package handler

import (
	"fmt"
	"go-printing/pkg/printer"
	tglib "go-printing/pkg/tg"
	"log"
	"strconv"

	tele "gopkg.in/telebot.v4"
)

func PrinterList(ctx tele.Context) error {
	printers, err := printer.GetAvailablePrinters()
	if err != nil {
		log.Println("printerList: ", err)
	}
	if len(printers) == 0 {
		kb := tglib.NewKeyboard()
		backBtn := tglib.NewInlineButton("🔙 Назад", "back_to_main").Build()
		kb.AddButton(backBtn).AddRow()
		return ctx.Send("Принтеры не найдены", kb.Build())
	}

	kb := tglib.NewKeyboard()
	for i, p := range printers {
		btn := tglib.NewInlineButton(p, fmt.Sprintf("select_printer_%d", i)).Build()
		kb.AddButton(btn).AddRow()
	}
	kb.AddButton(tglib.NewInlineButton("🔙 Назад", "back_to_main").Build()).AddRow()
	return ctx.Send("Список принтеров:", kb.Build())
}

func CancelPrint(ctx tele.Context) error {
	return ctx.Send("Печать отменена для выбранного документа.")
}

func PrinterSelected(ctx tele.Context, data string) error {
	printers, err := printer.GetAvailablePrinters()
	if err != nil {
		log.Println(err)
		return ctx.Edit("Ошибка при получении списка принтеров.")
	}
	id, err := strconv.Atoi(data)
	if err != nil {
		log.Println(err)
	}
	printerData.SetID(data)
	printerData.SetName(printers[id])
	return ctx.Edit(fmt.Sprintf("Вы выбрали принтер: %s", printers[id]))
}
