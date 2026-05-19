package tglib

import (
	"fmt"
	"strings"

	tele "gopkg.in/telebot.v4"
)

func MiddlewarePasport(b *Bot) *Middleware {
	return &Middleware{
		Bot:    b,
		Status: "",
	}
}

func (mw *Middleware) Use(middleware ...tele.MiddlewareFunc) {
	mw.Bot.Middlewares = append(mw.Bot.Middlewares, middleware...)
	mw.Bot.Bot.Use(middleware...)
}

func (mw *Middleware) Storager(next tele.HandlerFunc) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		ctx.Set("storage", mw.Bot.Storage)
		return next(ctx)
	}
}

func (mw *Middleware) Logger(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		fmt.Printf("Получено сообщение от %d: %s\n", c.Sender().ID, c.Text())
		return next(c)
	}
}

// AdminOnly с логированием
func (mw *Middleware) AdminOnly(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {

		user := c.Sender()

		if !mw.AdminList[user.ID] {
			// Логируем попытку доступа
			fmt.Printf("🚫 Попытка доступа к админ-команде: %s (@%s, ID: %d)\n",
				user.FirstName, user.Username, user.ID)

			// return c.Send("❌ Эта команда только для администраторов!")
			return next(c)
		}

		// Логируем успешный доступ
		fmt.Printf("✅ Админ доступ: %s (@%s) выполнил команду\n",
			user.FirstName, user.Username)

		return next(c)
	}
}

func (mw *Middleware) MessageDeleter(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if mw.Status != "form" {
			if mw.Status == "inbox" {
				return next(c)
			}
			c.Delete()
		}
		return next(c)
	}
}

func (mw *Middleware) FormMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		unknow := c.Get("storage")
		store := unknow.(*Storage)
		if store != nil {
			unstate := store.Get("formstate")
			if unstate != nil {
				state := unstate.(map[int64]*FormBuilder)
				if form, exists := state[c.Sender().ID]; exists {
					if form == nil {
						return next(c)
					}

					if !form.Status {
						mw.Status = ""
					}

				}
			}
		}
		return next(c)
	}
}

type CallbackHandler func(ctx tele.Context, data string) error

var routes = make(map[string]CallbackHandler)

// RegisterCallbackRoute - регистрирует обработчик для префикса
func RegisterCallbackRoute(prefix string, handler CallbackHandler) {
	routes[prefix] = handler
}

// CallbackRouterMiddleware - middleware для роутинга callback'ов
func (mw *Middleware) CallbackRouter(next tele.HandlerFunc) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		// Если это callback
		if ctx.Callback() != nil {
			data := strings.TrimSpace(ctx.Callback().Data)

			// Пропускаем системные callback'и формы
			if data == "form_back" || data == "form_complete" || data == "form_edit" {
				return next(ctx)
			}

			// Проверяем все зарегистрированные префиксы
			for prefix, handler := range routes {
				if strings.HasPrefix(data, prefix) {
					return handler(ctx, data)
				}
			}
		}

		// Если не callback или не нашли обработчик - идём дальше
		return next(ctx)
	}
}
