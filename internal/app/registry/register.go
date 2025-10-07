package registry

import (
	"github.com/lsariol/botsuite/internal/commands/catfacts"
	"github.com/lsariol/botsuite/internal/commands/help"
	"github.com/lsariol/botsuite/internal/commands/ping"
)

// Register all functions into the registry
func RegisterAll(r *Registry) {

	r.Register(ping.Ping{})
	r.Register(catfacts.CatFact{})
	r.Register(help.Help{})
}
