package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/LSariol/coveclient"
)

type Config struct {
	Twitch   *TwitchConfig
	Database *DatabaseConfig
}

// Holds never changing variables ONLY

func New() *Config {
	return &Config{
		Twitch:   NewTwitchConfig(),
		Database: NewDatabaseConfig(),
	}
}

func (c *Config) Initilize(cove *coveclient.Client) error {

	env := strings.ToUpper(os.Getenv("APP_ENV"))

	err := c.InitilizeDatabaseConfig(cove, env)
	if err != nil {
		return fmt.Errorf("InitilizeDatabaseConfig: %w", err)
	}

	err = c.InitilizeTwitchConfig(cove)
	if err != nil {
		return fmt.Errorf("InitilizeTwitchConfig: %w", err)
	}

	return nil
}
