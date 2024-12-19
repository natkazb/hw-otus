package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

const LIMIT = 10

type WordFrequency struct {
	word  string
	count int
}

var pattern = regexp.MustCompile(`[\p{L}][\p{L}0-9-]*|-{2,}`)

func Top10(s string) []string {
	words := pattern.FindAllString(s, -1)
	frequency := make(map[string]int)
	for _, v := range words {
		frequency[strings.ToLower(v)]++
	}

	pairs := make([]WordFrequency, 0, len(frequency))
	for word, count := range frequency {
		pairs = append(pairs, WordFrequency{word, count})
	}

	// Sort by frequency (desc) and alphabetically
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count == pairs[j].count {
			return pairs[i].word < pairs[j].word
		}
		return pairs[i].count > pairs[j].count
	})

	// Take LIMIT or less
	limit := min(LIMIT, len(pairs))
	result := make([]string, limit)
	for i := 0; i < limit; i++ {
		result[i] = pairs[i].word
	}

	return result
}
