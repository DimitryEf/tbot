package main

import (
	"tbot/config"
	bot2 "tbot/internal/bot"
	"tbot/internal/errors"
	wiki2 "tbot/internal/wiki"
)

func main() {
	cfg, err := config.NewConfig()
	errors.PanicIfErr(err)

	wiki := wiki2.NewWiki(cfg.WikiUrl)

	bot, err := bot2.NewBot(cfg, wiki)
	errors.PanicIfErr(err)

	err = bot.Start()
	errors.PanicIfErr(err)
}
