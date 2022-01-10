package entities

import (
	"fmt"
	"testing"

	"github.com/mertenb/jcard/internal/test"
)

func TestAddFburl(t *testing.T) {
	fburls := []struct {
		n     string
		param map[string][]string
		valid bool
	}{
		{"http://www.example.com/busy/janedoe", map[string][]string{"mediatype": {"application/xml-dtd", "application/zip"}}, true},
		{"http://www.example.com/busy/janedoe", map[string][]string{"pref": {"1"}}, true},
		{"ftp://example.com/busy/project-a.ifb", nil, true},
		{"noUri", nil, false},
		{"http://www.example.com/busy/janedoe", map[string][]string{"unknown": {"0"}}, false},
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for _, tupel := range fburls {
		t.Run(fmt.Sprintf("%v", tupel.n), func(t *testing.T) {
			j.RemoveAll("fburl")
			if err := j.AddFburl(tupel.n, tupel.param); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v. %w", !tupel.valid, tupel.valid, err)
			}
		})
	}
}

func TestAddCadadruri(t *testing.T) {
	caladruris := []struct {
		n     string
		param map[string][]string
		valid bool
	}{
		{"http://www.example.com/busy/janedoe", map[string][]string{"mediatype": {"application/xml-dtd", "application/zip"}}, true},
		{"https://example.com/calendar/jdoe", map[string][]string{"pref": {"1"}}, true},
		{"ftp://example.com/busy/project-a.ifb", nil, true},
		{"noUri", nil, false},
		{"http://www.example.com/busy/janedoe", map[string][]string{"unknown": {"0"}}, false},
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for _, tupel := range caladruris {
		t.Run(fmt.Sprintf("%v", tupel.n), func(t *testing.T) {
			j.RemoveAll("caladruri")
			if err := j.AddCaladruri(tupel.n, tupel.param); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v. %w", !tupel.valid, tupel.valid, err)
			}
		})
	}
}

func TestAddCaluri(t *testing.T) {
	caluris := []struct {
		n     string
		param map[string][]string
		valid bool
	}{
		{"ftp://ftp.example.com/calA.ics", map[string][]string{"mediatype": {"text/calendar"}}, true},
		{"http://cal.example.com/calA", map[string][]string{"pref": {"1"}}, true},
		{"ftp://example.com/busy/project-a.ifb", nil, true},
		{"noUri", nil, false},
		{"http://www.example.com/busy/janedoe", map[string][]string{"unknown": {"0"}}, false},
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for _, tupel := range caluris {
		t.Run(fmt.Sprintf("%v", tupel.n), func(t *testing.T) {
			j.RemoveAll("caluri")
			if err := j.AddCaluri(tupel.n, tupel.param); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v. %w", !tupel.valid, tupel.valid, err)
			}
		})
	}
}
