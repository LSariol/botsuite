package main

import (
	"context"
	"fmt"
	"log"

	"github.com/lsariol/botsuite/internal/adapters/twitch"
	"github.com/lsariol/botsuite/internal/app"
)

func main() {

	deps, err := app.NewDeps()
	if err != nil {
		log.Fatal("Error loading deps %w")
	}

	var twitchBot twitch.TwitchBot = twitch.NewTwitchBot(deps.HTTP, &deps.Config.Twitch)

	if err := twitchBot.Init(); err != nil {
		log.Fatal(err)
	}

	ctx := context.Context(context.Background())
	var registry app.Registry = *app.NewRegistry()

	go twitchBot.Read()

	go registry.Run(ctx)

	fmt.Println("Both running")

	go func() {
		for env := range twitchBot.Envelopes() {
			registry.Inbound() <- env
		}
	}()

	for resp := range registry.Outbound() {
		twitchBot.Chew(resp)
	}

	//Idea
	//twitchBot.Run() <- replaces Read()
	//run will start a thread for reading messages.
	//ctx, cancel := context.WithCancel(context.Background())
	//twitchBot.Run()
}
