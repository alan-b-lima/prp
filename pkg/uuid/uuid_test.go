// Copyright (C) 2025 Alan Barbosa Lima
//
// PRP is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// PRP is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
// or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public
// License for more details.
//
// You should have received a copy of the GNU General Public License
// along with PRP, located in LICENSE, at the root of the source
// tree. If not, see <https://www.gnu.org/licenses/>.

package uuid_test

import (
	"sync"
	"testing"

	. "github.com/alan-b-lima/prp/pkg/uuid"
)

func TestConcurrentUUIDGeneration(t *testing.T) {
	numBatches, batchSize := 123, 1999
	limit := numBatches * batchSize

	result := make([]UUID, limit)
	var wg sync.WaitGroup

	wg.Add(numBatches)
	for i := range numBatches {
		offset := i * batchSize
		r := result[offset : offset+batchSize]

		go func() {
			for i := range batchSize {
				uuid := NewUUIDv7()
				r[i] = uuid
			}

			wg.Done()
		}()
	}

	wg.Wait()

	set := make(map[UUID]struct{}, limit)
	for _, v := range result {
		set[v] = struct{}{}
	}

	if len(set) < limit {
		t.Error("equal UUIDs have been generated")
	}
}
