package client

// RequestPayload represents the data needed to build an HTTP request
type RequestPayload struct {
	Method      string
	URL         string
	Headers     map[string]string
	QueryParams map[string]string
	Body        []byte
}

// APIResponse represents the standard API response structure
type APIResponse[T any] struct {
	Success   bool   `json:"success"`
	Data      T      `json:"data,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
	ErrorMsg  string `json:"error_msg,omitempty"`
}

// TokenRefreshRequest represents a request to refresh an access token
type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// TokenRefreshResponse represents the response from a token refresh request
type TokenRefreshResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
