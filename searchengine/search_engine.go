package searchengine

import (
	"math"
	"sort"
	"strings"

	documents "go4search/documents"
	nlp "go4search/nlp"
)

type SearchEngine struct {
	Index         InvertedIndex
	Documents     []documents.Document
	TotalDocCount float64
	TotalDocLen   float64
	AvgDocLength  float64
	K1            float64
	B             float64
}

const SCORE_THRESHOLD = 0.5

const BM25_WEIGHT = 0.5
const TFIDF_WEIGHT = 0.5

func (se *SearchEngine) SetK1(k1 float64) {
	se.K1 = k1
}

func (se *SearchEngine) SetB(b float64) {
	se.B = b
}

func (se *SearchEngine) AddNewDocument(doc documents.Document) {
	se.Documents = append(se.Documents, doc)
	count := se.TotalDocCount
	docLength := se.TotalDocLen

	currentDocLength := float64(len(doc.Content))
	// check if docLength + currentDocLength is too large to avoid overflow
	if docLength+currentDocLength > math.MaxFloat64-100 {
		return
	}

	// increase the docLength
	docLength += currentDocLength
	count++
	count_f := float64(count)

	// update the search engine
	se.TotalDocCount = count_f
	se.TotalDocLen = docLength
	se.AvgDocLength = docLength / count_f
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

				// TF-IDF score * weight
				scores[docID] += tf * idf * TFIDF_WEIGHT
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

				// BM25 score
				score := idf * numerator / denominator
				// apply weight to the score
				scores[docID] += score * BM25_WEIGHT
			}
		}
	}

	return scores
}

func (se *SearchEngine) Search(query string, limit int) []documents.Document {
	// remove stopwords from the query
	cleaned_query := removeStopwords(query)

	// tokenize the query
	// tokens := strings.Fields(strings.ToLower(cleaned_query))
	tokens := nlp.Tokenize_Query(strings.ToLower(cleaned_query))

	// ranking with TF-IDF
	scores := se.CalculateTFIDFScore(tokens)
	// ranking with BM25
	scores_bm25 := se.CalculateBM25Score(tokens)

	// combine the scores from TF-IDF and BM25 for weighted ranking
	for docID, score := range scores_bm25 {
		scores[docID] += score
	}

	var results []documents.Document
	for docID, score := range scores {
		// filter out the results with score less than SCORE_THRESHOLD
		if score < SCORE_THRESHOLD {
			continue
		}

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

	// return the at most top N results
	if len(results) > limit {
		return results[:limit]
	}
	return results
}
