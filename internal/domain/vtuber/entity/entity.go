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
	Agencies            []Agency
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
	UpdatedAt           time.Time
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

// Agency is entity for agency.
type Agency struct {
	ID    int64
	Name  string
	Image string
}

// Channel is entity for channel.
type Channel struct {
	ID         string
	Name       string
	Type       ChannelType
	URL        string
	Image      string
	Subscriber int
	Videos     []Video
}

// Video is entity for video.
type Video struct {
	ID        string
	Title     string
	URL       string
	Image     string
	StartDate *time.Time
	EndDate   *time.Time
}

// SearchMode is search mode.
type SearchMode string

// Available search mode.
const (
	SearchModeAll    SearchMode = "all"
	SearchModeSimple SearchMode = "simple"
)

// GetAllRequest is get all request model.
type GetAllRequest struct {
	Mode               SearchMode
	Names              string
	Name               string
	OriginalName       string
	Nickname           string
	ExcludeActive      bool
	ExcludeRetired     bool
	StartDebutMonth    int
	EndDebutMonth      int
	StartDebutYear     int
	EndDebutYear       int
	StartRetiredMonth  int
	EndRetiredMonth    int
	StartRetiredYear   int
	EndRetiredYear     int
	Has2D              *bool
	Has3D              *bool
	CharacterDesigner  string
	Character2DModeler string
	Character3DModeler string
	InAgency           *bool
	Agency             string
	AgencyID           int64
	ChannelTypes       []ChannelType
	BirthdayDay        int
	StartBirthdayMonth int
	EndBirthdayMonth   int
	BloodTypes         []string
	Genders            []string
	Zodiacs            []string
	StartSubscriber    int
	EndSubscriber      int
	Sort               string
	Page               int
	Limit              int
}

// StatusCount is entity for status count.
type StatusCount struct {
	Active  int
	Retired int
}

// DebutRetireCount is entity for debut & retire count.
type DebutRetireCount struct {
	Month  int
	Year   int
	Debut  int
	Retire int
}

// ModelCount is entity for 2d & 3d model count.
type ModelCount struct {
	None      int
	Has2DOnly int
	Has3DOnly int
	Both      int
}

// InAgencyCount is entity for in agency count.
type InAgencyCount struct {
	InAgency    int
	NotInAgency int
}

// SubscriberCount is entity for subscriber count.
type SubscriberCount struct {
	Min   int
	Max   int
	Count int
}

// DesignerCount is entity for designer count.
type DesignerCount struct {
	Name  string
	Count int
}

// VideoCountByDate is entity for video count by date.
type VideoCountByDate struct {
	Day   int
	Hour  int
	Count int
}

// VideoCount is entity for video count.
type VideoCount struct {
	ID    int64
	Name  string
	Count int
}

// VideoDuration is entity for video duration.
type VideoDuration struct {
	ID       int64
	Name     string
	Duration int
}
