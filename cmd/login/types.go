package login

// Credentials represents the user's authentication credentials for login.
type Credentials struct {
	Token     string `json:"token"`                // OAuth or access token
	ExpiresAt string `json:"expires_at,omitempty"` // Expiry time (optional)
	Email     string `json:"email,omitempty"`      // User email (optional)
	Username  string `json:"username,omitempty"`   // Username (optional)
}
