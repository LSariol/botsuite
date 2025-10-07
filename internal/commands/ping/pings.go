package ping

import (
	"context"
	"fmt"
	"time"

	"github.com/lsariol/botsuite/internal/adapters/adapter"
	"github.com/lsariol/botsuite/internal/app"
)

type Ping struct{}

func (Ping) Name() string           { return "!ping" }
func (Ping) Aliases() []string      { return nil }
func (Ping) Description() string    { return "Latency check." }
func (Ping) Usage() string          { return "!ping" }
func (Ping) Timeout() time.Duration { return 3 * time.Second }

func (Ping) Execute(ctx context.Context, e adapter.Envelope, deps *app.Deps) (adapter.Response, error) {

	diff := time.Since(e.Timestamp)
	return adapter.Response{Text: fmt.Sprintf("pong! (Not currently accurate. Twitch doesnt return real timestamps on their event messages. %d ms)", diff.Milliseconds())}, nil

}
