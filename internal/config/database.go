package config

import (
	"fmt"

	"github.com/LSariol/coveclient"
)

type DatabaseConfig struct {
	ConnectionString string
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{}
}

func (c *Config) InitilizeDatabaseConfig(cove *coveclient.Client, env string) error {

	var dbString string
	var err error

	if env == "PROD" {
		dbString, err = cove.GetSecret("BOTSUITE_DB_CONNECTION_STRING")
		if err != nil {
			return fmt.Errorf("get secret BOTSUITE_DB_CONNECTION_STRING: %w", err)
		}
	} else {
		dbString, err = cove.GetSecret("BOTSUITE_DB_CONNECTION_STRING_DEV")
		if err != nil {
			return fmt.Errorf("get secret BOTSUITE_DB_CONNECTION_STRING_DEV: %w", err)
		}
	}

	c.Database.ConnectionString = dbString

	return nil
}
