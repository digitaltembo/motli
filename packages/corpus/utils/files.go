package utils

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
)

func DataDir() (string, error) {

	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("failed to find data directory")
	}
	return path.Join(path.Dir(path.Dir(file)), "data"), nil
}
func WikiExtractFile(language string) (string, error) {
	data, err := DataDir()
	if err != nil {
		return "", err
	}
	return path.Join(data, fmt.Sprintf("%s.jsonl.gz", language)), nil
}

func NgramFile(language string, ngram int) (string, error) {
	data, err := DataDir()
	if err != nil {
		return "", err
	}

	return path.Join(data, fmt.Sprintf("%s-%dgram.csv", language, ngram)), nil
}

func TileFile(language string, tileCount int) (string, error) {
	data, err := DataDir()
	if err != nil {
		return "", err
	}

	return path.Join(data, fmt.Sprintf("%s-%dtiles.json", language, tileCount)), nil
}

func FileExists(file string) bool {
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
