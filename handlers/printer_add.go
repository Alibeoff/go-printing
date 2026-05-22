package handler

import (
	"fmt"
	"go-printing/pkg/printer"
	tglib "go-printing/pkg/tg"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	tele "gopkg.in/telebot.v4"
)

var AdminID int64

type Job struct {
	Path     string
	Printer  string
	Name     string
	UserID   int64
	UserName string
}

var jobs = make(map[string]*Job)

func PrinterAddDocument(ctx tele.Context) error {
	if ctx.Message().Document == nil {
		return ctx.Send("❌ Пожалуйста, отправьте файл")
	}

	user := ctx.Sender()
	file := ctx.Message().Document

	if err := ensureFilesDir(); err != nil {
		return ctx.Send("❌ Ошибка сервера")
	}

	uniqueName := getUniqueFileName(file.FileName)
	filePath := filepath.Join("files", uniqueName)

	log.Printf("📥 Скачиваем файл: %s -> %s", file.FileName, filePath)

	err := ctx.Bot().Download(&file.File, filePath)
	if err != nil {
		log.Printf("❌ Ошибка скачивания: %v", err)
		return ctx.Send("❌ Ошибка при загрузке файла")
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("❌ Файл не скачался: %s", filePath)
		return ctx.Send("❌ Ошибка при загрузке файла")
	}

	if printerData.ID == "" {
		cleanupFile(filePath)
		return ctx.Send("❌ Сначала выберите принтер")
	}

	id, _ := strconv.Atoi(printerData.ID)
	printers, _ := printer.GetAvailablePrinters()
	printerName := printers[id]

	jobID := fmt.Sprintf("%d", time.Now().UnixNano())

	jobs[jobID] = &Job{
		Path:     filePath,
		Printer:  printerName,
		Name:     file.FileName,
		UserID:   user.ID,
		UserName: user.Username,
	}
	kb := tglib.NewKeyboard()
	approveBtn := tglib.NewInlineButton("✅ Напечатать", fmt.Sprintf("approve_%s", jobID)).Build()
	rejectBtn := tglib.NewInlineButton("❌ Отклонить", fmt.Sprintf("reject_%s", jobID)).Build()

	kb.AddButton(approveBtn)
	kb.AddButton(rejectBtn)

	// БЕЗ Markdown
	msg := fmt.Sprintf("📄 НОВЫЙ ФАЙЛ\nОт: @%s\nФайл: %s\nПринтер: %s",
		user.Username, file.FileName, printerName)

	_, err = ctx.Bot().Send(&tele.User{ID: AdminID}, &tele.Document{
		File:    file.File,
		Caption: msg,
	}, &tele.SendOptions{
		ReplyMarkup: kb.Build(),
	})

	if err != nil {
		log.Printf("❌ Ошибка отправки админу: %v", err)
		cleanupFile(filePath)
		delete(jobs, jobID)
		return ctx.Send("❌ Ошибка отправки администратору")
	}

	return ctx.Send("⏳ Ваш файл отправлен администратору на проверку")
}

func PrinterAddPhoto(ctx tele.Context) error {
	log.Printf("📸 Получено фото")

	if ctx.Message().Photo == nil {
		return ctx.Send("❌ Пожалуйста, отправьте фото")
	}

	user := ctx.Sender()

	// Photo - это одно фото
	photo := ctx.Message().Photo
	file := &photo.File

	if err := ensureFilesDir(); err != nil {
		return ctx.Send("❌ Ошибка сервера")
	}

	// Генерируем уникальное имя файла
	uniqueName := getUniqueFileName("photo.jpg")
	filePath := filepath.Join("files", uniqueName)

	log.Printf("📥 Скачиваем фото: %s", filePath)

	// Скачиваем фото
	err := ctx.Bot().Download(file, filePath)
	if err != nil {
		log.Printf("❌ Ошибка скачивания фото: %v", err)
		return ctx.Send("❌ Ошибка при загрузке фото")
	}

	// Проверяем что фото скачалось
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("❌ Фото не скачалось: %s", filePath)
		return ctx.Send("❌ Ошибка при загрузке фото")
	}

	if printerData.ID == "" {
		cleanupFile(filePath)
		return ctx.Send("❌ Сначала выберите принтер")
	}

	id, _ := strconv.Atoi(printerData.ID)
	printers, _ := printer.GetAvailablePrinters()
	printerName := printers[id]

	// Генерируем ID заказа
	jobID := fmt.Sprintf("%d", time.Now().UnixNano())

	jobs[jobID] = &Job{
		Path:     filePath,
		Printer:  printerName,
		Name:     "photo.jpg",
		UserID:   user.ID,
		UserName: user.Username,
	}

	log.Printf("✅ СОЗДАН ЗАКАЗ: %s", jobID)
	log.Printf("📁 Путь к фото: %s", filePath)
	log.Printf("👤 Пользователь: @%s", user.Username)

	// Создаем кнопки
	kb := tglib.NewKeyboard()
	approveBtn := tglib.NewInlineButton("✅ Напечатать", fmt.Sprintf("approve_%s", jobID)).Build()
	rejectBtn := tglib.NewInlineButton("❌ Отклонить", fmt.Sprintf("reject_%s", jobID)).Build()

	kb.AddButton(approveBtn)
	kb.AddButton(rejectBtn)

	// БЕЗ Markdown, обычный текст
	msg := fmt.Sprintf("🖼️ НОВОЕ ФОТО\nОт: @%s\nПринтер: %s",
		user.Username, printerName)

	// Отправляем админу
	_, err = ctx.Bot().Send(&tele.User{ID: AdminID}, &tele.Photo{
		File:    *file,
		Caption: msg,
	}, &tele.SendOptions{
		ReplyMarkup: kb.Build(), // Убрали ParseMode
	})

	if err != nil {
		log.Printf("❌ Ошибка отправки админу: %v", err)
		cleanupFile(filePath)
		delete(jobs, jobID)
		return ctx.Send("❌ Ошибка отправки администратору")
	}

	log.Printf("📨 Отправлено админу с кнопками")

	return ctx.Send("⏳ Ваше фото отправлено администратору на проверку")
}

func HandleApprove(ctx tele.Context, data string) error {
	log.Printf("🔥🔥🔥 HANDLE_APPROVE ВЫЗВАН! data: %s", data)

	// Обязательно отвечаем на callback
	ctx.Respond()

	job, exists := jobs[data]
	if !exists {
		log.Printf("❌ Заказ %s не найден", data)
		ctx.Send("❌ Заказ не найден")
		return nil
	}

	log.Printf("✅ Найден заказ: %+v", job)

	// Проверяем файл
	if _, err := os.Stat(job.Path); os.IsNotExist(err) {
		log.Printf("❌ Файл не найден: %s", job.Path)
		ctx.Edit("❌ Файл не найден")
		delete(jobs, data)
		return nil
	}

	// Печатаем
	log.Printf("🖨️ lp -d %s %s", job.Printer, job.Path)
	cmd := exec.Command("lp", "-d", job.Printer, job.Path)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("❌ Ошибка печати: %v, вывод: %s", err, string(output))
		ctx.Edit("❌ Ошибка печати")
		ctx.Bot().Send(&tele.User{ID: job.UserID}, fmt.Sprintf("❌ Ошибка при печати файла %s", job.Name))
	} else {
		log.Printf("✅ Печать успешна: %s", string(output))
		ctx.Edit(fmt.Sprintf("✅ %s отправлен на печать", job.Name))
		ctx.Bot().Send(&tele.User{ID: job.UserID}, fmt.Sprintf("✅ Ваш файл %s напечатан", job.Name))
	}

	// Чистим
	cleanupFile(job.Path)
	delete(jobs, data)

	return nil
}

func HandleReject(ctx tele.Context, data string) error {
	log.Printf("🔥🔥🔥 HANDLE_REJECT ВЫЗВАН! data: %s", data)

	// Обязательно отвечаем на callback
	ctx.Respond()

	job, exists := jobs[data]
	if !exists {
		log.Printf("❌ Заказ %s не найден", data)
		ctx.Send("❌ Заказ не найден")
		return nil
	}

	log.Printf("✅ Найден заказ для reject: %+v", job)

	ctx.Edit(fmt.Sprintf("❌ %s отклонен", job.Name))
	ctx.Bot().Send(&tele.User{ID: job.UserID}, fmt.Sprintf("❌ Ваш файл %s отклонен", job.Name))

	cleanupFile(job.Path)
	delete(jobs, data)

	return nil
}
