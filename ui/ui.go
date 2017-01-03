package ui

import (
	"github.com/gizak/termui"
	"fmt"
	"strings"
)

type UIEvent string

const (
	PAUSE    = UIEvent("/sys/kbd/p")
	CONTINUE = UIEvent("/sys/kbd/o")
	STEP     = UIEvent("/sys/kbd/s")
	QUIT     = UIEvent("/sys/kbd/q")
)

var (
	instructions = []string{
		"[s] Step",
		"[p] Pause execution",
		"[o] Continue execution",
		"[q] Close the VM",
	}
)

type UI struct{}

func (ui *UI) Run(ram []byte, registers []string) {
	if err := termui.Init(); err != nil {
		panic(err)
	}
	defer termui.Close()

	// Register default handlers
	ui.Handle(QUIT, func() {
		termui.StopLoop()
	})

	ui.RenderRegistersWidget(registers)
	ui.RenderStatusWidget("initialized")
	ui.RenderInstructionsWidget()
	ui.RenderRAMWidget(ram)
	ui.RenderExecutingCommand("")

	termui.Loop()
}

func (ui *UI) Handle(ev UIEvent, f func()) {
	termui.Handle(string(ev), func(termui.Event) { f() })
}

func (ui *UI) RenderRAMWidget(ram []byte) {
	ls := termui.NewList()

	rowCount := int(len(ram) / 16) + 1
	rows := make([]string, rowCount, rowCount)

	// Draw rows
	for i := 0; i < rowCount; i++ {
		// Draw cols
		colsCount := len(ram) - i*16
		if colsCount > 16 {
			colsCount = 16
		}

		colsCount += 1

		cols := make([]string, colsCount, colsCount)
		cols[0] = fmt.Sprintf("%06x:", i*16)
		for j := 1; j < colsCount; j++ {
			cols[j] = fmt.Sprintf("%02x", ram[i*16 + j - 1])
		}
		rows[i] = strings.Trim(strings.Join(cols[:], " "), " ")
	}

	ls.Items = rows[:]
	ls.ItemFgColor = termui.ColorYellow
	ls.BorderLabel = "RAM"
	ls.Height = 48
	ls.Width = 58
	ls.Y = 0
	ls.X = 30

	termui.Render(ls)
}
func (ui *UI) RenderStatusWidget(status string) {
	st := termui.NewPar(status)
	st.Height = 3
	st.Width = 30
	st.Y = 2 + len(instructions)
	st.BorderLabel = "Status"
	termui.Render(st)
}
func (ui *UI) RenderRegistersWidget(str []string) {
	ls := termui.NewList()
	ls.Items = str
	ls.ItemFgColor = termui.ColorYellow
	ls.BorderLabel = "Registers"
	ls.Height = 11
	ls.Width = 30
	ls.Y = 3 + 2 + len(instructions)

	termui.Render(ls)
}
func (ui *UI) RenderInstructionsWidget() {
	ls := termui.NewList()
	ls.Items = instructions
	ls.ItemFgColor = termui.ColorYellow
	ls.BorderLabel = "Instructions"
	ls.Height = 2 + len(instructions)
	ls.Width = 30
	ls.Y = 0

	termui.Render(ls)
}
func (ui *UI) RenderExecutingCommand(cmd string) {
	st := termui.NewPar(cmd)
	st.Height = 3
	st.Width = 30
	st.Y = 16 + len(instructions)
	st.BorderLabel = "Last executed"
	termui.Render(st)
}
