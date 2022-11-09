package telegram

import (
	"fmt"
	"time"
	"wikibot/pkg/config"

	tb "gopkg.in/tucnak/telebot.v2"
)

// StartBotAsync starts passive handling of commands
func StartBotAsync(cfg config.TelegramClientConfig) (*tb.Bot, error) {
	b, err := tb.NewBot(tb.Settings{
		Token:  cfg.BotToken,
		Poller: &tb.LongPoller{Timeout: 1 * time.Second},
	})

	if err != nil {
		return nil, err
	}

	// Persistent server
	b.Handle("/start", func(m *tb.Message) {
		b.Send(m.Sender, fmt.Sprintf("Hello I am the %s !", cfg.Name))
	})

	// Persistent server
	b.Handle("/help", func(m *tb.Message) {
		b.Send(m.Sender, "This bot is push based only, you cannot really interact with it")
	})

	go b.Start()

	return b, nil
}

// GetChat transforms a chat id int to a chat object
func GetChat(id int) *tb.Chat {
	return &tb.Chat{
		ID: int64(id),
	}
}
