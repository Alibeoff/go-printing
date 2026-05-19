package tglib

import (
	"fmt"

	"gopkg.in/telebot.v4"
)

// Предопределенные валидаторы
func TextName(input any) (bool, string) {
	text, ok := input.(string)
	if !ok {
		return false, "❌ Неверный формат данных"
	}
	if len(text) < 2 {
		return false, "⚠️ Имя должно содержать минимум 2 символа"
	}
	if len(text) > 50 {
		return false, "⚠️ Имя слишком длинное (макс. 50 символов)"
	}
	return true, ""
}

// Валидатор для ссылки
func Link(input any) (bool, string) {
	link, ok := input.(string)
	if !ok {
		return false, "❌ Неверный формат данных"
	}
	if link == "" {
		return false, "⚠️ Ссылка не может быть пустой"
	}
	if len(link) < 10 {
		return false, "⚠️ Ссылка слишком короткая (мин. 10 символов)"
	}
	if len(link) > 200 {
		return false, "⚠️ Ссылка слишком длинная (макс. 200 символов)"
	}

	// Проверяем наличие протокола
	hasProtocol := false
	protocols := []string{"http://", "https://"}
	for _, proto := range protocols {
		if len(link) >= len(proto) && link[:len(proto)] == proto {
			hasProtocol = true
			break
		}
	}
	if !hasProtocol {
		return false, "⚠️ Ссылка должна начинаться с http:// или https://"
	}

	// Вызываем дополнительную проверку
	return true, link
}

func NumberInRange(min, max int) Validator {
	return func(input any) (bool, string) {
		text, ok := input.(string)
		if !ok {
			return false, "❌ Неверный формат данных"
		}

		var num int
		_, err := fmt.Sscanf(text, "%d", &num)
		if err != nil {
			return false, "📏 Пожалуйста, введите число"
		}

		if num < min || num > max {
			return false, fmt.Sprintf("📏 Введите число от %d до %d", min, max)
		}

		return true, ""
	}
}

func CheckVoice(input any) (bool, string) {
	// В телеграме голосовое сообщение - это объект Voice
	_, ok := input.(*telebot.Voice)
	if !ok {
		return false, "🎤 Пожалуйста, отправьте голосовое сообщение"
	}
	return true, ""
}

func CheckPhoto(input any) (bool, string) {
	// В телеграме фото - это массив Photo
	_, ok := input.(telebot.Photo)
	if !ok {
		return false, "📸 Пожалуйста, отправьте фото"
	}
	return true, ""
}

func CallbackValidator(input any) (bool, string) {
	// fee := input.(map[string]any)
	// fee["calback"] = value()
	return false, ""
}
