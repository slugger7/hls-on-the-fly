package ffprobe

import (
	"fmt"
	"strconv"
	"strings"
)

type lineType int

const (
	PACKET lineType = iota
	DURATION
)

type significantType struct {
	ContentType lineType
	Value       float64
}

type frameWriter struct {
	ch chan significantType
}

func (fw frameWriter) Write(p []byte) (int, error) {
	ss := string(p)

	lines := strings.Split(ss, "\n")

	for _, s := range lines {
		comma := strings.Index(s, ",")
		if comma == -1 {
			continue
		}
		lineType := s[:comma]
		rest := s[comma+1:]
		comma = strings.Index(rest, ",")
		switch lineType {
		case "packet":
			ptsTimeString := rest[:comma]
			flags := rest[comma+1:]

			if strings.HasPrefix(flags, "K_") {
				ptsTime, err := strconv.ParseFloat(ptsTimeString, 64)
				if err != nil {
					fmt.Println("could not parse float:", ptsTimeString, err.Error())
					break
				}
				fw.ch <- significantType{PACKET, ptsTime}
			}
		case "format", "stream":
			d, err := strconv.ParseFloat(strings.Trim(rest, "\n"), 64)
			if err != nil {
				fmt.Println("could not parse float for duration:", rest, err.Error())
				break
			}

			fw.ch <- significantType{DURATION, d}
		}
	}

	return len(p), nil
}

type FrameProbe struct {
	Frames   []float64 `json:"frames"`
	Duration float64   `json:"duration"`
}

func (fw frameWriter) KeyFrameConsumer(ch chan<- FrameProbe) {
	l := []float64{}
	d := 0.0

	for val := range fw.ch {
		switch val.ContentType {
		case PACKET:
			l = append(l, val.Value)
		case DURATION:
			if val.Value > d {
				d = val.Value
			}
		}
	}

	ch <- FrameProbe{
		Frames:   l,
		Duration: d,
	}

	close(ch)
}
