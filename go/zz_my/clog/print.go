package clog

import (
	"fmt"
	"os"
)

type Level int

const (
	CharReset  = "\033[0m"
	CharRed    = "\033[31m"
	CharGreen  = "\033[32m"
	CharYellow = "\033[33m"
	CharBlue   = "\033[34m"
	CharPurple = "\033[35m"
	CharCyan   = "\033[36m"
	CharGray   = "\033[37m"
	CharWhite  = "\033[97m"
)
const (
	LevelDefault Level = 1 << iota
	LevelRed
	LevelGreen
	LevelYellow
	LevelBlue
	LevelPurple
	LevelCyan
	LevelGray
	LevelWhite
)

var colorMap = map[Level]string{
	LevelDefault: CharReset,
	LevelRed:     CharRed,
	LevelGreen:   CharGreen,
	LevelYellow:  CharYellow,
	LevelBlue:    CharBlue,
	LevelPurple:  CharPurple,
	LevelCyan:    CharCyan,
	LevelGray:    CharGray,
	LevelWhite:   CharWhite,
}

var defaultColorLevel = LevelDefault

func SetDefaultColorLevel(levels ...Level) {
	var level = LevelDefault
	for _, l := range levels {
		level |= l
	}
	defaultColorLevel = level
}
func Printf(level Level, format string, a ...any) {
	if defaultColorLevel != LevelDefault && (level <= 0 || defaultColorLevel&(level) <= 0) {
		return
	}

	msg := colorMap[level] + format + CharReset
	fmt.Fprintf(os.Stdout, msg, a...)
}

func Println(a ...any) {
	fmt.Println(a...)
}
