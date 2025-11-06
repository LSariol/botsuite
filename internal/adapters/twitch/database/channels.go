package twitchdb

import (
	"context"
	"fmt"
)

func (s *Store) AddChannel(ctx context.Context, ch ChannelInsert) error {
	query := `
	INSERT INTO botsuite.twitch_channels (user_id, username, access_token, refresh_token)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (user_id) DO UPDATE
		SET username = EXCLUDED.username,
			access_token = EXCLUDED.access_token,
			refresh_token = EXCLUDED.refresh_token,
			active = TRUE,
			times_joined = botsuite.twitch_channels.times_joined + 1,
			last_joined = NOW();`

	_, err := s.pool.Exec(ctx, query, ch.UserID, ch.Username, ch.AccessToken, ch.RefreshToken)
	if err == nil {
		err = s.LogEvent(ctx, ch, "grant")
	}

	return err
}

func (s *Store) RemoveChannel(ctx context.Context, channelID string) error {
	return nil
}

func (s *Store) BanChannel(ctx context.Context, channelID string) error {
	return nil
}

func (s *Store) GetAllChannels(ctx context.Context) (map[string]*TwitchChannel, error) {

	query := `
	SELECT tc.user_id, tc.username, ts.command_prefix
	FROM botsuite.twitch_channels AS tc
	LEFT JOIN botsuite.twitch_settings AS ts
	ON tc.user_id = ts.user_id
	WHERE tc.active = TRUE;
	`
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query getallchannels: %w", err)
	}
	defer rows.Close()

	channels := make(map[string]*TwitchChannel)

	for rows.Next() {
		var channel TwitchChannel
		if err := rows.Scan(&channel.ID, &channel.Username, &channel.CommandPrefix); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		channels[channel.ID] = &channel
	}

	return channels, nil
}

// Action can either be grant or revoke
func (s *Store) LogEvent(ctx context.Context, ch ChannelInsert, action string) error {
	query := `
	INSERT INTO botsuite.twitch_channel_events (user_id, username, action)
	VALUES ($1, $2, $3)
	`
	_, err := s.pool.Exec(ctx, query, ch.UserID, ch.Username, action)
	if err != nil {
		err = fmt.Errorf("error logging event: %w", err)
	}
	return err
}
