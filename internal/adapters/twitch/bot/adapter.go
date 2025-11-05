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
	twitchdb "github.com/lsariol/botsuite/internal/adapters/twitch/database"
	"github.com/lsariol/botsuite/internal/config"
)

// CLI URLS
// const (
// 	HelixBaseURL       = "https://api.twitch.tv/helix"
// 	EventSubAPIBaseURL = "http://127.0.0.1:8080/eventsub/subscriptions"
// 	EventSubWSURL      = "ws://127.0.0.1:8080/ws"
// )

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
	DB          *twitchdb.Store
	events      chan adapter.Envelope
}

func NewTwitchBot(client *http.Client, cfg *config.TwitchConfig, dbStore *twitchdb.Store) *TwitchClient {
	return &TwitchClient{
		HTTP:   client,
		Config: cfg,
		DB:     dbStore,
		events: make(chan adapter.Envelope, 100),
	}
}

func (c *TwitchClient) Initilize(ctx context.Context) error {
	if err := c.refreshTokens(); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	if err := c.loadChannels(); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	conn, sessionData, err := c.dialWebsocket(ctx, EventSubWSURL)
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}
	c.WS = conn
	c.SessionData.SessionID = sessionData.SessionID
	c.SessionData.KeepAliveTimeout = sessionData.KeepAliveTimeout

	var channelIDs []string
	for _, c := range c.SessionData.Channels {
		channelIDs = append(channelIDs, c.ID)
	}

	c.Join(ctx, channelIDs)

	return nil
}

func (c *TwitchClient) Run(ctx context.Context) error {

	if err := c.Initilize(ctx); err != nil {
		return fmt.Errorf("initilization error: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("context canceled, stopping read loop")
			return nil

		default:
			messageType, data, err := c.WS.ReadMessage()
			if err != nil {
				fmt.Printf("read error: %q\n", err)
				c.hardResetWS(ctx, EventSubWSURL)
				continue
			}

			switch {
			case messageType == websocket.TextMessage:
				if messageType == websocket.TextMessage {
					var event EventSubMessage
					if err := json.Unmarshal(data, &event); err != nil {
						fmt.Println(fmt.Errorf("json unmarshal: %w", err))
						fmt.Println(string(data))
					}

					c.ConsumeEvent(ctx, event)
					continue
				}

			case messageType == websocket.CloseMessage || messageType == -1:
				log.Println("Read error, should be wsarecv")
				log.Println("messagetype == websocket.CloseMessage || messageType == -1 ERROR PATH.")
				log.Printf("%d: %s", messageType, err.Error())
				c.hardResetWS(ctx, EventSubWSURL)
				continue

			default:
				log.Printf("unknown message type: %d", messageType)
				log.Println("broken input: " + string(data))
			}

		}
	}
}

// Adapter Functions
func (c *TwitchClient) Shutdown(ctx context.Context) error {
	return nil
}

func (c *TwitchClient) Restart(ctx context.Context) error {

	return nil
}

func (c *TwitchClient) OutBoundEvents() <-chan adapter.Envelope {
	return c.events
}

func (c *TwitchClient) ConsumeEvent(ctx context.Context, event EventSubMessage) {

	switch event.Metadata.MessageType {
	case "session_welcome":
		return

	case "session_keepalive":
		//Add a timer to see if the socket is alive and healthy
		return

	case "session_reconnect":

		fmt.Println("Attempting to reconnect")
		err := c.resetWS(ctx, event)

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

			log.Printf("%s: @%s %s\n", event.Payload.Event.BroadcasterUserName, event.Payload.Event.ChatterUserName, event.Payload.Event.Message.Text)

			if event.Payload.Event.Message.Text[0] == '!' {
				cmd, body := ConsumeMessage(event.Payload.Event.Message.Text)
				var envelope adapter.Envelope = pack(&event, cmd, body)

				c.events <- envelope
			}

		default:
			fmt.Printf("Unknown notification type '%s'\n", event.Payload.Subscription.Type)
			//log error to db
		}

	default:

	}
}

func ConsumeMessage(msg string) (string, string) {

	i := strings.IndexByte(msg, ' ')
	if i == -1 {
		return msg, ""
	} else {
		return msg[:i], msg[i+1:]
	}

}

func (c *TwitchClient) DeliverResponse(ctx context.Context, r adapter.Response) error {

	reqBody := chatMessageReq{
		BroadcasterID: r.ChannelID,
		SenderID:      c.Config.BotID,
		Message:       r.Text,
	}

	var out any
	if err := c.postHelixJSON(ctx, "/chat/messages", reqBody, &out); err != nil {
		if errors.Is(err, ErrMissingChannelBot) {
			c.Leave(ctx, []string{r.ChannelID})
		}
		return err
	}
	return nil

}

func (c *TwitchClient) Join(ctx context.Context, targetIDs []string) error {

	for _, channelID := range targetIDs {

		body := map[string]any{
			"type":    "channel.chat.message",
			"version": "1",
			"condition": map[string]string{
				"broadcaster_user_id": channelID,
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
			c.SessionData.Channels[channelID].SubscriptionID = sub.ID
		}
	}

	return nil
}

func (c *TwitchClient) Leave(ctx context.Context, targets []string) error {

	for _, channelID := range targets {

		subscriptionID := c.SessionData.Channels[channelID].SubscriptionID
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
	}

	return nil
}

func (c *TwitchClient) Health(ctx context.Context) adapter.HealthStatus {
	status := adapter.HealthStatus{}

	return status
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

func (c *TwitchClient) helixAPI(ctx context.Context, method string, path string, in any, out any) error {
	return nil
}

func (c *TwitchClient) postHelixJSON(ctx context.Context, path string, in any, out any) error {

	b, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, HelixBaseURL+path, bytes.NewReader(b))
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
	newEnvelope.Command = strings.ToLower(command)
	newEnvelope.Content = body
	newEnvelope.Timestamp = rawTime

	fmt.Printf("[%s] %s: @%s %s\n", timestamp, newEnvelope.ChannelName, newEnvelope.Username, body)

	return newEnvelope
}
