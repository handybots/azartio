package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/demget/clickrus"
	"github.com/demget/don"
	"github.com/sirupsen/logrus"

	"github.com/handybots/azartio/handler"
	"github.com/handybots/azartio/storage"
	"github.com/handybots/store"

	tele "gopkg.in/tucnak/telebot.v3"
	"gopkg.in/tucnak/telebot.v3/layout"
	"gopkg.in/tucnak/telebot.v3/middleware"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
	b.Use(lt.Middleware("ru"))
	b.Use(h.Validate)
	b.Use(h.ApplyBonuses)

	// Handlers
	//b.Handle(tele.OnText, h.OnText)
	b.Handle(tele.OnPinned, h.OnPinned)
	b.Handle("/start", h.OnStart)

	// Menu
	b.Handle(lt.Callback("play"), h.OnRoulette)
	b.Handle(lt.Callback("roll"), h.OnRouletteGo)
	b.Handle(lt.Callback("stats"), h.OnStats)
	b.Handle(lt.Callback("leaderboard"), h.OnLeaderboard)
	b.Handle(lt.Callback("perks"), h.OnPerks)
	b.Handle(lt.Callback("bonuses"), h.OnBonuses)

	// Aliases
	b.Handle("/roll", h.OnRoulette)
	b.Handle("/go", h.OnRouletteGo)
	b.Handle("/balance", h.OnStats)
	b.Handle("/leaderboard", h.OnLeaderboard)
	b.Handle("/perks", h.OnPerks)
	b.Handle("/bonus", h.OnBonuses)

	// Game
	b.Handle(lt.Callback("bet_r"), h.OnRouletteBet)
	b.Handle(lt.Callback("bet_g"), h.OnRouletteBet)
	b.Handle(lt.Callback("bet_b"), h.OnRouletteBet)

	// Balance
	b.Handle(lt.Callback("deposit"), h.OnDeposit)

	// Perks
	b.Handle(lt.Callback("perk"), h.OnPerk)
	b.Handle(lt.Callback("perks_back"), h.OnPerks)

	// Bonuses
	b.Handle(lt.Callback("bonus_daily"), h.OnBonusDaily)
	b.Handle(lt.Callback("bonus_sponsor"), h.OnBonusSponsor)

	// Admin
	b.Handle("/_balance", h.AdminBalance)
	b.Handle("/_perk", h.AdminPerk)

	// Payment receivers
	go store.Listen(lt.Int("store_id"), h.OnPayment)

	b.Start()
}

var clickHouseConfig = clickrus.Config{
	Addr:    os.Getenv("CLICKHOUSE_URL"),
	Columns: []string{"event", "user_id", "chat_id"},
	Table:   "azartio.logs",
}
