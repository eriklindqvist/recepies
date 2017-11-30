package lib

import "os"

func Getenv(name string, def string) string {
		val := os.Getenv(name)
		if val == "" {
			return def
		}
		return val
}
