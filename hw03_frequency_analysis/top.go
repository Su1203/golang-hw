package hw03frequencyanalysis

import (
	"slices"
	"strings"
)

type wordStat struct {
	word string

	count int
}

func Top10(s string) []string {
	result := []string{}

	words := strings.Fields(s)

	// посчитаем кол-во

	wordCounts := make(map[string]int)

	for _, word := range words {
		wordCounts[word]++
	}

	stats := make([]wordStat, 0, len(wordCounts))

	for word, count := range wordCounts {
		stats = append(stats, wordStat{word, count})
	}

	// сортировка

	slices.SortFunc(stats, func(a, b wordStat) int {
		if a.count != b.count {
			return b.count - a.count
		}

		if a.word < b.word {
			return -1
		}

		if a.word > b.word {
			return 1
		}

		return 0
	})

	coutRes := 10
	if len(stats) < coutRes {
		coutRes = len(stats)
	}

	for _, s := range stats[:coutRes] {
		result = append(result, s.word)
	}

	return result
}
