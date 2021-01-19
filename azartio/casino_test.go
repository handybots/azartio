package azartio

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCasino_RollMany(t *testing.T) {
	casino := &Casino{}
	for i := 0; i <= 1000000; i++ {
		sign := []string{Red, Clever, Black}[rand.Intn(3)]
		amount := rand.Int63n(100000000000000)
		userID := rand.Intn(10000000000)
		bet := NewBet(sign, amount, userID)
		results, err := casino.RollMany([]*Bet{bet})
		if err != nil {
			t.Error(err)
		}
		for _, result := range results {
			if result.Won {
				if sign != Clever {
					assert.Equal(t, amount*2, result.Amount)
				} else {
					assert.Equal(t, amount*12, result.Amount)
				}
			} else {
				assert.Equal(t, amount, result.Amount)
			}
			assert.Equal(t, bet, result.Bet)
		}
	}
}
