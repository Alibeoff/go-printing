package tglib

import (
	"fmt"
	"strings"

	tele "gopkg.in/telebot.v4"
)

// Handler регистрирует обработчики по вашей нотации
func (b *Bot) Handler(endpoint any, h tele.HandlerFunc) {
	switch v := endpoint.(type) {
	case string:
		// Обработка текстовых сообщений
		if !strings.HasPrefix(v, "/") && !strings.HasPrefix(v, "|") && !strings.HasPrefix(v, ":") {
			if v == "" {
				b.Bot.Handle(tele.OnText, func(ctx tele.Context) error {
					unknow := ctx.Get("storage")
					store := unknow.(*Storage)
					if store != nil {
						unstate := store.Get("formstate")
						if unstate != nil {
							state := unstate.(map[int64]*FormBuilder)
							if form, exists := state[ctx.Sender().ID]; exists {
								if !form.Status {
									return nil
								}
								return form.HandleInput(ctx)
							}
						}
					}
					return nil
				})
				b.Bot.Handle(tele.OnVoice, func(ctx tele.Context) error {
					unknow := ctx.Get("storage")
					store := unknow.(*Storage)
					if store != nil {
						unstate := store.Get("formstate")
						if unstate != nil {
							state := unstate.(map[int64]*FormBuilder)
							if form, exists := state[ctx.Sender().ID]; exists {
								if !form.Status {
									return nil
								}
								return form.HandleInput(ctx)
							}
						}
					}
					return nil
				})
				b.Bot.Handle(tele.OnPhoto, func(ctx tele.Context) error {
					unknow := ctx.Get("storage")
					store := unknow.(*Storage)
					if store != nil {
						unstate := store.Get("formstate")
						if unstate != nil {
							state := unstate.(map[int64]*FormBuilder)
							if form, exists := state[ctx.Sender().ID]; exists {
								if !form.Status {
									return nil
								}
								return form.HandleInput(ctx)
							}
						}
					}
					return nil
				})

				b.Bot.Handle(tele.OnCallback, func(ctx tele.Context) error {
					unknow := ctx.Get("storage")
					store := unknow.(*Storage)
					if store != nil {
						unstate := store.Get("formstate")
						if unstate != nil {
							state := unstate.(map[int64]*FormBuilder)
							if form, exists := state[ctx.Sender().ID]; exists {
								fmt.Println("CHECK:", form.Status)
								if !form.Status {
									return nil
								}

								return form.handleCallback(ctx)
							}
						}
					}
					return nil
				})
			} else {
				b.Bot.Handle(v, h)
			}
		} else if strings.HasPrefix(v, "/") { // Обработка команд
			b.Bot.Handle(v, h)
		} else if callbackData, ok := strings.CutPrefix(v, "|"); ok {
			b.Bot.Handle(&tele.Btn{Unique: callbackData}, h)
		} else if mediaType, ok := strings.CutPrefix(v, ":"); ok { // Обработка медиа типов (voice, photo, document и т.д.)
			switch mediaType {
			case "voice":
				b.Bot.Handle(tele.OnVoice, h)
			case "photo":
				b.Bot.Handle(tele.OnPhoto, h)
			case "document":
				b.Bot.Handle(tele.OnDocument, h)
			case "audio":
				b.Bot.Handle(tele.OnAudio, h)
			case "video":
				b.Bot.Handle(tele.OnVideo, h)
			case "sticker":
				b.Bot.Handle(tele.OnSticker, h)
			case "location":
				b.Bot.Handle(tele.OnLocation, h)
			case "contact":
				b.Bot.Handle(tele.OnContact, h)
			default:
				fmt.Printf("Unknown media type: %s\n", mediaType)
			}
		}
	default:
		b.Bot.Handle(v, h)
	}
}
