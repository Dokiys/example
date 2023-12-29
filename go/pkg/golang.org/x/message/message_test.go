package message

import (
	"fmt"
	"testing"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// https://www.alexedwards.net/blog/i18n-managing-translations
func TestMessage(t *testing.T) {
	// init
	message.SetString(language.Chinese, "%s went to %s.\n", "%s去了%s。\n")
	message.SetString(language.Chinese, "%s has been stolen.\n", "%s被偷走了。\n")
	message.SetString(language.Chinese, "How are you?\n", "你好吗?\n")
	message.SetString(language.AmericanEnglish, "%s went to %s.\n", "%s is in %s.\n")
	message.SetString(language.AmericanEnglish, "%s has been stolen.\n", "%s has been stolen.\n")

	p := message.NewPrinter(language.Chinese)
	p.Printf("%s went to %s.\n", "彼得", "英格兰")
	p.Printf("%s has been stolen.\n", "宝石")
	p.Printf("How are you?\n")
	fmt.Println()
	p = message.NewPrinter(language.AmericanEnglish)
	p.Printf("%s went to %s.\n", "Peter", "England")
	p.Printf("%s has been stolen.\n", "The Gem")
}
