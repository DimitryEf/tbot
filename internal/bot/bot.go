package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
	"tbot/config"
	"tbot/services"
)

type Bot struct {
	bot *tgbotapi.BotAPI
	cfg *config.Config
	sm  *services.ServiceManager
}

func NewBot(cfg *config.Config, sm *services.ServiceManager) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, err
	}
	bot.Debug = true
	return &Bot{
		bot: bot,
		cfg: cfg,
		sm:  sm,
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

		msg.Text, err = b.act(msg.Text)
		if err != nil {
			log.Printf("error: %v", err)
		}

		if msg.Text == "" {
			msg.Text = "не получилось..."
		}

		if len(msg.Text) > 4096 {
			msg.Text = msg.Text[:4096]
		}

		_, err = b.bot.Send(msg)
		if err != nil {
			log.Printf("error: %v", err)
		}
	}
	return nil
}

func (b *Bot) act(msgText string) (string, error) {
	if len(msgText) < 3 {
		return msgText, nil
	}

	for serviceTag, acts := range b.cfg.Stg.Services {
		for _, act := range acts {
			if strings.HasPrefix(msgText, act) {
				text := msgText[len(act):]
				return b.chooseAction(serviceTag, text)
			}
		}
	}
	return msgText, nil
}

func (b *Bot) chooseAction(serviceTag string, text string) (string, error) {
	switch serviceTag {
	case b.cfg.Stg.WikiStg.Tag:
		return b.takeAction(b.cfg.Stg.WikiStg.Tag, text)
	case b.cfg.Stg.NewtonStg.Tag:
		return b.takeAction(b.cfg.Stg.NewtonStg.Tag, text)
	}
	return text, nil
}

func (b *Bot) takeAction(tag string, text string) (string, error) {
	result, err := b.sm.Services[tag].Query(text)
	if err != nil {
		return "", err
	}
	return result, nil
}
