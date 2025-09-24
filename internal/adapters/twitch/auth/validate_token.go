package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func ValidateToken(token string, HTTP *http.Client) (*VerificationSuccess, error) {

	req, err := http.NewRequest("GET", "https://id.twitch.tv/oauth2/validate", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "OAuth "+token)

	resp, err := HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var fail VerificationFail
		if err := json.Unmarshal(body, &fail); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("unexpected status %d. error: %s", fail.Status, fail.Message)
	}

	var success VerificationSuccess
	if err := json.Unmarshal(body, &success); err != nil {
		return nil, err
	}

	return &success, nil

}
