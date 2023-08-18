package nlp

import (
	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/pretrained"
)

var tokenizers *tokenizer.Tokenizer

func Init_Tokenizer() {
	// Download and cache pretrained tokenizer. In this case `bert-base-multilingual-cased` from Huggingface
	// can be any model with `tokenizer.json` available. E.g. `tiiuae/falcon-7b`
	configFile, err := tokenizer.CachedPath("bert-base-multilingual-cased", "tokenizer.json")
	if err != nil {
		panic(err)
	}

	tk, err := pretrained.FromFile(configFile)
	if err != nil {
		panic(err)
	}
	tokenizers = tk
}

func Tokenize_Query(query string) []string {
	en, err := tokenizers.EncodeSingle(query)
	if err != nil {
		panic(err)
	}
	return en.Tokens
}
