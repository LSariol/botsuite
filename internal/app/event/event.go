package event

import "time"

type Envelope struct {
	Platform    string
	Username    string
	UserID      string
	ChannelName string
	ChannelID   string
	Command     string
	Content     string
	Timestamp   time.Time

	// Implement real time stamping
	// Ingress_ts        string
	// Router_in_ts      string
	// Dispatch_start_ts string
}

type Response struct {
	Platform    string
	Username    string
	UserID      string
	ChannelName string
	ChannelID   string
	Text        string

	//Temp timestamps
	TimeStart    time.Time
	TimeFinished time.Time

	//Timestamp variables
	// Ingress_ts        string // Time received from twitch (EventSubTimestamp)
	// Egress_ts         string // Time sent back to twitch
	// Router_in_ts      string // When router dequeues the envelope
	// Dispatch_start_ts string // Before dispatcher runs the handler chain
	// Dispatch_end_ts   string // After dispatcher runs the handler chain

	Success bool
	Error   bool
}
