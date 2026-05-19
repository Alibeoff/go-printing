package printer

import (
	"fmt"
	"os/exec"
	"strconv"
)

func (p *Printer) Do() {
	var s []string

	if len(p.Arguments.Color) != 0 {
		arg := fmt.Sprintf("ColorModel=%s", p.Arguments.Color)
		s = append(s, "-o", arg)
	}

	if p.Arguments.Orientation {
		s = append(s, "-o orientation-requested=4") // альбомная
	} else if false {
		s = append(s, "-o orientation-requested=3") // книжная
	}

	if p.Arguments.AutoPull {
		s = append(s, "-o", "fit-to-page")
	}

	if p.Arguments.Copies != 0 {
		arg := "-n"
		s = append(s, arg, strconv.Itoa(p.Arguments.Copies))
	}

	if len(p.Arguments.Scale) != 0 {
		arg := fmt.Sprintf("-o media=Custom.%s", p.Arguments.Scale)
		s = append(s, arg)
	}

	if len(p.Arguments.Format) != 0 {
		arg := fmt.Sprintf("-o media=%s", p.Arguments.Format)
		s = append(s, arg)
	}

	if p.Arguments.Double {
		s = append(s, "-o", "sides=two-sided-long-edge")
	} else if false {
		s = append(s, "-o", "sides=one-side")
	}

	if len(p.Arguments.Pages) != 0 {
		s = append(s, "-P", p.Arguments.Pages)
	}

	if len(p.Arguments.File) != 0 {
		s = append([]string{"-d", p.UsePrinter}, s...)
		s = append(s, p.Arguments.File)
	}

	cmd := exec.Command("lp", s...)
	cmd.Run()
}
