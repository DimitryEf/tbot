package main

import (
	"flag"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"tbot/config"
	bot2 "tbot/internal/bot"
	"tbot/internal/errors"
	"tbot/services"
	golang2 "tbot/services/golang"

	//markov2 "tbot/services/markov"
	newton2 "tbot/services/newton"
	playground2 "tbot/services/playground"
	wiki2 "tbot/services/wiki"
)

var (
	defaultSettingsFile = "settings.yaml"
	settingsFile        = flag.String("settings", defaultSettingsFile, "path to settings file, if not it will use the default settings")
	defaultToken        = "" //TODO write a bot token here or in program arg, or in .env file, or in env variables
	token               = flag.String("token", defaultToken, "path to settings file")
)

func main() {
	flag.Parse()

	// load config
	cfg, err := config.NewConfig(*settingsFile, *token)
	errors.PanicIfErr(err)

	// load db
	db, err := gorm.Open(sqlite.Open("/home/dim/projects/tbot/databases/golang.db"), &gorm.Config{})
	errors.PanicIfErr(err)

	// load services
	wiki := wiki2.NewWiki(cfg.Stg.WikiStg)
	newton := newton2.NewNewton(cfg.Stg.NewtonStg)
	playground := playground2.NewPlayground(cfg.Stg.PlaygroundStg)
	//markov := markov2.NewMarkov(cfg.Stg.MarkovStg)
	golang := golang2.NewGolang(cfg.Stg.GolangStg, db)
	serviceManager := services.NewServiceManager(wiki, newton, playground, golang)

	// load bot
	bot, err := bot2.NewBot(cfg, serviceManager)
	errors.PanicIfErr(err)

	// start bot to listen chat messages
	err = bot.Start()
	errors.PanicIfErr(err)
}
