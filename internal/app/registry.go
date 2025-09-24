package app

import (
	"context"
	"fmt"

	"github.com/lsariol/botsuite/internal/event"
)

type Registry struct {
	inbound  chan event.Envelope
	outbound chan event.Response
}

func NewRegistry() *Registry {
	return &Registry{
		inbound:  make(chan event.Envelope, 100),
		outbound: make(chan event.Response, 100),
	}
}

func (r *Registry) Inbound() chan<- event.Envelope {
	return r.inbound
}

func (r *Registry) Outbound() <-chan event.Response {
	return r.outbound
}

func (r *Registry) Run(ctx context.Context) {

	fmt.Println("Registry is running")
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

			fmt.Println("Registry has recieved the envelope")
			resp := r.disptch(env)
			r.outbound <- resp

		}
	}
}

func (r *Registry) disptch(envelope event.Envelope) event.Response {

	if envelope.Command == "!ping" {
		return event.Response{
			Platform:  envelope.Platform,
			ChannelID: envelope.ChannelID,
			Text:      "pong!",
		}
	}
	return event.Response{
		Platform:  envelope.Platform,
		ChannelID: envelope.ChannelID,
		Text:      "unknown command",
		Error:     true,
	}
}
