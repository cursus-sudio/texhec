package console

import "strings"

type console struct {
	drawnLines int
	print      func(string)
}

func newConsole() Console {
	return &console{
		drawnLines: 0,
		print:      func(s string) { print(s) },
	}
}

func clearCurrentLine() string { return "\033[2K" }
func goToPreviousLine() string { return "\033[1A" }

func (console *console) ClearConsole() {
	text := ""
	text += clearCurrentLine()
	for i := 0; i < console.drawnLines; i++ {
		text += goToPreviousLine()
		text += clearCurrentLine()
	}
	console.print(text)
	console.drawnLines = 0
}

func (console *console) LogToConsole(text string) {
	console.print(text)
	console.drawnLines += strings.Count(text, "\n")
}

func (console *console) ClearAndLogToConsole(text string) {
	flushed := ""
	flushed += clearCurrentLine()
	for i := 0; i < console.drawnLines; i++ {
		flushed += goToPreviousLine()
		flushed += clearCurrentLine()
	}
	flushed += text
	console.print(flushed)
	console.drawnLines = strings.Count(text, "\n")
}
