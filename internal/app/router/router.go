package router

import (
	"context"

	"github.com/lsariol/botsuite/internal/adapters/adapter"
	"github.com/lsariol/botsuite/internal/app"
	"github.com/lsariol/botsuite/internal/app/registry"
	"github.com/lsariol/botsuite/internal/commands"
)

type Router struct {
	inbound  chan adapter.Envelope
	outbound chan adapter.Response
	adapters map[string]adapter.Adapter
	registry *registry.Registry
	rootCtx  context.Context
}

func NewRouter(ctx context.Context, reg *registry.Registry) *Router {
	return &Router{
		inbound:  make(chan adapter.Envelope, 100),
		outbound: make(chan adapter.Response, 100),
		rootCtx:  ctx,
		registry: reg,
	}
}

func (r *Router) Inbound() chan<- adapter.Envelope {
	return r.inbound
}

func (r *Router) Outbound() <-chan adapter.Response {
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

			var resp adapter.Response = Dispatch(ctx, env, cmd, deps)
			r.outbound <- resp

		}
	}
}

func (r *Router) registryLookUp(envelope adapter.Envelope) (commands.Command, bool) {

	cmd, ok := r.registry.Get(envelope.Command)

	return cmd, ok
}
