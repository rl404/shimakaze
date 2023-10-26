package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	_nr "github.com/rl404/fairy/log/newrelic"
	nrCache "github.com/rl404/fairy/monitoring/newrelic/cache"
	nrPS "github.com/rl404/fairy/monitoring/newrelic/pubsub"
	_consumer "github.com/rl404/shimakaze/internal/delivery/consumer"
	agencyRepository "github.com/rl404/shimakaze/internal/domain/agency/repository"
	agencyMongo "github.com/rl404/shimakaze/internal/domain/agency/repository/mongo"
	bilibiliRepository "github.com/rl404/shimakaze/internal/domain/bilibili/repository"
	bilibiliClient "github.com/rl404/shimakaze/internal/domain/bilibili/repository/client"
	niconicoRepository "github.com/rl404/shimakaze/internal/domain/niconico/repository"
	niconicoClient "github.com/rl404/shimakaze/internal/domain/niconico/repository/client"
	nonVtuberRepository "github.com/rl404/shimakaze/internal/domain/non_vtuber/repository"
	nonVtuberMongo "github.com/rl404/shimakaze/internal/domain/non_vtuber/repository/mongo"
	publisherRepository "github.com/rl404/shimakaze/internal/domain/publisher/repository"
	publisherPubsub "github.com/rl404/shimakaze/internal/domain/publisher/repository/pubsub"
	twitchRepository "github.com/rl404/shimakaze/internal/domain/twitch/repository"
	twitchClient "github.com/rl404/shimakaze/internal/domain/twitch/repository/client"
	vtuberRepository "github.com/rl404/shimakaze/internal/domain/vtuber/repository"
	vtuberMongo "github.com/rl404/shimakaze/internal/domain/vtuber/repository/mongo"
	wikiaRepository "github.com/rl404/shimakaze/internal/domain/wikia/repository"
	wikiaClient "github.com/rl404/shimakaze/internal/domain/wikia/repository/client"
	youtubeRepository "github.com/rl404/shimakaze/internal/domain/youtube/repository"
	youtubeClient "github.com/rl404/shimakaze/internal/domain/youtube/repository/client"
	"github.com/rl404/shimakaze/internal/service"
	"github.com/rl404/shimakaze/internal/utils"
	"github.com/rl404/shimakaze/pkg/cache"
	"github.com/rl404/shimakaze/pkg/pubsub"
)

func consumer() error {
	// Get config.
	cfg, err := getConfig()
	if err != nil {
		return err
	}
	utils.Info("config initialized")

	// Init newrelic.
	nrApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(cfg.Newrelic.Name),
		newrelic.ConfigLicense(cfg.Newrelic.LicenseKey),
		newrelic.ConfigDistributedTracerEnabled(true),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		utils.Error(err.Error())
	} else {
		defer nrApp.Shutdown(10 * time.Second)
		utils.AddLog(_nr.NewFromNewrelicApp(nrApp, _nr.LogLevel(cfg.Log.Level)))
		utils.Info("newrelic initialized")
	}

	// Init in-memory.
	im, err := cache.New(cache.InMemory, "", "", time.Hour)
	if err != nil {
		return err
	}
	im = nrCache.New("inmemory", "inmemory", im)
	utils.Info("in-memory initialized")
	defer im.Close()

	// Init db.
	db, err := newDB(cfg.DB)
	if err != nil {
		return err
	}
	utils.Info("database initialized")
	defer db.Client().Disconnect(context.Background())

	// Init pubsub.
	ps, err := pubsub.New(pubsubType[cfg.PubSub.Dialect], cfg.PubSub.Address, cfg.PubSub.Password)
	if err != nil {
		return err
	}
	ps = nrPS.New(cfg.PubSub.Dialect, ps, nrApp)
	utils.Info("pubsub initialized")
	defer ps.Close()

	// Init wikia.
	var wikia wikiaRepository.Repository = wikiaClient.New()
	utils.Info("repository wikia initialized")

	// Init vtuber.
	var vtuber vtuberRepository.Repository = vtuberMongo.New(db, cfg.Cron.ActiveAge, cfg.Cron.RetiredAge)
	utils.Info("repository vtuber initialized")

	// Init non-vtuber.
	var nonVtuber nonVtuberRepository.Repository = nonVtuberMongo.New(db)
	utils.Info("repository non-vtuber initialized")

	// Init agency.
	var agency agencyRepository.Repository = agencyMongo.New(db, cfg.Cron.AgencyAge)
	utils.Info("repository agency initialized")

	// Init publisher.
	var publisher publisherRepository.Repository = publisherPubsub.New(ps, pubsubTopic)
	utils.Info("repository publisher initialized")

	// Init youtube.
	var youtube youtubeRepository.Repository = youtubeClient.New(cfg.Youtube.Key, cfg.Youtube.MaxAge)
	utils.Info("repository youtube initialized")

	// Init twitch.
	var twitch twitchRepository.Repository = twitchClient.New(im, cfg.Twitch.ClientID, cfg.Twitch.ClientSecret, cfg.Twitch.MaxAge)
	utils.Info("repository twitch initialized")

	// Init bilibili.
	var bilibili bilibiliRepository.Repository = bilibiliClient.New(cfg.Bilibili.MaxAge)
	utils.Info("repository bilibili initialized")

	// Init niconico.
	var niconico niconicoRepository.Repository = niconicoClient.New(cfg.Niconico.MaxAge)
	utils.Info("repository niconico initialized")

	// Init service.
	service := service.New(wikia, vtuber, nonVtuber, agency, publisher, youtube, twitch, bilibili, niconico)
	utils.Info("service initialized")

	// Init consumer.
	consumer := _consumer.New(service, ps, pubsubTopic)
	utils.Info("consumer initialized")
	defer consumer.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Start subscribe.
	if err := consumer.Subscribe(); err != nil {
		return err
	}

	utils.Info("consumer ready")
	<-sigChan

	return nil
}
