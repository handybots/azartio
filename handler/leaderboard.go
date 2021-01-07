package handler

import (
	"github.com/handybots/azartio/storage"
	"strconv"

	tele "gopkg.in/tucnak/telebot.v3"
)

type Leaderboard struct {
		Users []storage.User
		Chats map[int]*tele.Chat
	}

func (h handler) OnLeaderboard(c tele.Context) error {
	chats := make(map[int]*tele.Chat)

	users, err := h.db.Users.Leaderboard()
	if err != nil {
		return err
	}

	for _, user := range users {
		chat, err := h.b.ChatByID(strconv.Itoa(user.ID))
		if err != nil {
			return err
		}
		chats[user.ID] = chat
	}

	lb := Leaderboard{Users: users, Chats: chats}
	_, err = h.b.Send(c.Chat(), h.lt.Text(c, "leaderboard", lb))
	return nil
}
