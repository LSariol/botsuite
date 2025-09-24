package event

type Envelope struct {
	Platform    string
	Username    string
	UserID      string
	ChannelName string
	ChannelID   string
	Command     string
	Content     string
	Timestamp   string
}

type Response struct {
	Platform     string
	Username     string
	UserID       string
	ChannelName  string
	ChannelID    string
	Text         string
	TimeStart    string
	TimeFinished string
	Success      bool
	Error        bool
}
