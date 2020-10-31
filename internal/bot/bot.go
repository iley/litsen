package bot

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/tucnak/telebot.v2"

	"github.com/iley/litsen/internal/camera"
)

type Bot struct {
	tb *telebot.Bot
}

func New(token string) (*Bot, error) {
	tb, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		return nil, fmt.Errorf("could not initialize telebot %w", err)
	}

	b := &Bot{tb: tb}

	tb.Handle("/photo", b.takePhoto)

	return b, nil
}

func (b *Bot) Run() {
	log.Printf("running the bot")
	b.tb.Start()
}

func (b *Bot) takePhoto(m *telebot.Message) {
	b.tb.Notify(m.Sender, telebot.UploadingPhoto)
	log.Printf("taking a photo...")

	imagePath, err := camera.TakePhoto(".") // TODO: Make directory configurable.
	if err != nil {
		log.Printf("could not take a photo: %s", err)
		b.tb.Send(m.Sender, fmt.Sprintf("Could not take a photo: %s", err))
		return
	}

	img := &telebot.Photo{File: telebot.FromDisk(imagePath)}
	b.tb.Send(m.Sender, img)
}
