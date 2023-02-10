package env

import "os"

func Get(key, def string) string {
	v := os.Getenv(key)
	if len(v) == 0 {
		return def
	}
	return v
}
