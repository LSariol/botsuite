package adapter

import "context"

type Adapter interface {

	// Lifecycle
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Restart(ctx context.Context) error
	Close() error

	// Streams
	Events() <-chan Envelope

	// Outbound
	Deliver(ctx context.Context, r Response) error

	// Dynamic Membership
	Join(ctx context.Context, target string) error
	Leave(ctx context.Context, target string) error

	// Ops/Diagnostics
	Health(ctx context.Context) error
	Name() string
}
