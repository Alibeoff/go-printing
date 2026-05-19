package tglib

import (
	"fmt"
	"time"

	tele "gopkg.in/telebot.v4"
)

func NewBot(token string) *Bot {
	fmt.Println("NEW BOT TOKEN:", token)
	bot := &Bot{}

	settings := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(settings)
	if err != nil {
		panic(err)
	}

	bot.Storage = &Storage{
		Data: map[string]any{},
	}
	bot.Bot = b
	return bot
}
