# Go4Search

Golang based search engine implementation.

## Implemented Features

* Indexing
    * Inverted Index
* Search
    * TF-IDF
    * BM25
* Natural Language Processing
    * Subword Tokenization
    * Stopword Removal
    * Language Detection

## Profiling and Tracing

Basically, this application uses `net/http/pprof` for profiling and tracing.

For visualizing the profiling and tracing, open `http://localhost:6060/debug/pprof/` in your browser.

## ToDos

* [ ] Build index from reading and parsing raw text files
* [ ] Save data to file (dump and load)
* [ ] Levenshtein Distance Spell Correction
* [ ] Pseudo Relevance Feedback
* [ ] Query Expansion

## References

- [blurfx/mini-search-engine](https://github.com/blurfx/mini-search-engine)
- [System Design for Discovery](https://eugeneyan.com/writing/system-design-for-discovery/)
- [🤗 bert-base-multilingual-cased](https://huggingface.co/bert-base-multilingual-cased)
- [sugarme/tokenizer](https://pkg.go.dev/github.com/sugarme/tokenizer)
