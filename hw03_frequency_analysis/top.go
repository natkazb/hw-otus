package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

const LIMIT = 10

type WordFrequency struct {
	word  string
	count int
}

func Top10(s string) []string {
	words := strings.Fields(s)
	frequency := make(map[string]int)
	for _, v := range words {
		frequency[v]++
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
