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
		result      []string         = make([]string, 0)
		storage     map[string]int64 = make(map[string]int64)
		pairStorage []*pair          = make([]*pair, 0)
		topWord     int64            = 10
	)

	words := strings.Fields(text)
	for _, word := range words {
		if len(word) == 1 {
			word = strings.ToLower(word)
		}
		storage[word] = storage[word] + 1

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
		} else {
			return pairStorage[i].Count > pairStorage[j].Count
		}
	})

	if len(pairStorage) < 10 {
		topWord = int64(len(pairStorage))
	}

	for i := 0; i < int(topWord); i++ {
		result = append(result, pairStorage[i].Word)
	}

	return result
}
