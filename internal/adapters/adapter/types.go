package adapter

import "time"

type Envelope struct {
	Platform    string
	Username    string
	UserID      string
	ChannelName string
	ChannelID   string
	Command     string
	Args        []string
	Timestamp   time.Time
	RawMessage  string
	IsRegex     bool

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
	//Leaving text field empty will result in no response back to the adapter
	Text string

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
	// If send reply is false, we do not send back to the adapter.
	SuppressReply bool
}

type SystemEvent struct {
	Username     string
	UserID       string
	ChannelName  string
	ChannelID    string
	Setting      string
	SettingValue string
}

type HealthStatus struct {
	Name   string
	Status string
	Detail string
}
