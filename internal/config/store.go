package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func StoreTwitchConfig(cfg *TwitchConfig) error {

	d := cfg.ToDict()

	err := godotenv.Write(d, "configs/.env")
	if err != nil {
		return fmt.Errorf("store twitch config: %w", err)
	}

	return nil
}
