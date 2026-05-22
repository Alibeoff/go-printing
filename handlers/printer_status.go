package handler

import (
	"fmt"
	"go-printing/pkg/printer"
	"log"

	tele "gopkg.in/telebot.v4"
)

func PrinterStatus(ctx tele.Context) error {
	status, err := printer.CheckPrinterStatus(printerData.Name)
	if err != nil {
		log.Println(err)
	}

	ready := "Нет"
	if status.IsReady {
		ready = "Да"
	}
	printing := "Нет"
	if status.IsPrinting {
		printing = "Да"
	}

	text := fmt.Sprintf("Статус принтера %s:\nГотов к печати: %s\nПечатает: %s\nСообщение: %s\nКоличество заданий в очереди: %d", status.Name, ready, printing, status.Message, status.JobsCount)
	return ctx.Send(text)
}
