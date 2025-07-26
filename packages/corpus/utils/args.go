package utils

import "flag"

// Struct representing parsed command line args for the corpus tool
type Args struct {
	// Either parsed download command or nil, if we do not want to download
	Download *DownloadArgs
	// Either parsed analyze command or nil, if we do not want to analyze
	Analyze *AnalyzeArgs
}

// Struct representing parsed command line args for the download command in the corpus tool
type DownloadArgs struct {
	// Language to analyze. Currently supported "simple-en" and "en"
	Language string
}

// Struct representing parsed command line args for the analyze command in the corpus tool
type AnalyzeArgs struct {
	// Language to analyze. Currently supported "simple-en" and "en"
	Language string
	// Size of ngrams to analyze, defaults to 1 - meaning ["a", "b", "c",...].
	// 2 means ["aa", "ab", ...], and as you can see, grows by power of 26^n,
	// so not very scalable - but good for analyzing a complete set of letters/digrams/trigrams
	Ngrams int
	// Building on ngram analysis, generates a distribution of tiles over the language, using
	// the provided number of tiles - e.g. passing Tiles = 100 means it will compute
	// a distribution of 100 tiles
	Tiles int
}

// Parse command line arguments into the structured Args type
func ParseArgs() Args {
	a := Args{Download: &DownloadArgs{}, Analyze: &AnalyzeArgs{}}

	flag.StringVar(&a.Download.Language, "download", "", "Download specified language")
	flag.StringVar(&a.Analyze.Language, "analyze", "", "Analyze specified language")
	flag.IntVar(&a.Analyze.Ngrams, "ngrams", 1, "Analyze all ngrams in the dictionary for this language up to the provided length")
	flag.IntVar(&a.Analyze.Tiles, "tiles", 0, "Analyze this language and create a set of tiles")
	flag.Parse()

	if a.Download.Language == "" {
		a.Download = nil
	}
	if a.Analyze.Language == "" {
		a.Analyze = nil
	}
	return a
}
