package main

import (
	"fmt"
	printer "go-printing"
)

func main() {
	var p printer.Printer
	p.PrinterList()                               // собирает данные о доступных принтерах
	fmt.Println(p.Names)                          // Отражает список доступных принтеров
	p.Select(0)                                   // выбор определенного принтера
	p.Arguments.ColorOption(p.UsePrinter, "Gray") // p.UsePrinter принтер который выбрали в p.Select
	p.Arguments.PagesOption("2-4")                // argument который позволяет распечать только указанные страницы
	p.Arguments.FormatOption("A4")                // argument который отражает формат(размер листа) список поддерживаемых форматов можно вызвать функцией printer.PageScale(p.UsePrinter) - принимает в себя имя принтера
	p.Arguments.FocusOption(1)                    // argument который печатает конкретную страницу которую вы пропишите.
	p.Arguments.ScaleOption("150x375")            // argument который указывает кастомный размер листа в пикселях, обязательно через x
	p.Arguments.CopiesOption(2)                   // argument который указывает количество копий
	p.Arguments.DoubleOption(p.UsePrinter, true)  // argument который позволит принтеру печатать с двух сторон, если вы вызовете этот метод, но ваш принтер не умеет печатать с двух сторон, то принтер просто проигнорирует такую настройку.
	p.Arguments.AutoPullOption(true)              // argument который позволит растянуть контент по всей странице если это возможно
	p.Arguments.MessageOption()                   // argument который в данный момент никак не используется(задел на будущие обновления)
	p.Arguments.OrientationOption(true)           // argument который позволяет выбрать ориентацию листа, книжную или альбомную
	p.Arguments.FileOption("text.txt")            // argument обязательный - путь к файлу
	p.Do()                                        // главный метод который стоит полную команду и отправляет файл на печать

	// для взаимодействия с принтером главное записать метод PrinterList Select и FileOption.
	// следует отметить что не может быть FormatOption и ScaleOption одновременно, так же как PagesOption и FocusOption
	// Выше я описал полное взаимодействие с билдером, прописав каждый метод, так делать не стоит!
	// Теперь предлагаю ознакомиться с неким функционалом:
	printer.KillProcess(1)                // удаление из очереди на печать(id смотреть в списке lpstat)
	printer.ColorList(p.UsePrinter)       // позволяет узнать цветовую палитру принтера
	printer.ActivePrintList(p.UsePrinter) // позволяет узнать очередь печати для конкретного принтера.
	// Если у вас есть идеи как улучшить данный софт, отправляйте свои решения и идеи.
	// В следущем обновлении скорее всего будет изменена работа со структурами.
}
