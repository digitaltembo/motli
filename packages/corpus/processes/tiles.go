package processes

import (
	"encoding/json"
	"fmt"
	"math"
	"os"

	"github.com/digitaltembo/motli/packages/corpus/sources"
	"github.com/digitaltembo/motli/packages/corpus/utils"
)

// Runs an ngram analysis on the language over the corpus created by example
// in the wiktionary examples in the provided language, and creates a "fairish"
// distribution of tiles - currently specifically by looking at the number of
// occurrences of a given character throughout the entire corpus of example sentences
// in the wiktionary for the provided language
func TileSet(language sources.LanguageSourceId, tileCount int) (map[string]int, error) {
	analysis, err := AnalyzeNgrams(language, 1)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Hey so found %d\n", len(analysis))
	// search for a bucket size that gives the appropriate number of tiles
	min := 1000000
	max := 0
	for _, a := range analysis {
		fmt.Fprintf(os.Stderr, "%s - %d\n", a.Symbol, a.CorpusCounts.Count)
		if a.CorpusCounts.Count > max {
			max = a.CorpusCounts.Count
		}
		if a.CorpusCounts.Count < min {
			min = a.CorpusCounts.Count
		}
	}
	tiles := -1
	var tileMap map[string]int
	for tiles != tileCount {
		bucketSize := (max-min)/2 + min
		tileMap, tiles = tilesGivenBucketSize(analysis, bucketSize)
		if tiles < tileCount {
			max = bucketSize
		} else {
			min = bucketSize
		}

		if max-min < 2 {
			fmt.Fprintf(os.Stderr, "Could not match file size, found %d with a bucket size of %d", tiles, bucketSize)
			break
		}
	}
	outputFile, err := utils.TileFile(string(language), tileCount)
	if err != nil {
		return nil, err
	}
	asJson, err := json.MarshalIndent(tileMap, "", "  ")
	if err != nil {
		return nil, err
	}
	out, err := os.Create(outputFile)
	if err != nil {
		return nil, err
	}
	defer out.Close()
	fmt.Fprint(out, string(asJson))
	return tileMap, nil
}

// Count the tiles created by making one tile for every bucketsize appearances of the tile string
// throughout the entire corpus
func tilesGivenBucketSize(analysis []*Analysis, bucketSize int) (map[string]int, int) {
	tileMap := map[string]int{}
	count := 0
	for _, tile := range analysis {
		tileMap[tile.Symbol] = int(math.Ceil(float64(tile.CorpusCounts.Count) / float64(bucketSize)))
		count += tileMap[tile.Symbol]
	}
	return tileMap, count
}
