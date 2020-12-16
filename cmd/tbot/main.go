package main

import (
	"flag"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/http"
	"os"
	"tbot/config"
	bot2 "tbot/internal/bot"
	"tbot/internal/errors"
	"tbot/services"
	golang2 "tbot/services/golang"
	"time"

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
	//db, err := gorm.Open(sqlite.Open("databases/golang.db"), &gorm.Config{
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second,   // Slow SQL threshold
				LogLevel:      logger.Silent, // Log level
				Colorful:      false,         // Disable color
			},
		),
	})
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

	// heroku
	addr, err := determineListenAddress()
	errors.PanicIfErr(err)
	http.Handle("/", http.FileServer(http.Dir("./databases")))
	http.HandleFunc("/hello", hello)
	//go func() {
	log.Println("port", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
	//}()

	// start bot to listen chat messages
	err = bot.Start()
	errors.PanicIfErr(err)
}

func determineListenAddress() (string, error) {
	//return ":8081", nil
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World")
}
