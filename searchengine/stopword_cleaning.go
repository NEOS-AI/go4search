package searchengine

import (
	"github.com/bbalet/stopwords"

	nlp "go4search/nlp"
)

func removeStopwords(query string) string {
	language, isExist := nlp.DetectLanguage(query)
	if !isExist {
		return query
	}

	cleaned_query := stopwords.CleanString(query, language, true)
	if cleaned_query == "" {
		cleaned_query = query
	}
	return cleaned_query
}
