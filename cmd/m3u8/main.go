package main

import (
	"fmt"
	"hls-on-the-fly/internal/m3u8"
)

func main() {
	vid := "./tmp/vid.mp4"

	manifest, err := m3u8.CreateManifestForFile(vid, 5, "./cache")
	if err != nil {
		panic(err)
	}

	fmt.Println("manifest created", manifest)

	segments, _ := m3u8.ParseManifest("./cache/vid/vid.m3u8")

	for _, s := range segments {
		fmt.Println(s.Name, s.Duration, s.Start, s.Start+s.Duration)
	}
}
