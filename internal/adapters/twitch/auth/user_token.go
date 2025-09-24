package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/lsariol/botsuite/internal/config"
)

func GenerateUserAcessToken(userToken string, cfg *config.TwitchConfig, HTTP *http.Client) (UserData, error) {

	var userData UserData

	vals := url.Values{}
	vals.Set("client_id", cfg.AppClientID)
	vals.Set("client_secret", cfg.AppClientSecret)
	vals.Set("code", userToken)
	vals.Set("grant_type", "authorization_code")
	vals.Set("redirect_uri", "http://localhost:3000")
	urlVals := vals.Encode()

	req, err := http.NewRequest("POST", "https://id.twitch.tv/oauth2/token", strings.NewReader(urlVals))
	if err != nil {
		return userData, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := HTTP.Do(req)
	if err != nil {
		return userData, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return userData, fmt.Errorf("generateuseraccesstoken unexpected status %d", resp.StatusCode)
	}

	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return userData, err
	}

	userAccessToken, ok := data["access_token"].(string)
	if !ok {
		return userData, fmt.Errorf("access_token not found in response")
	}

	refreshToken, ok := data["refresh_token"].(string)
	if !ok {
		return userData, fmt.Errorf("access_token not found in response")
	}

	userData.UserAccessToken = userAccessToken
	userData.UserRefreshToken = refreshToken

	return userData, nil
}

func RefreshUserAccessToken(refreshToken string, cfg *config.TwitchConfig, HTTP *http.Client) (UserData, error) {

	var userData UserData
	data := url.Values{}
	data.Set("client_id", cfg.AppClientID)
	data.Set("client_secret", cfg.AppClientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", "https://id.twitch.tv/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return userData, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := HTTP.Do(req)
	if err != nil {
		return userData, err
	}
	defer resp.Body.Close()

	var fail RefreshFail
	var success RefreshSuccess

	if resp.StatusCode != http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(&fail); err != nil {
			return userData, fmt.Errorf("RefreshUserAccessToken: StatusCode not OK: Decode Error: %w", err)
		}
		return userData, fmt.Errorf("RefreshUserAccessToken: StatusCode not OK: Code %d: Message %s: Error: %s: %w", fail.Status, fail.Message, fail.Error, err)
	}

	if err := json.NewDecoder(resp.Body).Decode(&success); err != nil {
		return userData, fmt.Errorf("RefreshUserAccessToken: StatusCode OK: Decode Error: %w", err)
	}

	userData.UserAccessToken = success.AccessToken
	userData.UserRefreshToken = success.RefreshToken

	return userData, nil

}
