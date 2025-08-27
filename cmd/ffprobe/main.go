package main

import (
	"fmt"
	"hls-on-the-fly/internal/ffprobe"
)

func main() {
	p := "./tmp/vid-big.mp4"

	data, err := ffprobe.Frames(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Probe result:", data)
}
