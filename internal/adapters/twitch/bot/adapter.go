package twitch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lsariol/botsuite/internal/adapters/adapter"
	"github.com/lsariol/botsuite/internal/adapters/twitch/auth"
	"github.com/lsariol/botsuite/internal/config"
)

const (
	HelixBaseURL       = "https://api.twitch.tv/helix"
	EventSubAPIBaseURL = "https://api.twitch.tv/helix/eventsub/subscriptions"
	EventSubWSURL      = "wss://eventsub.wss.twitch.tv/ws"
)

var _ adapter.Adapter = (*TwitchClient)(nil)

type TwitchClient struct {
	HTTP        *http.Client
	WS          *websocket.Conn
	Config      *config.TwitchConfig
	SessionData SessionData
	events      chan adapter.Envelope
}

func NewTwitchBot(client *http.Client, cfg *config.TwitchConfig) *TwitchClient {
	return &TwitchClient{
		HTTP:   client,
		Config: cfg,
		events: make(chan adapter.Envelope, 100),
	}
}

func (c *TwitchClient) Run(ctx context.Context) error {

	if err := c.refreshTokens(); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	if err := c.loadChannels(); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	conn, sessionData, err := c.newWebsocketConn(EventSubWSURL)
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}
	c.WS = conn
	c.SessionData.SessionID = sessionData.SessionID
	c.SessionData.KeepAliveTimeout = sessionData.KeepAliveTimeout

	//for channel in channels
	for _, channel := range c.SessionData.Channels {
		if err := c.Join(ctx, channel.ID); err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("Connected to %s\n", channel.Username)
	}

	go c.listen()

	return nil
}

func (c *TwitchClient) Stop(ctx context.Context) error {
	return nil
}

func (c *TwitchClient) Restart(ctx context.Context) error {
	return nil
}

// Gracefully closes the websocket connection
func (c *TwitchClient) Close(ctx context.Context) error {

	return nil
}

func (c *TwitchClient) Events() <-chan adapter.Envelope {
	return c.events
}

func (c *TwitchClient) Deliver(ctx context.Context, r adapter.Response) error {

	reqBody := chatMessageReq{
		BroadcasterID: r.ChannelID,
		SenderID:      c.Config.BotID,
		Message:       r.Text,
	}

	var out any
	if err := c.postHelixJSON(ctx, "/chat/messages", reqBody, &out); err != nil {
		if errors.Is(err, ErrMissingChannelBot) {
			c.Leave(ctx, r.ChannelID)
		}
		return err
	}
	return nil

}

func (c *TwitchClient) Join(ctx context.Context, targetID string) error {

	body := map[string]any{
		"type":    "channel.chat.message",
		"version": "1",
		"condition": map[string]string{
			"broadcaster_user_id": targetID,
			"user_id":             c.Config.BotID,
		},
		"transport": map[string]string{
			"method":     "websocket",
			"session_id": c.SessionData.SessionID,
		},
	}

	buf, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", EventSubAPIBaseURL, bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("join: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Config.UserAccessToken) // MUST be a user token for WebSocket subs
	req.Header.Set("Client-Id", c.Config.AppClientID)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var respData EventSubJoinResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return err
	}

	for _, sub := range respData.Data {
		c.SessionData.Channels[targetID].SubscriptionID = sub.ID
	}

	return nil
}

func (c *TwitchClient) Leave(ctx context.Context, target string) error {

	subscriptionID := c.SessionData.Channels[target].SubscriptionID
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, EventSubAPIBaseURL+"?id="+subscriptionID, nil)
	if err != nil {
		return fmt.Errorf("unsubscribe: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Config.UserAccessToken) // MUST be a user token for WebSocket subs
	req.Header.Set("Client-Id", c.Config.AppClientID)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNoContent:
		fmt.Println("Successfully delete")

		//db.Purge(user)

	case http.StatusBadRequest:
		fmt.Println("400: Bad Delete")

	case http.StatusUnauthorized:
		fmt.Println("401: Unauthorized")

	case http.StatusNotFound:
		fmt.Println("404: Subscription Not found")

	}

	return nil
}

func (c *TwitchClient) Health(ctx context.Context) error {
	return nil
}

func (c *TwitchClient) Name() string {
	return "twitch"
}

// ___________ Private Helper Functions ____________

// Refreshes the bots User Access Token and Refresh Token
func (c *TwitchClient) refreshTokens() error {

	//Get new App Access Token and store it to config
	err := auth.RefreshAppAccessToken(c.Config, c.HTTP)
	if err != nil {
		return fmt.Errorf("RefreshTokens: %w", err)
	}

	newTokens, err := auth.RefreshUserAccessToken(c.Config.UserRefreshToken, c.Config, c.HTTP)
	if err != nil {
		return fmt.Errorf("RefreshTokens: %w", err)
	}

	c.Config.UserAccessToken = newTokens.UserAccessToken
	c.Config.UserRefreshToken = newTokens.UserRefreshToken

	if err := config.StoreTwitchConfig(c.Config); err != nil {
		return fmt.Errorf("RefreshToken: %w", err)
	}

	return nil
}

// Loads all channels into config
func (c *TwitchClient) loadChannels() error {

	//Get Users from DB
	userData, err := auth.LoadUserData()
	if err != nil {
		return fmt.Errorf("LoadChannels: %w", err)
	}

	channels := make(map[string]*TwitchChannel)
	for _, user := range userData {
		var channel TwitchChannel
		channel.ID = user.UserID
		channel.Username = user.Username
		channels[channel.ID] = &channel
	}

	c.SessionData.Channels = channels
	return nil
}

// Looping function to read chat messages in and parse them
func (c *TwitchClient) listen() {

	fmt.Println("TwitchBot is listening")

	for {
		messageType, data, err := c.WS.ReadMessage()
		if err != nil {
			log.Println(err)
		}

		if messageType == websocket.TextMessage {
			var event EventSubMessage
			if err := json.Unmarshal(data, &event); err != nil {
				fmt.Println(fmt.Errorf("json unmarshal: %w", err))
				fmt.Println(string(data))
			}

			c.handleEvent(event)
			continue
		}

		// If Socket closes randomly, or we read in an error
		if messageType == websocket.CloseMessage || messageType == -1 {
			conn, sessiondata, err := c.newWebsocketConn(EventSubAPIBaseURL)
			if err != nil {
				fmt.Println("establish new connection failure: %w", err)
			}

			c.WS = conn
			c.SessionData.SessionID = sessiondata.SessionID
			c.SessionData.KeepAliveTimeout = sessiondata.KeepAliveTimeout
			continue
		}

	}
}

func (c *TwitchClient) handleEvent(event EventSubMessage) {

	switch event.Metadata.MessageType {
	case "session_welcome":
		return

	case "session_keepalive":
		//Add a timer to see if the socket is alive and healthy
		return

	case "session_reconnect":

		fmt.Println("Attempting to reconnect")
		err := c.ReconnectWebsocket(event)

		if err != nil {
			fmt.Println("error in handleEvent session_reconnect %w", err)
			panic(err)
		}
		return

	case "revocation":

		fmt.Println("Revocation")
		//youll receive the message once and then no longer receive messages for the specified user and subscription type
		// check status field on how to handle

		switch event.Payload.Subscription.Status {

		case "user_removed":
			// user_removed -> user mentioned in the subscription no longer exists. ( Channel banned or whatver)
			fmt.Println("User Removed")
			return

		case "authorization_revoked":
			// authorization_revoked -> user revoked the authorization token that the subscription relied on (user removed bot permissions), remove user from subscription list
			fmt.Println("Auth Revoked")
			return

		case "version_removed":
			// version_removed -> the subscribed to subscription type and version is no longer supported
			fmt.Println("Version removed")
			return
		}
		return

	case "notification":
		//check if message is a command

		switch event.Payload.Subscription.Type {
		case "channel.chat.message":

			if event.Payload.Event.Message.Text[0] == '!' {
				cmd, body := parseCommand(event.Payload.Event.Message.Text)
				var envelope adapter.Envelope = pack(&event, cmd, body)

				c.events <- envelope
			}
			fmt.Printf("%s @%s: %s\n", event.Payload.Event.BroadcasterUserName, event.Payload.Event.ChatterUserName, event.Payload.Event.Message.Text)

		default:
			fmt.Printf("Unknown notification type '%s'\n", event.Payload.Subscription.Type)
			//log error to db
		}

	default:

	}
}

func (c *TwitchClient) helixAPI(ctx context.Context, method string, path string, in any, out any) error {
	return nil
}

func (c *TwitchClient) postHelixJSON(ctx context.Context, path string, in any, out any) error {

	const base = "https://api.twitch.tv/helix"

	b, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+path, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("posthelixjson: %w", err)
	}
	c.applyAuthHeaders(req)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		if out != nil && resp.ContentLength != 0 {
			if err := json.NewDecoder(resp.Body).Decode(out); err != nil && !errors.Is(err, io.EOF) {
				return fmt.Errorf("decode success body: %w", err)
			}
		}
		return nil
	}

	var herr helixError
	_ = json.NewDecoder(resp.Body).Decode(&herr)
	if resp.StatusCode == http.StatusUnauthorized || strings.EqualFold(herr.Error, "Unauthorized") {
		return ErrMissingChannelBot
	}

	if herr.Error != "" || herr.Message != "" {
		return fmt.Errorf("helix %d %s: %s", resp.StatusCode, herr.Error, herr.Message)
	}
	return fmt.Errorf("helix %d", resp.StatusCode)

}

func (c *TwitchClient) applyAuthHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+c.Config.AppAccessToken) // App Access Token for sending messages out
	req.Header.Set("Client-Id", c.Config.AppClientID)
	req.Header.Set("Content-Type", "application/json")
}

func pack(msg *EventSubMessage, command string, body string) adapter.Envelope {
	var newEnvelope adapter.Envelope

	rawTime, _ := time.Parse(time.RFC3339Nano, msg.Metadata.MessageTimestamp)
	var timestamp string = (rawTime).Format("2006-01-02 15:04:05.00")

	newEnvelope.Platform = "twitch"
	newEnvelope.Username = msg.Payload.Event.ChatterUserName
	newEnvelope.UserID = msg.Payload.Event.ChatterUserID
	newEnvelope.ChannelName = msg.Payload.Event.BroadcasterUserName
	newEnvelope.ChannelID = msg.Payload.Event.BroadcasterUserID
	newEnvelope.Command = command
	newEnvelope.Content = body
	newEnvelope.Timestamp = rawTime

	fmt.Printf("[%s] %s: @%s %s\n", timestamp, newEnvelope.ChannelName, newEnvelope.Username, body)

	return newEnvelope
}

func parseCommand(msg string) (string, string) {

	i := strings.IndexByte(msg, ' ')
	if i == -1 {
		return msg, ""
	} else {
		return msg[:i], msg[i+1:]
	}

}
