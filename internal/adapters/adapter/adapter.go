package adapter

import "context"

type Adapter interface {

	// Lifecycle
	Run(ctx context.Context) error
	Stop(ctx context.Context) error
	Restart(ctx context.Context) error
	Close(ctx context.Context) error

	// Streams
	Events() <-chan Envelope

	// Outbound
	Deliver(ctx context.Context, r Response) error

	// Dynamic Membership
	Join(ctx context.Context, targetID string) error
	Leave(ctx context.Context, targetID string) error

	// Ops/Diagnostics
	Health(ctx context.Context) error
	Name() string
}

// Easy Copy Paste to create a working class
// func (c *TwitchClient) Run(ctx context.Context) error { /* ... */ }
// func (c *TwitchClient) Stop(ctx context.Context) error  { /* ... */ }
// func (c *TwitchClient) Restart(ctx context.Context) error { /* ... */ }
// func (c *TwitchClient) Close(ctx context.Context) error { /* ... */ }
// func (c *TwitchClient) Events() <-chan adapters.Envelope { return c.events }
// func (c *TwitchClient) Deliver(ctx context.Context, r adapters.Response) error { /* ... */ }
// func (c *TwitchClient) Join(ctx context.Context, targetID string) error { /* ... */ }
// func (c *TwitchClient) Leave(ctx context.Context, targetID string) error { /* ... */ }
// func (c *TwitchClient) Health(ctx context.Context) error { return nil }
// func (c *TwitchClient) Name() string { return "twitch" }
