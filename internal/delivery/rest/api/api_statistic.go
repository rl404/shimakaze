package api

import (
	"net/http"

	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/utils"
)

// @summary Get vtuber count.
// @tags Statistic
// @produce json
// @success 200 {object} utils.Response{data=int}
// @failure 500 {object} utils.Response
// @router /statistics/vtubers/count [get]
func (api *API) handleGetVtuberCount(w http.ResponseWriter, r *http.Request) {
	count, code, err := api.service.GetVtuberCount(r.Context())
	utils.ResponseWithJSON(w, code, count, errors.Wrap(r.Context(), err))
}

// @summary Get agency count.
// @tags Statistic
// @produce json
// @success 200 {object} utils.Response{data=int}
// @failure 500 {object} utils.Response
// @router /statistics/agencies/count [get]
func (api *API) handleGetAgencyCount(w http.ResponseWriter, r *http.Request) {
	count, code, err := api.service.GetAgencyCount(r.Context())
	utils.ResponseWithJSON(w, code, count, errors.Wrap(r.Context(), err))
}

// @summary Get vtuber average active time.
// @tags Statistic
// @produce json
// @success 200 {object} utils.Response{data=float64}
// @failure 500 {object} utils.Response
// @router /statistics/vtubers/average-active-time [get]
func (api *API) handleGetVtuberAverageActiveTime(w http.ResponseWriter, r *http.Request) {
	avg, code, err := api.service.GetVtuberAverageActiveTime(r.Context())
	utils.ResponseWithJSON(w, code, avg, errors.Wrap(r.Context(), err))
}

// @summary Get vtuber status count.
// @tags Statistic
// @produce json
// @success 200 {object} utils.Response{data=service.vtuberStatusCount}
// @failure 500 {object} utils.Response
// @router /statistics/vtubers/count-by-status [get]
func (api *API) handleGetVtuberStatusCount(w http.ResponseWriter, r *http.Request) {
	cnt, code, err := api.service.GetVtuberStatusCount(r.Context())
	utils.ResponseWithJSON(w, code, cnt, errors.Wrap(r.Context(), err))
}

// @summary Get vtuber debut & retire count monthly.
// @tags Statistic
// @produce json
// @success 200 {object} utils.Response{data=[]service.vtuberDebutRetireCount}
// @failure 500 {object} utils.Response
// @router /statistics/vtubers/debut-retire-count-monthly [get]
func (api *API) handleGetVtuberDebutRetireCountMonthly(w http.ResponseWriter, r *http.Request) {
	cnt, code, err := api.service.GetVtuberDebutRetireCountMonthly(r.Context())
	utils.ResponseWithJSON(w, code, cnt, errors.Wrap(r.Context(), err))
}
