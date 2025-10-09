package twitch

import (
	"fmt"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

func (c *TwitchClient) ReconnectWebsocket(reconnectmsg EventSubMessage) error {

	reconnectUrl := reconnectmsg.Payload.Session.ReconnectURL
	newConn, err := c.newWebsocketConn(reconnectUrl)
	if err != nil {
		return fmt.Errorf("reconnectwebsocket error: %w", err)
	}

	var event EventSubMessage

	for {
		c.WS.SetReadDeadline(time.Now().Add(1 * time.Millisecond))
		if err := c.WS.ReadJSON(&event); err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// websocket connection is empty
				c.WS.Close()
				c.WS = newConn
				return nil
			}

			// Twitch has sent a close frame with 4004
			if websocket.IsCloseError(err, 4004) {
				c.WS.Close()
				c.WS = newConn

				c.SessionData.SessionID = reconnectmsg.Payload.Session.ID
				c.SessionData.KeepAliveTimeout = reconnectmsg.Payload.Session.KeepaliveTimeoutSeconds
				return nil
			}

			// unacounted for error
			return fmt.Errorf("UNACCOUNTED FOR ERROR IN RECONNECTWEBSOCKET(): %w", err)
		}

		c.handleEvent(event)
	}

}

func (c *TwitchClient) newWebsocketConn(URL string) (*websocket.Conn, error) {

	newConn, _, err := websocket.DefaultDialer.Dial(URL, nil)
	if err != nil {
		return nil, fmt.Errorf("newwebsocketconn dial error: %w", err)
	}

	var event EventSubMessage
	if err := newConn.ReadJSON(&event); err != nil {
		return nil, fmt.Errorf("newwebsocketconn read json: %w", err)
	}

	for event.Metadata.MessageType != "session_welcome" {
		c.handleEvent(event)

		if err := newConn.ReadJSON(&event); err != nil {
			return nil, fmt.Errorf("newwebsocketconn read json: %w", err)
		}
	}

	c.SessionData.SessionID = event.Payload.Session.ID
	c.SessionData.KeepAliveTimeout = event.Payload.Session.KeepaliveTimeoutSeconds
	return newConn, nil
}
