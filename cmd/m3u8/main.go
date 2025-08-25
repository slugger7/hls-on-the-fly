package main

import (
	"fmt"
	"hls-on-the-fly/internal/m3u8"
)

func main() {
	vid := "./tmp/vid.mp4"

	manifest, err := m3u8.CreateManifestForFile(vid)
	if err != nil {
		panic(err)
	}

	fmt.Println("manifest created", manifest)
}
