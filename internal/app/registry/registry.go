package registry

import (
	"fmt"

	"github.com/lsariol/botsuite/internal/commands"
)

type Registry struct {
	commands map[string]commands.Command
	readOnly map[string]string
}

func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]commands.Command),
		readOnly: make(map[string]string),
	}
}

func (r *Registry) Register(cmd commands.Command) {

	//Creating a helper function called add
	add := func(k string, t string) {
		if _, exists := r.commands[k]; exists {
			panic(fmt.Sprintf("duplicate command registration %s", k))
		}

		r.commands[k] = cmd
		r.readOnly[k] = t
	}

	add(cmd.Name(), "cmd")
	for _, a := range cmd.Aliases() {
		add(a, "cmd")
	}

	for _, t := range cmd.TriggerPhrases() {
		add(t, "trigger")
	}
}

func (r *Registry) Get(name string) (commands.Command, bool) {

	c, ok := r.commands[name]

	return c, ok
}

func (r *Registry) GetReadMap() map[string]string {

	return r.readOnly
}
