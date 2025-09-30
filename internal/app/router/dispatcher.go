package router

import (
	"context"
	"time"

	"github.com/lsariol/botsuite/internal/adapters/adapter"
	"github.com/lsariol/botsuite/internal/app"
	"github.com/lsariol/botsuite/internal/app/middleware"
	"github.com/lsariol/botsuite/internal/commands"
)

func Dispatch(parentCtx context.Context, envelope adapter.Envelope, cmd commands.Command, deps *app.Deps) adapter.Response {

	ctx := parentCtx

	if t := cmd.Timeout(); t > 0 {

		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(parentCtx, t)
		defer cancel()
	}

	final := chain(cmd.Execute, middleware.Recovery, middleware.Logging)

	response, err := final(ctx, envelope, deps)

	response.Platform = envelope.Platform
	response.ChannelID = envelope.ChannelID
	response.ChannelName = envelope.ChannelName
	response.Platform = envelope.Platform
	response.Username = envelope.Username
	response.UserID = envelope.UserID
	response.ChannelName = envelope.ChannelName
	response.ChannelID = envelope.ChannelID
	response.TimeStart = envelope.Timestamp
	response.TimeFinished = time.Now()

	if err != nil {
		response.Error = true
		response.Success = false
		return response
	}

	response.Success = true
	return response

}

func chain(next middleware.Handler, mws ...middleware.Middleware) middleware.Handler {
	wrapped := next
	for i := len(mws) - 1; i >= 0; i-- {
		wrapped = mws[i](wrapped)
	}
	return wrapped
}
