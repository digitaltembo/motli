/*
Corpus is a set of tools for creating and analyzing word corpora for the purposes of word games.

My goal here is a set of functionality for creating word lists and scoring algorithms
for tiles and words that is open-source and easily tunable for different games, play styles,
and languages.

Outputs data files in the data directory, including large wiktionary exported files that are
downloaded when run

Usage:

	corpus [--download [language]] [--analyze [language [--ngrams [int]] [--tiles [int]]

The flags are:

	--help
			Print the help text

	--download [language]
			Download the wikiextract file for the given language, storing the gzipped jsonl files
			in the data directory

	--analyze [language]
			Run analysis on the language, defaulting to an ngram analysis of size 1

	--analyse [language] --ngrams [int]
			Run ngram analysis of the provided size, storing thee results as a csv in the data directory

	--analyse [language] --tiles [int]
			Run analysis of ngram size of 1 and create a set of tiles of the provided size whose
			frequency corresponds to the frequency of the ngrams in that language's corpus, storing the
			results as a JSON file in the data directory
*/
package main

import (
	"fmt"

	"github.com/digitaltembo/motli/packages/corpus/processes"
	"github.com/digitaltembo/motli/packages/corpus/sources"
	"github.com/digitaltembo/motli/packages/corpus/utils"
)

func main() {
	args := utils.ParseArgs()

	if args.Download != nil {
		filename, err := sources.DownloadWikiExtract(sources.WikiExtractLanguage(args.Download.Language))
		if err != nil {
			fmt.Printf("Failed to download %s: %s\n", args.Download.Language, err.Error())
		} else {
			fmt.Printf("File downloaded at %s\n", filename)
		}
		return
	}

	if args.Analyze != nil {
		if args.Analyze.Tiles > 0 {
			_, err := processes.TileSet(
				sources.LanguageSourceId(args.Analyze.Language),
				args.Analyze.Tiles)

			if err != nil {
				fmt.Printf("Failed to get tiles %s: %s\n", args.Analyze.Language, err.Error())
			}
		} else {
			_, err := processes.AnalyzeNgrams(sources.LanguageSourceId(args.Analyze.Language), args.Analyze.Ngrams)

			if err != nil {
				fmt.Printf("Failed to analyze %s: %s\n", args.Analyze.Language, err.Error())
			}
		}
		return
	}

	fmt.Println("Didn't do anything")
}
