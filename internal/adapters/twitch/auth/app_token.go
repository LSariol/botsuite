package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/lsariol/botsuite/internal/config"
)

func RefreshAppAccessToken(cfg *config.TwitchConfig, HTTP *http.Client) error {

	var appAccessToken string

	vals := url.Values{}
	vals.Set("client_id", cfg.AppClientID)
	vals.Set("client_secret", cfg.AppClientSecret)
	vals.Set("grant_type", "client_credentials")
	urlVals := vals.Encode()

	req, err := http.NewRequest("POST", "https://id.twitch.tv/oauth2/token", strings.NewReader(urlVals))
	if err != nil {
		return fmt.Errorf("RefreshAppAccessToken: error with request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("RefreshAppAccessToken: response unexpected status: %d: %w", resp.StatusCode, err)
	}

	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("RefreshAppAccessToken: error with decode: %w", err)
	}

	appAccessToken, ok := data["access_token"].(string)

	if !ok {
		return fmt.Errorf("RefreshAppAccessToken: Access token not found in response")
	}

	cfg.AppAccessToken = appAccessToken

	return nil
}
