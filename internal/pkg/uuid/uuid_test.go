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

	. "github.com/alan-b-lima/prp/internal/pkg/uuid"
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
		t.Error("equal UUIDs have been generated")
	}
}
