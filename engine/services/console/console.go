package console

import (
	"strings"
)

type console struct {
	permanent          string
	previousDrawnLines int
	toPrint            string
	print              func(string)
}

func newConsole() Console {
	return &console{
		permanent:          "",
		previousDrawnLines: 0,
		print:              func(s string) { print(s) },
	}
}

func clearCurrentLine() string { return "\033[2K" }
func goToPreviousLine() string { return "\033[1A" }

func (console *console) PrintPermanent(text string) {
	console.permanent += text + "\n"
}

func (console *console) Print(text string) {
	console.toPrint += text
}

func (console *console) Flush() {
	flushed := ""
	flushed += clearCurrentLine()
	for i := 0; i < console.previousDrawnLines; i++ {
		flushed += goToPreviousLine()
		flushed += clearCurrentLine()
	}
	console.print(flushed + console.permanent + console.toPrint)
	console.permanent = ""
	flushed += console.toPrint
	console.previousDrawnLines = strings.Count(console.toPrint, "\n")
	console.toPrint = ""
}
