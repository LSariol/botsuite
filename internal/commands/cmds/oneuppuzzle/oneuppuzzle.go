package oneuppuzzle

import (
	"context"
	"fmt"
	"time"

	"github.com/lsariol/botsuite/internal/adapters/adapter"
	"github.com/lsariol/botsuite/internal/app/dependencies"
)

type OneUpPuzzle struct{}

func (OneUpPuzzle) Name() string      { return "oneuppuzzle" }
func (OneUpPuzzle) Aliases() []string { return []string{"oneup"} }
func (OneUpPuzzle) Regexes() []string {
	return []string{`Congratulations,\s*you finished puzzle #\d+\s+in\s+(?:\d+\s*hour[s]?,\s*)?(?:\d+\s*minute[s]?,\s*)?(?:\d+\s*second[s]?)`}
}
func (OneUpPuzzle) Description() string    { return "command for oneuppuzzle statistics." }
func (OneUpPuzzle) Usage() string          { return "oneuppuzzle <none> <help> <leaders> <username> <stats>" }
func (OneUpPuzzle) Timeout() time.Duration { return 3 * time.Second }

func (OneUpPuzzle) Execute(ctx context.Context, e adapter.Envelope, deps *dependencies.Deps) (adapter.Response, error) {

	if e.IsRegex {
		//store in DB
		game, err := validateGameEntry(e, deps)
		if err != nil {
			return adapter.Response{SuppressReply: true, Error: true}, err
		}

		if 3600 <= game.Details.TimeSeconds || game.Details.TimeSeconds <= 20 {
			err = storeFlaggedGame(game, deps)
			if err != nil {
				return adapter.Response{Text: "An error occured while storing your game. Please try again."}, err
			}
			return adapter.Response{Text: "Game has been recorded"}, nil
		} else {
			err = storeGame(game, deps)
			if err != nil {
				return adapter.Response{Text: "An error occured while storing your game. Please try again."}, err
			}
			return adapter.Response{Text: "Game has been recorded"}, nil
		}
	}

	switch e.Command {
	case "oneuppuzzle", "oneup":

		if len(e.Args) == 0 {
			return adapter.Response{Text: "https://www.oneuppuzzle.com"}, nil
		}

		switch e.Args[0] {
		case "help":
			return adapter.Response{Text: "Valid arguments for the !oneup command are: help, leaders, stats, or <username>"}, nil
		case "leaders":

			//if "!oneup leaders" - return top 3
			// if "!oneup leaders <int>" return top int players

			if len(e.Args) >= 2 {
				//Get X ammount of leaders
			}

			//return top 3 fastest averages
		case "stats":

			switch {
			case len(e.Args) == 1:
				cp, ok, err := getChannelStats(e.ChannelID, deps)
				if err != nil {
					return adapter.Response{Text: "An error occured while fetching channel stats. Please try again."}, nil
				}

				if !ok {
					return adapter.Response{Text: "This channel does not have any tracked games."}, nil
				}

				return adapter.Response{Text: fmt.Sprintf("Channel Stats: Total Games %d | Average %s | Fastest Game %s (%s) | Slowest Game %s (%s)", *cp.GamesCompleted, formatTime(*cp.Completions.AverageTime), formatTime(*cp.Completions.FastestTime), *cp.FastestUser, formatTime(*cp.Completions.SlowestTime), *cp.SlowestUser)}, nil

			case len(e.Args) == 2:

				up, ok, err := getUserStats(e.Args[1], deps)
				if err != nil {
					return adapter.Response{Text: "An error occured while fetching user stats. Please try again."}, nil
				}

				if !ok {
					return adapter.Response{Text: "This user does not have any tracked games."}, nil
				}

				return adapter.Response{Text: fmt.Sprintf("Stats for @%s: Total Games %d | Average  %s | Fastest Game %s | Slowest Game %s", e.Args[1], *up.GamesCompleted, formatTime(*up.Completions.AverageTime), formatTime(*up.Completions.FastestTime), formatTime(*up.Completions.SlowestTime))}, nil
			}

			// return stats from the bot
			// total games, channel average, slowest, fastest
		default:
			//try to look up player stats, otherwise
			return adapter.Response{Text: "Invalid Command"}, nil
		}
	}

	// help - return usage
	// leaders - return top 3 fastest averages
	// username - return total completed, average time, longest time, shortest time
	// stats - total games tracked, channel average, slowest, fastest
	return adapter.Response{Text: "Im not even sure how you managed to input a command to get here... but you did. "}, nil
}

func validateGameEntry(e adapter.Envelope, deps *dependencies.Deps) (OneUpGame, error) {

	// Attempt to parse the game entry. If Error in parsing, its invalid
	newGame, err := packageGame(e)
	if err != nil {
		return newGame, err
	}

	//Need some kind of check here to see if the game is the correct numbered game.

	// Validate the game ID as well as if the user has sumbitted this before
	recorded, err := isPuzzleRecorded(e.UserID, newGame.Details.PuzzleID, deps)
	if err != nil {
		return newGame, fmt.Errorf("ispuzzlerecorded: %w", err)
	}

	if recorded {
		return newGame, fmt.Errorf("recorded game")
	}

	return newGame, nil
}

func formatTime(seconds int) string {

	h := seconds / 3600
	m := (seconds % 3600) / 60
	s := seconds % 60

	if h == 0 {
		if m == 0 {
			return fmt.Sprintf("%d seconds", s)
		}
		return fmt.Sprintf("%d minutes, %d seconds", m, s)
	}

	return fmt.Sprintf("%d hours, %d minutes, %d seconds", h, m, s)

}
