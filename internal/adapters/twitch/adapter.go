package twitch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lsariol/botsuite/internal/adapters/twitch/auth"
	"github.com/lsariol/botsuite/internal/app/event"
	"github.com/lsariol/botsuite/internal/bot"
	"github.com/lsariol/botsuite/internal/config"
)

type TwitchClient struct {
	bot.Bot
	HTTP        *http.Client
	WS          *websocket.Conn
	Config      *config.TwitchConfig
	SessionData SessionData
	msgs        chan event.Envelope
}

func NewTwitchBot(client *http.Client, cfg *config.TwitchConfig) *TwitchClient {
	return &TwitchClient{
		HTTP:   client,
		Config: cfg,
		msgs:   make(chan event.Envelope, 100),
	}
}

func (t *TwitchClient) Run() {

}

// Initilizes TwitchClient object with values and establishes connections and subscriptions
func (t *TwitchClient) Init() error {

	if err := t.RefreshTokens(); err != nil {
		return fmt.Errorf("Init: %w", err)
	}

	if err := t.LoadChannels(); err != nil {
		return fmt.Errorf("Init: %w", err)
	}

	if err := t.EstablishConnection(); err != nil {
		return fmt.Errorf("Init: %w", err)
	}

	//for channel in channels
	for _, channel := range t.SessionData.Channels {
		if err := t.Subscribe(channel.ID); err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("Connected to %s\n", channel.Username)
	}

	return nil

}

// Connects to twitch Websocket
func (t *TwitchClient) EstablishConnection() error {

	url := "wss://eventsub.wss.twitch.tv/ws"

	// 	//Connect
	// 	//websocket <- calling package
	// 	//DefaultDialer <- a dialer with all fields set to default values
	// 	//Dial <- creates a new client connection
	// 	//returns: *websocket.conn, *http.response, err
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("dial error: %w", err)
	}

	t.WS = conn

	_, message, err := t.WS.ReadMessage()
	if err != nil {
		return fmt.Errorf("read error: %w", err)
	}

	var event EventSubMessage
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("json unmarshal error: %w", err)
	}

	t.SessionData.SessionID = event.Payload.Session.ID
	t.SessionData.KeepAliveTimeout = event.Payload.Session.KeepaliveTimeoutSeconds

	return nil
}

// Disconnects from a channelId
func (t *TwitchClient) Disconnect() {

}

// Sends msg to channelID
func (t *TwitchClient) SendChat(msg string, channelID string) error {

	payload := struct {
		BroadcasterID string `json:"broadcaster_id"`
		SenderID      string `json:"sender_id"`
		Message       string `json:"message"`
	}{
		BroadcasterID: channelID,
		SenderID:      t.Config.BotID,
		Message:       msg,
	}

	b, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "https://api.twitch.tv/helix/chat/messages", bytes.NewReader(b))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+t.Config.AppAccessToken)
	req.Header.Set("Client-Id", t.Config.AppClientID)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := t.HTTP.Do(req)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("sendchat failed: %s: %w", body, err)
	}

	return nil
}

// Looping function to read chat messages in and parse them
func (t *TwitchClient) Read() {

	fmt.Println("TwitchBot is reading")
	for {
		_, msg, err := t.WS.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		var message EventSubMessage
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("convert error: ", err)
		}

		switch message.Metadata.MessageType {
		case "session_welcome":
			continue
		case "session_keepalive":
			//Add a timer to see if the socket is alive and healthy
			continue
		case "session_reconnect":
			//connect to new url
			continue
		case "revocation":
			//youll receive the message once and then no longer receive messages for the specified user and subscription type
			// check status field on how to handle
			switch message.Payload.Subscription.Status {
			case "user_removed":
				// user_removed -> user mentioned in the subscription no longer exists. ( Channel banned or whatver)
			case "authorization_revoked":
				// authorization_revoked -> user revoked the authorization token that the subscription relied on (user removed bot permissions), remove user from subscription list
			case "version_removed":
				// version_removed -> the subscribed to subscription type and version is no longer supported
			}
			continue
		case "notification":
			//check if message is a command

			switch message.Payload.Subscription.Type {
			case "channel.chat.message":

				if message.Payload.Event.Message.Text[0] == '!' {
					cmd, body := parseCommand(message.Payload.Event.Message.Text)
					var envelope event.Envelope = Pack(&message, cmd, body)

					t.msgs <- envelope
				}
				fmt.Printf("%s @%s: %s\n", message.Payload.Event.BroadcasterUserName, message.Payload.Event.ChatterUserName, message.Payload.Event.Message.Text)

			default:
				fmt.Printf("Unknown notification type '%s'\n", message.Payload.Subscription.Type)
				//log error to db
			}
		default:

		}

	}
}

// Subscribes to a chat event
func (t *TwitchClient) Subscribe(channelID string) error {

	body := map[string]any{
		"type":    "channel.chat.message",
		"version": "1",
		"condition": map[string]string{
			"broadcaster_user_id": channelID,
			"user_id":             t.Config.BotID,
		},
		"transport": map[string]string{
			"method":     "websocket",
			"session_id": t.SessionData.SessionID,
		},
	}

	buf, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", "https://api.twitch.tv/helix/eventsub/subscriptions", bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("Subscribe: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+t.Config.UserAccessToken) // MUST be a user token for WebSocket subs
	req.Header.Set("Client-Id", t.Config.AppClientID)
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	return nil
}

// Closes all connections
func (t *TwitchClient) Close() {

}

// Loads all channels into config
func (t *TwitchClient) LoadChannels() error {

	//Get Users from DB
	userData, err := auth.LoadUserData()
	if err != nil {
		return fmt.Errorf("LoadChannels: %w", err)
	}

	var channels []TwitchChannel
	for _, user := range userData {
		var channel TwitchChannel
		channel.ID = user.UserID
		channel.Username = user.Username
		channels = append(channels, channel)
	}

	t.SessionData.Channels = channels
	return nil
}

// Handles sending out commands
func (t *TwitchClient) Command(msg string) error {

	return nil
}

func Pack(msg *EventSubMessage, command string, body string) event.Envelope {
	var newEnvelope event.Envelope

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

func (t *TwitchClient) RefreshTokens() error {

	//Get new App Access Token and store it to config
	err := auth.RefreshAppAccessToken(t.Config, t.HTTP)
	if err != nil {
		return fmt.Errorf("RefreshTokens: %w", err)
	}

	newTokens, err := auth.RefreshUserAccessToken(t.Config.UserRefreshToken, t.Config, t.HTTP)
	if err != nil {
		return fmt.Errorf("RefreshTokens: %w", err)
	}

	t.Config.UserAccessToken = newTokens.UserAccessToken
	t.Config.UserRefreshToken = newTokens.UserRefreshToken

	if err := config.StoreTwitchConfig(t.Config); err != nil {
		return fmt.Errorf("RefreshToken: %w", err)
	}

	return nil
}

func (t *TwitchClient) Envelopes() <-chan event.Envelope {
	return t.msgs
}

func (t *TwitchClient) Chew(msg event.Response) {

	fmt.Println("Twitch has received the response")

	t.SendChat(msg.Text, msg.ChannelID)
}
