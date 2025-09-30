package bot

import "github.com/lsariol/botsuite/internal/app/event"

type Bot interface {
	Run()
	Init() error
	Chew(msg event.Response)
}
