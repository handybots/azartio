package main

import (
	"database/sql"
	"log"
	"time"

	tele "gopkg.in/tucnak/telebot.v3"
)

func voteWatcher() {
	for {
		if err := updateVote(); err != nil {
			if err != tele.ErrTrueResult && err != sql.ErrNoRows {
				log.Println("voteWatcher:", err)
			}
		}
		time.Sleep(time.Hour)
	}
}

func updateVote() error {
	vote, err := db.LastVote()
	if err != nil {
		return err
	}
	if time.Now().Sub(vote.UpdatedAt).Hours() < 24 {
		return nil
	}
	if err := vote.SetDaysLeft(vote.DaysLeft - 1); err != nil {
		return err
	}

	log.Println("Vote", vote.ID, vote.DaysLeft, "days left")
	if err := vote.SetUpdatedAt(time.Now()); err != nil {
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

	lastIdea, ideas := ideas[vote.DaysLeft], ideas[:vote.DaysLeft]
	if err := lastIdea.SetUsed(false); err != nil {
		return err
	}

	if vote.DaysLeft > 1 {
		_, err = b.Edit(
			&tele.Message{InlineID: vote.MessageID},
			lt.TextLocale("ideas", "vote", ideas),
			ideasMarkup(ideas),
		)
		return err
	}

	if err := vote.SetDone(true); err != nil {
		return err
	}

	winner := ideas[0]
	if err := winner.SetDeleted(true); err != nil {
		return err
	}

	_, err = b.Edit(
		&tele.Message{InlineID: vote.MessageID},
		lt.TextLocale("ideas", "done", winner),
	)
	return err
}
