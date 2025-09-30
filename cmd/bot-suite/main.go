package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Create Registry
	var register *registry.Registry = registry.NewRegistry()
	registry.RegisterAll(register)

	// Create Router
	router := router.NewRouter(ctx, register)

	// Create TwitchClient
	var twitchClient twitch.TwitchClient = *twitch.NewTwitchBot(deps.HTTP, &deps.Config.Twitch)

	// Start Threads

	// Start Router
	go router.Run(ctx, deps)

	// Start TwitchBot Reading
	go twitchClient.Run(ctx)

	go func() {
		for env := range twitchClient.Events() {
			router.Inbound() <- env
		}
	}()

	go func() {
		for resp := range router.Outbound() {

			switch resp.Platform {
			case "twitch":
				twitchClient.Deliver(ctx, resp)
			case "discord":
				//discord Chew
				continue
			}

		}
	}()

	<-ctx.Done()

}
