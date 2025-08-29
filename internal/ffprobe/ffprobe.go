package ffprobe

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
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

	dat, err := json.Marshal(frames)
	if err != nil {
		fmt.Println("could not marshal frames for debug")
	} else {
		f, err := os.Create(path.Join(".", "cache", "vid", "debug.json"))
		if err != nil {
			fmt.Println("could not create debug log")
			return frames, nil
		}

		defer f.Close()

		if _, err := f.Write(dat); err != nil {
			fmt.Println("could not write data to debug log")
		}
	}

	return frames, nil
}
