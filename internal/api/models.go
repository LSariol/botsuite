package api

import "github.com/jackc/pgx/v5/pgxpool"

type Config struct {
	Address string
	Port    string
}

type APIStore struct {
	pool             *pgxpool.Pool
	connectionString string
}
