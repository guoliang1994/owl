package log

import (
	"fmt"
	"github.com/fatih/color"
)

const Print = 0
const PrintLn = 1

type Color string

const (
	Black  = "Black"
	White  = "White"
	Red    = "Red"
	Blue   = "Blue"
	Yellow = "Yellow"
	Green  = "Green"
)

// PrintLnBlue 行打印，蓝色
func PrintLnBlue(content ...any) {
	print(PrintLn, Blue, content)
}

// PrintLnYellow 行打印，黄色
func PrintLnYellow(content ...any) {
	print(PrintLn, Yellow, content)
}

// PrintLnRed 行打印，红色
func PrintLnRed(content ...any) {
	print(PrintLn, Red, content)
}

// PrintLnGreen 行打印，绿色
func PrintLnGreen(content ...any) {
	print(PrintLn, Green, content)
}

// PrintLnBlack 行打印，黑色
func PrintLnBlack(content ...any) {
	print(PrintLn, Black, content)
}

// PrintLnWhite 行打印，白色
func PrintLnWhite(content ...any) {
	print(PrintLn, White, content)
}

// PrintBlue 行内打印，蓝色
func PrintBlue(content ...any) {
	print(Print, Blue, content)
}

// PrintYellow 行内打印，黄色
func PrintYellow(content ...any) {
	print(Print, Yellow, content)
}

// PrintRed 行内打印，红色
func PrintRed(content ...any) {
	print(Print, Red, content)
}

// PrintGreen 行内打印，绿色
func PrintGreen(content ...any) {
	print(Print, Green, content)
}

// PrintBlack 行内打印，黑色
func PrintBlack(content ...any) {
	print(Print, Black, content)
}

// PrintWhite 行内打印，白色
func PrintWhite(content ...any) {
	print(Print, White, content)
}

func print(typ int, c string, content ...any) {
	var colorFun func(format string, a ...interface{}) string
	switch c {
	case White:
		colorFun = color.WhiteString
	case Black:
		colorFun = color.BlackString
	case Blue:
		colorFun = color.BlueString
	case Green:
		colorFun = color.GreenString
	case Red:
		colorFun = color.RedString
	case Yellow:
		colorFun = color.YellowString
	default:
		colorFun = color.WhiteString
	}

	switch typ {
	case PrintLn:
		fmt.Println(colorFun("%v", content))
	default:
		fmt.Print(colorFun("%v", content))
	}
}
