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

	twitch.AppClientID = os.Getenv("APP_CLIENT_ID")
	twitch.AppClientSecret = os.Getenv("APP_CLIENT_SECRET")
	twitch.AppAccessToken = os.Getenv("APP_ACCESS_TOKEN")
	twitch.UserAccessToken = os.Getenv("BOT_USER_ACCESS_TOKEN")
	twitch.UserRefreshToken = os.Getenv("BOT_USER_REFRESH_TOKEN")
	twitch.BotID = os.Getenv("BOT_ID")
	twitch.WebSocketURL = os.Getenv("TWITCH_WEBSOCKET_URL")
	twitch.APIURL = os.Getenv("TWITCH_API_URL")
	twitch.SubURL = os.Getenv("TWITCH_SUB_URL")

	cfg.Twitch = twitch

	return cfg, nil
}
