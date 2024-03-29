package internal

import "os"

func fullCanonica(path string) string {
	return os.Getenv("DOMAIN") + path
}

func getMinimo(x int, y int) int {
	if x < y {
		return x
	}
	return y
}
