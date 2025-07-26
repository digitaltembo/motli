package utils

import (
	"fmt"
	"time"

	"github.com/cavaliergopher/grab/v3"
)

var wikiextractFiles = map[string]string{
	// download wiktionary extracts via https://kaikki.org/dictionary/rawdata.html
	"en":        "https://kaikki.org/dictionary/raw-wiktextract-data.jsonl.gz",
	"simple-en": "https://kaikki.org/dictionary/downloads/simple/simple-extract.jsonl.gz",
}

func DownloadWikiExtract(language string) (string, error) {
	url, ok := wikiextractFiles[language]
	if !ok {
		return "", fmt.Errorf("invalid language to download")
	}
	target, err := WikiExtractFile(language)
	if err != nil {
		return "", err
	}
	if !FileExists(target) {
		err = downloadFile(target, url)
		if err != nil {
			return "", err
		}
	}
	return target, nil
}

func downloadFile(filepath string, url string) (err error) {
	fmt.Printf("Downloading file from %s to %s", url, filepath)

	client := grab.NewClient()
	req, err := grab.NewRequest(filepath, url)
	resp := client.Do(req)

	t := time.NewTicker(time.Second)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			fmt.Printf("%.02f%% complete, eta %s\n",
				resp.Progress()*100,
				resp.ETA().String())

		case <-resp.Done:
			if err := resp.Err(); err != nil {
				return err
			}
			return nil
		}
	}
}
