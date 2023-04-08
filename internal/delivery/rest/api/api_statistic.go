package api

import (
	"net/http"
	"strconv"

	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/service"
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

// @summary Get vtuber debut & retire count yearly.
// @tags Statistic
// @produce json
// @success 200 {object} utils.Response{data=[]service.vtuberDebutRetireCount}
// @failure 500 {object} utils.Response
// @router /statistics/vtubers/debut-retire-count-yearly [get]
func (api *API) handleGetVtuberDebutRetireCountYearly(w http.ResponseWriter, r *http.Request) {
	cnt, code, err := api.service.GetVtuberDebutRetireCountYearly(r.Context())
	utils.ResponseWithJSON(w, code, cnt, errors.Wrap(r.Context(), err))
}

// @summary Get vtuber 2d & 3d model count.
// @tags Statistic
// @produce json
// @success 200 {object} utils.Response{data=[]service.vtuberModelCount}
// @failure 500 {object} utils.Response
// @router /statistics/vtubers/model-count [get]
func (api *API) handleGetVtuberModelCount(w http.ResponseWriter, r *http.Request) {
	cnt, code, err := api.service.GetVtuberModelCount(r.Context())
	utils.ResponseWithJSON(w, code, cnt, errors.Wrap(r.Context(), err))
}

// @summary Get vtuber in agency count.
// @tags Statistic
// @produce json
// @success 200 {object} utils.Response{data=service.vtuberInAgencyCount}
// @failure 500 {object} utils.Response
// @router /statistics/vtubers/in-agency-count [get]
func (api *API) handleGetVtuberInAgencyCount(w http.ResponseWriter, r *http.Request) {
	cnt, code, err := api.service.GetVtuberInAgencyCount(r.Context())
	utils.ResponseWithJSON(w, code, cnt, errors.Wrap(r.Context(), err))
}

// @summary Get vtuber susbcriber count.
// @tags Statistic
// @produce json
// @param interval query integer false "interval" default(100000)
// @param max query integer false "max" default(5000000)
// @success 200 {object} utils.Response{data=[]service.vtuberSubscriberCount}
// @failure 500 {object} utils.Response
// @router /statistics/vtubers/subscriber-count [get]
func (api *API) handleGetVtuberSubscriberCount(w http.ResponseWriter, r *http.Request) {
	interval, _ := strconv.Atoi(r.URL.Query().Get("interval"))
	max, _ := strconv.Atoi(r.URL.Query().Get("max"))
	cnt, code, err := api.service.GetVtuberSubscriberCount(r.Context(), service.GetVtuberSubscriberCountRequest{
		Interval: interval,
		Max:      max,
	})
	utils.ResponseWithJSON(w, code, cnt, errors.Wrap(r.Context(), err))
}
