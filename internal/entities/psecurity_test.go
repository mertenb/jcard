package entities

import (
	"fmt"
	"testing"

	"github.com/mertenb/jcard/internal/test"
)

func TestAddKey(t *testing.T) {
	keys := []struct {
		n     string
		param map[string][]string
		stype string
	}{
		{"reiner text", map[string][]string{"pref": {"51"}}, "text"},
		{"http://www.example.com/keys/jdoe.cer", map[string][]string{"type": {"work"}}, "uri"},
		{"MEDIATYPE=application/pgp-keys:ftp://example.com/keys/jdoe", nil, "text"},
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for _, tupel := range keys {
		t.Run(fmt.Sprintf("%v", tupel.n), func(t *testing.T) {
			j.RemoveAll("key")
			if err := j.AddKey(tupel.n, tupel.param); (err != nil) || (tupel.stype != j.GetFirst("key").Type) {
				t.Errorf("Got %v, want %v. (%w)", j.GetFirst("key").Type, tupel.stype, err)
			}
		})
	}
}
