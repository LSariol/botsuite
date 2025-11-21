package config

import (
	"fmt"

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

	err := c.InitilizeDatabase(cove)
	if err != nil {
		return fmt.Errorf("InitilizeDatabase: %w", err)
	}

	err = c.InitilizeTwitch(cove)
	if err != nil {
		return fmt.Errorf("InitilizeTwitch: %w", err)
	}

	return nil
}
