package dependencies

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

func New() *Deps {
	return &Deps{}
}

func (d *Deps) Load() error {

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("NewDeps: %w", err)
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

	d.HTTP = client
	d.Config = cfg
	d.Logger = "*zap.Logger"
	d.DB = database.NewDatabase()
	d.BootTime = time.Now()

	return nil
}

func (d *Deps) RefreshTwitchUserTokens(userAccessToken string, refreshToken string) {

	d.Config.Twitch.UserAccessToken = userAccessToken
	d.Config.Twitch.UserRefreshToken = refreshToken
}
