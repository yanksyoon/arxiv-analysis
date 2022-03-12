package keyword

import (
	"math"
	"sort"
)

type TDIDFConfig struct {
	threshold float64
	topN      int
}

type TDIDFExtractor struct {
	cfg *TDIDFConfig

	wordCountMap  map[string]int
	totalNumWords int
}

var defaultTdIdfConfig = TDIDFConfig{
	threshold: 0,
	topN:      10,
}

func NewTdIdfExtractor(cfg *TDIDFConfig) *TDIDFExtractor {
	if cfg == nil {
		cfg = &defaultTdIdfConfig
	} else if cfg.threshold <= 0 && cfg.topN <= 0 {
		// log using default config
		cfg = &defaultTdIdfConfig
	}
	return &TDIDFExtractor{cfg: cfg}
}

func (e *TDIDFExtractor) GetKeyWords(corpus string, documentSet []document) []string {
	wordCountMap, totalNumWords := analyzeCorpus(corpus)
	e.wordCountMap = wordCountMap
	e.totalNumWords = totalNumWords
	termScores := map[string]float64{}

	for word := range wordCountMap {
		termScores[word] = e.getTfIdf(word, documentSet)
	}

	if e.cfg.threshold > 0 {
		termScores = filterByThreshold(termScores, e.cfg.threshold)
	}

	keywords := []string{}
	scores := []float64{}
	for term, score := range termScores {
		keywords = append(keywords, term)
		scores = append(scores, score)
	}

	sort.Slice(keywords, func(i, j int) bool {
		return scores[i] < scores[j]
	})

	return keywords[:e.cfg.topN]
}

func (e *TDIDFExtractor) getTfIdf(term string, documentSet []document) float64 {
	termFrequency := e.getTf(term)
	invDocFrequency := e.getIdf(term, documentSet)
	return termFrequency * invDocFrequency
}

func (e *TDIDFExtractor) getTf(term string) float64 {
	termFrequency := e.wordCountMap[term]
	return float64(termFrequency) / float64(e.totalNumWords)
}

func (e *TDIDFExtractor) getIdf(term string, documentSet []document) float64 {
	totalFrequency := 0
	for _, document := range documentSet {
		totalFrequency += document.GetTermFrequency(term)
	}
	return math.Log(float64(len(documentSet)) / float64(totalFrequency+1))
}

func filterByThreshold(scores map[string]float64, threshold float64) map[string]float64 {
	newScoresMap := map[string]float64{}
	for term, score := range scores {
		if score < threshold {
			continue
		}
		newScoresMap[term] = score
	}
	return newScoresMap
}
