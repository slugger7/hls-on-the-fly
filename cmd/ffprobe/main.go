package main

import (
	"fmt"
	"hls-on-the-fly/internal/ffprobe"
)

func main() {
	p := "./tmp/vid.mp4"

	data := ffprobe.FFProbe(p)

	fmt.Println("Probe result:", data)
}
