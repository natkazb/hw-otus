package hw02unpackstring

import (
	"errors"
	"strings"

	"github.com/rivo/uniseg" //nolint
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var result strings.Builder
	var prevStr string
	var prevRunes []rune
	isBeginShield := false

	gr := uniseg.NewGraphemes(s)
	for gr.Next() {
		currentStr := gr.Str()
		currentRunes := gr.Runes()

		// Handle numbers
		if isNumber(currentRunes) && !isBeginShield {
			if prevStr == "" || isNumber(prevRunes) {
				return "", ErrInvalidString
			}

			count := int(currentRunes[0] - '0')
			if count > 0 {
				result.WriteString(strings.Repeat(prevStr, count-1))
			}
		} else {
			// Check for next
			_, to := gr.Positions()
			gr2 := uniseg.NewGraphemes(s[to:])
			okNext := gr2.Next()
			if isShield(currentRunes) { // \
				isBeginShield = true
				if isShield(prevRunes) {
					result.WriteString(currentStr)
				}
				if okNext && !isNumber(gr2.Runes()) && !isShield(gr2.Runes()) {
					return "", ErrInvalidString
				}
			} else if !okNext || okNext && !isZero(gr2.Runes()) { // zero
				result.WriteString(currentStr)
				if isBeginShield && !isNumber(gr2.Runes()) {
					isBeginShield = false
				}
			}
		}
		prevStr = currentStr
		prevRunes = gr.Runes()
	}

	return result.String(), nil
}

func isNumber(r []rune) bool {
	return len(r) == 1 && r[0] >= 48 && r[0] <= 57
}

func isZero(r []rune) bool {
	return len(r) == 1 && r[0] == 48
}

func isShield(r []rune) bool {
	return len(r) == 1 && r[0] == 92
}
