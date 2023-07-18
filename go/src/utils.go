package src

import (
	"fmt"
	"unicode/utf8"
)

func fmtHan(len int, str string) string {
	var count int
	for _, char := range str {
		if utf8.RuneLen(char) > 1 {
			count++
		}
	}

	return fmt.Sprintf(fmt.Sprintf("%%%ds", len-count), str)
}
