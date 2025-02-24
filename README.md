# go-printing

`go-printing` — это обёртка над командами `lp` для печати на принтерах, написанная на Golang. Проект разработан для упрощения работы с командами печати в Linux, предоставляя удобный интерфейс для управления параметрами печати.

Для работы этого пакета необходимо установить утилиту CUPS (Common Unix Printing System) на вашем Linux-дистрибутиве.

## Установка CUPS

### Ubuntu/Debian
```bash
sudo apt update
sudo apt install cups
```
### Fedora
```bash
sudo dnf install cups 
``` 
### Arch
```bash
sudo pacman -S cups
```
### После установки CUPS проверьте наличие команды lp в терминале:
```bash
lp --help
```
Если команда выводит справку, значит CUPS установлен правильно.

## Get Starting

Для запуска печати существует множество методов, которые разделяются на
аргументы и настройки принтера. Чтобы присвоить аргументам значение, 
используйте следующие методы:
- FocusOption: текущая страница

- OrientationOption: альбомная или книжная ориентация

- AutoPullOption: автоматическое растягивание по масштабу страницы

- CopiesOption: количество копий для печати

- FormatOption: формат листа для печати (можно выбрать из списка или указать свои размеры)

- ScaleOption: самостоятельная регулировка масштаба текста или изображения

- DoubleOption: включение двусторонней печати

- ColorOption: выбор цветовой палитры для печати

- MessageOption: пока не используется (для будущих обновлений)

- PagesOption: добавление в очередь печати определённых страниц или перечисление страниц через тире

После заполнения необходимых аргументов для печати можно приступить к запуску. 
Вот пример кода на Golang:
```go   
package main 
import "printer" 

var p printer.Printer 
p.PrinterList()                               
p.Select(0)                                   
p.Arguments.FileOption("path-to-file")
p.Do() // главный метод который стоит полную команду и отправляет файл на печать
```
Тут мы явно указываем путь к файлу и используем функцию Do которая собирает все аргументы и отправляет файл указанный выше на печать.
Далее рассмотрим какой есть функционал у настроек аргументов.
## Аргументы
Для того чтобы в функцию Do передались аргументы их нужно заполнить, для этого существуют Options методы у printer.Printer.Arguments.

```go
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
```
Cледует отметить что не может быть FormatOption и ScaleOption одновременно, так же как PagesOption и FocusOption
Выше я описал полное взаимодействие с билдером, прописав каждый метод, так делать не стоит!

## Utils

Теперь предлагаю ознакомиться с неким функционалом:
```go
printer.KillProcess(1)                // удаление из очереди на печать(id смотреть в списке lpstat)
printer.ColorList(p.UsePrinter)       // позволяет узнать цветовую палитру принтера
printer.ActivePrintList(p.UsePrinter) // позволяет узнать очередь печати для конкретного принтера.
```
Если у вас есть идеи как улучшить данный софт, отправляйте свои решения и идеи.
В следущем обновлении скорее всего будет изменена работа со структурами.
