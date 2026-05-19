package sources

import "context"

// Source is the interface every feed source must implement.
// PullFeed runs until ctx is cancelled, reading from the source and
// pushing feed.Events to the channel handed to it via SetOutbound.
type Source interface {
	PullFeed(ctx context.Context)
}
