package tglib

import tele "gopkg.in/telebot.v4"

// NewInlineButton - создает новую inline-кнопку
func NewInlineButton(text, unique string) *InlineButtonBuilder {
	return &InlineButtonBuilder{
		button: &tele.InlineButton{
			Unique: unique,
			Text:   text,
		},
	}
}

func ContactButton(text string) *tele.ReplyMarkup {
	replyMarkup := &tele.ReplyMarkup{
		ReplyKeyboard: [][]tele.ReplyButton{
			{tele.ReplyButton{Text: text, Contact: true}},
		},
		ResizeKeyboard:  true, // Опционально: подогнать размер клавиатуры
		OneTimeKeyboard: true, // Опционально: скрыть после использования
	}
	return replyMarkup
}

// WithData - добавляет callback данные
func (b *InlineButtonBuilder) WithData(data string) *InlineButtonBuilder {
	b.button.Data = data
	return b
}

// WithURL - делает кнопку URL-кнопкой
func (b *InlineButtonBuilder) WithURL(url string) *InlineButtonBuilder {
	b.button.URL = url
	return b
}

// WithLogin - делает кнопкой для Web Apps авторизации
func (b *InlineButtonBuilder) WithLogin(loginURL string, requestWrite bool) *InlineButtonBuilder {
	b.button.Login = &tele.Login{
		URL:         loginURL,
		WriteAccess: requestWrite,
	}
	return b
}

// WithQuery - делает кнопкой для switch inline query
func (b *InlineButtonBuilder) WithQuery(query string) *InlineButtonBuilder {
	b.button.InlineQuery = query
	return b
}

// WithQueryChat- делает кнопкой для switch inline query в текущем чате
func (b *InlineButtonBuilder) WithQueryChat(query string) *InlineButtonBuilder {
	b.button.InlineQueryChat = query
	return b
}

// WithQueryChosenChat- делает кнопкой для switch inline query в текущем чате
func (b *InlineButtonBuilder) WithQueryChosenChat(query *tele.SwitchInlineQuery) *InlineButtonBuilder {
	b.button.InlineQueryChosenChat = query
	return b
}

// AsGame - делает кнопкой для запуска игры
func (b *InlineButtonBuilder) AsGame(query *tele.CallbackGame) *InlineButtonBuilder {
	b.button.CallbackGame = query
	return b
}

// WebApp- делает кнопкой для открытия веб приложения
func (b *InlineButtonBuilder) AsWebApp(query *tele.WebApp) *InlineButtonBuilder {
	b.button.WebApp = query
	return b
}

// AsPoll - делает кнопкой
func (b *InlineButtonBuilder) AsPay() *InlineButtonBuilder {
	b.button.Pay = true
	return b
}

func (b *InlineButtonBuilder) ContactButton(text string) *InlineButtonBuilder {
	return b
}

// Build - возвращает готовую кнопку
func (b *InlineButtonBuilder) Build() *tele.InlineButton {
	return b.button
}
