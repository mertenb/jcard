package entities

import (
	"testing"

	"github.com/mertenb/jcard/internal/test"
)

func TestAddKind(t *testing.T) {
	kinds := []struct {
		kind string

		valid bool
	}{
		{"individual", true},
		{"group", true},
		{"org", true},
		{"location", true},
		{"unknown", false}, // note: iana-token allowed, but not implement here.
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for _, tupel := range kinds {
		t.Run(tupel.kind, func(t *testing.T) {
			if err := j.AddKind(tupel.kind); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v", !tupel.valid, tupel.valid)
			}
		})
	}
}
