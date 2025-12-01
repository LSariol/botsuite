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
	Commands *CommandsConfig
}

// Holds never changing variables ONLY

func New() *Config {
	return &Config{
		Twitch:   NewTwitchConfig(),
		Database: NewDatabaseConfig(),
		Commands: NewCommandsConfig(),
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

	err = c.InitilizeCommandsConfig(cove)
	if err != nil {
		return fmt.Errorf("InitilizeCommandsConfig: %w", err)
	}

	return nil
}
