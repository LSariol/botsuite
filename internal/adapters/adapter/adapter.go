package adapter

import (
	"context"
)

type Adapter interface {

	// Lifecycle
	Run(ctx context.Context) error // Starts long lived loops; blocks until ctx cancel
	Shutdown(ctx context.Context)  // Gracefully stops the Adapter

	// This will be handled by the adapter orchestrator
	//Restart(ctx context.Context) error  // Re-Initilizes the Adapter

	// Streams

	//Not sure if necessary
	// inEvents() <-chan Envelope
	//OutboundEnvelopes() chan<- Envelope // Read only channel for Registry Router
	DeliverResponse(r Response)

	// Event Handling

	// Both of these are not needed here. We need type specific behavior,
	// ConsumeEvent(ctx context.Context, event twitch.EventSubMessage) error    // Takes events and decides if its one we care about. (command prefix, specific keywords, ect)
	// ConsumeMessage(ctx context.Context) error             // Parses the message from the event and packages it into a format readable by the router

	// Not needed?
	// EmitMessage() error                                    // Send Message to Registry Router

	// Not needed either?
	// ReceiveResponse() error                                // Receives outputs from the Registry Router
	// DeliverResponse(ctx context.Context, r Response) error // Send Response back to platform

	// Dynamic Membership
	// Join(ctx context.Context, targets []string) error  // Joins a single connection (channel, server, group)
	// Leave(ctx context.Context, targets []string) error // Leaves a single connection (channel, server, group)

	// Ops/Diagnostics
	Health(ctx context.Context) HealthStatus // Returns the current health statate of an adapter
	Name() string                            // Returns the name of the adapter
}

// Easy Copy Paste to create a working class

// func (c *TwitchClient) Initilize(ctx context.Context) error { /* ... */ }
// func (c *TwitchClient) Run(ctx context.Context) error { /* ... */ }
// func (c *TwitchClient) Shutdown(ctx context.Context) error { /* ... */ }
// func (c *TwitchClient) Restart(ctx context.Context) error { /* ... */ }
// func (c *TwitchClient) OutBoundEvents() <-chan Envelope
// func (c *TwitchClient) ConsumeEvent(ctx context.Context) error
// func (c *TwitchClient) ConsumeMessage(ctx context.Context) error
// func (c *TwitchClient) EmitMessage() error
// func (c *TwitchClient) ReceiveResponse() error
// func (c *TwitchClient) DeliverReponse(ctx context.Context, r Response) error
// func (c *TwitchClient) Join(ctx context.Context, targets ...string) error
// func (c *TwitchClient) Leave(ctx context.Context, targets ...string) error
// func (c *TwitchClient) Health(ctx context.Context) HealthStatus
// func (c *TwitchClient) Name() string
