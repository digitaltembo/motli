package main

import (
	"fmt"

	"github.com/digitaltembo/motli/packages/corpus/processes"
	"github.com/digitaltembo/motli/packages/corpus/utils"
)

func main() {
	args := utils.ParseArgs()

	if args.Download != nil {
		filename, err := utils.DownloadWikiExtract(args.Download.Language)
		if err != nil {
			fmt.Printf("Failed to download %s: %s\n", args.Download.Language, err.Error())
		} else {
			fmt.Printf("File downloaded at %s\n", filename)
		}
		return
	}

	if args.Analyze != nil {
		if args.Analyze.Tiles > 0 {
			_, err := processes.TileSet(args.Analyze.Language, args.Analyze.Tiles)

			if err != nil {
				fmt.Printf("Failed to get tiles %s: %s\n", args.Analyze.Language, err.Error())
			}
		} else {
			_, err := processes.AnalyzeNgrams(args.Analyze.Language, args.Analyze.Ngrams)

			if err != nil {
				fmt.Printf("Failed to analyze %s: %s\n", args.Analyze.Language, err.Error())
			}
		}
		return
	}

	fmt.Println("Didn't do anything")
}
