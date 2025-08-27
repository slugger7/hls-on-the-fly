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

type ffprobeRawStream struct {
	Index         int    `json:"index"`
	CodecName     string `json:"codec_name"`
	CodecLongName string `json:"codec_long_name"`
	Profile       string `json:"profile"`
	CodecType     string `json:"codec_type"`
}

type ffProbeRawData struct {
	Format  ffprobeRawFormat   `json:"format"`
	Streams []ffprobeRawStream `json:"streams"`
	Frames  []FFProbeFrame     `json:"frames"`
}

type FFProbeFrame struct {
	KeyFrame   int    `json:"key_frame"`
	PtsTime    string `json:"pts_time"`
	PktDtsTime string `json:"pkt_dts_time"`
	SteamIndex int    `json:"stream_index"`
	MediaType  string `json:"media_type"`
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

type FFProbeData struct {
	Format  FFProbeFormat
	Streams []FFProbeStream
	Frames  []FFProbeFrame `json:"frames"`
}

func (f *FFProbeStream) fromRaw(r *ffprobeRawStream) *FFProbeStream {
	f.Index = r.Index
	f.CodecName = r.CodecName
	f.CodecLongName = r.CodecLongName
	f.Profile = r.Profile
	f.CodecType = r.CodecType

	return f
}

func (f *FFProbeData) fromRaw(r *ffProbeRawData) *FFProbeData {
	f.Streams = make([]FFProbeStream, len(r.Streams))
	f.Format = *(&FFProbeFormat{}).fromRaw(&r.Format)
	f.Frames = r.Frames

	for i, m := range r.Streams {
		f.Streams[i] = *(&FFProbeStream{}).fromRaw(&m)
	}

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
		"-show_frames",
		"-show_entries", "frame=key_frame,pts_time,pkt_dts_time,media_type,stream_index:stream=index,codec_name,codec_long_name,profile,codec_type",
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

func Frames(p string) ([]float64, error) {
	cmd := exec.Command("ffprobe",
		"-fflags", "+genpts",
		"-v", "error",
		"-skip_frame", "nokey",
		"-show_entries", "packet=pts_time,flags",
		"-select_streams", "v",
		"-of", "csv",
		p,
	)

	fmt.Println(cmd.String())

	fw := frameWriter{
		ch: make(chan float64),
	}

	cmd.Stdout = fw

	ch := make(chan []float64)
	go fw.KeyFrameConsumer(ch)

	if err := cmd.Run(); err != nil {
		fmt.Println("error running frame command", err.Error())
	}
	close(fw.ch)

	frames := <-ch

	return frames, nil
}
