package ffprobe

import (
	"fmt"
	"os/exec"
)

func Frames(p string) (FrameProbe, error) {
	cmd := exec.Command("ffprobe",
		"-fflags", "+genpts",
		"-v", "error",
		"-skip_frame", "nokey",
		"-show_entries", "packet=pts_time,flags",
		"-show_entries", "format=duration",
		"-show_entries", "stream=duration",
		"-select_streams", "v",
		"-of", "csv",
		p,
	)

	fmt.Println(cmd.String())

	fw := frameWriter{
		ch: make(chan significantType),
	}

	cmd.Stdout = fw

	ch := make(chan FrameProbe)
	go fw.KeyFrameConsumer(ch)

	if err := cmd.Run(); err != nil {
		fmt.Println("error running frame command", err.Error())
	}
	close(fw.ch)

	frames := <-ch

	return frames, nil
}
