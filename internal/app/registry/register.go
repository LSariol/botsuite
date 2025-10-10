package registry

import (
	"github.com/lsariol/botsuite/internal/commands/catfacts"
	"github.com/lsariol/botsuite/internal/commands/help"
	"github.com/lsariol/botsuite/internal/commands/ping"
	"github.com/lsariol/botsuite/internal/commands/randomdog"
	"github.com/lsariol/botsuite/internal/commands/request"
	"github.com/lsariol/botsuite/internal/commands/uptime"
)

// Register all functions into the registry
func RegisterAll(r *Registry) {

	r.Register(catfacts.CatFact{})
	r.Register(help.Help{})
	r.Register(ping.Ping{})
	r.Register(randomdog.RandomDog{})
	r.Register(request.Request{})
	r.Register(uptime.UpTime{})

}
