package dependencies

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/LSariol/coveclient"
	"github.com/joho/godotenv"
	"github.com/lsariol/botsuite/internal/config"
	"github.com/lsariol/botsuite/internal/database"
)

type Deps struct {
	Config   *config.Config
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

	if os.Getenv("COVECLIENT_URL") == "" {
		err := godotenv.Load("configs/.env")
		if err != nil {
			return fmt.Errorf("error loading .env file: %w", err)
		}
	}

	cove := coveclient.New(os.Getenv("COVECLIENT_URL"), os.Getenv("COVECLIENT_SECRET"))

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
