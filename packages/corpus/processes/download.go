package processes

import (
	"fmt"

	"github.com/digitaltembo/motli/packages/corpus/utils"
)

var wikiextractFiles = map[string]string{
	"en":        "https://kaikki.org/dictionary/raw-wiktextract-data.jsonl.gz",
	"simple-en": "https://kaikki.org/dictionary/downloads/simple/simple-extract.jsonl.gz",
}

// Downloads wiktionary extracts produced by the https://github.com/tatuylonen/wiktextract project
// from the archive at https://kaikki.org/dictionary/rawdata.html. These are gzip-compressed JSONL
// files whose structure is documented in the above github
func DownloadWikiExtract(language string) (string, error) {
	url, ok := wikiextractFiles[language]
	if !ok {
		return "", fmt.Errorf("invalid language to download")
	}
	target, err := utils.WikiExtractFile(language)
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
