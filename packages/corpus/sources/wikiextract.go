package sources

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"bufio"
	"compress/gzip"
	"encoding/json"
	"os"

	"github.com/digitaltembo/motli/packages/corpus/utils"
)

type WikiExtractLanguage string

const (
	WikiExtractLanguage_En       = "we-en"
	WikiExtractLanguage_SimpleEn = "we-simple-en"
)

var wikiextractFiles = map[WikiExtractLanguage]string{
	"we-en":        "https://kaikki.org/dictionary/raw-wiktextract-data.jsonl.gz",
	"we-simple-en": "https://kaikki.org/dictionary/downloads/simple/simple-extract.jsonl.gz",
}

var wikiExtractAlphabets = map[WikiExtractLanguage]string{
	"we-en":        "abcdefghiklmnopqrstuvwxyz",
	"we-simple-en": "abcdefghiklmnopqrstuvwxyz",
}

var languageCode = map[WikiExtractLanguage]string{
	"we-en":        "en",
	"we-simple-en": "en",
}

// Downloads wiktionary extracts produced by the https://github.com/tatuylonen/wiktextract project
// from the archive at https://kaikki.org/dictionary/rawdata.html. These are gzip-compressed JSONL
// files whose structure is documented in the above github
func DownloadWikiExtract(language WikiExtractLanguage) (string, error) {
	url, ok := wikiextractFiles[language]
	if !ok {
		return "", fmt.Errorf("invalid language to download: %s", language)
	}
	target, err := utils.WikiExtractFile(string(language))
	if err != nil {
		return "", err
	}
	if !utils.FileExists(target) {
		err = utils.DownloadFile(target, url)
		if err != nil {
			return "", err
		}
	}
	return target, nil
}

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

type wikiExtractWordSource struct {
	words      map[string]*Word
	categories map[int]string
}

func newWikiExtractWordSource(language WikiExtractLanguage) (*wikiExtractWordSource, error) {
	w := wikiExtractWordSource{words: map[string]*Word{}, categories: map[int]string{}}
	invertedCats := map[string]int{}

	entryCh, err := ParseWikiExtract(language)
	if err != nil {
		return nil, err
	}

	analyzed := 0
	for {
		entry := <-entryCh
		if entry == nil {

			fmt.Fprintf(os.Stderr,
				"Finished reading dictionary, %d words\n", analyzed)
			break
		} else {
			analyzed++
			if analyzed%10000 == 0 {
				fmt.Fprintf(os.Stderr,
					"Added %d words to the %s dictionary (latest: %s)\n",
					analyzed, language, entry.Word)
			}
			cat, ok := invertedCats[entry.Pos]
			if !ok {
				cat = len(invertedCats)
				invertedCats[entry.Pos] = cat
				w.categories[cat] = entry.Pos
			}
			currentWord, ok := w.words[entry.Word]
			if ok {
				if !slices.Contains(currentWord.Categories, cat) {
					currentWord.Categories = append(currentWord.Categories, cat)
				}
			} else {
				w.words[entry.Word] = &Word{
					Word:       entry.Word,
					Categories: []int{cat},
				}
			}
		}
	}

	return &w, nil
}

func (w *wikiExtractWordSource) GetCategory(catId int) string {
	cat, ok := w.categories[catId]
	if !ok {
		return ""
	}
	return cat
}

func (w *wikiExtractWordSource) GetWord(s string) *Word {
	word, ok := w.words[strings.ToLower(s)]
	if !ok {
		return nil
	}
	return word
}
func (w *wikiExtractWordSource) GetWordList() []*Word {
	return slices.Collect(maps.Values((w.words)))
}

type wikiExtractLanguageSource struct {
	language WikiExtractLanguage
}

func newWikiExtractLanguageSource(language WikiExtractLanguage) *wikiExtractLanguageSource {
	return &wikiExtractLanguageSource{language: language}
}

func (w *wikiExtractLanguageSource) Alphabet() string {
	return wikiExtractAlphabets[w.language]
}

func (w *wikiExtractLanguageSource) Read() (chan *string, error) {

	inCh, err := ParseWikiExtract(w.language)
	if err != nil {
		return nil, err
	}
	outCh := make(chan *string)
	analyzed := 0
	go func(inCh chan *WeWord, outCh chan *string) {
		defer func(outCh chan *string) {
			close(outCh)
		}(outCh)

		for {
			entry := <-inCh
			if entry == nil {
				outCh <- nil
				break
			} else {
				analyzed++
				if analyzed%10000 == 0 {
					fmt.Fprintf(os.Stderr,
						"Added  %d examples from the %s dictionary (latest: %s)\n",
						analyzed, w.language, entry.Word)
				}
				for _, s := range entry.Senses {
					for _, e := range s.Examples {
						for _, w := range strings.Split(e.Text, " ") {
							outCh <- &w
						}
					}
				}
			}
		}
	}(inCh, outCh)
	return outCh, nil
}

// Downloads (if necessary) and parses the WikiExtract gzipped JSONL langauge file
func ParseWikiExtract(language WikiExtractLanguage) (chan *WeWord, error) {
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
	// Filter to only include words in the target language,
	// because by default wiktionary includes definitions in the entry language for words in all languages
	entryLanguage := languageCode[language]

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
			if err == nil && word.LangCode == string(entryLanguage) {
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
