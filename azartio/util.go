package azartio

import "strings"

func TranslateBet(s string) string{
	switch strings.ToLower(s) {
	case "к":
		return Red
	case "красное":
		return Red
	case "ч":
		return Black
	case "чёрное":
		return Black
	case "черное": // для дебилов
		return Black
	case "з":
		return Clever
	case "зелёное":
		return Clever
	case "зеленое": // для дебилов
		return Clever
	default:
		return ""
	}
}