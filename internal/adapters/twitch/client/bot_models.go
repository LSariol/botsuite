package twitch

import twitchdb "github.com/lsariol/botsuite/internal/adapters/twitch/database"

type SessionData struct {
	SessionID        string
	KeepAliveTimeout int
	Channels         map[string]*twitchdb.TwitchChannel
}
