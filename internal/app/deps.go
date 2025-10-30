package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/lsariol/botsuite/internal/config"
	"github.com/lsariol/botsuite/internal/database"
)

type Deps struct {
	Config   config.Config
	HTTP     *http.Client
	Logger   string
	DB       *database.Database
	BootTime time.Time
}

func NewDependencies() (*Deps, error) {

	var dependencies Deps

	cfg, err := config.Load()
	if err != nil {
		return &dependencies, fmt.Errorf("NewDeps: %w", err)
	}

	tr := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	dependencies.HTTP = client
	dependencies.Config = cfg
	dependencies.Logger = "*zap.Logger"
	dependencies.DB = database.NewDatabase()
	dependencies.BootTime = time.Now()

	return &dependencies, nil
}

func (d *Deps) RefreshTwitchUserTokens(userAccessToken string, refreshToken string) {

	d.Config.Twitch.UserAccessToken = userAccessToken
	d.Config.Twitch.UserRefreshToken = refreshToken
}
