package entities

import (
	"fmt"
	"testing"

	"github.com/mertenb/jcard/internal/test"
)

func TestAddTz(t *testing.T) {
	ns := []struct {
		n     string
		param map[string][]string
		stype string
	}{
		{"reiner text", map[string][]string{"mediatype": {"application/xml-dtd", "application/zip"}}, "text"},
		{"reiner text", map[string][]string{"pref": {"51"}}, "text"},
		{"-0630", map[string][]string{"pref": {"51"}}, "utc-offset"},
		{"telnet://192.0.2.16:80/", map[string][]string{"pref": {"51"}}, "uri"}, // don't know an example for 'uri timezone' ...
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for _, tupel := range ns {
		t.Run(fmt.Sprintf("%v", tupel.n), func(t *testing.T) {
			j.RemoveAll("tz")
			if err := j.AddTz(tupel.n, tupel.param); (err != nil) || (tupel.stype != j.GetFirst("tz").Type) {
				t.Errorf("Got %v, want %v. (%v)", j.GetFirst("tz").Type, tupel.stype, err)
			}
		})
	}
}

func TestAddGeo(t *testing.T) {
	geos := []struct {
		geo   string
		param map[string][]string
		valid bool
	}{
		{"geo:37.386013,-122.082932", map[string][]string{"mediatype": {"application/xml-dtd", "application/zip"}}, true},
		{"geo:37.386013,-122.082932", map[string][]string{"pref": {"51"}}, true},
		{"geo:37.386013,-122.082932", nil, true},
		{"invalid uri", map[string][]string{"pref": {"51"}}, false},
		{"", nil, false},
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for _, tupel := range geos {
		t.Run(fmt.Sprintf("%v", tupel.geo), func(t *testing.T) {
			if err := j.AddGeo(tupel.geo, tupel.param); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v", !tupel.valid, tupel.valid)
			}
		})

	}
}
