package bot

import (
	"fmt"
	"log"
	"regexp"
	"sync"
	"time"

	"gopkg.in/tucnak/telebot.v2"

	"github.com/iley/litsen/internal/camera"
)

var (
	watchCommandRe = regexp.MustCompile(`^/\w+\s+(\w+)`)
)

type Bot struct {
	tb *telebot.Bot

	watchesMu sync.Mutex
	watches   map[int64]*watch

	imagesDir string
}

type Settings struct {
	TelegramToken string
	UserWhitelist []string
	ImagesDir     string
}

func New(settings Settings) (*Bot, error) {
	poller := &telebot.LongPoller{Timeout: 10 * time.Second}

	authPoller := telebot.NewMiddlewarePoller(poller, func(upd *telebot.Update) bool {
		if upd.Message == nil {
			return true
		}
		return isWhitelisted(upd.Message.Sender.Username, settings.UserWhitelist)
	})

	tb, err := telebot.NewBot(telebot.Settings{
		Token:  settings.TelegramToken,
		Poller: authPoller,
	})

	if err != nil {
		return nil, fmt.Errorf("could not initialize telebot %w", err)
	}

	b := &Bot{
		tb:      tb,
		watches: make(map[int64]*watch),
	}

	if settings.ImagesDir == "" {
		b.imagesDir = "."
	} else {
		b.imagesDir = settings.ImagesDir
	}

	tb.Handle("/photo", b.handlePhoto)
	tb.Handle("/watch", b.handleWatch)
	tb.Handle("/stopwatch", b.handleStopWatch)

	return b, nil
}

func (b *Bot) Run() {
	log.Printf("running the bot")
	b.tb.Start()
}

func (b *Bot) handlePhoto(m *telebot.Message) {
	b.takePhoto(m.Chat)
}

func (b *Bot) handleWatch(m *telebot.Message) {
	groups := watchCommandRe.FindStringSubmatch(m.Text)
	if len(groups) != 2 {
		b.tb.Send(m.Chat, "Usage: /watch <interval>")
		return
	}

	periodStr := groups[1]

	period, err := time.ParseDuration(periodStr)
	if err != nil {
		b.tb.Send(m.Chat, fmt.Sprintf("Could not parse duration %s: %s", periodStr, err))
		return
	}

	log.Printf("staring a watch with period %s", period)

	recipient := m.Chat
	w := newWatch(period, func() {
		b.takePhoto(recipient)
	})

	b.watchesMu.Lock()
	b.watches[m.Chat.ID] = w
	b.watchesMu.Unlock()

	w.start()

	b.tb.Send(m.Chat, fmt.Sprintf("Started a watch with period %s", period))
}

func (b *Bot) handleStopWatch(m *telebot.Message) {
	b.watchesMu.Lock()
	w, found := b.watches[m.Chat.ID]
	delete(b.watches, m.Chat.ID)
	b.watchesMu.Unlock()

	if found {
		w.stop()
		b.tb.Send(m.Chat, "Stopped the watch")
	} else {
		b.tb.Send(m.Chat, "No watch currently running")
	}
}

func (b *Bot) takePhoto(recipient telebot.Recipient) {
	b.tb.Notify(recipient, telebot.UploadingPhoto)
	log.Printf("taking a photo...")

	imagePath, err := camera.TakePhoto(b.imagesDir)
	if err != nil {
		log.Printf("could not take a photo: %s", err)
		b.tb.Send(recipient, fmt.Sprintf("Could not take a photo: %s", err))
		return
	}

	log.Printf("saved a photo at %s", imagePath)

	img := &telebot.Photo{File: telebot.FromDisk(imagePath)}
	b.tb.Send(recipient, img)
}

func isWhitelisted(username string, whitelist []string) bool {
	for _, whitelisted := range whitelist {
		if username == whitelisted {
			return true
		}
	}
	log.Printf("unauthenticated user %s", username)
	return false
}
