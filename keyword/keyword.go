package keyword

import "strings"

type document interface {
	GetTermFrequency(term string) int
}

type Extractor interface {
	GetKeyWords(string, []document) []string
}

func analyzeCorpus(corpus string) (wordCountMap map[string]int, totalWordsCount int) {
	wordCountMap = map[string]int{}
	words := strings.Split(corpus, " ")
	for _, word := range words {
		wordCountMap[word] += 1
	}
	totalWordsCount = len(words)
	return wordCountMap, totalWordsCount
}
