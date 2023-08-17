package nlp

import (
	"github.com/pemistahl/lingua-go"
)

var languages = []lingua.Language{
	lingua.English,
	lingua.French,
	lingua.German,
	lingua.Spanish,
	lingua.Korean,
	lingua.Japanese,
	lingua.Chinese,
}
var detector = lingua.NewLanguageDetectorBuilder().FromLanguages(languages...).Build()

/**
 * Detect the language of a given text.
 *
 * @param text string
 * @return (string, bool) language, exists
 */
func DetectLanguage(text string) (string, bool) {
	language, exists := detector.DetectLanguageOf(text)
	return language.String(), exists
}
