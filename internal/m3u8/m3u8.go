package m3u8

import (
	"fmt"
	"hls-on-the-fly/internal/ffprobe"
	pathhelpers "hls-on-the-fly/internal/path_helpers"
	"os"
	"path"
	"strings"
)

func CreateManifestForFile(p string, hlsTime int) (string, error) {
	cacheDir := "./cache"

	if hlsTime == 0 {
		hlsTime = 5
	}

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

	previousFrame := 0.0
	i := 1
	for float64(i*hlsTime) <= probe.Duration {
		timestamp := float64(i * hlsTime)
		closestFrame := 0.0
		closestFrameIndex := 0
		for x, f := range probe.Frames {
			if timestamp < f {
				break
			}
			closestFrame = f
			closestFrameIndex = x
		}
		_ = closestFrameIndex

		// reduce our frames list soo we do not have to loop through all previous frames every time
		probe.Frames = probe.Frames[closestFrameIndex:]

		diff := closestFrame - previousFrame
		if _, err := f.WriteString(fmt.Sprintf("#EXTINF:%v,\n%v.%v.ts\n", diff, base, i)); err != nil {
			fmt.Println("could not write to manifest for: ", i, err.Error())
			return "", err
		}

		previousFrame = closestFrame
		i++
	}

	if previousFrame != probe.Duration {
		if _, err := f.WriteString(fmt.Sprintf("#EXTINF:%v,\n%v.%v.ts\n", probe.Duration-previousFrame, base, i)); err != nil {
			fmt.Println("could not write to manifest for: ", i, err.Error())
			return "", err
		}
	}

	if _, err := f.WriteString("\n#EXT-X-ENDLIST"); err != nil {
		fmt.Println("could not write end for file", err.Error())
		return "", err
	}

	return manifestPath, nil
}
