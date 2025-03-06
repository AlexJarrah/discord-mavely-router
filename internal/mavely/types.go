package mavely

// TokenData represents the token data retrieved via the Mavely API
type TokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	FetchedAt    int64  `json:"fetched_at"`
}
