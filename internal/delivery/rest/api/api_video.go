package api

import (
	"net/http"
	"strconv"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/service"
	"github.com/rl404/shimakaze/internal/utils"
)

// @summary Get videos.
// @tags Video
// @produce json
// @param start_date query string false "start date"
// @param end_date query string false "end date"
// @param is_finished query boolean false "is finished"
// @param sort query string false "sort" enums(video_start_date,-video_start_date) default(-video_start_date)
// @param page query integer false "page" default(1)
// @param limit query integer false "limit" default(20)
// @success 200 {object} utils.Response{data=[]service.video}
// @failure 400 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /videos [get]
func (api *API) handleGetVideos(w http.ResponseWriter, r *http.Request) {
	isFinished := utils.StrToPtrBool(r.URL.Query().Get("is_finished"))
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	sort := r.URL.Query().Get("sort")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	videos, pagination, code, err := api.service.GetVideos(r.Context(), service.GetVideosRequest{
		StartDate:  startDate,
		EndDate:    endDate,
		IsFinished: isFinished,
		Sort:       sort,
		Page:       page,
		Limit:      limit,
	})

	utils.ResponseWithJSON(w, code, videos, stack.Wrap(r.Context(), err), pagination)
}
