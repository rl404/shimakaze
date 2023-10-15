package entity

type messageType string

// Available message types.
const (
	TypeParseVtuber messageType = "parse-vtuber"
	TypeParseAgency messageType = "parse-agency"
)

// Message is pubsub message.
type Message struct {
	Type   messageType `json:"type"`
	ID     int64       `json:"id"`
	Forced bool        `json:"forced"`
}
