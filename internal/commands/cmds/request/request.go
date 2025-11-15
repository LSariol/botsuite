package request

import (
	"context"
	"time"

	"github.com/lsariol/botsuite/internal/adapters/adapter"
	"github.com/lsariol/botsuite/internal/app/dependencies"
)

type Request struct{}

func (Request) Name() string           { return "request" }
func (Request) Aliases() []string      { return nil }
func (Request) Regexes() []string      { return nil }
func (Request) Description() string    { return "Logs a request for things to be added to the bot." }
func (Request) Usage() string          { return "!request" }
func (Request) Timeout() time.Duration { return 3 * time.Second }

// Users can request things added to the bot. Will log it into the database.
func (Request) Execute(ctx context.Context, e adapter.Envelope, deps *dependencies.Deps) (adapter.Response, error) {

	return adapter.Response{Text: "Not set up yet. Leave a message in Mobasity's or BotMoba's chat and ill take a peak. The easier it is the sooner it can be implemented."}, nil

}
