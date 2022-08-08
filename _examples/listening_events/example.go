package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

var (
	token   = os.Getenv("disgo_token")
	guildID = snowflake.GetEnv("disgo_guild_id")
)

func main() {
	log.SetLevel(log.LevelInfo)
	log.Info("starting example...")
	log.Info("disgo version: ", disgo.Version)

	client, err := disgo.New(token,
		bot.WithDefaultGateway(),
		bot.WithEventListenerFunc(eventListenerFunc),
		bot.WithEventListenerChan(eventListenerChan()),
		bot.WithEventListeners(&events.ListenerAdapter{OnMessageCreate: eventListenerFunc}),
	)
	if err != nil {
		log.Fatal("error while building disgo instance: ", err)
		return
	}

	defer client.Close(context.TODO())

	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to gateway: ", err)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func eventListenerFunc(event *events.MessageCreate) {
	_, _ = event.Client().Rest().CreateMessage(event.ChannelID, discord.MessageCreate{
		Content: "pong",
	})
}

func eventListenerChan() chan<- *events.MessageCreate {
	c := make(chan *events.MessageCreate)
	go func() {
		defer close(c)
		for event := range c {
			if event.Message.Content == "ping" {
				_, _ = event.Client().Rest().CreateMessage(event.ChannelID, discord.MessageCreate{
					Content: "pong",
				})
			}
		}
	}()
	return c
}