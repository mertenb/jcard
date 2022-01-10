package entities

import (
	"fmt"
	"testing"

	"github.com/mertenb/jcard/internal/test"
)

func TestAddCategories(t *testing.T) {
	cats := []struct {
		categories []string
		param      map[string][]string
		valid      bool
	}{
		{[]string{"individual", "swimmer", "climber", "reader"}, nil, true},
		{[]string{"individual"}, map[string][]string{"pref": {"55"}, "pid": {"44"}, "altid": {"99"}}, true},
		{[]string{"invalid language param"}, map[string][]string{"language": {"de"}, "pref": {"55"}}, false},
		{[]string{}, map[string][]string{"pref": {"55"}}, true},
		{[]string{""}, map[string][]string{"pref": {"55"}}, true},
		{nil, map[string][]string{"pref": {"55"}}, true},
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for _, tupel := range cats {
		t.Run(fmt.Sprintf("%v", tupel.categories), func(t *testing.T) {
			if err := j.AddCategories(tupel.categories, tupel.param); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v", !tupel.valid, tupel.valid)
			}
		})
	}
}

func TestAddNote(t *testing.T) {
	notes := []struct {
		note  string
		param map[string][]string
		valid bool
	}{
		{"this is a note", nil, true},
		{"note", map[string][]string{"pref": {"55"}, "pid": {"44"}, "altid": {"99"}}, true},
		{"note with lang", map[string][]string{"language": {"de"}, "pref": {"55"}}, true},
		{"", map[string][]string{"pref": {"55"}}, true},
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for _, tupel := range notes {
		t.Run(fmt.Sprintf("%v", tupel.note), func(t *testing.T) {
			if err := j.AddNote(tupel.note, tupel.param); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v. %w", !tupel.valid, tupel.valid, err)
			}
		})
	}
}

func TestAddProdid(t *testing.T) {
	prodids := []struct {
		prodid string
		param  map[string][]string
		valid  bool
	}{
		{"this is a note", nil, true},
		{`-//ONLINE DIRECTORY//NONSGML Version 1//EN`, nil, true},
		{"no params allowed", map[string][]string{"pref": {"55"}}, false},
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for _, tupel := range prodids {
		t.Run(fmt.Sprintf("%v", tupel.prodid), func(t *testing.T) {
			if err := j.AddProdid(tupel.prodid, tupel.param); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v. %w", !tupel.valid, tupel.valid, err)
			}
		})
	}
}

func TestAddRev(t *testing.T) {
	revs := []struct {
		rev   string
		param map[string][]string
		valid bool
	}{
		{"19951031T222710Z", nil, true},
		{"19961022T140000Z", nil, true},
		{"only timestamp allowed", nil, false},
		{"19951031T222710Z", map[string][]string{"pref": {"55"}}, false},
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for _, tupel := range revs {
		t.Run(fmt.Sprintf("%v", tupel.rev), func(t *testing.T) {
			if err := j.AddRev(tupel.rev, tupel.param); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v. %w", !tupel.valid, tupel.valid, err)
			}
		})
	}
}

func TestAddSound(t *testing.T) {
	sound :=
		`data:audio/basic;base64,MIICajCCAdOgAwIBAgICBEUwDQYJKoZIh` +
			`AQEEBQAwdzELMAkGA1UEBhMCVVMxLDAqBgNVBAoTI05ldHNjYXBlIENvbW11bm` +
			`ljYXRpb25zIENvcnBvcmF0aW9uMRwwGgYDVQQLExNJbmZvcm1hdGlvbiBTeXN0` +
			`hhx4dbgYKAAA7`

	sounds := []struct {
		sound string
		param map[string][]string
		valid bool
	}{
		{"CID:JOHNQPUBLIC.part8.19960229T080000.xyzMail@example.com", map[string][]string{"type": {"home"}}, true},
		{sound, map[string][]string{"language": {"de"}, "pref": {"55"}}, true},
		{sound, map[string][]string{"pref": {"55"}}, true},
		{"This is NoValidUri", map[string][]string{"pref": {"55"}}, false}, // invalid uri
		{"http://www.example.com/pub/sound/jqpublic.mp6", nil, true},
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for i, tupel := range sounds {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			if err := j.AddSound(tupel.sound, tupel.param); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v. (%w)", !tupel.valid, tupel.valid, err)
			}
		})
	}
}

func TestAddUID(t *testing.T) {
	uids := []struct {
		uid          string
		param        map[string][]string
		expectedType string
	}{
		{"urn:uuid:f81d4fae-7dec-11d0-a765-00a0c91e6bf6", nil, "uri"},
		{"Ein text", nil, "text"},
		{"", nil, "text"},
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for _, tupel := range uids {
		t.Run(fmt.Sprintf("%v", tupel.uid), func(t *testing.T) {
			j.RemoveAll("uid")
			if err := j.AddUID(tupel.uid, tupel.param); (err != nil) || (j.GetFirst("uid").Type != tupel.expectedType) {
				t.Errorf("Got %v, want %v. %w", j.GetFirst("uid").Type, tupel.expectedType, err)
			}
		})
	}
}

func TestAddClientPIDMap(t *testing.T) {
	cpms := []struct {
		cpm   []interface{}
		param map[string][]string
		valid bool
	}{
		{[]interface{}{"1", "urn:uuid:f81d4fae-7dec-11d0-a765-00a0c91e6bf6"}, nil, true},
		{[]interface{}{"1", "urn:uuid:f6"}, nil, true},
		{[]interface{}{"-1", "urn:uuid:f81d4fae-7dec-11d0-a765-00a0c91e6bf6"}, nil, false},
		{[]interface{}{"urn:uuid:f81d4fae-7dec-11d0-a765-00a0c91e6bf6"}, nil, false},
		{nil, nil, false},
		{[]interface{}{"urn:uuid:f81d4fae-7dec-11d0-a765-00a0c91e6bf6", "abc"}, nil, false},
		{[]interface{}{"1", "urn:uuid:f81d4fae-7dec-11d0-a765-00a0c91e6bf6"}, map[string][]string{"type": {"home"}}, false},
		{[]interface{}{"abc", "urn:uuid:f81d4fae-7dec-11d0-a765-00a0c91e6bf6"}, nil, false},
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for _, tupel := range cpms {
		t.Run(fmt.Sprintf("%v", tupel.cpm), func(t *testing.T) {
			j.RemoveAll("clientpidmap")
			if err := j.AddClientPIDMap(tupel.cpm, tupel.param); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v. %w", !tupel.valid, tupel.valid, err)
			}
		})
	}
}

func TestAddURL(t *testing.T) {
	urls := []struct {
		url   string
		param map[string][]string
		valid bool
	}{

		{"http://www.example.com/pub/photos/jqpublic.gif", map[string][]string{"type": {"home"}}, true},
		{"http://example.org/restaurant.french/~chezchic.html", map[string][]string{"pref": {"55"}}, true},
		{"This is NoValidUri", map[string][]string{"pref": {"55"}}, false}, // invalid uri
		{"http://www.example.com/pub/photos/jqpublic.gif", nil, true},
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for _, tupel := range urls {
		t.Run(fmt.Sprintf("%v", tupel.url), func(t *testing.T) {
			j.RemoveAll("url")
			if err := j.AddURL(tupel.url, tupel.param); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v. %w", !tupel.valid, tupel.valid, err)
			}
		})
	}
}
