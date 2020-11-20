package bot

import (
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

		msg.Text, err = b.Act(msg.Text)
		if err != nil {
			log.Printf("error: %v", err)
		}

		_, err = b.bot.Send(msg)
		if err != nil {
			log.Printf("error: %v", err)
		}
	}
	return nil
}

func (b *Bot) Act(msgText string) (string, error) {
	if len(msgText) < 3 {
		return msgText, nil
	}

	for service, acts := range b.cfg.Stg.Services {
		for _, act := range acts {
			if strings.HasPrefix(msgText, act) {
				text := msgText[len(act):]
				return b.takeAction(service, text)
			}
		}
	}
	return msgText, nil
}

func (b *Bot) takeAction(service string, text string) (string, error) {
	switch service {
	case b.cfg.Stg.WikiStg.Tag:
		title, err := b.wiki.Query(text)
		if err != nil {
			return "", err
		}
		return title, nil
	}
	return text, nil
}
