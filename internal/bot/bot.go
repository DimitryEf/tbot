package bot

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
	"tbot/config"
	wiki2 "tbot/internal/wiki"
)

type Bot struct {
	bot  *tgbotapi.BotAPI
	cfg  *config.Config
	wiki *wiki2.Wiki
}

func NewBot(cfg *config.Config, wiki *wiki2.Wiki) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, err
	}
	bot.Debug = true
	return &Bot{
		bot:  bot,
		cfg:  cfg,
		wiki: wiki,
	}, nil
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID

		if strings.HasPrefix(msg.Text, "w ") {
			query := msg.Text[2:]
			title, err := b.wiki.Query(query)
			if err != nil {
				//log.Printf("error: %v", err)
				msg.Text = fmt.Sprintf("error: %v", err)
				b.bot.Send(msg)
				continue
			}
			msg.Text = title
			b.bot.Send(msg)
			continue
		}

		b.bot.Send(msg)
	}
	return nil
}
