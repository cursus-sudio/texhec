package console

type Console interface {
	ClearConsole()
	LogToConsole(string)
	ClearAndLogToConsole(string)
}
