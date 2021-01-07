package azartio

import (
	"errors"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Casino struct {
}

// doRoll makes roll for bet depending on picked n number.
// if n <= 45: wins red
// if n > 45 and n < 90: wins black
// if n > 90: wins clever
func (c *Casino) doRoll(bet *Bet, n int) (*RollResult, error) {
	amount := bet.Amount
	
	if amount <= 0 {
		return nil, errors.New("azartio: bet amount must be > 0")
	}

	if _, ok := Colors[bet.Sign]; !ok {
		return nil, errors.New("azartio: unknown sign")
	}

	wonSign, err := c.pickSign(n)
	if err != nil {
		return nil, err
	}

	if wonSign == bet.Sign {
		if bet.Sign == Clever {
			amount = amount * 12
		} else {
			amount = amount * 2
		}
	}
	
	return &RollResult{
		Amount:  amount,
		Won:     wonSign == bet.Sign,
		Bet:     *bet,
		WonSign: wonSign,
	}, nil
}

// pickSign picks a sign depending on picked n number.
func (c *Casino) pickSign(n int) (wonSign string, _ error) { // for test
	if n < 0 || n > 100 {
		return "", errors.New("azartio: pickSign argument must be < 100 and > 0")
	}
	
	switch {
	case n <= 45:
		wonSign = Red
	case n > 45 && n <= 90:
		wonSign = Black
	case n > 90:
		wonSign = Clever
	}
	
	return wonSign, nil
}

func (c *Casino) RollMany(bets []*Bet) (result []*RollResult, _ error) {
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

func (c *Casino) Roll(bet *Bet) (*RollResult, error) {
	n := rand.Intn(100)
	return c.doRoll(bet, n)
}

type Bet struct {
	Sign   string
	Amount int64
	UserID int
}

func NewBet(sign string, amount int64, userID int) *Bet {
	return &Bet{Sign: sign, Amount: amount, UserID: userID}
}

type RollResult struct {
	Amount  int64
	Won     bool
	Bet     Bet
	WonSign string
}

const (
	Clever = "g"
	Red    = "r"
	Black  = "b"
)

var Colors = map[string]string{
	Clever: "🍀",
	Red:    "🔴",
	Black:  "⚫️",
}
