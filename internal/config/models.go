package config

type Config struct {
	Twitch TwitchConfig
}

type TwitchConfig struct {
	AppClientID      string
	AppClientSecret  string
	AppAccessToken   string
	UserAccessToken  string
	UserRefreshToken string
	BotID            string
	WebSocketURL     string
	APIURL           string
	SubURL           string
}

func (cfg TwitchConfig) ToDict() map[string]string {

	d := make(map[string]string)

	d["TWITCH_APP_CLIENT_ID"] = cfg.AppClientID
	d["TWITCH_APP_CLIENT_SECRET"] = cfg.AppClientSecret
	d["TWITCH_APP_ACCESS_TOKEN"] = cfg.AppAccessToken
	d["TWITCH_BOT_USER_ACCESS_TOKEN"] = cfg.UserAccessToken
	d["TWITCH_BOT_USER_REFRESH_TOKEN"] = cfg.UserRefreshToken
	d["TWITCH_BOT_ID"] = cfg.BotID

	return d
}
