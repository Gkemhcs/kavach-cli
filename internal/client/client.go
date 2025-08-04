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

func DoAuthenticatedRequest[T any](payload RequestPayload, logger *utils.Logger, cfg *config.Config) (*ApiResponse[T], error) {
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

	var respBody ApiResponse[T]

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if cliErrors.IsConnectionError(err.Error()) {
			return nil, cliErrors.ErrConnectionFailed
		}
		return nil, err
	}
	if resp.StatusCode==403{
		return nil,cliErrors.ErrAccessDenied
	}
	defer resp.Body.Close()
	_ = json.NewDecoder(resp.Body).Decode(&respBody)

	if respBody.Success {
		return &respBody, nil
	}

	if !respBody.Success && respBody.ErrorCode != "expired_token" {
		return &respBody, nil
	}

	// making a request to refresh endpoint as the access token is expired
	refreshReqBody, _ := json.Marshal(TokenRefreshRequest{RefreshToken: tokenData.RefreshToken})
	refreshURL := fmt.Sprintf("%sauth/refresh", cfg.BackendEndpoint)
	refreshReq, _ := http.NewRequest("POST", refreshURL, bytes.NewReader(refreshReqBody))
	refreshReq.Header.Set("Content-Type", "application/json")
	refreshResp, err := http.DefaultClient.Do(refreshReq)
	if err != nil && cliErrors.IsConnectionError(err.Error()) {
		return nil, cliErrors.ErrConnectionFailed
	}
	type refreshTokenData struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	var refreshRespBody ApiResponse[refreshTokenData]
	if err := json.NewDecoder(refreshResp.Body).Decode(&refreshRespBody); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	if err != nil {
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
	payload.Headers["Authorization"] = "Bearer " + refreshRespBody.Data.Token
	retryReq, err := BuildRequest(payload)
	if err != nil {
		return nil, err
	}

	retryResp, err := http.DefaultClient.Do(retryReq)
	if retryResp.StatusCode==403{
		return nil,cliErrors.ErrAccessDenied
	}
	if err != nil {
		if cliErrors.IsConnectionError(err.Error()) {
			return nil, cliErrors.ErrConnectionFailed
		}
		return nil, err
	}
	defer retryResp.Body.Close()
	var retryRespBody ApiResponse[T]
	_ = json.NewDecoder(retryResp.Body).Decode(&retryRespBody)

	return &retryRespBody, nil

}

func getTokenData(logger *utils.Logger, cfg *config.Config) (*auth.TokenData, error) {
	tokenData, err := auth.LoadToken(logger, cfg)
	if err != nil {
		return nil, cliErrors.ErrNotLoggedIn
	}

	return tokenData, nil

}
