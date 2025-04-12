package random

import (
	"math/rand"
	"time"
)

func NewRandomString(size int) string {
	var r []rune = make([]rune, size)

	t := time.Now()

	seed := int64(t.Day()) + int64(t.Nanosecond())
	rnd := rand.New(rand.NewSource(seed))

	for i := range r {
		r[i] = rune(rnd.Int31n(122-97)) + 97
	}

	return string(r)
}
