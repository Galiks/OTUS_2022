package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type pair struct {
	Word  string
	Count int64
}

func Top10(text string) []string {
	var (
		result            = make([]string, 0)
		storage           = make(map[string]int64)
		pairStorage       = make([]*pair, 0)
		topWord     int64 = 10
	)

	words := strings.Fields(text)
	if len(words) == 1 {
		return []string{words[0]}
	}
	for _, word := range words {
		if len(word) == 1 {
			word = strings.ToLower(word)
		}
		storage[word]++
	}
	for k, v := range storage {
		pairStorage = append(pairStorage, &pair{
			Word:  k,
			Count: v,
		})
	}

	sort.Slice(pairStorage, func(i, j int) bool {
		if pairStorage[i].Count == pairStorage[j].Count {
			return pairStorage[i].Word < pairStorage[j].Word
		}
		return pairStorage[i].Count > pairStorage[j].Count
	})

	if len(pairStorage) < 10 {
		topWord = int64(len(pairStorage))
	}

	for i := 0; i < int(topWord); i++ {
		result = append(result, pairStorage[i].Word)
	}

	return result
}
