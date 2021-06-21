package azartio

import (
	"errors"
	"math/rand"
)

// Roulette is responsible for basic roulette game.
type Roulette struct{}

type Bet struct {
	Sign   string
	Amount int64
	UserID int
}

type RollResult struct {
	Amount int64
	Won    bool
	Bet    Bet
	Sign   string
}

func NewBet(sign string, amount int64, userID int) *Bet {
	return &Bet{Sign: sign, Amount: amount, UserID: userID}
}

// doRoll does roll for bet depending on the n number.
//
// if n <= 45: wins red
// if n > 45 and n < 90: wins black
// if n > 90: wins clever
//
func (c *Roulette) doRoll(bet Bet, n int) (*RollResult, error) {
	amount := bet.Amount

	if amount <= 0 {
		return nil, errors.New("azartio: bet amount must be > 0")
	}

	if _, ok := Colors[bet.Sign]; !ok {
		return nil, errors.New("azartio: unknown sign")
	}

	sign, err := c.pickSign(n)
	if err != nil {
		return nil, err
	}

	if sign == bet.Sign {
		if bet.Sign == Clever {
			amount = amount * 8
		} else {
			amount = amount * 2
		}
	}

	return &RollResult{
		Amount: amount,
		Won:    sign == bet.Sign,
		Bet:    bet,
		Sign:   sign,
	}, nil
}

// pickSign picks a sign depending on the n number.
func (c *Roulette) pickSign(n int) (sign string, _ error) {
	switch {
	case n <= 45:
		sign = Red
	case n > 45 && n <= 90:
		sign = Black
	case n > 90:
		sign = Clever
	default:
		return "", errors.New("azartio: bad n pick number")
	}

	return sign, nil
}

func (c *Roulette) RollMany(bets []Bet) (result []*RollResult, _ error) {
	n := rand.Intn(100)
	for _, bet := range bets {
		r, err := c.doRoll(bet, n)
		if err != nil {
			continue
		}
		result = append(result, r)
	}

	if len(result) == 0 {
		return nil, errors.New("azartio: result is empty")
	}

	return result, nil
}

func (c *Roulette) Roll(bet Bet) (*RollResult, error) {
	n := rand.Intn(100)
	return c.doRoll(bet, n)
}
