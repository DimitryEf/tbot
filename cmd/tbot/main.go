package main

import (
	"tbot/config"
	bot2 "tbot/internal/bot"
	"tbot/internal/errors"
)

func main() {
	cfg, err := config.NewConfig()
	errors.PanicIfErr(err)

	bot, err := bot2.NewBot(cfg)
	errors.PanicIfErr(err)

	err = bot.Start()
	errors.PanicIfErr(err)
}
