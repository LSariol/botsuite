package oneuppuzzle

import (
	"context"
	"time"

	"github.com/lsariol/botsuite/internal/adapters/adapter"
	"github.com/lsariol/botsuite/internal/app/dependencies"
)

type OneUpPuzzle struct{}

func (OneUpPuzzle) Name() string             { return "oneuppuzzle" }
func (OneUpPuzzle) Aliases() []string        { return []string{"oneup"} }
func (OneUpPuzzle) TriggerPhrases() []string { return []string{"congratulations, you finished puzzle"} }
func (OneUpPuzzle) Description() string      { return "command for oneuppuzzle statistics." }
func (OneUpPuzzle) Usage() string            { return "oneuppuzzle <none> <help> <leaders> <username> <stats>" }
func (OneUpPuzzle) Timeout() time.Duration   { return 3 * time.Second }

func (OneUpPuzzle) Execute(ctx context.Context, e adapter.Envelope, deps *dependencies.Deps) (adapter.Response, error) {

	//none - send link
	// help - return usage
	// leaders - return top 3 fastest averages
	// username - return total completed, average time, longest time, shortest time
	// stats - total games tracked, channel average, slowest, fastest

	return adapter.Response{Text: "Game has been recorded"}, nil

}
