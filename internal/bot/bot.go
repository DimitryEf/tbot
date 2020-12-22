package bot

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
	"tbot/config"
	"tbot/services"
)

const textMaxLength int = 4096

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
		//oneMsg.ReplyToMessageID = update.Message.MessageID

		for _, val := range b.cfg.Stg.Services["golang"] {
			if strings.HasPrefix(msg.Text, val) {
				msg.ParseMode = tgbotapi.ModeMarkdownV2
				break
			}
		}

		//oneMsg.Text = "\"Toggle Off\\\\n\""
		//b.bot.Send(oneMsg)
		//continue

		text, err := b.act(msg.Text)
		if err != nil {
			log.Printf("error: %v", err)
			text = fmt.Sprintf("error: %v", err)
		}

		if text == "" {
			text = "не получилось\\.\\.\\."
		}

		var msgs []tgbotapi.MessageConfig
		limit := 3
		offset := textMaxLength
		for {
			if len(text) > textMaxLength {
				if limit <= 0 {
					break
				}
				limit--
				oneMsg := msg
				oneMsg.Text = text[:textMaxLength]
				if oneMsg.Text[len(oneMsg.Text)-1] == '\\' && oneMsg.Text[len(oneMsg.Text)-2] != '\\' {
					oneMsg.Text = oneMsg.Text[:textMaxLength-1]
					offset--
				}
				text = text[offset:]
				if oneMsg.ParseMode == tgbotapi.ModeMarkdownV2 {
					oneMsg.Text = oneMsg.Text + "`"
					text = "`" + text
				}
				msgs = append(msgs, oneMsg)
				continue
			}
			msg.Text = text
			msgs = append(msgs, msg)
			break
		}

		for _, oneMsg := range msgs {
			_, err = b.bot.Send(oneMsg)
			if err != nil {
				log.Printf("error: %v", err)
				log.Println("text:", oneMsg.Text)
			}
		}
	}
	return nil
}

func (b *Bot) act(msgText string) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("recovered from ", r)
		}
	}()

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

	if strings.HasPrefix(msgText, "/help") {
		return b.cfg.Stg.HelpText, nil
	}
	return msgText, nil
}

func (b *Bot) chooseAction(serviceTag string, text string) (string, error) {
	if text == "" {
		return text, nil
	}
	switch serviceTag {
	case b.cfg.Stg.WikiStg.Tag:
		return b.takeAction(b.cfg.Stg.WikiStg.Tag, text)
	case b.cfg.Stg.NewtonStg.Tag:
		return b.takeAction(b.cfg.Stg.NewtonStg.Tag, text)
	case b.cfg.Stg.PlaygroundStg.Tag:
		return b.takeAction(b.cfg.Stg.PlaygroundStg.Tag, text)
	case b.cfg.Stg.MarkovStg.Tag:
		return b.takeAction(b.cfg.Stg.MarkovStg.Tag, text)
	case b.cfg.Stg.GolangStg.Tag:
		return b.takeAction(b.cfg.Stg.GolangStg.Tag, text)
	}
	return text, nil
}

func (b *Bot) takeAction(tag string, text string) (string, error) {
	if b.sm.Services[tag].IsReady() {
		result, err := b.sm.Services[tag].Query(text)
		if err != nil {
			return "", err
		}
		return result, nil
	}
	return "Not ready yet. Please, waiting some minutes", nil
}
