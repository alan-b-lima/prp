package hash

import (
	"math/rand/v2"
	"sync"
	"time"
)

var (
	source rand.Source
	mulock sync.Mutex
)

func init() {
	lo := uint64(time.Now().UnixNano())
	hi := lo * 0xF0551

	source = rand.NewPCG(lo, hi)
}

func Hash(data []byte) [32]byte {
	digest := [32]byte{}

	for i := range len(data) {
		for j := 0; j < 0x10; j += 4 {
			digest[(i+j)&0x1F] = digest[(i+j)&0x1F]*data[i] + data[i]
		}
	}

	return digest
}

func NewSalt() [16]byte {
	mulock.Lock()
	hi, lo := source.Uint64(), source.Uint64()
	mulock.Unlock()

	return [16]byte{
		0x00: byte(lo >> 0x00), byte(lo >> 0x08), byte(lo >> 0x10), byte(lo >> 0x18),
		0x04: byte(lo >> 0x20), byte(lo >> 0x28), byte(lo >> 0x30), byte(lo >> 0x38),
		0x08: byte(hi >> 0x00), byte(hi >> 0x08), byte(hi >> 0x10), byte(hi >> 0x18),
		0x0C: byte(hi >> 0x20), byte(hi >> 0x28), byte(hi >> 0x30), byte(hi >> 0x38),
	}
}
