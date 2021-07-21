package main

import (
	"database/sql"
	"github.com/handybots/azartio/cmd/ideas/database"
	tele "gopkg.in/tucnak/telebot.v3"
	"gopkg.in/tucnak/telebot.v3/layout"
	"gopkg.in/tucnak/telebot.v3/middleware"
	"log"
	"os"
	"strconv"
	"strings"
)

var admins = []tele.Recipient{
	&tele.User{ID: 1360834297},
	&tele.User{ID: 384327785},
}

var (
	lt *layout.Layout
	b  *tele.Bot
	db *database.DB
)

func init() {
	layout.AddFunc("inc", func(i int) int { return i + 1 })
}

func main() {
	var err error

	lt, err = layout.New("ideas.yml")
	if err != nil {
		log.Fatal(err)
	}

	b, err = tele.NewBot(lt.Settings())
	if err != nil {
		log.Fatal(err)
	}

	db, err = database.Open(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	b.OnError = func(err error, _ tele.Context) {
		if err != tele.ErrTrueResult {
			log.Println(err)
		}
	}

	b.Use(lt.Middleware("ideas"))
	b.Use(middleware.Whitelist(admins...))

	b.Handle("/idea", onIdea)
	b.Handle("/clear", onClear)
	b.Handle(lt.Callback("vote"), onVote)
	b.Handle(tele.OnQuery, onQuery)
	b.Handle(tele.OnChosenInlineResult, onInlineResult)

	go voteWatcher()

	b.Start()
}

func onIdea(c tele.Context) error {
	args := strings.Split(c.Message().Payload, " ")
	desc := strings.Split(c.Text(), "\n")[1:]

	idea := database.Idea{
		Emoji:       args[0],
		Title:       strings.Join(args[1:], " "),
		Description: strings.Join(desc, " "),
	}
	if _, err := db.InsertIdea(idea); err != nil {
		return err
	}

	return c.Send("Идея успешно добавлена")
}

func onClear(c tele.Context) error {
	ideas, err := db.IdeasByUsed(true)
	if err != nil {
		return err
	}

	for _, idea := range ideas {
		if err := idea.SetUsed(false); err != nil {
			return err
		}
	}

	return c.Send("Последнее голосование очищено")
}

func onVote(c tele.Context) error {
	defer c.Respond(&tele.CallbackResponse{
		Text: "Твой голос учтён, спасибо!",
	})

	userID := int64(c.Sender().ID)
	ideaID, _ := strconv.Atoi(c.Data())

	vote, err := db.LastVote()
	if err != nil {
		return err
	}

	voter, err := db.Voter(vote.ID, userID)
	if err == sql.ErrNoRows {
		_, err = db.InsertVoter(database.Voter{
			UserID: userID,
			VoteID: vote.ID,
			IdeaID: ideaID,
		})
	} else if err == nil {
		voter.IdeaID = ideaID
		_ = db.DeleteVoter(userID)
		_, err = db.InsertVoter(voter)
	}
	if err != nil {
		return err
	}

	voters, err := db.VotersByVoteID(vote.ID)
	if err != nil {
		return err
	}

	ideas, err := fetchIdeas(voters, vote.Ideas)
	if err != nil {
		return err
	}

	ideas = ideas[:vote.DaysLeft]
	return c.Edit(
		lt.Text(c, "vote", ideas),
		ideasMarkup(ideas),
	)
}

func onQuery(c tele.Context) error {
	if c.Data() == "" {
		return nil
	}

	res := &tele.ArticleResult{Title: "Отправить пост"}
	res.SetReplyMarkup(lt.Markup(c, "loading").InlineKeyboard)
	res.SetContent(&tele.InputTextMessageContent{
		Text:           lt.Text(c, "vote"),
		ParseMode:      tele.ModeHTML,
		DisablePreview: true,
	})

	return c.Answer(&tele.QueryResponse{
		Results:   tele.Results{res},
		CacheTime: -1,
	})
}

func onInlineResult(c tele.Context) error {
	messageID := c.ChosenInlineResult().MessageID

	days, err := strconv.Atoi(c.Data())
	if err != nil {
		return nil
	}

	free, err := db.IdeasByDeleted(false)
	if err != nil {
		return err
	}
	if len(free) < days {
		days = len(free)
	}

	var (
		ideas   []Idea
		ideaIDs []int32
	)
	for _, idea := range free[:days] {
		if err := idea.SetUsed(true); err != nil {
			return err
		}

		ideas = append(ideas, Idea{Idea: idea})
		ideaIDs = append(ideaIDs, int32(idea.ID))
	}

	_, err = db.InsertVote(database.Vote{
		DaysLeft:  days,
		Ideas:     ideaIDs,
		MessageID: messageID,
	})
	if err != nil {
		return err
	}

	_, err = b.Edit(
		&tele.Message{InlineID: messageID},
		lt.Text(c, "vote", ideas),
		ideasMarkup(ideas),
	)
	return err
}

func ideasMarkup(ideas []Idea) *tele.ReplyMarkup {
	var total int
	for _, idea := range ideas {
		total += idea.Votes
	}

	var btns []tele.Btn
	for _, idea := range ideas {
		if total == 0 {
			idea.Percentage = 0
		} else {
			idea.Percentage = int(float64(idea.Votes) / float64(total) * 100)
		}

		btns = append(btns, *lt.ButtonLocale("ideas", "vote", idea))
	}

	markup := b.NewMarkup()
	markup.Inline(markup.Split(1, btns...)...)
	return markup
}
