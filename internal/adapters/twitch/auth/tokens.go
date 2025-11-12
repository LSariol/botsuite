package auth

import (
	"sync"

	twitchdb "github.com/lsariol/botsuite/internal/adapters/twitch/database"
)

type SafeTwitchTokens struct {
	mu     sync.RWMutex
	tokens twitchdb.TwitchTokens
}

func NewSafeTwitchTokens(tok twitchdb.TwitchTokens) *SafeTwitchTokens {
	return &SafeTwitchTokens{
		tokens: tok,
	}
}

func (s *SafeTwitchTokens) UserAccessToken() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tokens.UserAccessToken
}

func (s *SafeTwitchTokens) UserRefreshToken() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tokens.UserRefreshToken
}

func (s *SafeTwitchTokens) Tokens() (access, refresh string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tokens.UserAccessToken, s.tokens.UserRefreshToken
}

func (s *SafeTwitchTokens) SetTokens(access string, refresh string) {
	s.mu.Lock()
	s.tokens.UserAccessToken = access
	s.tokens.UserRefreshToken = refresh
	s.mu.Unlock()
}
