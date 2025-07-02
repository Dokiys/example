package color

import "fmt"

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"
)

var mute = false

func Mute() {
	mute = true
	if mute {
		return
	}
}

func Printf(format string, a ...any) {
	if mute {
		return
	}
	fmt.Printf(format, a...)
}

func Println(a ...any) {
	if mute {
		return
	}
	fmt.Println(a...)
}
func PrintfRed(format string, a ...any) {
	if mute {
		return
	}
	fmt.Printf(Red+format+Reset, a...)
}
func PrintfGreen(format string, a ...any) {
	if mute {
		return
	}
	fmt.Printf(Green+format+Reset, a...)
}
func PrintfYellow(format string, a ...any) {
	if mute {
		return
	}
	fmt.Printf(Yellow+format+Reset, a...)
}
func PrintfBlue(format string, a ...any) {
	if mute {
		return
	}
	fmt.Printf(Blue+format+Reset, a...)
}
func PrintfPurple(format string, a ...any) {
	if mute {
		return
	}
	fmt.Printf(Purple+format+Reset, a...)
}
func PrintfCyan(format string, a ...any) {
	if mute {
		return
	}
	fmt.Printf(Cyan+format+Reset, a...)
}
func PrintfGray(format string, a ...any) {
	if mute {
		return
	}
	fmt.Printf(Gray+format+Reset, a...)
}
func PrintfWhite(format string, a ...any) {
	fmt.Printf(White+format+Reset, a...)
}
