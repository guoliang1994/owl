package log

import "owl/contract"

type ConsoleImpl struct {
}

func (c ConsoleImpl) Emergency(content ...interface{}) {
	PrintLnRed(contract.EMERGENCY, content)
}

func (c ConsoleImpl) Alert(content ...interface{}) {
	PrintLnRed(contract.ALERT, content)
}

func (c ConsoleImpl) Critical(content ...interface{}) {
	PrintLnYellow(contract.CRITICAL, content)
}

func (c ConsoleImpl) Error(content ...interface{}) {
	PrintLnRed(contract.ERROR, content)
}

func (c ConsoleImpl) Warning(content ...interface{}) {
	PrintLnYellow(contract.WARNING, content)
}

func (c ConsoleImpl) Notice(content ...interface{}) {
	PrintLnBlue(contract.NOTICE, content)
}

func (c ConsoleImpl) Info(content ...interface{}) {
	PrintLnWhite(contract.INFO, content)
}

func (c ConsoleImpl) Debug(content ...interface{}) {
	PrintLnWhite(contract.DEBUG, content)
}
