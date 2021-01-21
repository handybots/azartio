package handler

import (
	"strconv"

	"github.com/handybots/azartio/storage"

	tele "gopkg.in/tucnak/telebot.v3"
)

// Not finished yet
// TODO: emoji
func (h *handler) OnLeaderboard(c tele.Context) error {
	type Leaderboard struct {
		Users []storage.User
		Chats map[int]*tele.Chat
	}
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
	h.b.Send(c.Chat(), h.lt.Text(c, "leaderboard", Leaderboard{Users: users, Chats: chats}))
	return nil
}
