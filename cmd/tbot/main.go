package main

import (
	"flag"
	"tbot/config"
	bot2 "tbot/internal/bot"
	"tbot/internal/errors"
	wiki2 "tbot/internal/wiki"
)

var (
	defaultSettingsFile = "settings.yaml"
	settingsFile        = flag.String("settings", defaultSettingsFile, "the path to settings file, if not its will use default")
	defaultToken        = "" //TODO write bot token here or in program arg, or in .env file, or in env variables
	token               = flag.String("token", defaultToken, "the path to settings file")
)

func main() {
	flag.Parse()

	cfg, err := config.NewConfig(*settingsFile, *token)
	errors.PanicIfErr(err)

	wiki := wiki2.NewWiki(cfg.Stg.WikiStg)

	bot, err := bot2.NewBot(cfg, wiki)
	errors.PanicIfErr(err)

	err = bot.Start()
	errors.PanicIfErr(err)
}
