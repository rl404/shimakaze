package entity

type messageType string

// Available message type.
const (
	TypeParseVtuber messageType = "parse-vtuber"
	TypeParseAgency messageType = "parse-agency"
)

// Message is entity for message.
type Message struct {
	Type messageType `json:"type"`
	Data []byte      `json:"data"`
}

// ParseVtuberRequest is parse vtuber request model.
type ParseVtuberRequest struct {
	ID     int64 `json:"id"`
	Forced bool  `json:"forced"`
}

// ParseAgencyRequest is parse agency request model.
type ParseAgencyRequest struct {
	ID     int64 `json:"id"`
	Forced bool  `json:"forced"`
}
