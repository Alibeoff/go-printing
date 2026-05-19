package tglib

import tele "gopkg.in/telebot.v4"

// NewKeyboard - создает новую клавиатуру
func NewKeyboard() *KeyboardBuilder {
	return &KeyboardBuilder{
		rows: make([][]tele.InlineButton, 0),
	}
}

// AddButton - добавляет кнопку в текущий ряд
func (k *KeyboardBuilder) AddButton(btn *tele.InlineButton) *KeyboardBuilder {
	k.currentRow = append(k.currentRow, *btn)
	return k
}

// AddRow - завершает текущий ряд и начинает новый
func (k *KeyboardBuilder) AddRow() *KeyboardBuilder {
	if len(k.currentRow) > 0 {
		k.rows = append(k.rows, k.currentRow)
		k.currentRow = make([]tele.InlineButton, 0)
	}
	return k
}

// AddButtons - добавляет несколько кнопок в ряд
func (k *KeyboardBuilder) AddButtons(buttons ...*tele.InlineButton) *KeyboardBuilder {
	for _, btn := range buttons {
		k.currentRow = append(k.currentRow, *btn)
	}
	return k
}

// Build - возвращает готовую клавиатуру
func (k *KeyboardBuilder) Build() *tele.ReplyMarkup {
	// Добавляем последний ряд, если он не пустой
	if len(k.currentRow) > 0 {
		k.rows = append(k.rows, k.currentRow)
	}

	return &tele.ReplyMarkup{
		InlineKeyboard: k.rows,
	}
}

// BuildWithOptions - возвращает клавиатуру с дополнительными опциями
func (k *KeyboardBuilder) BuildWithOptions(resize, oneTime, selective bool) *tele.ReplyMarkup {
	if len(k.currentRow) > 0 {
		k.rows = append(k.rows, k.currentRow)
	}

	return &tele.ReplyMarkup{
		InlineKeyboard:  k.rows,
		ResizeKeyboard:  resize,
		OneTimeKeyboard: oneTime,
		Selective:       selective,
	}
}
