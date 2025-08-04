package client

type RequestPayload struct {
	Method      string
	URL         string
	Headers     map[string]string
	QueryParams map[string]string
	Body        []byte
}

type ApiResponse[T any] struct {
	Success   bool   `json:"success"`
	Data      T      `json:"data,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
	ErrorMsg  string `json:"error_msg,omitempty"`
}

// Helper types and functions (copied from org_client.go, should be moved here)
type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type TokenRefreshResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
