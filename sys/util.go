package sys

import (
	"math/rand"
	"time"
)

func Random(a int, z int) int {
	var (
		min int = a
		max int = z
	)
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	return rng.Intn(max-min+1) + min
}
