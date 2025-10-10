package hash_test

import (
	crand "crypto/rand"
	"math/rand/v2"
	"testing"

	. "github.com/alan-b-lima/prp/pkg/hash"
)

func TestHashComparison(t *testing.T) {
	const set_size = 100

	const max_len = 100
	const min_len = 10

	buf := make([]byte, 0, max_len)

	for range set_size {
		password := buf[:rand.IntN(max_len-min_len)+min_len]
		crand.Read(password)

		hash := Hash(password)
		if !Compare(hash[:], password) {
			t.Errorf("%x should have compared to true with its hash", password)
		}
	}
}
