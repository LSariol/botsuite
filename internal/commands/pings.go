package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/lsariol/botsuite/internal/app"
	"github.com/lsariol/botsuite/internal/app/event"
)

type Ping struct{}

func (Ping) Name() string           { return "!ping" }
func (Ping) Aliases() []string      { return nil }
func (Ping) Description() string    { return "Latency check." }
func (Ping) Usage() string          { return "!ping" }
func (Ping) Timeout() time.Duration { return 3 * time.Second }

func (Ping) Execute(ctx context.Context, e event.Envelope, deps *app.Deps) (event.Response, error) {

	diff := time.Since(e.Timestamp)
	fmt.Println(e.Timestamp)
	fmt.Println(time.Now())
	return event.Response{Text: fmt.Sprintf("pong! (%d ms)", diff.Milliseconds())}, nil

}
