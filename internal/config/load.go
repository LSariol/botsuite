package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Load() (Config, error) {

	var cfg Config
	var twitch TwitchConfig

	err := godotenv.Load("configs/.env")
	if err != nil {
		return cfg, fmt.Errorf("error loading .env file: %w", err)
	}

	twitch.AppClientID = os.Getenv("TWITCH_APP_CLIENT_ID")
	twitch.AppClientSecret = os.Getenv("TWITCH_APP_CLIENT_SECRET")
	twitch.AppAccessToken = os.Getenv("TWITCH_APP_ACCESS_TOKEN")
	twitch.UserAccessToken = os.Getenv("TWITCH_BOT_USER_ACCESS_TOKEN")
	twitch.UserRefreshToken = os.Getenv("TWITCH_BOT_USER_REFRESH_TOKEN")
	twitch.BotID = os.Getenv("TWITCH_BOT_ID")

	cfg.Twitch = twitch

	return cfg, nil
}
