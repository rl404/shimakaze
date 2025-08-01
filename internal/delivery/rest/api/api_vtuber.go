package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rl404/fairy/errors/stack"
	historyEntity "github.com/rl404/shimakaze/internal/domain/channel_stats_history/entity"
	"github.com/rl404/shimakaze/internal/domain/vtuber/entity"
	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/service"
	"github.com/rl404/shimakaze/internal/utils"
)

// @summary Get vtuber data.
// @tags Vtuber
// @produce json
// @param mode query string false "mode" enums(all, simple) default(all)
// @param names query string false "names"
// @param name query string false "name"
// @param original_name query string false "original name"
// @param nickname query string false "nickname"
// @param exclude_active query boolean false "exclude active"
// @param exclude_retired query boolean false "exclude retired"
// @param debut_day query integer false "start debut day"
// @param start_debut_month query integer false "start debut month"
// @param end_debut_month query integer false "end debut month"
// @param start_debut_year query integer false "start debut year"
// @param end_debut_year query integer false "end debut year"
// @param start_retired_month query integer false "start retired month"
// @param end_retired_month query integer false "end retired month"
// @param start_retired_year query integer false "start retired year"
// @param end_retired_year query integer false "end retired year"
// @param has_2d query boolean false "has 2d model"
// @param has_3d query boolean false "has 3d model"
// @param character_designer query string false "character designer"
// @param character_2d_modeler query string false "character 2d modeler"
// @param character_3d_modeler query string false "character 3d modeler"
// @param in_agency query boolean false "in agency"
// @param agency query string false "agency"
// @param agency_id query integer false "agency id"
// @param language_id query integer false "language id"
// @param channel_types query string false "channel types"
// @param birthday_day query integer false "birthday day"
// @param start_birthday_month query integer false "start birthday month"
// @param end_birthday_month query integer false "end birthday month"
// @param blood_types query string false "blood types"
// @param genders query string false "gender types"
// @param zodiacs query string false "zodiac types"
// @param start_subscriber query integer false "start subscriber"
// @param end_subscriber query integer false "end subscriber"
// @param start_video_count query integer false "start video count"
// @param end_video_count query integer false "end video count"
// @param sort query string false "sort" enums(name,-name,debut_date,-debut_date,retirement_date,-retirement_date,subscriber,-subscriber,monthly_subscriber,-monthly_subscriber,video_count,-video_count,average_video_length,-average_video_length,total_video_length,-total_video_length) default(name)
// @param page query integer false "page" default(1)
// @param limit query integer false "limit" default(20)
// @success 200 {object} utils.Response{data=[]service.vtuber}
// @failure 400 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /vtubers [get]
func (api *API) handleGetVtubers(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")
	names := r.URL.Query().Get("names")
	name := r.URL.Query().Get("name")
	originalName := r.URL.Query().Get("original_name")
	nickname := r.URL.Query().Get("nickname")
	excludeActive, _ := strconv.ParseBool(r.URL.Query().Get("exclude_active"))
	excludeRetired, _ := strconv.ParseBool(r.URL.Query().Get("exclude_retired"))
	debutDay, _ := strconv.Atoi(r.URL.Query().Get("debut_day"))
	startDebutMonth, _ := strconv.Atoi(r.URL.Query().Get("start_debut_month"))
	endDebutMonth, _ := strconv.Atoi(r.URL.Query().Get("end_debut_month"))
	startDebutYear, _ := strconv.Atoi(r.URL.Query().Get("start_debut_year"))
	endDebutYear, _ := strconv.Atoi(r.URL.Query().Get("end_debut_year"))
	startRetiredMonth, _ := strconv.Atoi(r.URL.Query().Get("start_retired_month"))
	endRetiredMonth, _ := strconv.Atoi(r.URL.Query().Get("end_retired_month"))
	startRetiredYear, _ := strconv.Atoi(r.URL.Query().Get("start_retired_year"))
	endRetiredYear, _ := strconv.Atoi(r.URL.Query().Get("end_retired_year"))
	has2D := utils.StrToPtrBool(r.URL.Query().Get("has_2d"))
	has3D := utils.StrToPtrBool(r.URL.Query().Get("has_3d"))
	characterDesigner := r.URL.Query().Get("character_designer")
	character2DModeler := r.URL.Query().Get("character_2d_modeler")
	character3DModeler := r.URL.Query().Get("character_3d_modeler")
	inAgency := utils.StrToPtrBool(r.URL.Query().Get("in_agency"))
	agency := r.URL.Query().Get("agency")
	agencyID, _ := strconv.ParseInt(r.URL.Query().Get("agency_id"), 10, 64)
	languageID, _ := strconv.ParseInt(r.URL.Query().Get("language_id"), 10, 64)
	channelTypes := utils.StrToStrSlice(r.URL.Query().Get("channel_types"))
	birthdayDay, _ := strconv.Atoi(r.URL.Query().Get("birthday_day"))
	startBirthdayMonth, _ := strconv.Atoi(r.URL.Query().Get("start_birthday_month"))
	endBirthdayMonth, _ := strconv.Atoi(r.URL.Query().Get("end_birthday_month"))
	bloodTypes := utils.StrToStrSlice(r.URL.Query().Get("blood_types"))
	genders := utils.StrToStrSlice(r.URL.Query().Get("genders"))
	zodiacs := utils.StrToStrSlice(r.URL.Query().Get("zodiacs"))
	startSubscriber, _ := strconv.Atoi(r.URL.Query().Get("start_subscriber"))
	endSubscriber, _ := strconv.Atoi(r.URL.Query().Get("end_subscriber"))
	startVideoCount, _ := strconv.Atoi(r.URL.Query().Get("start_video_count"))
	endVideoCount, _ := strconv.Atoi(r.URL.Query().Get("end_video_count"))
	sort := r.URL.Query().Get("sort")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	vtubers, pagination, code, err := api.service.GetVtubers(r.Context(), service.GetVtubersRequest{
		Mode:               entity.SearchMode(mode),
		Names:              names,
		Name:               name,
		OriginalName:       originalName,
		Nickname:           nickname,
		ExcludeActive:      excludeActive,
		ExcludeRetired:     excludeRetired,
		DebutDay:           debutDay,
		StartDebutMonth:    startDebutMonth,
		EndDebutMonth:      endDebutMonth,
		StartDebutYear:     startDebutYear,
		EndDebutYear:       endDebutYear,
		StartRetiredMonth:  startRetiredMonth,
		EndRetiredMonth:    endRetiredMonth,
		StartRetiredYear:   startRetiredYear,
		EndRetiredYear:     endRetiredYear,
		Has2D:              has2D,
		Has3D:              has3D,
		CharacterDesigner:  characterDesigner,
		Character2DModeler: character2DModeler,
		Character3DModeler: character3DModeler,
		InAgency:           inAgency,
		Agency:             agency,
		AgencyID:           agencyID,
		LanguageID:         languageID,
		ChannelTypes:       entity.StrsToChannelTypes(channelTypes),
		BirthdayDay:        birthdayDay,
		StartBirthdayMonth: startBirthdayMonth,
		EndBirthdayMonth:   endBirthdayMonth,
		BloodTypes:         bloodTypes,
		Genders:            genders,
		Zodiacs:            zodiacs,
		StartSubscriber:    startSubscriber,
		EndSubscriber:      endSubscriber,
		StartVideoCount:    startVideoCount,
		EndVideoCount:      endVideoCount,
		Sort:               sort,
		Page:               page,
		Limit:              limit,
	})

	utils.ResponseWithJSON(w, code, vtubers, stack.Wrap(r.Context(), err), pagination)
}

// @summary Get vtuber data.
// @tags Vtuber
// @produce json
// @param id path integer true "wikia id"
// @success 200 {object} utils.Response{data=service.vtuber}
// @failure 400 {object} utils.Response
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /vtubers/{id} [get]
func (api *API) handleGetVtuberByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, stack.Wrap(r.Context(), err, errors.ErrInvalidID))
		return
	}

	vtuber, code, err := api.service.GetVtuberByID(r.Context(), id)
	utils.ResponseWithJSON(w, code, vtuber, stack.Wrap(r.Context(), err))
}

// @summary Get vtuber channel histories.
// @tags Vtuber
// @produce json
// @param id path integer true "wikia id"
// @param start_date query string false "start date"
// @param end_date query string false "end date"
// @param group query string false "group" enums(DAILY,MONTHLY,YEARLY) default(MONTHLY)
// @success 200 {object} utils.Response{data=service.vtuberChannelHistory}
// @failure 400 {object} utils.Response
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /vtubers/{id}/channel-history [get]
func (api *API) handleGetVtuberChannelHistory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, stack.Wrap(r.Context(), err, errors.ErrInvalidID))
		return
	}

	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	group := r.URL.Query().Get("group")

	histories, code, err := api.service.GetVtuberChannelHistoriesByID(r.Context(), service.GetVtuberChannelHistoriesRequest{
		ID:        id,
		StartDate: startDate,
		EndDate:   endDate,
		Group:     historyEntity.Group(group),
	})

	utils.ResponseWithJSON(w, code, histories, stack.Wrap(r.Context(), err))
}

// @summary Get all vtuber images.
// @tags Vtuber
// @produce json
// @param shuffle query boolean false "shuffle"
// @param limit query integer false "limit"
// @success 200 {object} utils.Response{data=[]service.vtuberImage}
// @failure 500 {object} utils.Response
// @router /vtubers/images [get]
func (api *API) handleGetVtuberImages(w http.ResponseWriter, r *http.Request) {
	shuffle, _ := strconv.ParseBool(r.URL.Query().Get("shuffle"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	images, code, err := api.service.GetVtuberImages(r.Context(), shuffle, limit)
	utils.ResponseWithJSON(w, code, images, stack.Wrap(r.Context(), err))
}

// @summary Get vtuber family trees.
// @tags Vtuber
// @produce json
// @success 200 {object} utils.Response{data=service.vtuberFamilyTree}
// @failure 500 {object} utils.Response
// @router /vtubers/family-trees [get]
func (api *API) handleGetVtuberFamilyTrees(w http.ResponseWriter, r *http.Request) {
	tree, code, err := api.service.GetVtuberFamilyTrees(r.Context())
	utils.ResponseWithJSON(w, code, tree, stack.Wrap(r.Context(), err))
}

// @summary Get vtuber agency trees.
// @tags Vtuber
// @produce json
// @success 200 {object} utils.Response{data=service.vtuberAgencyTree}
// @failure 500 {object} utils.Response
// @router /vtubers/agency-trees [get]
func (api *API) handleGetVtuberAgencyTrees(w http.ResponseWriter, r *http.Request) {
	tree, code, err := api.service.GetVtuberAgencyTrees(r.Context())
	utils.ResponseWithJSON(w, code, tree, stack.Wrap(r.Context(), err))
}

// @summary Get vtuber character designers.
// @tags Vtuber
// @produce json
// @success 200 {object} utils.Response{data=[]string}
// @failure 500 {object} utils.Response
// @router /vtubers/character-designers [get]
func (api *API) handleGetVtuberCharacterDesigners(w http.ResponseWriter, r *http.Request) {
	designers, code, err := api.service.GetVtuberCharacterDesigners(r.Context())
	utils.ResponseWithJSON(w, code, designers, stack.Wrap(r.Context(), err))
}

// @summary Get vtuber character 2D modelers.
// @tags Vtuber
// @produce json
// @success 200 {object} utils.Response{data=[]string}
// @failure 500 {object} utils.Response
// @router /vtubers/2d-modelers [get]
func (api *API) handleGetVtuberCharacter2DModelers(w http.ResponseWriter, r *http.Request) {
	modelers, code, err := api.service.GetVtuberCharacter2DModelers(r.Context())
	utils.ResponseWithJSON(w, code, modelers, stack.Wrap(r.Context(), err))
}

// @summary Get vtuber character 3D modelers.
// @tags Vtuber
// @produce json
// @success 200 {object} utils.Response{data=[]string}
// @failure 500 {object} utils.Response
// @router /vtubers/3d-modelers [get]
func (api *API) handleGetVtuberCharacter3DModelers(w http.ResponseWriter, r *http.Request) {
	modelers, code, err := api.service.GetVtuberCharacter3DModelers(r.Context())
	utils.ResponseWithJSON(w, code, modelers, stack.Wrap(r.Context(), err))
}
