package twitch

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

func (c *TwitchClient) resetWS(reconnectmsg EventSubMessage) error {

	reconnectUrl := reconnectmsg.Payload.Session.ReconnectURL
	newConn, sessionData, err := c.dialWebsocket(reconnectUrl)
	if err != nil {
		return fmt.Errorf("reconnectwebsocket error: %w", err)
	}

	//Send a close frame to twitch to prevent more messages being sent here.
	_ = c.WS.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "reconnecting"), time.Now().Add(time.Second))

	var event EventSubMessage
	for {
		if err := c.WS.ReadJSON(&event); err != nil {

			if c.WS != nil {
				c.WS.Close()
				c.WS = nil
			}

			c.WS = newConn
			c.SessionData.SessionID = sessionData.SessionID
			c.SessionData.KeepAliveTimeout = sessionData.KeepAliveTimeout

			return nil
		}

		c.handleEvent(event)
	}

}

func (c *TwitchClient) dialWebsocket(URL string) (*websocket.Conn, SessionData, error) {

	var sD SessionData

	newConn, _, err := websocket.DefaultDialer.Dial(URL, nil)
	if err != nil {
		return nil, sD, fmt.Errorf("newwebsocketconn dial error: %w", err)
	}

	var event EventSubMessage
	if err := newConn.ReadJSON(&event); err != nil {
		return nil, sD, fmt.Errorf("newwebsocketconn read json: %w", err)
	}

	for event.Metadata.MessageType != "session_welcome" {
		c.handleEvent(event)

		if err := newConn.ReadJSON(&event); err != nil {
			return nil, sD, fmt.Errorf("newwebsocketconn read json: %w", err)
		}
	}

	sD = SessionData{
		SessionID:        event.Payload.Session.ID,
		KeepAliveTimeout: event.Payload.Session.KeepaliveTimeoutSeconds,
	}

	return newConn, sD, nil
}

func (c *TwitchClient) hardResetWS(ctx context.Context, URL string) error {

	newConn, _, err := websocket.DefaultDialer.Dial(URL, nil)
	if err != nil {
		newConn.Close()
		return fmt.Errorf("newwebsocketconn dial error: %w", err)
	}

	var event EventSubMessage
	if err := newConn.ReadJSON(&event); err != nil {
		newConn.Close()
		return fmt.Errorf("newwebsocketconn read json: %w", err)
	}

	if c.WS != nil {
		_ = c.WS.Close()
		c.WS = nil
	}

	c.SessionData.KeepAliveTimeout = event.Payload.Session.KeepaliveTimeoutSeconds
	c.SessionData.SessionID = event.Payload.Session.ID
	c.WS = newConn
	c.refreshTokens()
	c.JoinAllChannels(ctx)
	return nil
}
