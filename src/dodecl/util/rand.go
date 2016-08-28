package util

import "math/rand"

var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

// RandID generates a random ID given a length.
func RandID(length int) string {
	buf := make([]rune, length)
	for i := range buf {
		idx := rand.Intn(len(letters))
		buf[i] = letters[idx]
	}

	return string(buf)
}
