package processes

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"os"
)

// Structure of a wikiextract entry as defined in
// https://github.com/tatuylonen/wiktextract?tab=readme-ov-file#format-of-the-extracted-word-entries
// with a number of fields ignored for our uses
type WeWord struct {
	// The word itself
	Word string `json:"word"`
	// The part of speach, e.g. "noun", "verb", "adj", "adv", "pron", "determiner", "prep", etc
	Pos string `json:"pos"`
	// Wiktionary language code corresponding to lang key (e.g., en)
	LangCode string `json:"lang_code"`
	// List of word senses (dictionaries) for this word/part-of-speech (see below)
	Senses []WeSense `json:"senses"`
}

// Structure of a specific sense of a word within a wikiextract entry WeWord,
// with a number of fields ignored for our purposes
type WeSense struct {
	// List (probably of size 1) of "glosses", or definitions of the word in a particular sense
	Glosses []string `json:"glosses"`
	// Example sentences for the word
	Examples []WeExample `json:"examples"`
}

// Structure of an example of a word within a wikiextract word sesnse WeSense,
// with a small amount of annotative fields ignored for our purposes
type WeExample struct {
	// text of the example
	Text string `json:"text"`
}

// Downloads (if necessary) and parses the WikiExtract gzipped JSONL langauge file
func ParseWikiExtract(language string) (chan *WeWord, error) {
	WikiExtractFile, err := DownloadWikiExtract(language)
	if err != nil {
		return nil, err
	}
	// reads a gzipped jsonl file from wikiextract
	rawf, err := os.Open(WikiExtractFile)
	if err != nil {
		return nil, err
	}
	rawContents, err := gzip.NewReader(rawf)
	if err != nil {
		return nil, err
	}
	bufferedContents := bufio.NewReader(rawContents)
	ch := make(chan *WeWord)
	go func(ch chan *WeWord, contents *bufio.Reader) {
		defer func(ch chan *WeWord) {
			close(ch)
		}(ch)
		for {
			line, err := contents.ReadBytes('\n')
			if err != nil {
				// stop reading when we encounter any error and send a nil
				ch <- nil
				return
			}
			word, err := parseObj(line)
			if err == nil {
				ch <- word
			} // ignore regular errors
		}
	}(ch, bufferedContents)
	return ch, nil
}

func parseObjIntoWeWord(obj []byte, word *WeWord) error {
	if err := json.Unmarshal(obj, word); err != nil {
		return err
	}
	return nil
}

func parseObj(obj []byte) (*WeWord, error) {
	parsedWord := WeWord{}
	if err := parseObjIntoWeWord(obj, &parsedWord); err != nil {
		return nil, err
	}
	return &parsedWord, nil
}
