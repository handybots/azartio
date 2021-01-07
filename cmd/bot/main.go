package main

import (
	"log"
	"os"

	"github.com/demget/clickrus"
	"github.com/demget/don"
	"github.com/sirupsen/logrus"

	"github.com/handybots/azartio/handler"
	"github.com/handybots/azartio/storage"

	tele "gopkg.in/tucnak/telebot.v3"
	"gopkg.in/tucnak/telebot.v3/layout"
	"gopkg.in/tucnak/telebot.v3/middleware"
)

func main() {
	layout.AddFunc("increment", increment)

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
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	dons, err := don.Parse("don.yml")
	if err != nil {
		log.Fatal(err)
	}

	ch, err := clickrus.NewHook(clickHouseConfig)
	if err != nil {
		log.Println(err)
	} else {
		logger.AddHook(ch)
	}

	h := handler.New(handler.Handler{
		Layout: lt,
		Bot:    b,
		DB:     db,
		Dons:   dons,
	})

	// Middleware
	b.OnError = h.OnError
	b.Use(middleware.Logger(logger, h.LoggerFields))
	b.Use(lt.Middleware("ru", h.LocaleFunc))

	// Handlers
	b.Handle(tele.OnText, h.OnText)
	b.Handle(tele.OnPinned, h.OnPinned)
	b.Handle("/start", h.OnStart)
	b.Handle(lt.Callback("red_bet"), h.OnBet, h.Validate())
	b.Handle(lt.Callback("green_bet"), h.OnBet, h.Validate())
	b.Handle(lt.Callback("black_bet"), h.OnBet, h.Validate())
	b.Handle(lt.Callback("roulette"), h.OnRoulette, h.Validate())
	b.Handle(lt.Callback("play"), h.OnRoulette, h.Validate())
	b.Handle(lt.Callback("reply_balance"), h.OnBalance, h.Validate())
	b.Handle(lt.Callback("perks"), h.OnPerks, h.Validate())
	b.Handle("/perks", h.OnPerks, h.Validate())
	b.Handle(lt.Callback("back_to_perks"), h.OnPerks, h.Validate())
	b.Handle(lt.Callback("perk"), h.OnPerk, h.Validate())
	b.Handle(lt.Callback("bonus"), h.OnBonuses, h.Validate())
	b.Handle(lt.Callback("leaderboard"), h.OnLeaderboard)
	b.Handle("/leaderboard", h.OnLeaderboard)
	b.Handle("/roulette", h.OnRoulette)
	b.Handle(lt.Callback("balance"), h.OnBalance, h.Validate())
	b.Handle("/bonus", h.OnBonus, h.Validate())
	b.Handle("б", h.OnBalance, h.Validate())
	b.Handle("Б", h.OnBalance, h.Validate())
	b.Handle("/go", h.OnGo, h.Validate())
	b.Handle(lt.Callback("roll"), h.OnGo, h.Validate())
	b.Handle(lt.Callback("participate"), h.OnParticipate, h.Validate())
	b.Handle("/contest", h.OnMakeContest, h.Validate())

	b.Start()
}

var clickHouseConfig = clickrus.Config{
	Addr:    os.Getenv("CLICKHOUSE_URL"),
	Columns: []string{"event", "user_id"},
	Table:   "azartio.logs",
}
