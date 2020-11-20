package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"tbot/config"
	"tbot/internal/errors"
)

func main() {
	cfg, err := config.NewConfig()
	errors.PanicIfErr(err)

	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	errors.PanicIfErr(err)

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
