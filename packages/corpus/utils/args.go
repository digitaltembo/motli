package utils

import "flag"

type Args struct {
	Download *DownloadArgs
	Analyze  *AnalyzeArgs
}
type DownloadArgs struct {
	Language string
}

type AnalyzeArgs struct {
	Language string
	Ngrams   int
	Tiles    int
}

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
