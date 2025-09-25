package middleware

import (
	"context"

	"github.com/lsariol/botsuite/internal/app"
	"github.com/lsariol/botsuite/internal/app/event"
)

type Handler func(ctx context.Context, env event.Envelope, deps *app.Deps) (event.Response, error)
type Middleware func(Handler) Handler
