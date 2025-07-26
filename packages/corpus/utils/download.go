package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/cavaliergopher/grab/v3"
)

// Downloads a file from the provided url to the filepath, using
// the grab library, and printing progress
func DownloadFile(filepath string, url string) (err error) {
	fmt.Printf("Downloading file from %s to %s", url, filepath)

	client := grab.NewClient()
	req, err := grab.NewRequest(filepath, url)
	resp := client.Do(req)

	t := time.NewTicker(time.Second)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			fmt.Fprintf(os.Stderr,
				"%.02f%% complete, eta %s\n",
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
