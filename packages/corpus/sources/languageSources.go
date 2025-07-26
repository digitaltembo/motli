package sources

import "fmt"

// Source of words as used in language
type LanguageSource interface {
	Alphabet() string
	// Returns a channel streaming individual words, whch may have punctuation etc.
	Read() (chan *string, error)
}

type LanguageSourceId string

const (
	// Language sourced from wikiextract simple-english dictionary examples with default filters
	LanguageSourceId_SimpleEn = "we-simple-en"
	// Language sourced from all simple-english dictionary exampes
	LanguageSourceId_SimpleEnAll = "we-simple-en-all"

	// Language source from the wikiextract English dictionary with words in the simple-english dictionary
	LanguageSourceId_SimpleEnFromEnExamples = "we-simple-en-from-en"
	// Langauge sourced from the wikiextract English dictionary examples with default filters
	LanguageSourceId_En = "we-en"
	// Language source from the wikiextract English dictionary examples
	LanguageSourceId_EnAll = "we-en-all"
)

func GetLanguageSource(srcId LanguageSourceId) (LanguageSource, error) {
	switch srcId {
	// All from simple-english examples
	case LanguageSourceId_SimpleEnAll:
		return newWikiExtractLanguageSource(WikiExtractLanguage_SimpleEn), nil

	// All simple-english words from simple-english examples which meet the
	// ReasonableEnglishWord criteria
	case LanguageSourceId_SimpleEn:
		ws, err := newWikiExtractWordSource(WikiExtractLanguage_SimpleEn)
		if err != nil {
			return nil, err
		}
		filterWordSource := FilterWordSource(ws, ReasonableEnglishWord)
		return FilterLanguageSource(
			newWikiExtractLanguageSource(WikiExtractLanguage_SimpleEn),
			filterWordSource), nil

	// All simple-english words from the full english examples
	case LanguageSourceId_SimpleEnFromEnExamples:
		ws, err := newWikiExtractWordSource(WikiExtractLanguage_SimpleEn)
		if err != nil {
			return nil, err
		}
		return FilterLanguageSource(newWikiExtractLanguageSource(WikiExtractLanguage_En), ws), nil

	// All English words from the English examples that meet the
	// Reasonable English Word criteria
	case LanguageSourceId_En:
		ws, err := newWikiExtractWordSource(WikiExtractLanguage_En)
		if err != nil {
			return nil, err
		}
		filterWordSource := FilterWordSource(ws, ReasonableEnglishWord)
		return FilterLanguageSource(
			newWikiExtractLanguageSource(WikiExtractLanguage_En),
			filterWordSource), nil

	// All words from the English examples
	case LanguageSourceId_EnAll:
		return newWikiExtractLanguageSource(WikiExtractLanguage_En), nil

	default:
		return nil, fmt.Errorf("unsupported language source")
	}
}

type filteredLanguageSource struct {
	ls LanguageSource
	ws WordSource
}

func FilterLanguageSource(ls LanguageSource, ws WordSource) *filteredLanguageSource {
	return &filteredLanguageSource{ls: ls, ws: ws}
}
func (fls *filteredLanguageSource) Alphabet() string { return fls.ls.Alphabet() }
func (fls *filteredLanguageSource) Read() (chan *string, error) {
	inCh, err := fls.ls.Read()
	if err != nil {
		return nil, err
	}
	outCh := make(chan *string)
	go func(inCh chan *string, outCh chan *string, ws WordSource) {
		defer func(ch chan *string) {
			close(ch)
		}(outCh)
		for {
			word := <-inCh
			if word == nil || ws.GetWord(*word) != nil {
				outCh <- word
			}
		}

	}(inCh, outCh, fls.ws)

	return outCh, nil
}
