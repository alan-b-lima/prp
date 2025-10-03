// Copyright (C) 2025 Alan Barbosa Lima.
//
// PRP is licensed under the GNU General Public License
// version 3. You should have received a copy of the
// license, located in LICENSE, at the root of the source
// tree. If not, see <https://www.gnu.org/licenses/>.

package uuid_test

import (
	"crypto/rand"
	"sync"
	"testing"

	. "github.com/alan-b-lima/prp/pkg/uuid"
)

func TestConcurrentUUIDGeneration(t *testing.T) {
	const numBatches, batchSize = 123, 1999
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

func TestInversabilityBetweenStringAndFromString(t *testing.T) {
	const numTests = 1000

	for range numTests {
		var uuid UUID
		rand.Read(uuid[:])

		str := uuid.String()
		if uuid2, err := FromString(str); err != nil {
			t.Error(err)
		} else if uuid != uuid2 {
			t.Errorf("%x and %x should be equal", uuid, uuid2)
		}
	}
}
