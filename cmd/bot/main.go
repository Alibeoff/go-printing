package main

import (
	handler "go-printing/handlers"
	tglib "go-printing/pkg/tg"
	"log"
	"os"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v4"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN is required")
	}

	adminIDStr := os.Getenv("ADMIN_ID")
	if adminIDStr == "" {
		log.Fatal("ADMIN_ID is required")
	}

	adminID, err := strconv.ParseInt(adminIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Invalid ADMIN_ID: %v", err)
	}
	handler.AdminID = adminID

	b := tglib.NewBot(token)
	w := tglib.MiddlewarePasport(b)
	// Регистрируем админа в мидлваре
	w.AdminList = make(map[int64]bool)
	w.AdminList[adminID] = true

	w.Use(w.Storager, w.FormMiddleware, w.CallbackRouter, w.MessageDeleter)

	registerCallbackRoutes()

	b.Handler("/start", start)
	b.Handler("Принтеры", handler.PrinterList)
	b.Handler("Cтатус принтера", handler.PrinterStatus)
	b.Handler("Очередь печати", handler.PrinterQueue)
	b.Handler("Отменить печать", handler.CancelPrint)
	b.Handler(":document", handler.PrinterAddDocument)
	b.Handler("", func(ctx tele.Context) error { return nil })
	b.Bot.Handle(tele.OnPhoto, handler.PrinterAddPhoto)

	b.Bot.Start()
}

func start(ctx tele.Context) error {
	menu := [][]string{
		{"Принтеры", "Cтатус принтера"},
		{"Очередь печати", "Отменить печать"},
	}
	opts := tglib.KeyboardReplyBuilder(menu)

	return ctx.Send("Привет! Я бот для управления принтерами. Выберите действие из меню.", opts)
}

func registerCallbackRoutes() {
	tglib.RegisterCallbackRoute("select_printer_", func(ctx tele.Context, data string) error {
		alias := strings.TrimPrefix(data, "select_printer_")
		return handler.PrinterSelected(ctx, alias)
	})

	tglib.RegisterCallbackRoute("delete_queue_", func(ctx tele.Context, data string) error {
		alias := strings.TrimPrefix(data, "delete_queue_")
		return handler.DeleteQueue(ctx, alias)
	})

	// Роуты для админских кнопок
	tglib.RegisterCallbackRoute("approve_", func(ctx tele.Context, data string) error {
		alias := strings.TrimPrefix(data, "approve_")
		log.Printf("🔥 APPROVE CALLBACK: %s", data)
		return handler.HandleApprove(ctx, alias)
	})

	tglib.RegisterCallbackRoute("reject_", func(ctx tele.Context, data string) error {
		alias := strings.TrimPrefix(data, "reject_")
		log.Printf("🔥 REJECT CALLBACK: %s", alias)
		return handler.HandleReject(ctx, alias)
	})
}
