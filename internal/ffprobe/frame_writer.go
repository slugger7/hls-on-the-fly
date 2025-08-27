package ffprobe

import (
	"fmt"
	"strconv"
	"strings"
)

type frameWriter struct {
	ch chan float64
}

func (fw frameWriter) Write(p []byte) (int, error) {
	s := string(p)

	comma := strings.Index(s, ",")
	lineType := s[:comma]
	rest := s[comma+1:]
	if lineType == "packet" {
		comma = strings.Index(rest, ",")
		ptsTimeString := rest[:comma]
		flags := rest[comma+1:]

		if strings.HasPrefix(flags, "K_") {
			ptsTime, err := strconv.ParseFloat(ptsTimeString, 64)
			if err != nil {
				fmt.Println("could not parse float:", ptsTimeString, err.Error())
			}
			fw.ch <- ptsTime
		}
	}

	return len(p), nil
}

func (fw frameWriter) KeyFrameConsumer(ch chan<- []float64) {
	l := []float64{}

	for val := range fw.ch {
		l = append(l, val)
	}

	ch <- l

	close(ch)
}
