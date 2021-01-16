package azartio

import "strings"

// Translates s and returns bet sign
// Example:
//	 TranslateBet("красное") returns Red
//	 TranslateBet("к") returns Red
//
// if can not translate returns an empty string
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