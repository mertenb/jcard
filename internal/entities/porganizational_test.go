package entities

import (
	"fmt"
	"testing"

	"github.com/mertenb/jcard/internal/test"
)

func TestAddTitle(t *testing.T) {
	titles := []struct {
		title string
		param map[string][]string
		valid bool
	}{
		{"Research Scientist", map[string][]string{"language": {"de"}, "pref": {"55"}}, true},
		{"Big winner", map[string][]string{"language": {"de"}, "pref": {"55"}, "pid": {"44"}, "altid": {"99"}}, true},
		{"Big looser", map[string][]string{"mediatype": {"application/xml-dtd", "application/zip"}}, false}, // no mediatype allowed
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for _, tupel := range titles {
		t.Run(tupel.title, func(t *testing.T) {
			if err := j.AddTitle(tupel.title, tupel.param); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v", !tupel.valid, tupel.valid)
			}
		})
	}
}

func TestAddRole(t *testing.T) {
	roles := []struct {
		role  string
		param map[string][]string
		valid bool
	}{
		{"Project Leader", map[string][]string{"language": {"de"}, "pref": {"55"}}, true},
		{"Project Leader", map[string][]string{"language": {"de"}, "pref": {"55"}, "pid": {"44"}, "altid": {"99"}}, true},
		{"Big looser", map[string][]string{"mediatype": {"application/xml-dtd", "application/zip"}}, false}, // no mediatype allowed
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for _, tupel := range roles {
		t.Run(tupel.role, func(t *testing.T) {
			if err := j.AddRole(tupel.role, tupel.param); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v", !tupel.valid, tupel.valid)
			}
		})
	}
}
func TestAddLogo(t *testing.T) {
	image :=
		`data:image/gif;base64,R0lGODdhMAAwAPAAAAAAAP///ywAAAAAMAAw` +
			`AAAC8IyPqcvt3wCcDkiLc7C0qwyGHhSWpjQu5yqmCYsapyuvUUlvONmOZtfzgFz` +
			`ByTB10QgxOR0TqBQejhRNzOfkVJ+5YiUqrXF5Y5lKh/DeuNcP5yLWGsEbtLiOSp` +
			`a/TPg7JpJHxyendzWTBfX0cxOnKPjgBzi4diinWGdkF8kjdfnycQZXZeYGejmJl` +
			`ZeGl9i2icVqaNVailT6F5iJ90m6mvuTS4OK05M0vDk0Q4XUtwvKOzrcd3iq9uis` +
			`F81M1OIcR7lEewwcLp7tuNNkM3uNna3F2JQFo97Vriy/Xl4/f1cf5VWzXyym7PH` +
			`hhx4dbgYKAAA7`

	ns := []struct {
		name  string
		logo  string
		param map[string][]string
		valid bool
	}{
		{"url", "http://www.example.com/pub/photos/jqpublic.gif", map[string][]string{"type": {"home"}}, true},
		{"image", image, map[string][]string{"language": {"de"}, "pref": {"55"}}, true},
		{"image2", image, map[string][]string{"pref": {"55"}}, true},
		{"invalidUri", "This is NoValidUri", map[string][]string{"pref": {"55"}}, false}, // invalid uri
		{"validurl_nil", "http://www.example.com/pub/photos/jqpublic.gif", nil, true},
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for _, tupel := range ns {
		t.Run(fmt.Sprintf("%v", tupel.name), func(t *testing.T) {
			if err := j.AddLogo(tupel.logo, tupel.param); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v. (%v)", !tupel.valid, tupel.valid, err)
			}
		})
	}
}

func TestAddOrg(t *testing.T) {
	orgs := []struct {
		org   []string
		param map[string][]string
		valid bool
	}{
		{[]string{"ABC, Inc.", "North American Division", "Marketing"}, map[string][]string{"language": {"de"}, "pref": {"55"}}, true},
		{[]string{"ABC, Inc.", "North American Division", "Marketing"}, map[string][]string{"language": {"de"}, "pref": {"55"}, "sort-as": {"ABC, Inc.", "North American Division", "Marketing"}}, true},
		{[]string{"ABC, Inc.", "North American Division", "Marketing"}, map[string][]string{"language": {"de"}, "pref": {"55"}, "sort-as": {"ABC, Inc.", "North American Division"}}, true},
		{[]string{"ABC, Inc.", "North American Division", "Marketing"}, map[string][]string{"unknown": {"de"}, "pref": {"55"}, "sort-as": {"ABC, Inc.", "North American Division"}}, false},
		{[]string{"ABC, Inc."}, map[string][]string{"sort-as": {"ABC, Inc.", "North American Division"}}, false}, // too many sort params
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for i, tupel := range orgs {
		t.Run(fmt.Sprintf("%v_nr_%v", tupel.org, i), func(t *testing.T) {
			if err := j.AddOrg(tupel.org, tupel.param); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v", !tupel.valid, tupel.valid)
			}
		})
	}
}

func TestAddMember(t *testing.T) {
	members := []struct {
		uri   string
		param map[string][]string
		valid bool
	}{
		{"mailto:subscriber1@example.com", nil, true},
		{"xmpp:subscriber2@example.com", nil, true},
		{"tel:+1-418-555-5555", nil, true},
		{"mailto:subscriber1@example.com", map[string][]string{"type": {"home"}}, false},
		{"urn:uuid:b8767877-b4a1-4c70-9acc-505d3819e519", nil, true},
		{"no_uri", nil, false},
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	j.AddKind("group")
	for i, tupel := range members {
		t.Run(fmt.Sprintf("%v_nr_%v", tupel.uri, i), func(t *testing.T) {
			if err := j.AddMember(tupel.uri, tupel.param); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v", !tupel.valid, tupel.valid)
			}
		})
	}
}

func TestAddRelated(t *testing.T) {
	related := []struct {
		uriOrText string
		param     map[string][]string
		valid     bool
	}{
		{"urn:uuid:f81d4fae-7dec-11d0-a765-00a0c91e6bf6", map[string][]string{"mediatype": {"application/xml-dtd", "application/zip"}}, true},
		{"urn:uuid:f81d4fae-7dec-11d0-a765-00a0c91e6bf6", map[string][]string{"language": {"de"}}, false},
		{"language is only a text attributue", map[string][]string{"language": {"de"}}, true},
		{"http://example.com/directory/jdoe.vcf", map[string][]string{"mediatype": {"application/xml-dtd", "application/zip"}, "type": {"sweetheart", "emergency"}}, true},
		{"text, no mediatype allowed", map[string][]string{"mediatype": {"application/xml-dtd", "application/zip"}, "type": {"sweetheart", "emergency"}}, false},
		{"http://example.com/directory/jdoe.vcf", nil, true},
		{"http://example.com/directory/jdoe.vcf", map[string][]string{"type": {"sweetheart", "emergency", "agent"}}, true},
		{"nourl only text", map[string][]string{"type": {"notype"}}, false},
		{"nourl only text", map[string][]string{"type": {"agent"}}, true},
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for i, tupel := range related {
		t.Run(fmt.Sprintf("%v_nr_%v", tupel.uriOrText, i), func(t *testing.T) {
			if err := j.AddRelated(tupel.uriOrText, tupel.param); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v, %v", !tupel.valid, tupel.valid, err)
			}
		})
	}
}
