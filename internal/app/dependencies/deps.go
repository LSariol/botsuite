package dependencies

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/LSariol/coveclient"
	"github.com/joho/godotenv"
	"github.com/lsariol/botsuite/internal/config"
	"github.com/lsariol/botsuite/internal/database"
	"github.com/lsariol/botsuite/internal/runtime/settings"
)

type Deps struct {
	Config   *config.Config
	Settings *settings.Store
	HTTP     *http.Client
	Logger   string
	Cove     *coveclient.Client
	DB       *database.Database
	CTX      context.Context
	BootTime time.Time
}

func New(ctx context.Context) *Deps {
	return &Deps{
		CTX: ctx,
	}
}

func (d *Deps) Initilize() error {

	cfg := config.New()

	_ = godotenv.Load(".env")

	env := strings.ToUpper(os.Getenv("APP_ENV"))
	prod := env == "PROD"

	var cove *coveclient.Client
	if prod {
		cove = coveclient.New(os.Getenv("COVE_CLIENT_URL"), os.Getenv("COVE_CLIENT_SECRET"))
	} else {
		cove = coveclient.New(os.Getenv("COVE_CLIENT_URL_DEV"), os.Getenv("COVE_CLIENT_SECRET_DEV"))
	}

	err := cfg.Initilize(cove)
	if err != nil {
		return fmt.Errorf("initilize config: %w", err)
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

	d.Cove = cove
	d.HTTP = client
	d.Config = cfg
	d.Logger = "*zap.Logger"
	d.DB = database.NewDatabase(cfg.Database)
	d.BootTime = time.Now()

	return nil
}
