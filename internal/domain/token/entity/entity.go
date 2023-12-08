package entity

// CreateAccessTokenRequest is request model for create access token.
type CreateAccessTokenRequest struct {
	UserID      int64
	Username    string
	AccessUUID  string
	RefreshUUID string
}

// CreateRefreshTokenRequest is request model for create refresh token.
type CreateRefreshTokenRequest struct {
	UserID      int64
	Username    string
	RefreshUUID string
}

// Payload is entity for token payload.
type Payload struct {
	UserID      int64
	RefreshUUID string
}
