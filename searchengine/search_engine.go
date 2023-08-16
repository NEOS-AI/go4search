package searchengine

import (
	"math"
	"sort"
	"strings"

	documents "go4search/documents"
)

type SearchEngine struct {
	Index        InvertedIndex
	Documents    []documents.Document
	AvgDocLength float64
	K1           float64
	B            float64
}

/**
 * Calculate the TF-IDF score for each document.
 * Iterate all the tokens (extracted from user input query), and calculate the TF-IDF score for each document.
 *
 * @param tokens A slice of tokens
 * @return map[int]float64
 */
func (se *SearchEngine) CalculateTFIDFScore(tokens []string) map[int]float64 {
	scores := make(map[int]float64)

	// iterate all tokens in the query
	for _, token := range tokens {
		if docSet, ok := se.Index[token]; ok {
			idf := math.Log(float64(len(se.Documents)) / float64(len(docSet)))

			// iterate all document that contains the token
			for _, docID := range docSet {
				tf := float64(strings.Count(strings.ToLower(se.Documents[docID].Content), token))
				scores[docID] += tf * idf
			}
		}
	}

	return scores
}

/**
 * Calculate the BM25 score for each document.
 * Iterate all the tokens (extracted from user input query), and calculate the BM25 score for each document.
 *
 * @param tokens A slice of tokens
 * @return map[int]float64
 */
func (se *SearchEngine) CalculateBM25Score(tokens []string) map[int]float64 {
	scores := make(map[int]float64)

	// iterate all tokens in the query
	for _, token := range tokens {
		if docSet, ok := se.Index[token]; ok {
			idf := math.Log(float64(len(se.Documents)-len(docSet))+0.5) / (float64(len(docSet)) + 0.5)

			// iterate all document that contains the token
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
	// remove stopwords from the query
	cleaned_query := removeStopwords(query)

	// tokenize the query
	tokens := strings.Fields(strings.ToLower(cleaned_query))

	scores := se.CalculateTFIDFScore(tokens)
	// or, to use bm25 scoring algorithm:
	// scores_bm25 := se.CalculateBM25Score(tokens)

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

	// sort the results by score in descending order
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// return the at most top 10 results
	if len(results) > 10 {
		return results[:10]
	}
	return results
}
