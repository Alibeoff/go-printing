package tglib

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gopkg.in/telebot.v4"
	tele "gopkg.in/telebot.v4"
)

// Новый билдер формы
func NewFormBuilder(userID int64) *FormBuilder {
	return &FormBuilder{
		userID: userID,
		data:   make(map[string]any),
		steps:  make([]FormStep, 0),
	}
}

func (fb *FormBuilder) SaveReplyMarkup(markup *tele.ReplyMarkup) {
	fb.originalMarkup = markup
}

// Добавляем метод для установки колбэка
func (fb *FormBuilder) SetOnComplete(callback func(map[string]any)) *FormBuilder {
	fb.OnComplete = callback
	return fb
}

// Добавить кнопку отмены
func (fb *FormBuilder) AddCancelButton(ctx tele.Context) error {
	fb.hasCancelButton = true
	markup := &tele.ReplyMarkup{}
	// Добавляем кнопку отмены
	cancel := KeyboardReplyBuilder([][]string{{"❌ Отмена"}})
	markup.ReplyKeyboard = cancel.ReplyKeyboard

	ctx.Bot().Send(ctx.Chat(), fb.steps[fb.currentStep].Title, &telebot.SendOptions{
		ReplyMarkup: markup,
	})

	return nil
}

// Добавить кнопку "Назад"
func (fb *FormBuilder) AddBackInlineButton() {
	fb.hasBackButton = true
}

func (fb *FormBuilder) AddFinalMessage() {
	fb.hasFinalMessage = true
}

// Добавить кнопку завершения
func (fb *FormBuilder) AddReturnInlineButton() {
	fb.hasReturnButton = true
}

// Добавить шаги формы
func (fb *FormBuilder) AddSteps(steps ...FormStep) {
	fb.steps = append(fb.steps, steps...)
}

// Запустить форму
func (fb *FormBuilder) Start(ctx tele.Context) error {
	if len(fb.steps) == 0 {
		return ctx.Send("❌ Форма не содержит шагов")
	}

	// Начинаем с первого шага
	fb.currentStep = 0
	return fb.sendCurrentStep(ctx)
}

// Отправить текущий шаг
func (fb *FormBuilder) sendCurrentStep(ctx tele.Context) error {

	if fb.currentStep >= len(fb.steps) {
		return fb.completeForm(ctx)
	}

	step := fb.steps[fb.currentStep]

	// Создаем клавиатуру для этого шага
	markup := &tele.ReplyMarkup{}
	var rows []tele.Row

	// Добавляем кнопки для выбора (например, для пола)
	if step.Buttons == nil {
		if step.Type == QGender {
			rows = append(rows, markup.Row(
				tele.Btn{Text: "👨 Мужской", Data: "male"},
				tele.Btn{Text: "👩 Женский", Data: "female"},
			))
		}

	} else if step.Buttons != nil {
		if step.Type == QButton {
			fmt.Println(markup.InlineKeyboard)
			if markup.InlineKeyboard == nil {
				markup.InlineKeyboard = step.Buttons.InlineKeyboard
			}

			if fb.hasBackButton && fb.currentStep > 0 {
				backBtn := telebot.InlineButton{
					Unique: "form_back",
					Text:   "🔙 Назад",
				}
				markup.InlineKeyboard = append(markup.InlineKeyboard, []telebot.InlineButton{backBtn})
			}
		}

	}

	// Добавляем кнопку "Назад"
	if fb.hasBackButton && fb.currentStep > 0 {
		rows = append(rows, markup.Row(tele.Btn{
			Text: "🔙 Назад",
			Data: "form_back",
		}))
	}

	if markup.InlineKeyboard == nil {
		markup.Inline(rows...)
	}

	if fb.hasCancelButton && fb.currentStep == 0 {
		fb.hasCancelButton = false
		return nil
	}
	// if fb.state == "edit" {
	// 	return fb.showEditMenu(ctx)
	// }
	return ctx.Send(step.Title, markup)
}

// Обработать ответ пользователя
func (fb *FormBuilder) HandleInput(ctx tele.Context) error {
	if fb.currentStep >= len(fb.steps) {
		return fb.completeForm(ctx)
	}

	step := fb.steps[fb.currentStep]

	// Проверяем, не нажата ли inline-кнопка
	if ctx.Callback() != nil {
		return fb.handleCallback(ctx)
	}

	// Обрабатываем данные в зависимости от типа
	var input any
	var errorMsg string

	if ctx.Text() == "❌ Отмена" {
		return fb.cancelForm(ctx)
	}

	switch step.Type {
	case QText:
		input = ctx.Text()

	case QVoice:
		if ctx.Message().Voice != nil {
			input = ctx.Message().Voice
		}

	case QPhoto:
		if ctx.Message().Photo != nil {
			input = *ctx.Message().Photo
		}

	case QGender:
		return nil
	}

	// Валидация
	if step.Validate != nil {
		var ok bool
		ok, errorMsg = step.Validate(input)
		if !ok {
			if errorMsg == "" {
				errorMsg = "🚫 Пожалуйста, проверьте введенные данные"
			}
			return ctx.Send(errorMsg)
		}
	}

	// Сохраняем данные
	fb.data[step.Field] = input

	// Красивый лог
	fb.logStep(step, input)

	// Переходим к следующему шагу
	fb.currentStep++
	return fb.sendCurrentStep(ctx)
}

// Обработать callback от inline-кнопок
func (fb *FormBuilder) handleCallback(ctx tele.Context) error {
	data := ctx.Callback().Data
	data = strings.TrimSpace(data)
	// cals := fb.steps[fb.currentStep].CallBacks

	if strings.HasPrefix(data, "call") {
		step := fb.steps[fb.currentStep]
		fb.data[step.Field] = data

		switch data {
		case "call_sql":
			fb.logStep(step, "sql")
		case "call_mongo":
			fb.logStep(step, "mongo")
		case "call_pg":
			fb.logStep(step, "pg")
		}

		// Переходим дальше
		fb.currentStep++
		return fb.sendCurrentStep(ctx)
	}

	if strings.HasPrefix(data, "edit_") {
		stepnum := strings.TrimPrefix(data, "edit_")
		num, _ := strconv.Atoi(stepnum)
		fb.currentStep = num
		return fb.sendCurrentStep(ctx)
	}

	switch data {
	case "form_back":
		fmt.Println("FORM BACK WORKING")
		return fb.goBack(ctx)

	case "male", "female":
		// Сохраняем выбор пола
		step := fb.steps[fb.currentStep]
		fb.data[step.Field] = data

		// Лог
		genderText := "👨 Мужской"
		if data == "female" {
			genderText = "👩 Женский"
		}
		fb.logStep(step, genderText)

		// Переходим дальше
		fb.currentStep++
		return fb.sendCurrentStep(ctx)

	case "form_complete":
		fb.state = ""
		return fb.finalizeForm(ctx)

	case "form_edit":
		fb.state = "edit"
		return fb.showEditMenu(ctx)
	}

	return nil
}

// Вернуться на шаг назад
func (fb *FormBuilder) goBack(ctx tele.Context) error {
	if fb.currentStep > 0 {
		fb.currentStep--

		// Удаляем данные этого шага
		prevStep := fb.steps[fb.currentStep]
		delete(fb.data, prevStep.Field)

		// Отправляем шаг заново
		ctx.Respond()
		return fb.sendCurrentStep(ctx)
	}

	return ctx.Send("📍 Вы на первом шаге формы")
}

// Отменить форму
func (fb *FormBuilder) cancelForm(ctx tele.Context) error {
	fb.Status = false
	fb.state = ""

	st := ctx.Get("storage")
	store := st.(*Storage)
	unstate := store.Get("formstate")
	if unstate != nil {
		state := unstate.(map[int64]*FormBuilder)
		fb.Status = true
		state[ctx.Sender().ID] = fb
		store.Set("formstate", nil)
	}
	ctx.Respond(&tele.CallbackResponse{
		Text: "❌ Форма отменена",
	})

	// Восстанавливаем оригинальную клавиатуру если была
	if fb.originalMarkup != nil {
		ctx.Send("✅ Вы вернулись в обычный режим", fb.originalMarkup)
	} else {
		ctx.Send("✅ Форма отменена")
	}

	// Лог
	fmt.Printf("\n❌ ФОРМА ОТМЕНЕНА\n")
	fmt.Printf("👤 Пользователь: %s\n", ctx.Sender().Username)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━")

	return nil
}

// Завершить форму
func (fb *FormBuilder) completeForm(ctx tele.Context) error {
	fb.currentStep = 0

	// Если есть кнопка подтверждения, показываем её
	if fb.hasReturnButton {
		markup := &tele.ReplyMarkup{}
		markup.Inline(
			markup.Row(
				tele.Btn{Text: "✅ Подтвердить", Data: "form_complete"},
				tele.Btn{Text: "✏️ Редактировать", Data: "form_edit"},
			),
		)
		summary := fb.buildSummary()
		return ctx.Send(summary, markup)
	}

	return fb.finalizeForm(ctx)
}

// Показать меню редактирования
func (fb *FormBuilder) showEditMenu(ctx tele.Context) error {
	markup := &tele.ReplyMarkup{}
	var rows []tele.Row

	// Кнопки для редактирования каждого поля
	for i, step := range fb.steps {
		rows = append(rows, markup.Row(tele.Btn{
			Text: fmt.Sprintf("✏️ %s", step.Title[:min(20, len(step.Title))]),
			Data: fmt.Sprintf("edit_%d", i),
		}))
	}

	rows = append(rows, markup.Row(
		tele.Btn{Text: "✅ Завершить", Data: "form_complete"},
	))

	markup.Inline(rows...)

	return ctx.Send("📝 Выберите поле для редактирования:", markup)
}

// Завершить форму и вывести результат
func (fb *FormBuilder) finalizeForm(ctx tele.Context) error {
	// Красивый вывод в консоль
	fmt.Println("\n✨ ФОРМА ЗАПОЛНЕНА ✨")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for _, step := range fb.steps {
		value := fb.data[step.Field]
		fmt.Printf("%s: %v\n", step.Field, value)
	}

	fmt.Printf("👤 Пользователь: %s\n", ctx.Sender().Username)
	fmt.Printf("🕒 Время: %s\n", time.Now().Format("15:04:05"))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// ВЫЗЫВАЕМ КОЛБЭК ЕСЛИ ОН УСТАНОВЛЕН
	if fb.OnComplete != nil {
		fb.OnComplete(fb.data)
	}

	// Отправляем пользователю
	summary := fb.buildSummary()
	ctx.Respond(&tele.CallbackResponse{
		Text: "✅ Форма отправлена!",
	})

	// Восстанавливаем оригинальную клавиатуру
	st := ctx.Get("storage")
	store := st.(*Storage)
	unstate := store.Get("formstate")
	if unstate != nil {
		state := unstate.(map[int64]*FormBuilder)
		fb.Status = false
		delete(state, ctx.Sender().ID)
		store.Set("formstate", state)
	}

	markup := &tele.ReplyMarkup{}
	if fb.originalMarkup != nil {
		markup = fb.originalMarkup
	}

	if photo, ok := fb.data["photo"].(tele.Photo); ok {
		fmt.Println("📸 Отправляем фото")
		photo.Caption = summary
		return ctx.Send(&photo, markup)
	} else if fb.hasFinalMessage {
		return ctx.Send(summary, markup)
	} else {
		return nil
	}
}

// Построить сводку данных
func (fb *FormBuilder) buildSummary() string {
	var builder strings.Builder
	builder.WriteString("📋 *Сводка данных:*\n\n")

	for _, step := range fb.steps {
		value := fb.data[step.Field]
		var displayValue string

		switch v := value.(type) {
		case string:
			displayValue = v
		case *tele.Voice:
			displayValue = "🎤 Голосовое сообщение"
		case tele.Photo:
			displayValue = "📸 Фото"
		default:
			displayValue = fmt.Sprintf("%v", v)
		}

		builder.WriteString(fmt.Sprintf("• *%s:* %s\n", step.Field, displayValue))
	}

	return builder.String()
}

// Логирование шага
func (fb *FormBuilder) logStep(step FormStep, input any) {
	fmt.Printf("✅ [Шаг %d/%d] %s: %v\n",
		fb.currentStep+1, len(fb.steps), step.Field, input)
}

// Глобальное хранилище активных форм
var activeForms = make(map[int64]*FormBuilder)

// Основной хендлер для использования
func (fb *FormBuilder) HandleForm(ctx tele.Context, steps ...FormStep) error {
	// Сохраняем форму
	activeForms[ctx.Sender().ID] = fb
	unknow := ctx.Get("storage")
	store := unknow.(*Storage)
	if store != nil {
		unstate := store.Get("formstate")
		if unstate != nil {
			state := unstate.(map[int64]*FormBuilder)
			fb.Status = true
			state[ctx.Sender().ID] = fb
			store.Set("formstate", state)
		} else {
			state := make(map[int64]*FormBuilder)
			fb.Status = true
			state[ctx.Sender().ID] = fb
			store.Set("formstate", state)
		}
	}

	// Запускаем
	return fb.Start(ctx)
}

// Хендлер для обработки ввода
func HandleFormInput(ctx tele.Context) error {
	userID := ctx.Sender().ID

	if form, exists := activeForms[userID]; exists {
		return form.HandleInput(ctx)
	}

	return nil
}

// Хендлер для callback
func HandleFormCallback(ctx tele.Context) error {
	userID := ctx.Sender().ID

	if form, exists := activeForms[userID]; exists {
		return form.handleCallback(ctx)
	}

	return nil
}

func (fb *FormBuilder) BuildFormHandlers(bot *Bot) error {
	steps := fb.steps
	for _, step := range steps {
		switch step.Type {
		case QText:
			bot.Bot.Handle(tele.OnText, HandleFormInput)
		case QVoice:
			bot.Bot.Handle(tele.OnVoice, HandleFormInput)
		case QPhoto:
			bot.Bot.Handle(tele.OnPhoto, HandleFormInput)
		case QGender:
			bot.Bot.Handle(tele.OnCallback, HandleFormCallback)
		case QButton:
			bot.Bot.Handle(tele.OnCallback, HandleFormCallback)
		}
	}
	return nil
}
