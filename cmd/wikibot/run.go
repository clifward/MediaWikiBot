package main

import (
	"context"
	"fmt"
	"strings"
	"wikibot/pkg/config"
	"wikibot/pkg/telegram"
	"wikibot/pkg/wikicollector"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func run(ctx context.Context, cfg *config.Config) error {
	cfgBytes, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	logrus.Info("config\n", string(cfgBytes))

	bot, err := telegram.StartBotAsync(cfg.TelegramClientConfig)
	if err != nil {
		return err
	}

	wc := wikicollector.NewWikiClient(cfg.WikiClientConfig)

	raCh := make(chan []wikicollector.RecentAction, 10)
	done := make(chan struct{}, 1)
	go func() {
		// Send message when there are new recent actions
		tc := telegram.GetChat(cfg.TelegramClientConfig.ChanID)
		admin := telegram.GetChat(cfg.TelegramClientConfig.Admin)
		messages := make([]string, 0)
		for recentActions := range raCh {
			for _, recentAction := range recentActions {
				messages = append(messages, recentAction.String())
				messages = append(messages, "\n"+`_______________________`+"\n\n")
			}
			finalMessage := strings.Join(messages[:len(messages)-1], "")
			if len(finalMessage) == 0 {
				return
			}
			_, err = bot.Send(tc, finalMessage, "HTML")
			if err != nil {
				logrus.Error(err)
				bot.Send(admin, fmt.Sprint(err))
				bot.Send(admin, fmt.Sprint(finalMessage))
			}
		}
		done <- struct{}{}

	}()
	wc.StartFetchingActions(ctx, cfg.Mode, raCh)
	close(raCh)
	<-done
	bot.Stop()

	return nil
}
