package registry

import (
	"fmt"

	"github.com/lsariol/botsuite/internal/commands"
)

type Registry struct {
	byName map[string]commands.Command
}

func NewRegistry() *Registry {
	return &Registry{
		byName: make(map[string]commands.Command),
	}
}

func (r *Registry) Register(cmd commands.Command) {

	//Creating a helper function called add
	add := func(k string) {
		if _, exists := r.byName[k]; exists {
			panic(fmt.Sprintf("duplicate command registration %s", k))
		}

		r.byName[k] = cmd
	}

	add(cmd.Name())
	for _, a := range cmd.Aliases() {
		add(a)
	}
}

func (r *Registry) Get(name string) (commands.Command, bool) {
	c, ok := r.byName[name]
	return c, ok
}

func (r *Registry) GetAll() map[string]commands.Command {
	return r.byName
}
