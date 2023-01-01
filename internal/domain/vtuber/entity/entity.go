package entity

import "time"

// Vtuber is entity for vtuber.
type Vtuber struct {
	ID                  int64
	Name                string
	Image               string
	OriginalNames       []string
	Nicknames           []string
	Caption             string
	DebutDate           *time.Time
	RetirementDate      *time.Time
	Has2D               bool
	Has3D               bool
	CharacterDesigners  []string
	Character2DModelers []string
	Character3DModelers []string
	Agencies            []string
	Affiliations        []string
	Channels            []Channel
	SocialMedias        []string
	OfficialWebsites    []string
	Gender              string
	Age                 *float64
	Birthday            *time.Time
	Height              *float64
	Weight              *float64
	BloodType           string
	ZodiacSign          string
	Emoji               string
}

// ChannelType is channel types.
type ChannelType string

// Available channel types.
const (
	ChannelYoutube  ChannelType = "YOUTUBE"
	ChannelTwitch   ChannelType = "TWITCH"
	ChannelBilibili ChannelType = "BILIBILI"
	ChannelNiconico ChannelType = "NICONICO"
	ChannelOther    ChannelType = "OTHER"
)

// Channel is entity for channel.
type Channel struct {
	Name       string
	Type       ChannelType
	URL        string
	Image      string
	Subscriber int
	Videos     []Video
}

// Video is entity for video.
type Video struct {
	Title     string
	URL       string
	Image     string
	StartDate time.Time
	EndDate   time.Time
}
