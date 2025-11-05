package twitchdb

import (
	"context"
	"fmt"
)

type Channel struct {
	UserID       string
	Username     string
	AccessToken  string
	RefreshToken string
}

func (s *Store) AddChannel(ctx context.Context, ch Channel) error {
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

func (s *Store) RemoveChannel() error {
	return nil
}

func (s *Store) BanChannel() error {
	return nil
}

func (s *Store) GetAllChannels(ctx context.Context) error {

	query := `
	SELECT user_id, username
	FROM botsuite.twitch_channels
	WHERE active = TRUE;
	`
	println(query)
	return nil
}

func (s *Store) LogEvent(ctx context.Context, ch Channel, action string) error {
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
