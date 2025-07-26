package processes

import (
	"encoding/csv"
	"fmt"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/digitaltembo/motli/packages/corpus/utils"
)

var lettersByLanguage = map[string]string{
	"en":        "abcdefghiklmnopqrstuvwxyz",
	"simple-en": "abcdefghiklmnopqrstuvwxyz",
}

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
	Word string
	// counts of this word accross a sample linguistic use of the corpus
	CorpusCounts Counts
	// counts of this word for each word
	WordCounts Counts
}

// Representation of the Analysis as a line of a CSV
func (a Analysis) toString() string {
	return fmt.Sprintf("%s,%s,%s", a.Word, a.CorpusCounts.toString(), a.WordCounts.toString())
}

// Parse the Analysis from a csv line that has been split into string parts
func (a *Analysis) read(records []string) {
	a.Word = records[0]
	a.CorpusCounts.readAtOffset(records, 1)
	a.WordCounts.readAtOffset(records, 5)
}

// Analyze the frequency of ngrams in the corpus of example sentences
// for all words of language in the wiktionary, and save the output as a csv
// TODO(): Allow for other corpora? The dictionary examples may be biased
func AnalyzeNgrams(language string, n int) ([]*Analysis, error) {
	outputFile, err := utils.NgramFile(language, n)
	if err != nil {
		return nil, err
	}

	if !utils.FileExists(outputFile) {
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
		fmt.Fprintf(os.Stderr, "Already analyzed!")
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
			if len(record) == 9 {
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
func analyze(out *os.File, language string, words []string) ([]*Analysis, error) {
	wordMap := map[string]*Analysis{}
	for _, word := range words {
		wordMap[word] = &Analysis{Word: word}
	}

	entryLanguage := language
	if language == "simple-en" {
		entryLanguage = "en"
	}
	entryCh, err := ParseWikiExtract(language)
	if err != nil {
		return nil, err
	}

	analyzed := 0
	for {
		entry := <-entryCh

		if entry == nil {
			break
		} else if entry.LangCode == entryLanguage {
			analyzed++
			if analyzed%1000 == 0 {
				fmt.Fprintf(os.Stderr, "Analyzed %d words (latest: %s)\n", analyzed, entry.Word)
			}
			for _, word := range words {
				updateAnalysis(word, wordMap[word], entry)
			}
		}
	}
	fmt.Fprintln(out, "string,corpusCount,corpusMulti,corpusPrefix,corpusSuffix,wordCount,wordMulti,wordPrefix,wordSuffix")

	for _, word := range words {
		fmt.Fprintln(out, wordMap[word].toString())
	}
	return slices.Collect(maps.Values(wordMap)), nil
}

// Return the list of ngrams up to the count
// ie ngrams('en', 1) => ["a", "b", ...]
// ngrams('en', 2) => ["aa", "ab", ...]
func ngrams(language string, n int) ([]string, error) {
	letters, ok := lettersByLanguage[language]
	runes := []rune(letters)
	if !ok {
		return nil, fmt.Errorf("invalid language")
	}
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

func updateAnalysis(word string, analysis *Analysis, entry *WeWord) {
	updateCount(&analysis.WordCounts, word, entry.Word)
	for _, s := range entry.Senses {
		for _, e := range s.Examples {
			for _, w := range strings.Split(e.Text, " ") {
				updateCount(&analysis.CorpusCounts, word, w)
			}
		}
	}
}

func updateCount(counts *Counts, word string, inWord string) {
	inWord = strings.ToLower(inWord)
	if strings.HasPrefix(inWord, word) {
		counts.Prefix++
	}
	if strings.HasSuffix(inWord, word) {
		counts.Suffix++
	}
	i := strings.Index(inWord, word)
	li := strings.LastIndex(inWord, word)
	if i >= 0 {
		counts.Count++
		if i != li {
			counts.Multi++
		}
	}
}
