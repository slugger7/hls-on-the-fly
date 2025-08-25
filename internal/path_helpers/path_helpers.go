package pathhelpers

import "strings"

func GetNameWithoutExtension(f string) string {
	if pos := strings.LastIndexByte(f, '.'); pos != -1 {
		return f[:pos]
	}

	return f
}
