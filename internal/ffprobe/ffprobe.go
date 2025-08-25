package ffprobe

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type FFProbeRawFormat struct {
	Filename string `json:"filename"`
	Duration string `json:"duration"`
	Size     string `json:"size"`
	Bitrate  int    `json:"bitrate"`
}

type FFProbeStream struct {
	Index         int    `json:"index"`
	CodecName     string `json:"codec_name"`
	CodecLongName string `json:"codec_long_name"`
	Profile       string `json:"profile"`
	CodecType     string `json:"codec_type"`
}

type FFProbeData struct {
	Format  FFProbeRawFormat `json:"format"`
	Streams []FFProbeStream  `json:"streams"`
}

func FFProbe(p string) *FFProbeData {
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
		return nil
	}

	var js FFProbeData
	if err := json.Unmarshal(out, &js); err != nil {
		fmt.Println("Could not unmarshal ffprobe data:", err.Error())
		return nil
	}

	return &js
}
