package twitchdb

type ChannelInsert struct {
	UserID       string
	Username     string
	AccessToken  string
	RefreshToken string
}

type TwitchChannel struct {
	ID             string
	Username       string
	CommandPrefix  string
	SubscriptionID string
}
