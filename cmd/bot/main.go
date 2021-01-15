package main

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/handybots/azartio/handler"
	"github.com/handybots/azartio/storage"
	tele "gopkg.in/tucnak/telebot.v3"
	"gopkg.in/tucnak/telebot.v3/layout"
	"gopkg.in/tucnak/telebot.v3/middleware"
)

func main() {
	lt, err := layout.New("bot.yml")
	if err != nil {
		log.Fatal(err)
	}


	b, err := tele.NewBot(lt.Settings())
	if err != nil {
		log.Fatal(err)
	}

	db, err := storage.Open(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	//ch, err := clickrus.NewHook(clickHouseConfig)
	//if err != nil {
	//	log.Fatal(err)
	//}

	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	//logger.AddHook(ch)

	h := handler.New(handler.Handler{
		Layout: lt,
		Bot:    b,
		DB:     db,
	})

	// Middleware
	b.OnError = h.OnError
	b.Use(middleware.Logger(logger, h.LoggerFields))
	b.Use(lt.Middleware("ru", h.LocaleFunc))
	b.Use()
	// Handlers
	b.Handle("/start", h.OnStart)
	b.Handle(lt.Callback("lang"), h.OnLang)
	b.Handle(lt.Callback("bet"), h.OnBet, h.Validate())

	b.Handle("/menu", h.OnMenu)
	b.Handle(lt.Callback("balance"), h.OnBalance, h.Validate())

	b.Handle("б",  h.OnBalance, h.Validate())
	b.Handle("Б",  h.OnBalance, h.Validate())

	b.Start()
}

/*var clickHouseConfig = clickrus.Config{
	Addr:    os.Getenv("CLICKHOUSE_URL"),
	Columns: []string{"event", "user_id"},
	Table:   "bot.logs",
}*/
