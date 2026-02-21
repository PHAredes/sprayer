package job

import (
	"testing"
	"testing/quick"
)

func TestJob_Deduplication_Property(t *testing.T) {
	f := func(id1, id2 string) bool {
		j1 := Job{ID: id1, Title: "A"}
		j2 := Job{ID: id2, Title: "B"}
		jobs := []Job{j1, j2}
		
		deduped := Dedup()(jobs)
		
		if id1 == id2 {
			return len(deduped) == 1
		} else {
			return len(deduped) == 2
		}
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
