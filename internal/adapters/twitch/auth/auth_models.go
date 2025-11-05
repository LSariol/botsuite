package auth

import (
	"encoding/json"
	"fmt"
	"os"
)

type RefreshSuccess struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}

type RefreshFail struct {
	Error   string `json:"error"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type VerificationFail struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type VerificationSuccess struct {
	ClientId  string   `json:"client_id"`
	Login     string   `json:"login"`
	Scopes    []string `json:"scopes"`
	UserId    string   `json:"user_id"`
	ExpiresIn int      `json:"expires_in"`
}

type UserData struct {
	UserAccessToken  string `json:"user_access_token"`
	UserRefreshToken string `json:"user_refresh_token"`
	UserID           string `json:"user_id"`
	Username         string `json:"username"`
}

func (u UserData) Store() error {

	existingUserData, err := LoadUserData()
	if err != nil {
		return err
	}

	existingUserData = append(existingUserData, u)

	if err := SaveUserData(existingUserData); err != nil {
		return err
	}
	return nil
}

func LoadUserData() ([]UserData, error) {

	var userData []UserData

	data, err := os.ReadFile("internal/database/users.json")
	if err != nil {
		if os.IsNotExist(err) {
			return []UserData{}, nil
		}
		return nil, fmt.Errorf("read file: %w", err)
	}

	if err := json.Unmarshal(data, &userData); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return userData, nil
}

func SaveUserData(userData []UserData) error {

	data, err := json.MarshalIndent(userData, "", "  ")
	if err != nil {
		return fmt.Errorf("masrhsal: %w", err)
	}

	if err := os.WriteFile("internal/database/users.json", data, 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil

}
