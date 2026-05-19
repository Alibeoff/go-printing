package printer

import (
	"fmt"
	"log"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type Printer struct {
	Names      []string
	UsePrinter string
	Arguments  Command
}

type PrintJob struct {
	Rank  string
	JobID int
	Files string
}

// этот метод
func GetAvailablePrinters() ([]string, error) {
	// Выполняем команду lpstat -a
	out, err := exec.Command("lpstat", "-a").Output()
	if err != nil {
		return nil, err
	}

	// Разделяем вывод на строки
	lines := strings.Split(string(out), "\n")

	// Извлекаем имена принтеров
	printers := make([]string, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) > 0 {
			printers = append(printers, parts[0])
		}
	}

	return printers, nil
}

func (p *Printer) PrinterList() {
	printers, err := GetAvailablePrinters()
	if err != nil {
		log.Fatal(err)
	}
	for _, text := range printers {
		p.Names = append(p.Names, text)
	}
}

func (p *Printer) Select(id int) string {
	p.UsePrinter = p.Names[id]
	return p.UsePrinter
}

func DuplexBool(printerName string) (bool, error) {
	ppdPath := fmt.Sprintf("/etc/cups/ppd/%s.ppd", printerName)

	// Проверяем, существует ли PPD-файл
	if _, err := exec.Command("ls", ppdPath).CombinedOutput(); err != nil {
		return false, fmt.Errorf("PPD-файл для принтера '%s' не найден: %v", printerName, err)
	}

	// Выполняем команду grep для поиска строки "Duplex" в PPD-файле
	cmd := exec.Command("grep", "-i", "duplex", ppdPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Если grep не нашёл строку, это не обязательно ошибка
		if strings.Contains(string(output), "No such file or directory") {
			return false, fmt.Errorf("PPD-файл для принтера '%s' не найден", printerName)
		}
		return false, nil // Duplex не поддерживается
	}

	// Если в выводе есть строки с "Duplex", принтер поддерживает двустороннюю печать
	return true, nil
}

func PageScale(printerName string) ([]string, error) {
	// Выполняем команду lpoptions для получения настроек принтера
	cmd := exec.Command("lpoptions", "-p", printerName, "-l")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении команды lpoptions: %v", string(output))
	}

	// Ищем строку, содержащую PageSize
	lines := strings.Split(string(output), "\n")
	var pageSizesLine string
	for _, line := range lines {
		if strings.Contains(line, "PageSize") {
			pageSizesLine = line
			break
		}
	}

	if pageSizesLine == "" {
		return nil, fmt.Errorf("опция PageSize не найдена в настройках принтера")
	}

	// Извлекаем допустимые размеры страниц
	// Пример строки: "PageSize/Media Size: Letter Legal A4 *A5 Custom.523x523mm"
	parts := strings.Split(pageSizesLine, ":")
	if len(parts) < 2 {
		return nil, fmt.Errorf("неверный формат строки PageSize")
	}

	// Разделяем размеры страниц
	sizes := strings.Fields(parts[1])

	// Убираем звёздочку (*), которая обозначает размер по умолчанию
	for i, size := range sizes {
		sizes[i] = strings.TrimPrefix(size, "*")
	}

	return sizes, nil
}

func ColorList(printerName string) ([]string, error) {
	// Выполняем команду lpoptions для получения настроек принтера
	cmd := exec.Command("lpoptions", "-p", printerName, "-l")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении команды lpoptions: %v", string(output))
	}

	// Ищем строки, содержащие ColorModel или OutputMode
	lines := strings.Split(string(output), "\n")
	var colorOptions []string

	for _, line := range lines {
		if strings.Contains(line, "ColorModel") || strings.Contains(line, "OutputMode") {
			// Извлекаем часть строки после двоеточия
			parts := strings.Split(line, ":")
			if len(parts) < 2 {
				continue
			}

			// Разделяем варианты цвета
			options := strings.Fields(parts[1])
			for _, option := range options {
				// Убираем звёздочку (*), которая обозначает вариант по умолчанию
				option = strings.TrimPrefix(option, "*")
				colorOptions = append(colorOptions, option)
			}
		}
	}

	if len(colorOptions) == 0 {
		return nil, fmt.Errorf("варианты цвета печати не найдены")
	}

	return colorOptions, nil
}

func Unic(input string) (string, error) {
	input = strings.ReplaceAll(input, " ", "")
	parts := strings.Split(input, ",")

	var numbers []int

	for _, part := range parts {
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return "", fmt.Errorf("неверный формат диапазона: %s", part)
			}

			start, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				return "", fmt.Errorf("неверное число в диапазоне: %s", rangeParts[0])
			}
			end, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				return "", fmt.Errorf("неверное число в диапазоне: %s", rangeParts[1])
			}

			if start > end {
				return "", fmt.Errorf("неверный диапазон: %s (начальное число больше конечного)", part)
			}

			for i := start; i <= end; i++ {
				numbers = append(numbers, i)
			}
		} else {
			num, err := strconv.Atoi(part)
			if err != nil {
				return "", fmt.Errorf("неверное число: %s", part)
			}
			numbers = append(numbers, num)
		}
	}

	uniqueNumbers := removeDuplicates(numbers)
	sort.Ints(uniqueNumbers)

	var result []string
	for _, num := range uniqueNumbers {
		result = append(result, strconv.Itoa(num))
	}

	return strings.Join(result, ","), nil
}

func removeDuplicates(numbers []int) []int {
	unique := make(map[int]bool)
	var result []int
	for _, num := range numbers {
		if !unique[num] {
			unique[num] = true
			result = append(result, num)
		}
	}
	return result
}

func ActivePrintList(printerName string) ([]PrintJob, error) {
	cmd := exec.Command("lpq", "-P", printerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении команды lpq: %v", string(output))
	}

	lines := strings.Split(string(output), "\n")

	var jobs []PrintJob
	for _, line := range lines[2:] {
		if strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		jobID, err := strconv.Atoi(fields[2])
		if err != nil {
			continue
		}

		job := PrintJob{
			Rank:  fields[0],
			JobID: jobID,
			Files: fields[3],
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func KillProcess(id int) {
	cmd := exec.Command("cancel", strconv.Itoa(id))
	cmd.Run()
}
