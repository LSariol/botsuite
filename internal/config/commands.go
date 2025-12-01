package config

import (
	"fmt"

	"github.com/LSariol/coveclient"
)

type CommandsConfig struct {
	ChatGPTKey string
}

func NewCommandsConfig() *CommandsConfig {
	return &CommandsConfig{}
}

func (c *Config) InitilizeCommandsConfig(cove *coveclient.Client) error {

	chatGPTKey, err := cove.GetSecret("OPEN_AI_API_KEY")
	if err != nil {
		return fmt.Errorf("get secret OPEN_AI_API_KEY: %w", err)
	}

	c.Commands.ChatGPTKey = chatGPTKey

	return nil
}
