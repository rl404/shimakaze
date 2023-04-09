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

// @summary Get vtuber subscriber count.
// @tags Statistic
// @produce json
// @param interval query integer false "interval" default(100000)
// @param max query integer false "max" default(5000000)
// @success 200 {object} utils.Response{data=[]service.vtuberSubscriberCount}
// @failure 400 {object} utils.Response
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

// @summary Get vtuber character designer count.
// @tags Statistic
// @produce json
// @param top query integer false "top count" default(10)
// @success 200 {object} utils.Response{data=[]service.vtuberDesignerCount}
// @failure 400 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /statistics/vtubers/designer-count [get]
func (api *API) handleGetVtuberDesignerCount(w http.ResponseWriter, r *http.Request) {
	top, _ := strconv.Atoi(r.URL.Query().Get("top"))
	cnt, code, err := api.service.GetVtuberDesignerCount(r.Context(), service.GetVtuberDesignerCountRequest{Top: top})
	utils.ResponseWithJSON(w, code, cnt, errors.Wrap(r.Context(), err))
}

// @summary Get vtuber character 2d modeler count.
// @tags Statistic
// @produce json
// @param top query integer false "top count" default(10)
// @success 200 {object} utils.Response{data=[]service.vtuberDesignerCount}
// @failure 400 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /statistics/vtubers/2d-modeler-count [get]
func (api *API) handleGetVtuber2DModelerCount(w http.ResponseWriter, r *http.Request) {
	top, _ := strconv.Atoi(r.URL.Query().Get("top"))
	cnt, code, err := api.service.GetVtuber2DModelerCount(r.Context(), service.GetVtuberDesignerCountRequest{Top: top})
	utils.ResponseWithJSON(w, code, cnt, errors.Wrap(r.Context(), err))
}

// @summary Get vtuber character 3d modeler count.
// @tags Statistic
// @produce json
// @param top query integer false "top count" default(10)
// @success 200 {object} utils.Response{data=[]service.vtuberDesignerCount}
// @failure 400 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /statistics/vtubers/3d-modeler-count [get]
func (api *API) handleGetVtuber3DModelerCount(w http.ResponseWriter, r *http.Request) {
	top, _ := strconv.Atoi(r.URL.Query().Get("top"))
	cnt, code, err := api.service.GetVtuber3DModelerCount(r.Context(), service.GetVtuberDesignerCountRequest{Top: top})
	utils.ResponseWithJSON(w, code, cnt, errors.Wrap(r.Context(), err))
}

// @summary Get vtuber average video count.
// @tags Statistic
// @produce json
// @success 200 {object} utils.Response{data=float64}
// @failure 500 {object} utils.Response
// @router /statistics/vtubers/average-video-count [get]
func (api *API) handleGetVtuberAverageVideoCount(w http.ResponseWriter, r *http.Request) {
	avg, code, err := api.service.GetVtuberAverageVideoCount(r.Context())
	utils.ResponseWithJSON(w, code, avg, errors.Wrap(r.Context(), err))
}

// @summary Get vtuber average video duration.
// @tags Statistic
// @produce json
// @success 200 {object} utils.Response{data=float64}
// @failure 500 {object} utils.Response
// @router /statistics/vtubers/average-video-duration [get]
func (api *API) handleGetVtuberAverageVideoDuration(w http.ResponseWriter, r *http.Request) {
	avg, code, err := api.service.GetVtuberAverageVideoDuration(r.Context())
	utils.ResponseWithJSON(w, code, avg, errors.Wrap(r.Context(), err))
}

// @summary Get vtuber video count by date.
// @tags Statistic
// @produce json
// @param hourly query boolean false "hourly"
// @param daily query boolean false "daily"
// @success 200 {object} utils.Response{data=[]service.vtuberVideoCountByDate}
// @failure 400 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /statistics/vtubers/video-count-by-date [get]
func (api *API) handleGetVtuberVideoCountByDate(w http.ResponseWriter, r *http.Request) {
	hourly, _ := strconv.ParseBool(r.URL.Query().Get("hourly"))
	daily, _ := strconv.ParseBool(r.URL.Query().Get("daily"))
	cnt, code, err := api.service.GetVtuberVideoCountByDate(r.Context(), hourly, daily)
	utils.ResponseWithJSON(w, code, cnt, errors.Wrap(r.Context(), err))
}

// @summary Get vtuber video count.
// @tags Statistic
// @produce json
// @param top query integer false "top count" default(10)
// @success 200 {object} utils.Response{data=[]service.vtuberVideoCount}
// @failure 400 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /statistics/vtubers/video-count [get]
func (api *API) handleGetVtuberVideoCount(w http.ResponseWriter, r *http.Request) {
	top, _ := strconv.Atoi(r.URL.Query().Get("top"))
	cnt, code, err := api.service.GetVtuberVideoCount(r.Context(), service.GetVtuberVideoCountRequest{Top: top})
	utils.ResponseWithJSON(w, code, cnt, errors.Wrap(r.Context(), err))
}

// @summary Get vtuber video duration.
// @tags Statistic
// @produce json
// @param top query integer false "top count" default(10)
// @success 200 {object} utils.Response{data=[]service.vtuberVideoDuration}
// @failure 400 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /statistics/vtubers/video-duration [get]
func (api *API) handleGetVtuberVideoDuration(w http.ResponseWriter, r *http.Request) {
	top, _ := strconv.Atoi(r.URL.Query().Get("top"))
	cnt, code, err := api.service.GetVtuberVideoDuration(r.Context(), service.GetVtuberVideoDurationRequest{Top: top})
	utils.ResponseWithJSON(w, code, cnt, errors.Wrap(r.Context(), err))
}

// @summary Get vtuber birthday count.
// @tags Statistic
// @produce json
// @success 200 {object} utils.Response{data=[]service.vtuberBirthdayCount}
// @failure 500 {object} utils.Response
// @router /statistics/vtubers/birthday-count [get]
func (api *API) handleGetVtuberBirthdayCount(w http.ResponseWriter, r *http.Request) {
	cnt, code, err := api.service.GetVtuberBirthdayCount(r.Context())
	utils.ResponseWithJSON(w, code, cnt, errors.Wrap(r.Context(), err))
}
