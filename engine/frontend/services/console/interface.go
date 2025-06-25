package console

type Console interface {
	ClearConsole()
	LogPermanentlyToConsole(string)
	LogToConsole(string)
	ClearAndLogToConsole(string)
}
