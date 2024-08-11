package main

import (
	"context"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	_nr "github.com/rl404/fairy/log/newrelic"
	nrPS "github.com/rl404/fairy/monitoring/newrelic/pubsub"
	"github.com/rl404/shimakaze/internal/delivery/cron"
	agencyRepository "github.com/rl404/shimakaze/internal/domain/agency/repository"
	agencyMongo "github.com/rl404/shimakaze/internal/domain/agency/repository/mongo"
	languageRepository "github.com/rl404/shimakaze/internal/domain/language/repository"
	languageMongo "github.com/rl404/shimakaze/internal/domain/language/repository/mongo"
	nonVtuberRepository "github.com/rl404/shimakaze/internal/domain/non_vtuber/repository"
	nonVtuberMongo "github.com/rl404/shimakaze/internal/domain/non_vtuber/repository/mongo"
	publisherRepository "github.com/rl404/shimakaze/internal/domain/publisher/repository"
	publisherPubsub "github.com/rl404/shimakaze/internal/domain/publisher/repository/pubsub"
	vtuberRepository "github.com/rl404/shimakaze/internal/domain/vtuber/repository"
	vtuberMongo "github.com/rl404/shimakaze/internal/domain/vtuber/repository/mongo"
	wikiaRepository "github.com/rl404/shimakaze/internal/domain/wikia/repository"
	wikiaClient "github.com/rl404/shimakaze/internal/domain/wikia/repository/client"
	"github.com/rl404/shimakaze/internal/service"
	"github.com/rl404/shimakaze/internal/utils"
	"github.com/rl404/shimakaze/pkg/pubsub"
)

func cronFill() error {
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

	// Init language.
	var language languageRepository.Repository = languageMongo.New(db)
	utils.Info("repository language initialized")

	// Init publisher.
	var publisher publisherRepository.Repository = publisherPubsub.New(ps, pubsubTopic)
	utils.Info("repository publisher initialized")

	// Init service.
	service := service.New(wikia, vtuber, nonVtuber, agency, language, publisher, nil, nil, nil, nil, nil, nil, nil, nil)
	utils.Info("service initialized")

	// Run cron.
	utils.Info("filling missing data...")
	if err := cron.New(service, nrApp).Fill(cfg.Cron.FillLimit); err != nil {
		return err
	}

	utils.Info("done")
	return nil
}
