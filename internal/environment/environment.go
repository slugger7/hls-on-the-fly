package environment

import (
	"fmt"
	"os"
	"strconv"
)

type Env struct {
	Port    string
	Cache   string
	Media   string
	HlsTime int
}

func GetEnv() *Env {
	env := Env{
		Port:    getStringOrDefault("PORT", "8080"),
		Cache:   getStringOrDefault("CACHE_DIR", "cache"),
		Media:   getStringOrDefault("MEDIA_DIR", "tmp"),
		HlsTime: getIntOrDefault("HLS_TIME", 5),
	}

	return &env
}

func getIntOrDefault(key string, def int) int {
	val, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		val = def
	}

	fmt.Println(key, val)

	return val
}

func getStringOrDefault(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		val = def
	}

	fmt.Println(key, val)

	return val
}
