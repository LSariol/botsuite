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

	d["APP_CLIENT_ID"] = cfg.AppClientID
	d["APP_CLIENT_SECRET"] = cfg.AppClientSecret
	d["APP_ACCESS_TOKEN"] = cfg.AppAccessToken
	d["BOT_USER_ACCESS_TOKEN"] = cfg.UserAccessToken
	d["BOT_USER_REFRESH_TOKEN"] = cfg.UserRefreshToken
	d["BOT_ID"] = cfg.BotID
	d["TWITCH_WEBSOCKET_URL"] = cfg.WebSocketURL
	d["TWITCH_API_URL"] = cfg.APIURL
	d["TWITCH_SUB_URL"] = cfg.SubURL

	return d
}
