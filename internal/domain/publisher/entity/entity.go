package entity

// Available topic.
const (
	TopicParseVtuber = "shimakaze-pubsub-parse-vtuber"
	TopicParseAgency = "shimakaze-pubsub-parse-agency"
)

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
