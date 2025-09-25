package middleware

import (
	"context"

	"github.com/lsariol/botsuite/internal/app"
	"github.com/lsariol/botsuite/internal/app/event"
)

func Logging(next Handler) Handler {
	return func(ctx context.Context, env event.Envelope, deps *app.Deps) (event.Response, error) {
		resp, err := next(ctx, env, deps)
		if err != nil {
			resp.Error = true
			return resp, err
		}

		return resp, err
	}
}

func Recovery(next Handler) Handler {
	return func(ctx context.Context, env event.Envelope, deps *app.Deps) (event.Response, error) {
		resp, err := next(ctx, env, deps)
		if err != nil {
			resp.Error = true
			return resp, err
		}

		return resp, err
	}
}
