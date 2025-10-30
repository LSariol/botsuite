package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool       *pgxpool.Pool
	ConnString string
}

func NewDatabase() *Database {

	return &Database{
		ConnString: os.Getenv("BOTSUITE_DATABASE_URL"),
	}
}

func (d *Database) Connect(ctx context.Context) error {

	connString := fmt.Sprintf(d.ConnString, os.Getenv("BOTSUITE_DATABASE_USERNAME"), os.Getenv("BOTSUITE_DATABASE_PASSWORD"), os.Getenv("BOTSUITE_DATABASE_NAME"))
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return fmt.Errorf("ping database: %w", err)
	}

	d.Pool = pool
	d.ConnString = connString
	return nil
}

func (d *Database) Close() {
	if d.Pool != nil {
		d.Pool.Close()
	}
}
