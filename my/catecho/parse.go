package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	flag.Parse()
	defer os.Exit(2)

	if len(flag.Args()) > 0 {
		catecho(strings.Join(flag.Args(), " "))
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		catecho(scanner.Text())
	}
}

func catecho(str string) {
	fmt.Print(str)
}
