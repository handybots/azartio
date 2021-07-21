package main

import (
	"sort"

	"github.com/handybots/azartio/cmd/ideas/database"
)

type Idea struct {
	database.Idea
	Votes      int
	Percentage int
}

func fetchIdeas(voters []database.Voter, ids []int32) ([]Idea, error) {
	ideaVotes := make(map[int]int)
	for _, voter := range voters {
		ideaVotes[voter.IdeaID] += 1
	}

	var ideas []Idea
	for _, id := range ids {
		id := int(id)

		idea, err := db.IdeaByID(id)
		if err != nil {
			return nil, err
		}

		ideas = append(ideas, Idea{
			Idea:  idea,
			Votes: ideaVotes[id],
		})
	}

	sort.Slice(ideas, func(i, j int) bool {
		return ideas[i].Votes > ideas[j].Votes
	})

	return ideas, nil
}
