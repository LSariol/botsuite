package auth

import (
	"context"
	"fmt"
	"net/http"

	twitchdb "github.com/lsariol/botsuite/internal/adapters/twitch/database"
	"github.com/lsariol/botsuite/internal/config"
)

type AuthClient struct {
	HTTP   *http.Client
	DB     *twitchdb.Store
	Config *config.TwitchConfig
	Tokens *SafeTwitchTokens
}

func New(db *twitchdb.Store, cfg *config.TwitchConfig, http *http.Client) *AuthClient {
	return &AuthClient{
		DB:     db,
		Config: cfg,
		HTTP:   http,
		Tokens: NewSafeTwitchTokens(twitchdb.TwitchTokens{}),
	}
}

func (c *AuthClient) Initilize(ctx context.Context) error {

	tokens, err := c.DB.GetTokens(ctx)
	if err != nil {
		return fmt.Errorf("get tokens: %w", err)
	}

	c.Tokens.SetTokens(tokens.UserAccessToken, tokens.UserRefreshToken)

	return nil
}

func (c *AuthClient) Shutdown(ctx context.Context) error {

	if err := c.DB.StoreTokens(ctx, c.Tokens.UserAccessToken(), c.Tokens.UserRefreshToken()); err != nil {
		return fmt.Errorf("store tokens: %w", err)
	}

	return nil
}
