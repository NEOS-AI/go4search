package searchengine

import (
	"math"
	"sort"
	"strings"

	documents "go4search/documents"
)

type InvertedIndex map[string][]int

type SearchEngine struct {
	Index        InvertedIndex
	Documents    []documents.Document
	AvgDocLength float64
	K1, B        float64
}

func BuildInvertedIndex(documents []documents.Document) InvertedIndex {
	index := make(InvertedIndex)

	for _, doc := range documents {
		tokens := strings.Fields(strings.ToLower(doc.Content))

		for _, token := range tokens {
			if _, ok := index[token]; !ok {
				index[token] = make([]int, 0)
			}
			index[token] = append(index[token], doc.ID)
		}
	}

	return index
}

func (se *SearchEngine) CalculateTFIDFScore(tokens []string) map[int]float64 {
	scores := make(map[int]float64)

	for _, token := range tokens {
		if docSet, ok := se.Index[token]; ok {
			idf := math.Log(float64(len(se.Documents)) / float64(len(docSet)))
			for _, docID := range docSet {
				tf := float64(strings.Count(strings.ToLower(se.Documents[docID].Content), token))
				scores[docID] += tf * idf
			}
		}
	}

	return scores
}

func (se *SearchEngine) CalculateBM25Score(tokens []string) map[int]float64 {
	scores := make(map[int]float64)

	for _, token := range tokens {
		if docSet, ok := se.Index[token]; ok {
			idf := math.Log(float64(len(se.Documents)-len(docSet))+0.5) / (float64(len(docSet)) + 0.5)
			for _, docID := range docSet {
				tf := float64(strings.Count(strings.ToLower(se.Documents[docID].Content), token))
				dl := float64(len(strings.Fields(strings.ToLower(se.Documents[docID].Content))))
				numerator := (se.K1 + 1) * tf * (se.K1 + 1) / (tf + se.K1*(1.0-se.B+se.B*dl/se.AvgDocLength))
				denominator := tf + se.K1*(1.0-se.B+se.B*dl/se.AvgDocLength)
				scores[docID] += idf * numerator / denominator
			}
		}
	}

	return scores
}

func (se *SearchEngine) Search(query string) []documents.Document {
	tokens := strings.Fields(strings.ToLower(query))
	scores := se.CalculateTFIDFScore(tokens)
	// or, to use bm25 scoring algorithm:
	// scores := se.CalculateBM25Score(tokens)

	var results []documents.Document
	for docID, score := range scores {
		results = append(
			results,
			documents.Document{
				ID:      docID,
				Content: se.Documents[docID].Content,
				Score:   score,
			},
		)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	if len(results) > 10 {
		results = results[:10]
	}
	return results
}
