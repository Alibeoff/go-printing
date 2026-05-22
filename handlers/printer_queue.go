package handler

import (
	"fmt"
	"go-printing/pkg/printer"
	tglib "go-printing/pkg/tg"
	"log"
	"strconv"

	tele "gopkg.in/telebot.v4"
)

func PrinterQueue(ctx tele.Context) error {
	q, err := printer.ActivePrintList(printerData.Name)

	if err != nil {
		log.Println(err)
	}
	kb := tglib.NewKeyboard()
	for _, job := range q {
		text := "❌ " + job.Files
		btn := tglib.NewInlineButton(text, fmt.Sprintf("delete_queue_%d", job.JobID)).Build()
		kb.AddButton(btn).AddRow()
	}
	kb.AddButton(tglib.NewInlineButton("🔙 Назад", "back_to_main").Build()).AddRow()
	return ctx.Send("Список очереди:", kb.Build())
}

func DeleteQueue(ctx tele.Context, data string) error {
	id, err := strconv.Atoi(data)

	if err != nil {
		log.Println(err)
	}
	fmt.Println(id, data)
	printer.KillProcess(id)
	return ctx.Edit("Задание удалено из очереди.")
}
