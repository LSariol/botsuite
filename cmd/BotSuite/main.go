package main

import (
	"context"
	"log"
	"os/signal"

	"github.com/lsariol/botsuite/internal/adapters/twitch"
	"github.com/lsariol/botsuite/internal/app"
	"github.com/lsariol/botsuite/internal/app/registry"
	"github.com/lsariol/botsuite/internal/app/router"
)

func main() {

	// Create Dependencies
	deps, err := app.NewDependencies()
	if err != nil {
		log.Fatal("Error loading deps %w")
	}

	ctx, stop := signal.NotifyContext(context.Background())
	defer stop()

	// Create Registry
	var register *registry.Registry = registry.NewRegistry()
	registry.RegisterAll(register)

	// Create Router
	router := router.NewRouter(ctx, register)

	// Create TwitchClient
	var twitchBot twitch.TwitchClient = *twitch.NewTwitchBot(deps.HTTP, &deps.Config.Twitch)
	if err := twitchBot.Init(); err != nil {
		log.Fatal(err)
	}

	// Start Threads

	// Start Router
	go router.Run(ctx, deps)

	// Start TwitchBot Reading
	go twitchBot.Read()

	go func() {
		for env := range twitchBot.Envelopes() {
			router.Inbound() <- env
		}
	}()

	go func() {
		for resp := range router.Outbound() {

			switch resp.Platform {
			case "twitch":
				twitchBot.Chew(resp)
			case "discord":
				//discord Chew
				continue
			}

		}
	}()

	select {}
}
