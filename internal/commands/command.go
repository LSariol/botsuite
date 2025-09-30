package commands

import (
	"context"
	"time"

	"github.com/lsariol/botsuite/internal/adapters/adapter"
	"github.com/lsariol/botsuite/internal/app"
)

type Command interface {
	Name() string
	Aliases() []string
	Description() string
	Usage() string
	Timeout() time.Duration
	Execute(ctx context.Context, e adapter.Envelope, deps *app.Deps) (adapter.Response, error)
}
