package main

import (
	"fmt"
	"hls-on-the-fly/internal/ffmpeg"
	"hls-on-the-fly/internal/ffprobe"
	"path"
)

func main() {
	vid := "./tmp/vid.mp4"
	cacheDir := "./cache/vid"
	hlsTime := 5

	probeData, err := ffprobe.Frames(vid)
	if err != nil {
		panic(err)
	}

	for i := 0; float64(i*hlsTime) <= probeData.Duration; i++ {
		out, err := ffmpeg.HLSChunk(hlsTime, hlsTime*i, vid, path.Join(cacheDir, fmt.Sprintf("vid.%v.ts", i)))
		if err != nil {
			fmt.Println("i", i)
			panic(err)
		}

		fmt.Println(out)
	}
}
