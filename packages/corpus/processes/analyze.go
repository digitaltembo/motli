package processes

import (
	"encoding/csv"
	"fmt"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/digitaltembo/motli/packages/corpus/sources"
	"github.com/digitaltembo/motli/packages/corpus/utils"
)

// Chunk of info on the count of occurrences of a word in a context in a corpus
type Counts struct {
	// number of times it happens at least once in a word
	Count int
	// number of times it happens more than once in a word
	Multi int
	// number of times it happens at the start of a word
	Prefix int
	// number of times it happens at the end of a word
	Suffix int
}

// Count representation in a portion of a line of a CSV
func (c Counts) toString() string {
	return fmt.Sprintf("%d,%d,%d,%d", c.Count, c.Multi, c.Prefix, c.Suffix)
}

// Parse count representation from a portion of a line of a CSV
func (c *Counts) readAtOffset(records []string, offset int) {
	c.Count, _ = strconv.Atoi(records[offset])
	c.Multi, _ = strconv.Atoi(records[offset+1])
	c.Prefix, _ = strconv.Atoi(records[offset+2])
	c.Suffix, _ = strconv.Atoi(records[offset+3])
}

// Analysis of the occurrences of a 'word'/ngram/letter in a corpus
type Analysis struct {
	// the 'word'/ngram/letter we are counting
	Symbol string
	// counts of this word accross a sample linguistic use of the corpus
	CorpusCounts Counts
}

// Representation of the Analysis as a line of a CSV
func (a Analysis) toString() string {
	return fmt.Sprintf("%s,%s", a.Symbol, a.CorpusCounts.toString())
}

// Parse the Analysis from a csv line that has been split into string parts
func (a *Analysis) read(records []string) {
	a.Symbol = records[0]
	a.CorpusCounts.readAtOffset(records, 1)
}

// Analyze the frequency of ngrams in the corpus of example sentences
// for all words of language in the wiktionary, and save the output as a csv
// TODO(): Allow for other corpora? The dictionary examples may be biased
func AnalyzeNgrams(languageId sources.LanguageSourceId, n int) ([]*Analysis, error) {
	outputFile, err := utils.NgramFile(string(languageId), n)
	if err != nil {
		return nil, err
	}

	if !utils.FileExists(outputFile) {
		language, err := sources.GetLanguageSource(languageId)
		if err != nil {
			return nil, err
		}

		output, err := os.Create(outputFile)
		if err != nil {
			return nil, err
		}
		defer output.Close()

		ngrams, err := ngrams(language, n)
		if err != nil {
			return nil, err
		}
		return analyze(output, language, ngrams)
	} else {
		fmt.Fprintf(os.Stderr, "Already analyzed!\n")
		f, err := os.Open(outputFile)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		csvReader := csv.NewReader(f)
		records, err := csvReader.ReadAll()

		if err != nil {
			return nil, err
		}
		analysis := []*Analysis{}
		for _, record := range records[1:] {
			if len(record) > 4 {
				newAnalysis := Analysis{}
				newAnalysis.read(record)
				analysis = append(analysis, &newAnalysis)
			}
		}
		return analysis, nil
	}
}

// Generate an analysis of the frequency of occurrences of words as substrings in the
// wiktionary example sentences of the provided language, and save it as a csv
func analyze(out *os.File, languageSource sources.LanguageSource, symbols []string) ([]*Analysis, error) {
	symbolMap := map[string]*Analysis{}
	for _, symbol := range symbols {
		symbolMap[symbol] = &Analysis{Symbol: symbol}
	}

	wordCh, err := languageSource.Read()
	if err != nil {
		return nil, err
	}

	analyzed := 0
	for {
		word := <-wordCh

		if word == nil {
			break
		} else {
			analyzed++
			if analyzed%100000 == 0 {
				fmt.Fprintf(os.Stderr, "Analyzed %d words (latest: %s)\n", analyzed, *word)
			}
			for _, symbol := range symbols {
				updateCount(&symbolMap[symbol].CorpusCounts, symbol, *word)
			}
		}
	}
	fmt.Fprintln(out, "string,corpusCount,corpusMulti,corpusPrefix,corpusSuffix")

	for _, symbol := range symbols {
		fmt.Fprintln(out, symbolMap[symbol].toString())
	}
	return slices.Collect(maps.Values(symbolMap)), nil
}

// Return the list of ngrams up to the count
// ie ngrams('en', 1) => ["a", "b", ...]
// ngrams('en', 2) => ["aa", "ab", ...]
func ngrams(languageSource sources.LanguageSource, n int) ([]string, error) {
	letters := languageSource.Alphabet()
	runes := []rune(letters)
	ngrams := []string{}
	for _, r := range runes {
		ngrams = append(ngrams, string(r))
	}

	for i := 1; i < n; i++ {
		newNgrams := []string{}
		// newNgrams := make([]string, len(ngrams))
		// copy(newNgrams, ngrams)
		for _, ngram := range ngrams {
			for _, r := range runes {
				newNgrams = append(newNgrams, ngram+string(r))
			}
		}
		ngrams = newNgrams
	}
	return ngrams, nil
}
func updateCount(counts *Counts, symbol string, inWord string) {
	inWord = strings.ToLower(inWord)
	if strings.HasPrefix(inWord, symbol) {
		counts.Prefix++
	}
	if strings.HasSuffix(inWord, symbol) {
		counts.Suffix++
	}
	i := strings.Index(inWord, symbol)
	li := strings.LastIndex(inWord, symbol)
	if i >= 0 {
		counts.Count++
		if i != li {
			counts.Multi++
		}
	}
}
