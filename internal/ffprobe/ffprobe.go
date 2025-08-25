package ffprobe

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
)

type ffprobeRawFormat struct {
	Filename string `json:"filename"`
	Duration string `json:"duration"`
	Size     string `json:"size"`
	Bitrate  int    `json:"bitrate"`
}

type FFProbeFormat struct {
	Filename string
	Duration float64
	Size     int
	Bitrate  int
}

type FFProbeStream struct {
	Index         int    `json:"index"`
	CodecName     string `json:"codec_name"`
	CodecLongName string `json:"codec_long_name"`
	Profile       string `json:"profile"`
	CodecType     string `json:"codec_type"`
}

type ffProbeRawData struct {
	Format  ffprobeRawFormat `json:"format"`
	Streams []FFProbeStream  `json:"streams"`
}

type FFProbeData struct {
	Format  FFProbeFormat
	Streams []FFProbeStream
}

func (f *FFProbeData) fromRaw(r *ffProbeRawData) *FFProbeData {
	f.Streams = r.Streams
	f.Format = *(&FFProbeFormat{}).fromRaw(&r.Format)

	return f
}

func (f *FFProbeFormat) fromRaw(r *ffprobeRawFormat) *FFProbeFormat {
	f.Duration, _ = strconv.ParseFloat(r.Duration, 64)
	f.Size, _ = strconv.Atoi(r.Size)

	f.Filename = r.Filename
	f.Bitrate = r.Bitrate

	return f
}

func FFProbe(p string) (*FFProbeData, error) {
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		p,
	)

	fmt.Println(cmd.String())

	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Colud not get output for command: ", err.Error())
		return nil, err
	}

	var js ffProbeRawData
	if err := json.Unmarshal(out, &js); err != nil {
		fmt.Println("Could not unmarshal ffprobe data:", err.Error())
		return nil, err
	}

	return (&FFProbeData{}).fromRaw(&js), nil
}
