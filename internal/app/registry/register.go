package registry

import "github.com/lsariol/botsuite/internal/commands"

//Register all functions into the registry
func RegisterAll(r *Registry) {

	r.Register(commands.Ping{})
}
