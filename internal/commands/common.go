package commands

import "github.com/lsariol/botsuite/internal/adapters/adapter"

func SuppresedReply() (adapter.Response, error) {
	return adapter.Response{SuppressReply: true, Error: true}, nil
}
