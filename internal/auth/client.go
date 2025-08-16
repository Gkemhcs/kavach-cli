package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
)

// LoginWithGitHub performs the device flow login with GitHub and saves the access token.
// Handles device code retrieval, polling for token, and error handling/logging.
func LoginWithGitHub(logger *utils.Logger, cfg *config.Config) error {
	defer logger.Close()
	deviceCodeURL := cfg.DeviceCodeURL
	deviceTokenURL := cfg.DeviceTokenURL

	// 1. Start device flow
	req, _ := http.NewRequest("POST", deviceCodeURL, nil)
	req.Header.Set("Content-Type", "application/json")
	logger.LogRequest("POST " + deviceCodeURL)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if cliErrors.IsConnectionError(err.Error()) {
			return cliErrors.ErrUnReachableBackend
		}
		logger.Error("Failed to initiate device flow", err, map[string]interface{}{"cmd": "login"})
		return err
	}
	defer resp.Body.Close()

	var data struct {
		DeviceCode              string `json:"device_code"`
		UserCode                string `json:"user_code"`
		VerificationURI         string `json:"verification_uri"`
		VerificationURIComplete string `json:"verification_uri_complete"`
		ExpiresIn               int    `json:"expires_in"`
		Interval                int    `json:"interval"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		logger.Error("Failed to parse device code response", err, map[string]interface{}{"cmd": "login"})
		return err
	}

	logger.LogResponse(fmt.Sprintf("DeviceCode: %s, UserCode: %s, VerificationURI: %s", data.DeviceCode, data.UserCode, data.VerificationURI))
	fmt.Printf("🔑 To authenticate, visit %s and enter the code: %s\n\n", data.VerificationURI, data.UserCode)
	logger.Info("Prompted user to authenticate via browser", map[string]interface{}{"cmd": "login", "verification_uri": data.VerificationURI, "user_code": data.UserCode})

	fmt.Println("🔄 Waiting for you to complete authentication in your browser...")
	fmt.Println("⏱️  You have 2 minutes to complete authentication...")

	// 2. Poll for token (max 2 minutes)
	var tokenResp struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		User         struct {
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"user"`
	}

	maxAttempts := 60 // 2 minutes at 2s interval

	for i := 0; i < maxAttempts; i++ {
		body := map[string]string{"device_code": data.DeviceCode}
		b, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", deviceTokenURL, strings.NewReader(string(b)))
		req.Header.Set("Content-Type", "application/json")
		logger.LogRequest("POST " + deviceTokenURL + " Body: " + string(b))
		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			_ = json.NewDecoder(resp.Body).Decode(&tokenResp)
			logger.LogResponse(fmt.Sprintf("Token: %s, RefreshToken: %s, User: %v", tokenResp.Token, tokenResp.RefreshToken, tokenResp.User))
			if tokenResp.Token != "" {
				break
			}
		} else if err == nil {
			// Check for backend timeout error in response body
			var errResp struct {
				Error string `json:"error"`
			}
			_ = json.NewDecoder(resp.Body).Decode(&errResp)
			if errResp.Error == "device_authorization_timeout" {
				logger.Error("Backend device authorization timed out after 2 minutes", nil, map[string]interface{}{"cmd": "login"})
				return fmt.Errorf("device_authorization_timeout_backend")
			}
		}
		// Wait before polling again
		logger.Debug("Polling for device token", map[string]interface{}{"cmd": "login", "attempt": i + 1})
		time.Sleep(2 * time.Second)
	}

	if tokenResp.Token == "" {
		fmt.Println("❌ ⏰ Login timeout!")
		fmt.Println("❌ Unable to login within 2 minutes. Please try again.")
		logger.Error("Device authorization timed out after 2 minutes", nil, map[string]interface{}{"cmd": "login"})
		return fmt.Errorf("device_authorization_timeout")
	}

	fmt.Println("✅ Authentication completed successfully!")
	logger.Info("Device authorized, saving token", map[string]interface{}{"cmd": "login", "user": tokenResp.User.Email})
	return SaveToken(TokenData{
		AccessToken:  tokenResp.Token,
		RefreshToken: tokenResp.RefreshToken,
		Email:        tokenResp.User.Email,
		Name:         tokenResp.User.Username,
	}, logger, cfg)
}

// LogoutUser logs out the user by deleting the credentials file and logging the result.
func LogoutUser(logger *utils.Logger, cfg *config.Config) error {
	defer logger.Close()
	fmt.Println("🚪 Logging out...")
	logger.Info("Logout command started", map[string]interface{}{"cmd": "logout"})
	// Call the logout logic
	err := Logout(logger, cfg)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		logger.Warn("Logout attempted but credentials file did not exist", map[string]interface{}{"cmd": "logout"})
		return nil
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Logout failed. Please try again.")
		logger.Error("Logout failed", err, map[string]interface{}{"cmd": "logout"})
		return err
	}
	fmt.Println("✔️ Logout successful.")
	logger.Info("User successfully logged out", map[string]interface{}{"cmd": "logout"})
	return nil
}
