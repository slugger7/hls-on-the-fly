package m3u8

import (
	"fmt"
	"hls-on-the-fly/internal/ffprobe"
	pathhelpers "hls-on-the-fly/internal/path_helpers"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

type Segment struct {
	Name     string
	Duration float64
	Start    float64
}

func ParseManifest(p string) ([]Segment, error) {
	dat, err := os.ReadFile(p)
	if err != nil {
		fmt.Println("colud not open playlist file", p, err.Error())
		return nil, err
	}

	tags := strings.Split(string(dat), "#")

	segments := []Segment{}
	start := 0.0
	for _, t := range tags {
		if strings.HasPrefix(t, "EXTINF") {
			t = strings.ReplaceAll(t, "\n", "")
			colonIndex := strings.Index(t, ":")
			commaIndex := strings.Index(t, ",")
			durationString := t[colonIndex+1 : commaIndex]
			name := t[commaIndex+1:]

			duration, err := strconv.ParseFloat(durationString, 64)
			if err != nil {
				fmt.Println("could not parse duration", durationString, err.Error())
				return nil, err
			}

			segments = append(segments, Segment{
				Name:     name,
				Duration: duration,
				Start:    start,
			})

			start += duration
		}
	}

	return segments, nil
}

func CreateManifestForFile(p string, hlsTime int, cacheDir string) (string, error) {
	probe, err := ffprobe.Frames(p)
	if err != nil {
		fmt.Println("Could not p:", p, err.Error())
		return "", err
	}

	fmt.Println(probe.Duration, len(probe.Frames))

	lines := []string{
		"#EXTM3U",
		"#EXT-X-VERSION:3",
		fmt.Sprintf("#EXT-X-TARGETDURATION:%v", hlsTime),
		"#EXT-X-MEDIA-SEQUENCE:0",
		"\n",
	}

	data := strings.Join(lines, "\n")

	base := pathhelpers.GetNameWithoutExtension(path.Base(p))

	mediaFolder := path.Join(cacheDir, base)

	if err := os.MkdirAll(mediaFolder, 0777); err != nil {
		fmt.Println("could not create directory for media: ", err.Error())
		return "", err
	}

	manifestPath := path.Join(mediaFolder, fmt.Sprintf("%v.m3u8", base))

	f, err := os.Create(manifestPath)
	if err != nil {
		fmt.Println("could not create manifest file: ", err.Error())
		return "", nil
	}

	defer f.Close()

	if _, err := f.WriteString(data); err != nil {
		fmt.Println("could not write head of manifest:", err.Error())
		return "", err
	}

	// have to append the duration as "the last keyframe"
	segments := generateSegmentsForManifest(hlsTime, probe.Duration, probe.Frames, func(i int) string {
		return fmt.Sprintf("%v.%v.ts", base, i)
	})

	for _, s := range segments {
		if _, err := f.WriteString(fmt.Sprintf("#EXTINF:%v,\n%v\n", s.Duration, s.Name)); err != nil {
			fmt.Println("could not write to manifest for: ", s.Name, err.Error())
			return "", err
		}
	}

	if _, err := f.WriteString("\n#EXT-X-ENDLIST"); err != nil {
		fmt.Println("could not write end for file", err.Error())
		return "", err
	}

	return manifestPath, nil
}

func generateSegmentsForManifest(hlsTime int, duration float64, frames []float64, nameFunc func(int) string) []Segment {
	if len(frames) == 0 {
		return []Segment{}
	}

	segments := []Segment{}

	previousFrame := 0.0
	latestSegmentFrame := 0.0
	segmentCounter := 0
	for _, f := range frames {
		if float64(hlsTime) >= f-latestSegmentFrame {
			previousFrame = f
			continue
			// find a frame that is just over the hlsTime
		}

		dur, exact := decimal.NewFromFloat(previousFrame).Sub(decimal.NewFromFloat(latestSegmentFrame)).Float64()
		if !exact {
			fmt.Println("Not exact")
		}
		segments = append(segments, Segment{
			Name:     nameFunc(segmentCounter),
			Start:    latestSegmentFrame,
			Duration: dur,
		})
		segmentCounter++
		latestSegmentFrame = previousFrame
	}

	dur, exact := decimal.NewFromFloat(duration).Sub(decimal.NewFromFloat(latestSegmentFrame)).Float64()
	if !exact {
		fmt.Println("Not exact")
	}
	if latestSegmentFrame != duration {
		segments = append(segments, Segment{
			Name:     nameFunc(segmentCounter),
			Start:    latestSegmentFrame,
			Duration: dur,
		})
	}

	return segments
}
