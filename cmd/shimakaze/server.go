package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/newrelic/go-agent/v3/newrelic"
	_nr "github.com/rl404/fairy/log/newrelic"
	nrCache "github.com/rl404/fairy/monitoring/newrelic/cache"
	nrPS "github.com/rl404/fairy/monitoring/newrelic/pubsub"
	"github.com/rl404/shimakaze/internal/delivery/rest/api"
	"github.com/rl404/shimakaze/internal/delivery/rest/ping"
	"github.com/rl404/shimakaze/internal/delivery/rest/swagger"
	agencyRepository "github.com/rl404/shimakaze/internal/domain/agency/repository"
	agencyCache "github.com/rl404/shimakaze/internal/domain/agency/repository/cache"
	agencyMongo "github.com/rl404/shimakaze/internal/domain/agency/repository/mongo"
	channelStatsHistoryRepository "github.com/rl404/shimakaze/internal/domain/channel_stats_history/repository"
	channelStatsHistoryCache "github.com/rl404/shimakaze/internal/domain/channel_stats_history/repository/cache"
	channelStatsHistoryMongo "github.com/rl404/shimakaze/internal/domain/channel_stats_history/repository/mongo"
	languageRepository "github.com/rl404/shimakaze/internal/domain/language/repository"
	languageCache "github.com/rl404/shimakaze/internal/domain/language/repository/cache"
	languageMongo "github.com/rl404/shimakaze/internal/domain/language/repository/mongo"
	nonVtuberRepository "github.com/rl404/shimakaze/internal/domain/non_vtuber/repository"
	nonVtuberCache "github.com/rl404/shimakaze/internal/domain/non_vtuber/repository/cache"
	nonVtuberMongo "github.com/rl404/shimakaze/internal/domain/non_vtuber/repository/mongo"
	publisherRepository "github.com/rl404/shimakaze/internal/domain/publisher/repository"
	publisherPubsub "github.com/rl404/shimakaze/internal/domain/publisher/repository/pubsub"
	ssoRepository "github.com/rl404/shimakaze/internal/domain/sso/repository"
	ssoClient "github.com/rl404/shimakaze/internal/domain/sso/repository/client"
	tierListRepository "github.com/rl404/shimakaze/internal/domain/tier_list/repository"
	tierListCache "github.com/rl404/shimakaze/internal/domain/tier_list/repository/cache"
	tierListMongo "github.com/rl404/shimakaze/internal/domain/tier_list/repository/mongo"
	tokenRepository "github.com/rl404/shimakaze/internal/domain/token/repository"
	tokenToken "github.com/rl404/shimakaze/internal/domain/token/repository/cache"
	userRepository "github.com/rl404/shimakaze/internal/domain/user/repository"
	userCache "github.com/rl404/shimakaze/internal/domain/user/repository/cache"
	userMongo "github.com/rl404/shimakaze/internal/domain/user/repository/mongo"
	vtuberRepository "github.com/rl404/shimakaze/internal/domain/vtuber/repository"
	vtuberCache "github.com/rl404/shimakaze/internal/domain/vtuber/repository/cache"
	vtuberMongo "github.com/rl404/shimakaze/internal/domain/vtuber/repository/mongo"
	wikiaRepository "github.com/rl404/shimakaze/internal/domain/wikia/repository"
	wikiaClient "github.com/rl404/shimakaze/internal/domain/wikia/repository/client"
	"github.com/rl404/shimakaze/internal/service"
	"github.com/rl404/shimakaze/internal/utils"
	"github.com/rl404/shimakaze/pkg/cache"
	"github.com/rl404/shimakaze/pkg/http"
	"github.com/rl404/shimakaze/pkg/pubsub"
)

func server() error {
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

	// Init cache.
	c, err := cache.New(cacheType[cfg.Cache.Dialect], cfg.Cache.Address, cfg.Cache.Password, cfg.Cache.Time)
	if err != nil {
		return err
	}
	c = nrCache.New(cfg.Cache.Dialect, cfg.Cache.Address, c)
	utils.Info("cache initialized")
	defer c.Close()

	// Init in-memory.
	im, err := cache.New(cache.InMemory, "", "", 5*time.Second)
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
	var vtuber vtuberRepository.Repository
	vtuber = vtuberMongo.New(db, cfg.Cron.ActiveAge, cfg.Cron.RetiredAge)
	vtuber = vtuberCache.New(c, vtuber)
	vtuber = vtuberCache.New(im, vtuber)
	utils.Info("repository vtuber initialized")

	// Init non-vtuber.
	var nonVtuber nonVtuberRepository.Repository
	nonVtuber = nonVtuberMongo.New(db)
	nonVtuber = nonVtuberCache.New(c, nonVtuber)
	nonVtuber = nonVtuberCache.New(im, nonVtuber)
	utils.Info("repository non-vtuber initialized")

	// Init agency.
	var agency agencyRepository.Repository
	agency = agencyMongo.New(db, cfg.Cron.AgencyAge)
	agency = agencyCache.New(c, agency)
	agency = agencyCache.New(im, agency)
	utils.Info("repository agency initialized")

	// Init language.
	var language languageRepository.Repository
	language = languageMongo.New(db)
	language = languageCache.New(c, language)
	language = languageCache.New(im, language)
	utils.Info("repository language initialized")

	// Init channel stats history.
	var channelStatsHistory channelStatsHistoryRepository.Repository
	channelStatsHistory = channelStatsHistoryMongo.New(db)
	channelStatsHistory = channelStatsHistoryCache.New(c, channelStatsHistory)
	channelStatsHistory = channelStatsHistoryCache.New(im, channelStatsHistory)
	utils.Info("repository channel-stats-history initialized")

	// Init publisher.
	var publisher publisherRepository.Repository = publisherPubsub.New(ps, pubsubTopic)
	utils.Info("repository publisher initialized")

	// Init sso.
	var sso ssoRepository.Repository = ssoClient.New(cfg.SSO.ClientID, cfg.SSO.ClientSecret, cfg.SSO.RedirectURL)
	utils.Info("repository sso initialized")

	// Init user.
	var user userRepository.Repository
	user = userMongo.New(db)
	user = userCache.New(c, user)
	user = userCache.New(im, user)
	utils.Info("repository user initialized")

	// Init token.
	var token tokenRepository.Repository = tokenToken.New(c, cfg.JWT.AccessSecret, cfg.JWT.AccessExpired, cfg.JWT.RefreshSecret, cfg.JWT.RefreshExpired)
	utils.Info("repository token initialized")

	// Init tier-list.
	var tierList tierListRepository.Repository
	tierList = tierListMongo.New(db)
	tierList = tierListCache.New(c, tierList)
	tierList = tierListCache.New(im, tierList)
	utils.Info("repository tier-list initialized")

	// Init service.
	service := service.New(wikia, vtuber, nonVtuber, agency, language, channelStatsHistory, publisher, nil, nil, nil, nil, sso, user, token, tierList)
	utils.Info("service initialized")

	// Init web server.
	httpServer := http.New(http.Config{
		Port:            cfg.HTTP.Port,
		ReadTimeout:     cfg.HTTP.ReadTimeout,
		WriteTimeout:    cfg.HTTP.WriteTimeout,
		GracefulTimeout: cfg.HTTP.GracefulTimeout,
	})
	utils.Info("http server initialized")

	r := httpServer.Router()
	r.Use(middleware.RealIP)
	r.Use(utils.Recoverer)
	utils.Info("http server middleware initialized")

	// Register ping route.
	ping.New().Register(r)
	utils.Info("http route ping initialized")

	// Register swagger route.
	swagger.New().Register(r)
	utils.Info("http route swagger initialized")

	// Register api route.
	api.New(service, cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret).Register(r, nrApp)
	utils.Info("http route api initialized")

	// Run web server.
	httpServerChan := httpServer.Run()
	utils.Info("http server listening at :%s", cfg.HTTP.Port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case err := <-httpServerChan:
		if err != nil {
			return err
		}
	case <-sigChan:
	}

	return nil
}
