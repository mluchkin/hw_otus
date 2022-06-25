package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type counter struct {
	Word  string
	Count int
}

func Top10(text string) []string {
	maxWords := 10
	words := make(map[string]int)
	for _, word := range strings.Fields(text) {
		words[word]++
	}

	slice := make([]counter, 0, len(words))
	for word, count := range words {
		slice = append(slice, counter{word, count})
	}

	sort.SliceStable(slice, func(i, j int) bool {
		return slice[i].Count > slice[j].Count || (slice[i].Count == slice[j].Count && slice[i].Word < slice[j].Word)
	})

	if len(slice) < 10 {
		maxWords = len(slice)
	}

	result := make([]string, 0, maxWords)

	for _, counter := range slice[:maxWords] {
		result = append(result, counter.Word)
	}

	return result
}
