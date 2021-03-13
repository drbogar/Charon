package main

import (
	"log"
	"os"

	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/gateway"
)

func main() {
	var token = os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatalln("No $BOT_TOKEN given.")
	}

	commands := &Commands{}

	bot.Run(token, commands, func(ctx *bot.Context) error {
		ctx.HasPrefix = bot.NewPrefix("!", "~", "Charon, ")
		ctx.EditableCommands = true
		ctx.AddIntents(gateway.IntentGuildMessageReactions)

		return nil
	})
}
