package middleware

import (
	"context"

	"github.com/lsariol/botsuite/internal/adapters/adapter"
	"github.com/lsariol/botsuite/internal/app"
)

func Logging(next Handler) Handler {
	return func(ctx context.Context, env adapter.Envelope, deps *app.Deps) (adapter.Response, error) {
		resp, err := next(ctx, env, deps)
		if err != nil {
			resp.Error = true
			return resp, err
		}

		return resp, err
	}
}

func Recovery(next Handler) Handler {
	return func(ctx context.Context, env adapter.Envelope, deps *app.Deps) (adapter.Response, error) {
		resp, err := next(ctx, env, deps)
		if err != nil {
			resp.Error = true
			return resp, err
		}

		return resp, err
	}
}
