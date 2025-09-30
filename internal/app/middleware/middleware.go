package middleware

import (
	"context"

	"github.com/lsariol/botsuite/internal/adapters/adapter"
	"github.com/lsariol/botsuite/internal/app"
)

type Handler func(ctx context.Context, env adapter.Envelope, deps *app.Deps) (adapter.Response, error)
type Middleware func(Handler) Handler
