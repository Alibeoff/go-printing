package printer

import (
	"fmt"
	"strings"
)

type Command struct {
	Focus       int    // на какой странице сейчас
	Orientation bool   // альбомная или обычная
	AutoPull    bool   // aвтоматическое растягивание по масштабу страницы
	Copies      int    //-n аргумет  | количество копий печати
	Format      string // a4 a5 a6 и тд  или конкретные размеры: 720x1080
	Scale       string // самостоятельная регулировка маштабом текста или картинки
	Double      bool   // -o sides=two-sided-{short\long}-edge -альбомная двусторонняя\ книжная | функционал двусторонней печати(если позволяет принтер)
	Color       string // выбор цветовой палитры для печати
	Message     string
	Pages       string // функционал определенных страниц печати.
	File        string
}

func (c *Command) FileOption(file string) *Command {
	c.File = file
	return c
}

func (c *Command) FocusOption(page int) *Command {
	c.Focus = page
	return c
}

func (c *Command) OrientationOption(let bool) *Command {
	c.Orientation = let
	return c
}

func (c *Command) AutoPullOption(let bool) *Command {
	c.AutoPull = let
	return c
}

func (c *Command) CopiesOption(s int) *Command {
	c.Copies = s
	return c
}

func (c *Command) FormatOption(input string) *Command {
	c.Format = input
	return c
}

func (c *Command) ScaleOption(input string) *Command {
	input = strings.ReplaceAll(input, " ", "")
	c.Scale = input
	return c
}

func (c *Command) DoubleOption(pu string, let bool) *Command {
	b, err := DuplexBool(pu)
	if err != nil {
		fmt.Println(err)
	}
	if b {
		c.Double = let
	}
	return c
}

func (c *Command) ColorOption(pu string, input string) *Command {
	c.Color = input
	return c
}

func (c *Command) MessageOption() *Command {
	return c
}

func (c *Command) PagesOption(input string) *Command {
	out, err := Unic(input)
	if err != nil {
		fmt.Println("Ошибка:", err)
	}
	c.Pages = out

	return c
}
