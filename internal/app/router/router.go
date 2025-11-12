package router

import (
	"context"

	"github.com/lsariol/botsuite/internal/adapters/adapter"
	"github.com/lsariol/botsuite/internal/app/dependencies"
	"github.com/lsariol/botsuite/internal/app/registry"
	"github.com/lsariol/botsuite/internal/commands"
)

type Router struct {
	inbound         chan adapter.Envelope
	adapterRegistry map[string]adapter.Adapter
	registry        *registry.Registry
	rootCtx         context.Context
}

func NewRouter(ctx context.Context, reg *registry.Registry) *Router {
	return &Router{
		inbound:         make(chan adapter.Envelope, 100),
		adapterRegistry: make(map[string]adapter.Adapter),
		registry:        reg,
		rootCtx:         ctx,
	}
}

func (r *Router) Inbound() chan<- adapter.Envelope {
	return r.inbound
}

func (r *Router) Run(ctx context.Context, deps *dependencies.Deps) {

	for {
		select {
		case <-ctx.Done():
			return
		case env, ok := <-r.inbound:
			if !ok {
				return
			}

			cmd, ok := r.registryLookUp(env)
			if !ok {
				continue
			}

			r.Dispatch(ctx, env, cmd, deps)

		}
	}
}

func (r *Router) registryLookUp(envelope adapter.Envelope) (commands.Command, bool) {

	cmd, ok := r.registry.Get(envelope.Command)

	return cmd, ok
}

func (r *Router) RegisterAdapter(a adapter.Adapter) error {

	r.adapterRegistry[a.Name()] = a
	return nil
}
