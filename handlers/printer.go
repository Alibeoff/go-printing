package handler

import (
	"go-printing/pkg/printer"
	"sync"
)

type PrinterInfo struct {
	Name   string
	ID     string
	Status string
	Queue  []string
}

type PendingJob struct {
	UserID      int64
	UserName    string
	FilePath    string
	FileName    string
	PrinterID   int
	PrinterName string
	ChatID      int64
	MessageID   int // ID сообщения с запросом
}

var printerData PrinterInfo
var pt printer.Printer

// Хранилище ожидающих заказов
var pendingJobs = make(map[string]*PendingJob) // key: уникальный ID заказа
var pendingMutex sync.RWMutex

func (p *PrinterInfo) SetName(name string) {
	p.Name = name
}

func (p *PrinterInfo) SetID(id string) {
	p.ID = id
}

func (p *PrinterInfo) SetStatus(status string) {
	p.Status = status
}

func (p *PrinterInfo) SetQueue(queue []string) {
	p.Queue = queue
}
