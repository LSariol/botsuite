package database

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

func (d *Database) AddChannel(ctx context.Context, ch Channel) error {
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

	_, err := d.Pool.Exec(ctx, query, ch.UserID, ch.Username, ch.AccessToken, ch.RefreshToken)
	if err == nil {
		err = d.LogEvent(ctx, ch, "grant")
	}

	return err
}

func RemoveChannel() error {
	return nil
}

func BanChannel() error {
	return nil
}

func GetAllChannels() error {

	return nil
}

func (d *Database) LogEvent(ctx context.Context, ch Channel, action string) error {
	query := `
	INSERT INTO botsuite.twitch_channel_events (user_id, username, action)
	VALUES ($1, $2, $3)
	`
	_, err := d.Pool.Exec(ctx, query, ch.UserID, ch.Username, action)
	if err != nil {
		err = fmt.Errorf("error logging event: %w", err)
	}
	return err
}
