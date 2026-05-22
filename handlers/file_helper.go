package handler

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// ensureFilesDir создает папку files в корне проекта
func ensureFilesDir() error {
	dir := "files"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// getUniqueFileName генерирует уникальное имя файла
func getUniqueFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	name := originalName[:len(originalName)-len(ext)]
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s_%d%s", name, timestamp, ext)
}

// cleanupFile удаляет файл
func cleanupFile(filePath string) {
	if err := os.Remove(filePath); err != nil {
		log.Printf("⚠️ Ошибка удаления файла %s: %v", filePath, err)
	} else {
		log.Printf("✅ Файл удален: %s", filePath)
	}
}

// generateJobID генерирует уникальный ID заказа
func generateJobID() string {
	return fmt.Sprintf("job_%d", time.Now().UnixNano())
}
