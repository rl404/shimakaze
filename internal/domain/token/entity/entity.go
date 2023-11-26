package entity

// CreateAccessTokenRequest is request model for create access token.
type CreateAccessTokenRequest struct {
	UserID      int64
	AccessUUID  string
	RefreshUUID string
}

// CreateRefreshTokenRequest is request model for create refresh token.
type CreateRefreshTokenRequest struct {
	UserID      int64
	RefreshUUID string
}

// Payload is entity for token payload.
type Payload struct {
	UserID      int64
	RefreshUUID string
}
