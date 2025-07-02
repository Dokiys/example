package color

import (
	"log/slog"
	"testing"
)

func TestPrint(t *testing.T) {
	PrintfRed("Red\n")
	PrintfGreen("Green\n")
	PrintfYellow("Yellow\n")
	PrintfBlue("Blue\n")
	PrintfPurple("Purple\n")
	PrintfCyan("Cyan\n")
	PrintfGray("Gray\n")
	PrintfWhite("White\n")

	Mute()
	Printf("No Print")
}

func TestName(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Debug("123")
	// TODO[Dokiy] to be continued! (2025/7/2)
	// 自定义slog.SetDefault()
	// 利用slog的level来当颜色输出
	// slog.Log()
}
