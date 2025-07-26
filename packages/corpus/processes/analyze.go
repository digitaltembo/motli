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

func (c Counts) toString() string {
	return fmt.Sprintf("%d,%d,%d,%d", c.Count, c.Multi, c.Prefix, c.Suffix)
}

func (c *Counts) readAtOffset(records []string, offset int) {
	c.Count, _ = strconv.Atoi(records[offset])
	c.Multi, _ = strconv.Atoi(records[offset+1])
	c.Prefix, _ = strconv.Atoi(records[offset+2])
	c.Suffix, _ = strconv.Atoi(records[offset+3])
}

type Analysis struct {
	Word string
	// counts of this word accross a sample linguistic use of the corpus
	CorpusCounts Counts
	// counts of this word for each word
	WordCounts Counts
}

func (a Analysis) toString() string {
	return fmt.Sprintf("%s,%s,%s", a.Word, a.CorpusCounts.toString(), a.WordCounts.toString())
}
func (a *Analysis) read(records []string) {
	a.Word = records[0]
	a.CorpusCounts.readAtOffset(records, 1)
	a.WordCounts.readAtOffset(records, 5)
}

func Analyze(out *os.File, language string, words []string) ([]*Analysis, error) {
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
		return Analyze(output, language, ngrams)
	} else {
		fmt.Println("Already analyzed!")
		f, err := os.Open(outputFile)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		csvReader := csv.NewReader(f)
		records, err := csvReader.ReadAll()
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
