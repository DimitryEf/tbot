package bot

import (
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

		for _, val := range b.cfg.Stg.Services["golang"] {
			if strings.HasPrefix(msg.Text, val) {
				msg.ParseMode = tgbotapi.ModeMarkdownV2
				break
			}
		}

		//oneMsg.Text = "\"Toggle Off\\\\n\""
		//b.bot.Send(oneMsg)
		//continue

		resp, err := b.sm.Act(msg.Text)
		if err != nil {
			log.Printf("error: %v", err)
			//resp.Text = fmt.Sprintf("error: %v", err)
		}

		if resp.Text == "" {
			resp.Text = "не получилось"
		}

		msg.ParseMode = resp.ParseMod

		msgs := pagination(resp.Text, msg)

		for i, oneMsg := range msgs {
			if i == len(msgs)-1 {
				oneMsg.Text += "\n\n[if I fall asleep wake me up](https://myapp20200522.herokuapp.com/hello)"
			}
			_, err = b.bot.Send(oneMsg)
			if err != nil {
				log.Printf("error: %v", err)
				log.Println("text:", oneMsg.Text)
			}
		}
	}
	return nil
}

func pagination(text string, msg tgbotapi.MessageConfig) []tgbotapi.MessageConfig {
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
	return msgs
}
