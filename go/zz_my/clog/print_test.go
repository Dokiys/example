package clog

import (
	"testing"
)

func TestPrint(t *testing.T) {
	Printf(LevelRed, "Red\n")
	Printf(LevelGreen, "Green\n")
	Printf(LevelYellow, "Yellow\n")
	Printf(LevelBlue, "Blue\n")
	Printf(LevelPurple, "Purple\n")
	Printf(LevelCyan, "Cyan\n")
	Printf(LevelGray, "Gray\n")
	Printf(LevelWhite, "White\n")
	Println("---")
	SetDefaultColorLevel(LevelRed + LevelBlue)
	Printf(LevelRed, "Red\n")
	Printf(LevelGreen, "Green\n")
	Printf(LevelYellow, "Yellow\n")
	Printf(LevelBlue, "Blue\n")
	Printf(LevelPurple, "Purple\n")
	Printf(LevelCyan, "Cyan\n")
	Printf(LevelGray, "Gray\n")
	Printf(LevelWhite, "White\n")
}
