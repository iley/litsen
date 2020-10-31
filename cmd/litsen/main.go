package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/iley/litsen/internal/bot"
)

func main() {
	token := flag.String("token", "", "Telegram API token")
	flag.Parse()

	if *token == "" {
		fmt.Fprint(os.Stderr, "missing --token")
		os.Exit(1)
	}

	b, err := bot.New(*token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running the bot: %s", err)
		os.Exit(1)
	}

	b.Run()
}
