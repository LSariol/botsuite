package middleware

import (
	"context"
	"fmt"

	"github.com/lsariol/botsuite/internal/app"
	"github.com/lsariol/botsuite/internal/app/event"
)

func Logging(next Handler) Handler {
	return func(ctx context.Context, env event.Envelope, deps *app.Deps) (event.Response, error) {
		fmt.Println("Logging: IN")
		resp, err := next(ctx, env, deps)
		if err != nil {
			resp.Error = true
			return resp, err
		}
		fmt.Println("Logging: OUT")
		return resp, err
	}
}

func Recovery(next Handler) Handler {
	return func(ctx context.Context, env event.Envelope, deps *app.Deps) (event.Response, error) {
		fmt.Println("Recovery: IN")
		resp, err := next(ctx, env, deps)
		if err != nil {
			resp.Error = true
			return resp, err
		}
		fmt.Println("Recovery: OUT")

		return resp, err
	}
}
