package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/iley/litsen/internal/bot"
)

func main() {
	token := flag.String("token", "", "Telegram API token")
	whitelistStr := flag.String("whitelist", "", "Telegram user whiltelist")
	imagesDir := flag.String("images", "", "Directory to store images")
	flag.Parse()

	if *token == "" {
		fmt.Fprintln(os.Stderr, "missing --token")
		os.Exit(1)
	}

	if *whitelistStr == "" {
		fmt.Fprintln(os.Stderr, "missing --whitelist")
		os.Exit(1)
	}

	whitelist := strings.Split(*whitelistStr, ",")

	b, err := bot.New(bot.Settings{
		TelegramToken: *token,
		UserWhitelist: whitelist,
		ImagesDir:     *imagesDir,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running the bot: %s", err)
		os.Exit(1)
	}

	b.Run()
}
