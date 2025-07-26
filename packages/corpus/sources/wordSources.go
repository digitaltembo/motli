package sources

import (
	"fmt"
	"regexp"
)

type WordSourceId string

const (
	// Words sourced from wikiextract simple-english dictionary with default filters
	WordSourceId_WeSimpleEn = "we-simple-en"
	// All words from the wikiextract simple-english dictionary
	WordSourceId_WeSimpleEnAll = "we-simple-en-all"

	// Words source from wikiextract English dictionary with default filters
	WordSourceId_WeEn = "we-en"
	// All English words from the wikiextract English dictionary
	WordSourceId_WeEnAll = "we-en-all"
)

type Word struct {
	Word       string `json:"w"`
	Categories []int  `json:"cats,omitempty"`
	Freq       int    `json:"freq"`
}
type WordSource interface {
	GetCategory(catId int) string
	GetWord(w string) *Word
	GetWordList() []*Word
}

func GetWordSource(srcId WordSourceId) (WordSource, error) {
	switch srcId {
	case WordSourceId_WeSimpleEnAll:
		return newWikiExtractWordSource(WikiExtractLanguage_SimpleEn)
	case WordSourceId_WeEnAll:
		return newWikiExtractWordSource(WikiExtractLanguage_En)
	default:
		return nil, fmt.Errorf("unsupported word source")
	}
}

type filteredWordSource struct {
	s      WordSource
	filter func(*Word) bool
}

func (f *filteredWordSource) GetCategory(catId int) string {
	return f.s.GetCategory(catId)
}
func (f *filteredWordSource) GetWord(w string) *Word {
	word := f.s.GetWord(w)
	if word == nil {
		return nil
	}
	if !f.filter(word) {
		return nil
	}
	return word
}

func (f *filteredWordSource) GetWordList() (ret []*Word) {
	for _, s := range f.s.GetWordList() {
		if f.filter(s) {
			ret = append(ret, s)
		}
	}
	return
}

func FilterWordSource(source WordSource, filter func(*Word) bool) WordSource {
	return &filteredWordSource{
		s:      source,
		filter: filter,
	}
}

var reasonableEnglishRegex *regexp.Regexp = nil

// a simple filter for reasonable words to play,
// requiring at least 2 letters and no punction or spaces
func ReasonableEnglishWord(w *Word) bool {
	if reasonableEnglishRegex == nil {
		reasonableEnglishRegex, _ = regexp.Compile("^[a-zA-Z][a-zA-Z]+$")
	}
	if reasonableEnglishRegex.Match([]byte(w.Word)) {
		return true
	}
	return false
}
