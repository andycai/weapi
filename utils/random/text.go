package random

import (
	"math/rand"
	"time"
)

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyz")
var numberRunes = []rune("0123456789")

func randRunes(n int, source []rune) string {
	// r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// r.Intn(len(source))
	b := make([]rune, n)
	for i := range b {
		b[i] = source[rand.Intn(len(source))]
	}
	return string(b)
}

func RandText(n int) string {
	return randRunes(n, letterRunes)
}

func RandNumberText(n int) string {
	return randRunes(n, numberRunes)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
