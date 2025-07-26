package processes

import (
	"encoding/json"
	"fmt"
	"math"
	"os"

	"github.com/digitaltembo/motli/packages/corpus/utils"
)

func TileSet(language string, tileCount int) (map[string]int, error) {
	analysis, err := AnalyzeNgrams(language, 1)
	if err != nil {
		return nil, err
	}
	// search for a bucket size that gives the appropriate number of tiles
	min := 1000000
	max := 0
	for _, a := range analysis {
		fmt.Fprintf(os.Stderr, "%s - %d\n", a.Word, a.CorpusCounts.Count)
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
			fmt.Fprintf(os.Stderr, "Could not match file size :(")
			break
		}
	}
	outputFile, err := utils.TileFile(language, tileCount)
	if err != nil {
		return nil, err
	}
	asJson, err := json.Marshal(tileMap)
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

func tilesGivenBucketSize(analysis []*Analysis, bucketSize int) (map[string]int, int) {
	tileMap := map[string]int{}
	count := 0
	for _, tile := range analysis {
		tileMap[tile.Word] = int(math.Ceil(float64(tile.CorpusCounts.Count) / float64(bucketSize)))
		count += tileMap[tile.Word]
	}
	return tileMap, count
}
