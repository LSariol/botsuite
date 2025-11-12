package bot

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lsariol/botsuite/internal/adapters/adapter"
	"github.com/lsariol/botsuite/internal/adapters/twitch/eventsub"
)

func (c *TwitchClient) ingestLoop(ctx context.Context, in <-chan eventsub.EventSubMessage) error {

	fmt.Println("[TwitchBot] Running")
	for {
		select {
		case <-ctx.Done():
			log.Println("context canceled, stopped ingest loop")
			return nil

		case msg, ok := <-in:
			if !ok {
				fmt.Println("twitchClient inboundEvents channel closed")
				return nil
			}

			//Log message inbound
			c.printE(msg)
			go c.pack(msg)

			//Check leading words. If words are found send off

		}
	}
}

func (c *TwitchClient) pack(msg eventsub.EventSubMessage) {

	m := msg.Payload.Event.Message.Text
	prefix := c.ChannelSettings[msg.Payload.Event.BroadcasterUserID].CommandPrefix
	cmd, args, ok := c.parseCMD(m, prefix)
	if !ok {
		return
	}

	var newEnvelope adapter.Envelope
	rawTime, _ := time.Parse(time.RFC3339Nano, msg.Metadata.MessageTimestamp)
	newEnvelope.Platform = "twitch"
	newEnvelope.Username = msg.Payload.Event.ChatterUserName
	newEnvelope.UserID = msg.Payload.Event.ChatterUserID
	newEnvelope.ChannelName = msg.Payload.Event.BroadcasterUserName
	newEnvelope.ChannelID = msg.Payload.Event.BroadcasterUserID
	newEnvelope.Command = strings.ToLower(cmd)
	newEnvelope.Content = args
	newEnvelope.Timestamp = rawTime

	c.outEnvelopes <- newEnvelope
}

func (c *TwitchClient) parseCMD(msg string, prefix string) (string, []string, bool) {

	msg = strings.TrimSpace(msg)

	if strings.HasPrefix(msg, prefix) {
		withoutPrefix := strings.TrimPrefix(strings.ToLower(msg), prefix)
		for cmd, kind := range c.Commands {
			if kind == "cmd" && strings.HasPrefix(withoutPrefix, cmd) {
				if len(withoutPrefix) == len(cmd) || withoutPrefix[len(cmd)] == ' ' {
					remainder := strings.TrimSpace(withoutPrefix[len(cmd):])
					if remainder != "" {
						args := strings.Fields(remainder)
						return cmd, args, true
					}
					return cmd, nil, true
				}
			}
		}
	}

	for cmd, kind := range c.Commands {
		if kind == "trigger" && strings.HasPrefix(strings.ToLower(msg), cmd) {
			remainder := strings.TrimSpace(msg[len(cmd):])
			if remainder != "" {
				args := strings.Fields(remainder)
				return cmd, args, true
			}
			return cmd, nil, true
		}
	}
	return "", nil, false
}

func (c *TwitchClient) printE(e eventsub.EventSubMessage) {

	now := timestamp()

	log.Printf("[%s] %-10.10s: @%-10.10s - %s",
		now,                                 // formatted as YYYY-MM-DD-HH-MM-SS-msms
		e.Payload.Event.BroadcasterUserName, // max width 10 chars, padded right
		e.Payload.Event.ChatterUserName,     // max width 10 chars, padded right
		e.Payload.Event.Message.Text,        // remainder of line
	)
}

func timestamp() string {

	return time.Now().Format("2006-01-02 15-04-05.000")
}
