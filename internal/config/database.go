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

func (c *Config) InitilizeDatabase(cove *coveclient.Client) error {

	dbString, err := cove.GetSecret("BOTSUITE_DB_CONNECTION_STRING")
	if err != nil {
		return fmt.Errorf("get secret BOTSUITE_DB_CONNECTION_STRING: %w", err)
	}

	c.Database.ConnectionString = dbString

	return nil
}
