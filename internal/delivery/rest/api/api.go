package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/log"
	"github.com/rl404/fairy/monitoring/newrelic/middleware"
	"github.com/rl404/shimakaze/internal/service"
	"github.com/rl404/shimakaze/internal/utils"
)

// API contains all functions for api endpoints.
type API struct {
	service       service.Service
	accessSecret  string
	refreshSecret string
}

// New to create new api endpoints.
func New(service service.Service, accessSecret, refreshSecret string) *API {
	return &API{
		service:       service,
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
	}
}

// Register to register api routes.
func (api *API) Register(r chi.Router, nrApp *newrelic.Application) {
	r.Route("/", func(r chi.Router) {
		r.Use(middleware.NewHTTP(nrApp))
		r.Use(log.HTTPMiddlewareWithLog(utils.GetLogger(0), log.APIMiddlewareConfig{Error: true}))
		r.Use(log.HTTPMiddlewareWithLog(utils.GetLogger(1), log.APIMiddlewareConfig{
			RequestHeader:  true,
			RequestBody:    true,
			ResponseHeader: true,
			ResponseBody:   true,
			RawPath:        true,
			Error:          true,
		}))
		r.Use(utils.Recoverer)

		r.Get("/wikia/image/*", api.maxConcurrent(api.handleGetWikiaImage, 5))

		r.Post("/auth/callback", api.handleAuthCallback)
		r.Post("/auth/token/refresh", api.jwtAuth(api.handleTokenRefresh, tokenRefresh))
		r.Post("/auth/logout", api.jwtAuth(api.handleLogout))

		r.Get("/user", api.jwtAuth(api.handleGetProfile))

		r.Get("/tier-lists", api.handleGetTierLists)
		r.Put("/tier-lists", api.jwtAuth(api.handleUpsertTierListByID))
		r.Get("/tier-lists/{id}", api.handleGetTierListByID)
		r.Delete("/tier-lists/{id}", api.jwtAuth(api.handleDeleteTierListByID))

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
		r.Get("/statistics/vtubers/subscriber-count", api.handleGetVtuberSubscriberCount)
		r.Get("/statistics/vtubers/designer-count", api.handleGetVtuberDesignerCount)
		r.Get("/statistics/vtubers/2d-modeler-count", api.handleGetVtuber2DModelerCount)
		r.Get("/statistics/vtubers/3d-modeler-count", api.handleGetVtuber3DModelerCount)
		r.Get("/statistics/vtubers/average-video-count", api.handleGetVtuberAverageVideoCount)
		r.Get("/statistics/vtubers/average-video-duration", api.handleGetVtuberAverageVideoDuration)
		r.Get("/statistics/vtubers/video-count-by-date", api.handleGetVtuberVideoCountByDate)
		r.Get("/statistics/vtubers/video-count", api.handleGetVtuberVideoCount)
		r.Get("/statistics/vtubers/video-duration", api.handleGetVtuberVideoDuration)
		r.Get("/statistics/vtubers/birthday-count", api.handleGetVtuberBirthdayCount)
		r.Get("/statistics/vtubers/average-height", api.handleGetVtuberAverageHeight)
		r.Get("/statistics/vtubers/average-weight", api.handleGetVtuberAverageWeight)
		r.Get("/statistics/vtubers/blood-type-count", api.handleGetVtuberBloodTypeCount)
		r.Get("/statistics/vtubers/channel-type-count", api.handleGetVtuberChannelTypeCount)
		r.Get("/statistics/vtubers/gender-count", api.handleGetVtuberGenderCount)
		r.Get("/statistics/vtubers/zodiac-count", api.handleGetVtuberZodiacCount)

		r.Get("/statistics/agencies/count", api.handleGetAgencyCount)
	})
}
