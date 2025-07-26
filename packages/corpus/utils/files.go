package utils

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
)

// Directory of all output files and stored input files for processing steps
func DataDir() (string, error) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("failed to find data directory")
	}
	// Note: when moving this file around, we need to update the path
	return path.Join(path.Dir(path.Dir(file)), "data"), nil
}

// Path to the wikiextract gzipped jsonl file of the provided language
func WikiExtractFile(language string) (string, error) {
	data, err := DataDir()
	if err != nil {
		return "", err
	}
	return path.Join(data, fmt.Sprintf("%s.jsonl.gz", language)), nil
}

// Path to the csv file of analysis of ngrams of size ngram in the provided language
func NgramFile(language string, ngram int) (string, error) {
	data, err := DataDir()
	if err != nil {
		return "", err
	}

	return path.Join(data, fmt.Sprintf("%s-%dgram.csv", language, ngram)), nil
}

// Path to the json file of the suggested tile distribution in the provided language
func TileFile(language string, tileCount int) (string, error) {
	data, err := DataDir()
	if err != nil {
		return "", err
	}

	return path.Join(data, fmt.Sprintf("%s-%dtiles.json", language, tileCount)), nil
}

// helper to determine whether file exists
func FileExists(file string) bool {
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
