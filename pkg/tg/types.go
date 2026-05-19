package tglib

import tele "gopkg.in/telebot.v4"

type Bot struct {
	Bot         *tele.Bot
	Middlewares []tele.MiddlewareFunc
	Storage     *Storage
}

type Middleware struct {
	Bot       *Bot
	Status    string
	AdminList map[int64]bool
}

type Storage struct {
	Data map[string]any
}

type ChatInfo struct {
	ID       int64
	UserName string
	UserTag  string
}

// InlineButtonBuilder - билдер для inline-кнопок
type InlineButtonBuilder struct {
	button *tele.InlineButton
}

// KeyboardBuilder - билдер для inline-клавиатур
type KeyboardBuilder struct {
	rows       [][]tele.InlineButton
	currentRow []tele.InlineButton
}

// Типы полей формы
const (
	QText    = "text"
	QVoice   = "voice"
	QPhoto   = "photo"
	QGender  = "gender"
	QConfirm = "confirm"
	QButton  = "button"
)

type Validator func(any) (bool, string)

type FormStep struct {
	Title     string
	Field     string
	Type      string
	Validate  Validator
	Buttons   *tele.ReplyMarkup // Для кнопок (например, выбор пола)
	CallBacks []*tele.Callback
}

// Билдер формы
type FormBuilder struct {
	Status          bool
	state           string
	userID          int64
	steps           []FormStep
	currentStep     int
	data            map[string]any
	hasCancelButton bool
	hasBackButton   bool
	hasReturnButton bool
	hasFinalMessage bool
	originalMarkup  *tele.ReplyMarkup
	OnComplete      func(map[string]any)
}
