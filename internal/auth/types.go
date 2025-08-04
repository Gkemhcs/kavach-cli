package auth

// TokenData represents the user's authentication tokens and profile info.
type TokenData struct {
	AccessToken  string `json:"access_token"`  // OAuth access token
	RefreshToken string `json:"refresh_token"` // OAuth refresh token
	Name         string `json:"name"`          // User's display name
	Email        string `json:"email"`         // User's email address
}
