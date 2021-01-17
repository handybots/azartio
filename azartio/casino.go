package azartio

import (
	"errors"
	tele "gopkg.in/tucnak/telebot.v3"
	"math/rand"
	"time"
)

func init(){
	rand.Seed(time.Now().UnixNano())
}


type Casino struct{

}

// doRoll makes roll for bet depends on n
// if n <= 45: wins red
// if n > 45 and n < 90: wins black
// if n > 90: wins clever
func (c *Casino) doRoll(bet *Bet, n int) (*RollResult, error){
	if bet.Amount <= 0 {
		return nil, errors.New("casino: bet amount must be > 0")
	}
	amount := bet.Amount
	if _, ok := Colors[bet.Sign]; !ok{
		return nil, errors.New("casino: unknown sign")
	}
	wonSign, err := c.pickSign(n)
	if err != nil{
		return nil, err
	}
	if wonSign == bet.Sign {
		if bet.Sign == Clever{
			amount = amount * 12
		}else{
			amount = amount * 2
		}

	}
	return &RollResult{
		Amount: amount,
		Won: wonSign == bet.Sign,
		Bet: *bet,
	}, nil
}

// pickSign picks a sign depends on n
func (c *Casino) pickSign(n int) (wonSign string, _ error) { // for test
	if n < 0 || n > 100{
		return "", errors.New("casino: pickSign argument must be < 100 and > 0")
	}
	switch{
	case n <= 45:
		wonSign = Red
	case n > 45 && n <= 90:
		wonSign = Black
	case n > 90:
		wonSign = Clever
	}
	return wonSign, nil
}



func (c *Casino) RollMany(bets []*Bet) (result []*RollResult, _ error){
	n := rand.Intn(100)
	for _, bet := range bets{
		r, err := c.doRoll(bet, n)
		if err != nil{
			continue
		}
		result = append(result, r)
	}
	if len(result) == 0 {
		return nil, errors.New("casino: result is empty")
	}
	return result, nil
}

func (c *Casino) Roll(bet *Bet) (*RollResult, error){
	n := rand.Intn(100)
	return c.doRoll(bet, n)
}

// Bet :)
// the identifier is needed to pass it to the RollResult and not lose whose it is
type Bet struct {
	Sign string
	Amount int64
	User *tele.User
}

func NewBet(sign string, amount int64, user *tele.User) *Bet {
	return &Bet{Sign: sign, Amount: amount, User: user}
}



// RollResult represents a result of Casino.Roll
type RollResult struct {
	Amount int64
    Won bool
	Bet Bet
}

// Signs
const (
	Clever = "g"
	Red = "r"
	Black = "b" // nigger
)

var Colors = map[string]string{
	Clever: "🍀",
	Red: "🔴",
	Black: "⚫️",
}