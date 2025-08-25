package ffmpeg

import (
	"fmt"
	"os/exec"
	"path"
	"strconv"
)

func FFMpeg(arg ...string) *exec.Cmd {
	cmd := exec.Command("ffmpeg", arg...)

	fmt.Println(cmd.String())

	return cmd
}

func HLSChunk(duration, offset int, input, destination string) (string, error) {
	cmd := FFMpeg(
		"-hwaccel", "vaapi",
		"-vaapi_device", "/dev/dri/renderD128",
		"-hwaccel_output_format", "vaapi",
		"-ss", strconv.Itoa(offset),
		"-t", strconv.Itoa(duration),
		"-i", path.Clean(input),
		"-c:v", "h264_vaapi",
		"-c:a", "aac",
		"-f", "mpegts", path.Clean(destination),
	)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	fmt.Println(out)

	return destination, nil
}
