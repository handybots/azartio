package handler

import (
	"github.com/demget/don"
	"github.com/handybots/azartio/azartio"
	"github.com/handybots/azartio/storage"

	tele "gopkg.in/tucnak/telebot.v3"
	"gopkg.in/tucnak/telebot.v3/layout"
)

func New(c Handler) handler {
	return handler{
		lt: c.Layout,
		b:  c.Bot,
		db: c.DB,
		c:  c.Casino,
		d:  c.Dons,
	}
}

type (
	Handler struct {
		Layout *layout.Layout
		Bot    *tele.Bot
		DB     *storage.DB
		Casino *azartio.Casino
		Dons   *don.Dons
	}
	handler struct {
		lt *layout.Layout
		b  *tele.Bot
		db *storage.DB
		c  *azartio.Casino
		d  *don.Dons
	}
)