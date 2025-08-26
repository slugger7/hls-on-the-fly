package pathhelpers

import (
	"fmt"
	"strconv"
	"strings"
)

func GetNameWithoutExtension(f string) string {
	if pos := strings.LastIndexByte(f, '.'); pos != -1 {
		return f[:pos]
	}

	return f
}

func GetChunkNumber(f string) (int, error) {
	chunks := strings.Split(f, ".")
	if len(chunks) < 3 {
		return 0, fmt.Errorf("no chunk in video file")
	}

	chunk := chunks[len(chunks)-2]

	res, err := strconv.Atoi(chunk)
	if err != nil {
		fmt.Println("could not parse chunk: ", chunk, err.Error())
		return 0, err
	}

	return res, nil
}
