package client

import (
	"bytes"
	"encoding/json"

	"fmt"
	"net/http"
	"net/url"

	"github.com/Gkemhcs/kavach-cli/internal/auth"
	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
)

// BuildRequest creates an HTTP request from the given payload
func BuildRequest(p RequestPayload) (*http.Request, error) {
	// Attach query params to URL
	reqURL, err := url.Parse(p.URL)
	if err != nil {
		return nil, err
	}
	q := reqURL.Query()
	for key, val := range p.QueryParams {
		q.Set(key, val)
	}
	reqURL.RawQuery = q.Encode()

	// Build request
	req, err := http.NewRequest(p.Method, reqURL.String(), bytes.NewReader(p.Body))
	if err != nil {
		return nil, err
	}
	for key, val := range p.Headers {
		req.Header.Set(key, val)
	}
	return req, nil
}

// isAuthenticationError checks if the API response indicates an authentication failure
func isAuthenticationError[T any](respBody *APIResponse[T]) bool {
	if respBody == nil {
		return false
	}

	// Check for authentication-related error codes
	authErrorCodes := []string{
		"invalid_token",
		"expired_token",
		"unauthorized",
		"authentication_failed",
	}

	for _, code := range authErrorCodes {
		if respBody.ErrorCode == code {
			return true
		}
	}

	return false
}

// DoAuthenticatedRequest performs an authenticated HTTP request with automatic token refresh
func DoAuthenticatedRequest[T any](payload RequestPayload, logger *utils.Logger, cfg *config.Config) (*APIResponse[T], error) {
	tokenData, err := getTokenData(logger, cfg)
	if err != nil {
		return nil, err
	}

	if payload.Headers == nil {
		payload.Headers = make(map[string]string)
	}
	payload.Headers["Authorization"] = "Bearer " + tokenData.AccessToken
	req, err := BuildRequest(payload)
	if err != nil {
		return nil, err
	}

	var respBody APIResponse[T]

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if cliErrors.IsConnectionError(err.Error()) {
			return nil, cliErrors.ErrConnectionFailed
		}
		return nil, err
	}

	// Handle authentication errors
	if resp.StatusCode == 401 {
		defer resp.Body.Close()
		_ = json.NewDecoder(resp.Body).Decode(&respBody)

		// Check for specific authentication error codes
		if respBody.ErrorCode == "invalid_token" || respBody.ErrorCode == "expired_token" {
			return nil, cliErrors.ErrInvalidToken
		}

		// Generic 401 error
		return nil, cliErrors.ErrInvalidToken
	}

	if resp.StatusCode == 403 {
		return nil, cliErrors.ErrAccessDenied
	}

	defer resp.Body.Close()
	_ = json.NewDecoder(resp.Body).Decode(&respBody)

	if respBody.Success {
		return &respBody, nil
	}

	// Check if the error is due to expired token
	if !respBody.Success && respBody.ErrorCode == "expired_token" {
		// making a request to refresh endpoint as the access token is expired
		refreshReqBody, _ := json.Marshal(TokenRefreshRequest{RefreshToken: tokenData.RefreshToken})
		refreshURL := fmt.Sprintf("%sauth/refresh", cfg.BackendEndpoint)
		refreshReq, _ := http.NewRequest("POST", refreshURL, bytes.NewReader(refreshReqBody))
		refreshReq.Header.Set("Content-Type", "application/json")
		refreshResp, err := http.DefaultClient.Do(refreshReq)
		if err != nil && cliErrors.IsConnectionError(err.Error()) {
			return nil, cliErrors.ErrConnectionFailed
		}

		// If refresh fails, return authentication error
		if err != nil {
			return nil, cliErrors.ErrInvalidToken
		}

		type refreshTokenData struct {
			Token        string `json:"token"`
			RefreshToken string `json:"refresh_token"`
		}
		var refreshRespBody APIResponse[refreshTokenData]
		if err := json.NewDecoder(refreshResp.Body).Decode(&refreshRespBody); err != nil {
			defer refreshResp.Body.Close()
			return nil, cliErrors.ErrInvalidToken
		}

		// Check if refresh was successful
		if !refreshRespBody.Success {
			defer refreshResp.Body.Close()
			return nil, cliErrors.ErrInvalidToken
		}

		defer refreshResp.Body.Close()

		err = auth.SaveToken(auth.TokenData{
			AccessToken:  refreshRespBody.Data.Token,
			RefreshToken: refreshRespBody.Data.RefreshToken,
			Name:         tokenData.Name,
			Email:        tokenData.Email,
		}, logger, cfg)
		if err != nil {
			return nil, err
		}

		// Retry the original request with new token
		payload.Headers["Authorization"] = "Bearer " + refreshRespBody.Data.Token
		retryReq, err := BuildRequest(payload)
		if err != nil {
			return nil, err
		}

		retryResp, err := http.DefaultClient.Do(retryReq)
		if err != nil {
			if cliErrors.IsConnectionError(err.Error()) {
				return nil, cliErrors.ErrConnectionFailed
			}
			return nil, err
		}

		// Handle authentication errors in retry
		if retryResp.StatusCode == 401 {
			defer retryResp.Body.Close()
			return nil, cliErrors.ErrInvalidToken
		}

		if retryResp.StatusCode == 403 {
			defer retryResp.Body.Close()
			return nil, cliErrors.ErrAccessDenied
		}

		defer retryResp.Body.Close()
		var retryRespBody APIResponse[T]
		_ = json.NewDecoder(retryResp.Body).Decode(&retryRespBody)

		return &retryRespBody, nil
	}

	// For other errors, return the response body as is
	if !respBody.Success && isAuthenticationError(&respBody) {
		return nil, cliErrors.ErrInvalidToken
	}

	return &respBody, nil
}

func getTokenData(logger *utils.Logger, cfg *config.Config) (*auth.TokenData, error) {
	tokenData, err := auth.LoadToken(logger, cfg)
	if err != nil {
		return nil, cliErrors.ErrNotLoggedIn
	}

	return tokenData, nil

}
