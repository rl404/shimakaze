package api

import (
	"github.com/go-chi/chi"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/log"
	"github.com/rl404/fairy/monitoring/newrelic/middleware"
	"github.com/rl404/shimakaze/internal/service"
	"github.com/rl404/shimakaze/internal/utils"
)

// API contains all functions for api endpoints.
type API struct {
	service service.Service
}

// New to create new api endpoints.
func New(service service.Service) *API {
	return &API{
		service: service,
	}
}

// Register to register api routes.
func (api *API) Register(r chi.Router, nrApp *newrelic.Application) {
	r.Route("/", func(r chi.Router) {
		r.Use(middleware.NewHTTP(nrApp))
		r.Use(log.MiddlewareWithLog(utils.GetLogger(0), log.MiddlewareConfig{Error: true}))
		r.Use(log.MiddlewareWithLog(utils.GetLogger(1), log.MiddlewareConfig{
			RequestHeader:  true,
			RequestBody:    true,
			ResponseHeader: true,
			ResponseBody:   true,
			RawPath:        true,
			Error:          true,
		}))
		r.Use(utils.Recoverer)

		r.Get("/wikia/image/*", api.maxConcurrent(api.handleGetWikiaImage, 5))

		r.Get("/vtubers", api.handleGetVtubers)
		r.Get("/vtubers/{id}", api.handleGetVtuberByID)
		r.Get("/vtubers/images", api.handleGetVtuberImages)
		r.Get("/vtubers/family-trees", api.handleGetVtuberFamilyTrees)
		r.Get("/vtubers/agency-trees", api.handleGetVtuberAgencyTrees)
		r.Get("/vtubers/character-designers", api.handleGetVtuberCharacterDesigners)
		r.Get("/vtubers/character-2d-modelers", api.handleGetVtuberCharacter2DModelers)
		r.Get("/vtubers/character-3d-modelers", api.handleGetVtuberCharacter3DModelers)

		r.Get("/agencies", api.handleGetAgencies)
		r.Get("/agencies/{id}", api.handleGetAgencyByID)

		r.Get("/statistics/vtubers/count", api.handleGetVtuberCount)
		r.Get("/statistics/vtubers/average-active-time", api.handleGetVtuberAverageActiveTime)
		r.Get("/statistics/vtubers/status-count", api.handleGetVtuberStatusCount)
		r.Get("/statistics/vtubers/debut-retire-count-monthly", api.handleGetVtuberDebutRetireCountMonthly)
		r.Get("/statistics/vtubers/debut-retire-count-yearly", api.handleGetVtuberDebutRetireCountYearly)
		r.Get("/statistics/vtubers/model-count", api.handleGetVtuberModelCount)
		r.Get("/statistics/vtubers/in-agency-count", api.handleGetVtuberInAgencyCount)

		r.Get("/statistics/agencies/count", api.handleGetAgencyCount)
	})
}
