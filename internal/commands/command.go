package commands

import (
	"context"
	"time"

	"github.com/lsariol/botsuite/internal/app"
	"github.com/lsariol/botsuite/internal/app/event"
)

type Command interface {
	Name() string
	Aliases() []string
	Description() string
	Usage() string
	Timeout() time.Duration
	Execute(ctx context.Context, e event.Envelope, deps *app.Deps) (event.Response, error)
}
