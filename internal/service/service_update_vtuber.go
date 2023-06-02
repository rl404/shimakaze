package service

import (
	"context"
	_errors "errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/rl404/shimakaze/internal/domain/agency/entity"
	vtuberEntity "github.com/rl404/shimakaze/internal/domain/vtuber/entity"
	wikiaEntity "github.com/rl404/shimakaze/internal/domain/wikia/entity"
	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/utils"
)

func (s *service) updateVtuber(ctx context.Context, id int64) (int, error) {
	// Call wikia api.
	page, code, err := s.wikia.GetPageByID(ctx, id)
	if err != nil {
		if code == http.StatusNotFound {
			// Delete existing vtuber.
			if code, err := s.vtuber.DeleteByID(ctx, id); err != nil {
				return code, errors.Wrap(ctx, err)
			}

			// Insert to non-vtuber.
			if code, err := s.nonVtuber.Create(ctx, id); err != nil {
				return code, errors.Wrap(ctx, err)
			}
			return http.StatusOK, nil
		}
		return code, errors.Wrap(ctx, err)
	}

	// Non-vtuber page.
	if s.isNonVtuberPage(*page) {
		// Delete existing vtuber.
		if code, err := s.vtuber.DeleteByID(ctx, id); err != nil {
			return code, errors.Wrap(ctx, err)
		}

		// Insert to non-vtuber.
		if code, err := s.nonVtuber.Create(ctx, id); err != nil {
			return code, errors.Wrap(ctx, err)
		}

		return http.StatusOK, nil
	}

	// Fill vtuber data.
	vtuber := vtuberEntity.WikiaPageToVtuber(*page)

	// Get image.
	vtuber.Image = s.getVtuberImage(ctx, id)

	// Get agencies.
	agencyMap := s.getAgencyMap(ctx)
	agencyFromAffiliation := s.getAgencyFromAffiliation(vtuber.Affiliations, agencyMap)

	// Get categories.
	category := s.getVtuberCategory(ctx, id, agencyMap)
	vtuber.Has2D = category.has2D
	vtuber.Has3D = category.has3D
	vtuber.Agencies = s.mergeAgencies(agencyFromAffiliation, category.agencies)
	vtuber.CharacterDesigners = category.charDesigner
	vtuber.Character2DModelers = category.char2DModeler
	vtuber.Character3DModelers = category.char3DModeler

	// Fill channel data.
	vtuber.Channels = s.fillChannelData(ctx, vtuber.Channels)

	// Update data.
	if code, err := s.vtuber.UpdateByID(ctx, id, vtuber); err != nil {
		return code, errors.Wrap(ctx, err)
	}

	return http.StatusOK, nil
}

func (s *service) isNonVtuberPage(page wikiaEntity.Page) bool {
	return strings.Contains(page.Content, "#REDIRECT") ||
		!strings.Contains(page.Content, "{{Character\n|") ||
		strings.Contains(page.Title, "/Gallery") ||
		strings.Contains(page.Title, "/Discography")
}

func (s *service) getVtuberImage(ctx context.Context, id int64) string {
	pageImage, _, err := s.wikia.GetPageImageByID(ctx, id)
	if err != nil {
		errors.Wrap(ctx, err)
		return ""
	}
	return pageImage.Image
}

func (s *service) getAgencyMap(ctx context.Context) map[string]vtuberEntity.Agency {
	agencies, _, _, err := s.agency.GetAll(ctx, entity.GetAllRequest{Page: 1, Limit: -1})
	if err != nil {
		errors.Wrap(ctx, err)
		return nil
	}

	agencyMap := make(map[string]vtuberEntity.Agency)
	for _, a := range agencies {
		agencyMap[strings.ToLower(a.Name)] = vtuberEntity.Agency{
			ID:    a.ID,
			Name:  a.Name,
			Image: a.Image,
		}
	}

	return agencyMap
}

type vtuberCategory struct {
	has2D         bool
	has3D         bool
	agencies      []vtuberEntity.Agency
	charDesigner  []string
	char2DModeler []string
	char3DModeler []string
}

func (s *service) getVtuberCategory(ctx context.Context, id int64, agencyMap map[string]vtuberEntity.Agency) (vtuberCategory vtuberCategory) {
	// Loop and map categories.
	var lastTitle string
	limitPerPage := 500
	for {
		pageCategories, nextTitle, _, err := s.wikia.GetPageCategories(ctx, id, limitPerPage, lastTitle)
		if err != nil {
			errors.Wrap(ctx, err)
			return
		}

		lastTitle = nextTitle

		for _, pageCategory := range pageCategories {
			split := strings.Split(pageCategory.Title, ":")
			if len(split) < 2 {
				continue
			}

			category := strings.Join(split[1:], ":")

			if category == "2D" || category == "Live2D" {
				vtuberCategory.has2D = true
			}

			if category == "3D" {
				vtuberCategory.has3D = true
			}

			if v, ok := agencyMap[strings.ToLower(category)]; ok {
				vtuberCategory.agencies = append(vtuberCategory.agencies, v)
			}

			if designedBy := strings.Split(category, "Designed by "); len(designedBy) > 1 {
				vtuberCategory.charDesigner = append(vtuberCategory.charDesigner, designedBy[1])
			}

			if modeled2DBy := strings.Split(category, "Live2D by "); len(modeled2DBy) > 1 {
				vtuberCategory.char2DModeler = append(vtuberCategory.char2DModeler, modeled2DBy[1])
			}

			if modeled3DBy := strings.Split(category, "3D by "); len(modeled3DBy) > 1 {
				vtuberCategory.char3DModeler = append(vtuberCategory.char3DModeler, modeled3DBy[1])
			}
		}

		if len(pageCategories) == 0 || lastTitle == "" {
			return
		}
	}
}

func (s *service) getAgencyFromAffiliation(affiliations []string, agencyMap map[string]vtuberEntity.Agency) []vtuberEntity.Agency {
	var res []vtuberEntity.Agency
	for _, a := range affiliations {
		if v, ok := agencyMap[strings.ToLower(a)]; ok {
			res = append(res, v)
		}
	}
	return res
}

func (s *service) mergeAgencies(a1, a2 []vtuberEntity.Agency) []vtuberEntity.Agency {

	agencyMap := make(map[int64]vtuberEntity.Agency)
	for _, a := range a1 {
		agencyMap[a.ID] = a
	}

	for _, a := range a2 {
		agencyMap[a.ID] = a
	}

	var a3 []vtuberEntity.Agency
	for _, a := range agencyMap {
		a3 = append(a3, a)
	}

	return a3
}

func (s *service) fillChannelData(ctx context.Context, channels []vtuberEntity.Channel) []vtuberEntity.Channel {
	for i, channel := range channels {
		switch channel.Type {
		case vtuberEntity.ChannelYoutube:
			channels[i] = s.fillYoutubeChannel(ctx, channels[i])
			channels[i] = s.fillYoutubeVideos(ctx, channels[i])
		case vtuberEntity.ChannelTwitch:
			channels[i] = s.fillTwitchChannel(ctx, channels[i])
			channels[i] = s.fillTwitchVideo(ctx, channels[i])
		case vtuberEntity.ChannelBilibili:
			// channels[i] = s.fillBilibiliChannel(ctx, channels[i])
			// channels[i] = s.fillBilibiliVideo(ctx, channels[i])
		case vtuberEntity.ChannelNiconico:
			channels[i] = s.fillNiconicoChannel(ctx, channels[i])
			channels[i] = s.fillNiconicoVideo(ctx, channels[i])
		}
	}

	return channels
}

func (s *service) fillYoutubeChannel(ctx context.Context, channel vtuberEntity.Channel) vtuberEntity.Channel {
	channelID := utils.GetLastPathFromURL(channel.URL)
	if channelID == "" {
		return channel
	}

	ch, _, err := s.youtube.GetChannelByID(ctx, channelID)
	if err == nil {
		channel.ID = ch.ID
		channel.Name = ch.Name
		channel.Image = ch.Image
		channel.Subscriber = ch.Subscriber
		return channel
	}

	// URL not contain channel id.
	if !_errors.Is(err, errors.ErrChannelNotFound) {
		errors.Wrap(ctx, err)
	}

	channelID, _, err = s.youtube.GetChannelIDByURL(ctx, channel.URL)
	if err != nil {
		errors.Wrap(ctx, err)
		return channel
	}

	ch, _, err = s.youtube.GetChannelByID(ctx, channelID)
	if err != nil {
		errors.Wrap(ctx, err)
		return channel
	}

	channel.ID = ch.ID
	channel.Name = ch.Name
	channel.Image = ch.Image
	channel.Subscriber = ch.Subscriber

	return channel
}

func (s *service) fillYoutubeVideos(ctx context.Context, channel vtuberEntity.Channel) vtuberEntity.Channel {
	if channel.ID == "" {
		return channel
	}

	videoIDs, _, err := s.youtube.GetVideoIDsByChannelID(ctx, channel.ID)
	if err != nil {
		errors.Wrap(ctx, err)
		return channel
	}

	videos, _, err := s.youtube.GetVideosByIDs(ctx, videoIDs)
	if err != nil {
		errors.Wrap(ctx, err)
		return channel
	}

	res := make([]vtuberEntity.Video, len(videos))
	for i, v := range videos {
		res[i] = vtuberEntity.Video{
			ID:        v.ID,
			Title:     v.Title,
			URL:       fmt.Sprintf("https://www.youtube.com/watch?v=%s", v.ID),
			Image:     v.Image,
			StartDate: v.StartDate,
			EndDate:   v.EndDate,
		}
	}

	channel.Videos = res

	return channel
}

func (s *service) fillTwitchChannel(ctx context.Context, channel vtuberEntity.Channel) vtuberEntity.Channel {
	username := utils.GetLastPathFromURL(channel.URL)
	if username == "" {
		return channel
	}

	user, _, err := s.twitch.GetUser(ctx, username)
	if err != nil {
		errors.Wrap(ctx, err)
		return channel
	}

	channel.ID = user.ID
	channel.Name = user.Name
	channel.Image = user.Image

	follower, _, err := s.twitch.GetFollowerCount(ctx, user.ID)
	if err != nil {
		errors.Wrap(ctx, err)
		return channel
	}

	channel.Subscriber = follower

	return channel
}

func (s *service) fillTwitchVideo(ctx context.Context, channel vtuberEntity.Channel) vtuberEntity.Channel {
	if channel.ID == "" {
		return channel
	}

	videos, _, err := s.twitch.GetVideos(ctx, channel.ID)
	if err != nil {
		errors.Wrap(ctx, err)
		return channel
	}

	res := make([]vtuberEntity.Video, len(videos))
	for i, v := range videos {
		res[i] = vtuberEntity.Video{
			ID:        v.ID,
			Title:     v.Title,
			URL:       v.URL,
			Image:     v.Image,
			StartDate: v.StartDate,
			EndDate:   v.EndDate,
		}
	}

	channel.Videos = res

	return channel
}

func (s *service) fillBilibiliChannel(ctx context.Context, channel vtuberEntity.Channel) vtuberEntity.Channel {
	userID := utils.GetLastPathFromURL(channel.URL)
	if userID == "" {
		return channel
	}

	user, _, err := s.bilibili.GetUser(ctx, userID)
	if err != nil {
		errors.Wrap(ctx, err)
		return channel
	}

	channel.ID = user.ID
	channel.Name = user.Name
	channel.Image = user.Image

	follower, _, err := s.bilibili.GetFollowerCount(ctx, user.ID)
	if err != nil {
		errors.Wrap(ctx, err)
		return channel
	}

	channel.Subscriber = follower

	return channel
}

func (s *service) fillBilibiliVideo(ctx context.Context, channel vtuberEntity.Channel) vtuberEntity.Channel {
	if channel.ID == "" {
		return channel
	}

	videos, _, err := s.bilibili.GetVideos(ctx, channel.ID)
	if err != nil {
		errors.Wrap(ctx, err)
		return channel
	}

	res := make([]vtuberEntity.Video, len(videos))
	for i, v := range videos {
		res[i] = vtuberEntity.Video{
			ID:        v.ID,
			Title:     v.Title,
			URL:       fmt.Sprintf("https://www.bilibili.com/video/%s", v.ID),
			Image:     v.Image,
			StartDate: v.StartDate,
			EndDate:   v.EndDate,
		}
	}

	channel.Videos = res

	return channel
}

func (s *service) fillNiconicoChannel(ctx context.Context, channel vtuberEntity.Channel) vtuberEntity.Channel {
	if channel.URL == "" {
		return channel
	}

	user, _, err := s.niconico.GetUser(ctx, channel.URL)
	if err != nil {
		errors.Wrap(ctx, err)
		return channel
	}

	channel.ID = user.ID
	channel.Name = user.Name
	channel.Image = user.Image
	channel.Subscriber = user.Subscriber

	return channel
}

func (s *service) fillNiconicoVideo(ctx context.Context, channel vtuberEntity.Channel) vtuberEntity.Channel {
	if channel.ID == "" {
		return channel
	}

	videos, _, err := s.niconico.GetVideos(ctx, channel.ID)
	if err != nil {
		errors.Wrap(ctx, err)
		return channel
	}

	res := make([]vtuberEntity.Video, len(videos))
	for i, v := range videos {
		res[i] = vtuberEntity.Video{
			ID:        v.ID,
			Title:     v.Title,
			URL:       fmt.Sprintf("https://www.nicovideo.jp/watch/%s", v.ID),
			Image:     v.Image,
			StartDate: v.StartDate,
			EndDate:   v.EndDate,
		}
	}

	channel.Videos = res

	return channel
}
