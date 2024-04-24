package searchengine

import (
	"strings"

	documents "go4search/documents"
	nlp "go4search/nlp"
)

type InvertedIndex map[string][]int

/**
 * Build an inverted index by tokenizing the documents and storing the doucment IDs to key-value store.
 * The key is the token, and the value is a slice of document IDs.
 *
 * @param documents A slice of documents
 * @return InvertedIndex
 */
func BuildInvertedIndex(documents []documents.Document, use_tokenizer bool) InvertedIndex {
	index := make(InvertedIndex)

	// iterate all documents
	for _, doc := range documents {
		// support both pre-trained sentence-piece tokenizer and simple whitespace tokenizer
		var tokens []string
		if use_tokenizer {
			tokens = nlp.Tokenize_Query(strings.ToLower(doc.Content))
		} else {
			tokens = strings.Fields(strings.ToLower(doc.Content))
		}

		// iterate all tokens in the document, and store the document ID to the key-value store
		for _, token := range tokens {
			if _, ok := index[token]; !ok {
				index[token] = make([]int, 0)
			}
			index[token] = append(index[token], doc.ID)
		}
	}

	return index
}
