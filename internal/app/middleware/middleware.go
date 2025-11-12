package middleware

import (
	"context"

	"github.com/lsariol/botsuite/internal/adapters/adapter"
	"github.com/lsariol/botsuite/internal/app/dependencies"
)

type Handler func(ctx context.Context, env adapter.Envelope, deps *dependencies.Deps) (adapter.Response, error)
type Middleware func(Handler) Handler
