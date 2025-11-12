package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	twitchbot "github.com/lsariol/botsuite/internal/adapters/twitch/bot"
	twitchdb "github.com/lsariol/botsuite/internal/adapters/twitch/database"

	"github.com/lsariol/botsuite/internal/app/dependencies"
	"github.com/lsariol/botsuite/internal/app/registry"
	"github.com/lsariol/botsuite/internal/app/router"
)

func main() {

	// Create Dependencies

	deps := dependencies.New()

	err := deps.Load()
	if err != nil {
		log.Fatal("Error loading deps %w")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Create Database Connection
	if err := deps.DB.Connect(ctx); err != nil {
		log.Fatal(err)
	}

	// Create Registry
	var register *registry.Registry = registry.NewRegistry()
	registry.RegisterAll(register)

	// Create Router
	var router *router.Router = router.NewRouter(ctx, register)

	// Create TwitchClient
	var twitchDBStore *twitchdb.Store = twitchdb.NewStore(deps.DB.Pool, deps.DB.ConnString)
	var twitchClient *twitchbot.TwitchClient = twitchbot.New(deps.HTTP, &deps.Config.Twitch, twitchDBStore, router.Inbound(), register.GetReadMap())

	// // Start Threads

	// Register Adapters with the router
	router.RegisterAdapter(twitchClient)
	// Start Router
	go router.Run(ctx, deps)

	// Start TwitchBot Reading
	go func() {
		err := twitchClient.Run(ctx)
		if err != nil {
			log.Println(err)
			stop()
		}
	}()

	// time.Sleep(3 * time.Second)
	// RunLiveTests(ctx, twitchClient)

	<-ctx.Done()

}

// func RunLiveTests(ctx context.Context, c *twitchbot.TwitchClient) {

// 	r := adapter.Response{
// 		ChannelID: "42217464",
// 		Text:      newMessage(10),
// 	}
// 	r2 := adapter.Response{
// 		ChannelID: "228428175",
// 		Text:      newMessage(10),
// 	}
// 	r3 := adapter.Response{
// 		ChannelID: "410570760",
// 		Text:      newMessage(10),
// 	}

// 	c.DeliverResponse(ctx, r)
// 	c.DeliverResponse(ctx, r2)
// 	c.DeliverResponse(ctx, r3)

// }

// func newMessage(length int) string {
// 	// each byte = 2 hex chars, so divide by 2 (round up)
// 	byteLen := (length + 1) / 2
// 	b := make([]byte, byteLen)

// 	if _, err := io.ReadFull(rand.Reader, b); err != nil {
// 		panic(err)
// 	}

// 	hexStr := hex.EncodeToString(b)

// 	// trim to exact requested length (in case of odd number)
// 	if len(hexStr) > length {
// 		hexStr = hexStr[:length]
// 	}

// 	return hexStr
// }

// func rateTesting(ctx context.Context, c twitchbot.TwitchClient) {

// 	var i int = 1
// 	var messages int = 10
// 	for i <= messages {
// 		r := adapter.Response{
// 			ChannelID: "42217464",
// 			Text:      newMessage(500-4) + " (" + strconv.Itoa(i) + ")",
// 		}

// 		err := c.DeliverResponse(ctx, r)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		i++
// 		time.Sleep(1 * time.Second)
// 	}

// }
