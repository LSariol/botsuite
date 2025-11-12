package twitchdb

type ChannelInsert struct {
	UserID       string
	Username     string
	AccessToken  string
	RefreshToken string
}

type TwitchChannel struct {
	UserID         string
	Username       string
	Role           string
	SessionID      string
	SubscriptionID string
}

type TwitchTokens struct {
	UserAccessToken  string
	UserRefreshToken string
}

type TwitchChannelSettings struct {
	UserID        string
	Username      string
	CommandPrefix string
}
