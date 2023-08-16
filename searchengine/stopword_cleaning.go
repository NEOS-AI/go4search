package searchengine

import (
	stopwords "github.com/bbalet/stopwords"
)

func removeStopwords(query string) string {
	cleaned_query := stopwords.CleanString(query, "en", true)
	if cleaned_query == "" {
		cleaned_query = query
	}
	return cleaned_query
}
