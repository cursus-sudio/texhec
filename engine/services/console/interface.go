package console

type Console interface {
	PrintPermanent(string)
	Print(string)
	Flush()
}
