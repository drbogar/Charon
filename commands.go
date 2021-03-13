package main

import (
	"fmt"
	"log"
	"time"

	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
)

// Commands containes the all commands of this bot
type Commands struct {
	Ctx *bot.Context
}

// Ping command responses the latency.
func (c *Commands) Ping(e *gateway.MessageCreateEvent) error {
	now := time.Now()
	msg, err := c.Ctx.SendMessageReply(e.ChannelID, "Loading...", nil, e.ID)
	if err != nil {
		return err
	}
	cmdTime := e.Timestamp
	if e.EditedTimestamp.IsValid() {
		cmdTime = e.EditedTimestamp
	}
	c.Ctx.DeleteMessage(e.ChannelID, msg.ID)
	resp := fmt.Sprintf("🏓\nLatency of message receiving is %v.\nLatency of message sending is %v.", now.Sub(time.Time(cmdTime)), time.Time(msg.Timestamp).Sub(now))
	embed := discord.Embed{
		Description: resp,
		Type:        discord.NormalEmbed,
		Color:       discord.DefaultEmbedColor,
	}
	c.Ctx.SendEmbedReply(e.ChannelID, embed, e.ID)
	return nil
}

// StartGame start a new game
func (c *Commands) StartGame(e *gateway.MessageCreateEvent) error {
	_, err := c.Ctx.SendMessageReply(e.ChannelID, "Okay, let's play!", nil, e.ID)
	if err != nil {
		return err
	}
	whoWantPlayText := "Who want to play? (Add reaction!)"
	msg, err := c.Ctx.SendMessage(e.ChannelID, whoWantPlayText, nil)
	if err != nil {
		return err
	}

	reactionAddChan, _ := c.Ctx.ChanFor(func(v interface{}) bool {
		ev, ok := v.(*gateway.MessageReactionAddEvent)
		if !ok {
			return false
		}
		return ev.MessageID == msg.ID
	})
	reactionRemoveChan, _ := c.Ctx.ChanFor(func(v interface{}) bool {
		ev, ok := v.(*gateway.MessageReactionRemoveEvent)
		if !ok {
			return false
		}
		return ev.MessageID == msg.ID
	})
	go handleWannaPlayReactions(reactionAddChan, reactionRemoveChan, c.Ctx, whoWantPlayText)
	return nil
}

func handleWannaPlayReactions(add <-chan interface{}, remove <-chan interface{}, Ctx *bot.Context, whoWantPlayText string) {
	type wannaPlay struct {
		player Player
		count  int
	}
	wannaPlays := make(map[discord.UserID]wannaPlay)
	for add != nil && remove != nil {
		select {
		case e, ok := <-add:
			if !ok {
				add = nil
				continue
			}
			ev, ok := e.(*gateway.MessageReactionAddEvent)
			if !ok {
				log.Println("Type error")
			}
			player := Player{Member: *ev.Member}
			if _, ok := wannaPlays[ev.Member.User.ID]; !ok {
				wannaPlays[ev.Member.User.ID] = wannaPlay{player: Player{Member: *ev.Member}, count: 1}
			} else {
				wannaPlay := wannaPlays[ev.Member.User.ID]
				wannaPlay.count++
				wannaPlays[ev.Member.User.ID] = wannaPlay
			}
			log.Println("Reaction added by " + player.GetNick())
			wannaText := whoWantPlayText
			for _, wannaPlay := range wannaPlays {
				wannaText += "\n- " + wannaPlay.player.GetNick()
			}
			Ctx.EditText(ev.ChannelID, ev.MessageID, wannaText)
		case e, ok := <-remove:
			if !ok {
				add = nil
				continue
			}
			ev, ok := e.(*gateway.MessageReactionRemoveEvent)
			if !ok {
				log.Println("Type error")
			}
			wannaPlay := wannaPlays[ev.UserID]
			wannaPlay.count--
			wannaPlays[ev.UserID] = wannaPlay
			if wannaPlay.count < 1 {
				delete(wannaPlays, ev.UserID)
			}
			log.Println("Reaction removed by " + wannaPlay.player.GetNick())
			wannaText := whoWantPlayText
			for _, wannaPlay := range wannaPlays {
				wannaText += "\n- " + wannaPlay.player.GetNick()
			}
			Ctx.EditText(ev.ChannelID, ev.MessageID, wannaText)
		}
	}
}
