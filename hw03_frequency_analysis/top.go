package hw03frequencyanalysis

import (
	"sort"
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

	sort.Slice(stats, func(i, j int) bool {
		if stats[i].count != stats[j].count {
			return stats[i].count > stats[j].count // По убыванию частоты
		}
		return stats[i].word < stats[j].word // По алфавиту при равенстве
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
