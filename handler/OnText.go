package handler

import (
	"regexp"

	tele "gopkg.in/tucnak/telebot.v3"
)

var betRegexp = regexp.MustCompile(`(\d+)\s*(?:на\s*)?(к(?:расное)?|ч(?:[её]рное)?|з(?:ел[её]ное)?)$`)

func (h *handler) OnText(c tele.Context) error {
	if c.Message().Text == h.lt.Text(c, "my_balance") {
		return h.OnBalance(c)
	} else if betRegexp.MatchString(c.Message().Text) {
		return h.OnBet(c)
	}
	return nil
}
