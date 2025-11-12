package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/lsariol/botsuite/internal/adapters/adapter"
)

func (c *ChatClient) SendChat(ctx context.Context, r adapter.Response) (*SendChatMessageResponse, *HelixError, error) {

	const url = "https://api.twitch.tv/helix/chat/messages"

	body := SendChatBody{
		BroadcasterID: r.ChannelID,
		SenderID:      c.Config.BotID,
		Message:       r.Text,
	}

	payload, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Config.AppAccessToken)
	req.Header.Set("Client-Id", c.Config.AppClientID)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var helixErr HelixError
		if err := json.Unmarshal(respBody, &helixErr); err != nil {
			return nil, nil, fmt.Errorf("decode helix error: %w", err)
		}

		return nil, &helixErr, nil

	}
	// Handle successful 200 response
	var data SendChatMessageResponse
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, nil, fmt.Errorf("decode success response: %w", err)
	}

	return &data, nil, nil
}
