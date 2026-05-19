package feed

const (
	EventRSS          EventType = "RSS"
	EventNotification EventType = "notification"
)

type EventType string

// Event is the universal message passed from a feed Source to the Router.
type Event struct {
	Type        EventType
	Subscribers []Subscriber
	// Payload carries event-type-specific data.
	// Type-assert on EventType to access the concrete type.
	// EventNotification → NotificationPayload
	Payload any
}

// Subscriber identifies a platform channel that should receive this event.
type Subscriber struct {
	Platform    string
	ChannelID   string
	ChannelName string
}

// NotificationPayload is the Payload type for EventNotification events.
// It lives here (not in the notification source package) so the router can
// reference it without creating an import cycle.
type NotificationPayload struct {
	Message string
	UserID  string
	GUID    string
}
