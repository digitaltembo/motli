package processes

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"os"

	"github.com/digitaltembo/motli/packages/corpus/utils"
)

type WeWord struct {
	Word     string    `json:"word"`
	Pos      string    `json:"pos"`
	LangCode string    `json:"lang_code"`
	Senses   []WeSense `json:"senses"`
}

type WeSense struct {
	// definitions?
	Glosses  []string    `json:"glosses"`
	Examples []WeExample `json:"examples"`
}

type WeExample struct {
	Text string `json:"text"`
}

// Downloads (if necessary) and parses the WikiExtract gzipped JSONL langauge file
func ParseWikiExtract(language string) (chan *WeWord, error) {
	WikiExtractFile, err := utils.DownloadWikiExtract(language)
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
