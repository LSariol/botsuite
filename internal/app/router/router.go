package router

import (
	"context"

	"github.com/lsariol/botsuite/internal/app"
	"github.com/lsariol/botsuite/internal/app/event"
	"github.com/lsariol/botsuite/internal/app/registry"
	"github.com/lsariol/botsuite/internal/commands"
)

type Router struct {
	inbound  chan event.Envelope
	outbound chan event.Response
	registry *registry.Registry
	rootCtx  context.Context
}

func NewRouter(ctx context.Context, reg *registry.Registry) *Router {
	return &Router{
		inbound:  make(chan event.Envelope, 100),
		outbound: make(chan event.Response, 100),
		rootCtx:  ctx,
		registry: reg,
	}
}

func (r *Router) Inbound() chan<- event.Envelope {
	return r.inbound
}

func (r *Router) Outbound() <-chan event.Response {
	return r.outbound
}

func (r *Router) Run(ctx context.Context, deps *app.Deps) {

	for {
		select {
		case <-ctx.Done():
			close(r.outbound)
			return
		case env, ok := <-r.inbound:
			if !ok {
				close(r.outbound)
				return
			}

			cmd, ok := r.registryLookUp(env)
			if !ok {
				continue
			}

			var resp event.Response = Dispatch(ctx, env, cmd, deps)
			r.outbound <- resp

		}
	}
}

func (r *Router) registryLookUp(envelope event.Envelope) (commands.Command, bool) {

	cmd, ok := r.registry.Get(envelope.Command)

	return cmd, ok
}
