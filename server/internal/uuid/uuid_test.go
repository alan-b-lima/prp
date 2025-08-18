package uuid_test

import (
	"sync"
	"testing"

	. "github.com/alan-b-lima/prp/server/internal/uuid"
)

func TestConcurrentUUIDGeneration(t *testing.T) {
	limit := 10000

	result := make(chan UUID, limit)
	var wg sync.WaitGroup

	wg.Add(limit)
	for range limit {
		go func() {
			uuid := NewUUIDv7()
			
			result <- uuid
			wg.Done()
		}()
	}

	wg.Wait()
	close(result)

	set := make(map[UUID]struct{}, limit)
	for v := range result {
		set[v] = struct{}{}
	}

	if len(set) < limit {
		t.Error("equal UUID's have been generated")
	}
}
