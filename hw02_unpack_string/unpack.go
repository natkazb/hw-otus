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

	gr := uniseg.NewGraphemes(s)
	for gr.Next() {
		currentStr := gr.Str()
		currentRunes := gr.Runes()

		// Handle numbers
		if isNumber(currentRunes) {
			if prevStr == "" || isNumber(prevRunes) {
				return "", ErrInvalidString
			}

			count := int(currentRunes[0] - '0')
			if count > 0 {
				result.WriteString(strings.Repeat(prevStr, count-1))
			}
		} else {
			// Check for zero
			_, to := gr.Positions()
			gr2 := uniseg.NewGraphemes(s[to:])
			okNext := gr2.Next()
			if !okNext || okNext && !isZero(gr2.Runes()) {
				result.WriteString(currentStr)
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
