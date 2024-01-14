package entities

import (
	"fmt"
	"testing"

	"github.com/mertenb/jcard/internal/test"
)

// ALLOWED: type-param / language-param / altid-param / pid-param / pref-param / any-param
func TestAddFn(t *testing.T) {
	fns := []struct {
		fn    string
		param map[string][]string
		valid bool
	}{
		{"Mr. John Q. Public\\, Esq.", map[string][]string{"language": {"de"}, "pref": {"55"}}, true},
		{"Sir Arthur Dent", map[string][]string{"language": {"de"}, "pref": {"55"}, "pid": {"44"}, "altid": {"99"}}, true},
		{"Mr. John Q. Public\\, Esq.", map[string][]string{"mediatype": {"application/xml-dtd", "application/zip"}, "pref": {"55"}}, false},
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for _, tupel := range fns {
		t.Run(tupel.fn, func(t *testing.T) {
			if err := j.AddFn(tupel.fn, tupel.param); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v", !tupel.valid, tupel.valid)
			}
		})
	}
}

func TestAddN(t *testing.T) {
	ns := []struct {
		n     interface{}
		param map[string][]string
		valid bool
	}{
		{[]string{"Stevenson", "John", "Philip", "Paul", "Dr.", "Jr.", "M.D.", "A.C.P."}, map[string][]string{"language": {"de"}}, true},
		{[]string{"van der Harten", "Rene", "J.", "Sir", "R.D.O.N."}, map[string][]string{"language": {"de"}, "sort-as": {"Harten", "Rene"}}, true},
		{[]string{"van der Harten"}, map[string][]string{"language": {"de"}, "sort-as": {"Harten"}}, true},
		{[]string{"van der Harten"}, map[string][]string{"language": {"de"}, "sort-as": {"Harten", "Rene"}}, false}, // one name, but two 'sort-as'
		{"van der Harten", map[string][]string{"language": {"de"}, "sort-as": {"Harten"}}, false},                   // name must be string array
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for _, tupel := range ns {
		t.Run(fmt.Sprintf("%v", tupel.n), func(t *testing.T) {
			if err := j.AddN(tupel.n, tupel.param); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v. (%v)", !tupel.valid, tupel.valid, err)
			}
		})
	}
}

func TestAddNickname(t *testing.T) {
	ns := []struct {
		nickname string
		param    map[string][]string
		valid    bool
	}{
		{"Seppel", map[string][]string{"language": {"de"}}, true},
		{"Schnaggy", map[string][]string{"language": {"de"}, "pref": {"55"}}, true},
		{"Welli", map[string][]string{"language": {"de"}, "sort-as": {"Harten", "Rene"}}, false}, // no 'sort-as' allowed
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for _, tupel := range ns {
		t.Run(fmt.Sprintf("%v", tupel.nickname), func(t *testing.T) {
			if err := j.AddNickname(tupel.nickname, tupel.param); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v. (%v)", !tupel.valid, tupel.valid, err)
			}
		})
	}
}

func TestAddPhoto(t *testing.T) {
	image :=
		`data:image/gif;base64,R0lGODdhMAAwAPAAAAAAAP///ywAAAAAMAAw` +
			`AAAC8IyPqcvt3wCcDkiLc7C0qwyGHhSWpjQu5yqmCYsapyuvUUlvONmOZtfzgFz` +
			`ByTB10QgxOR0TqBQejhRNzOfkVJ+5YiUqrXF5Y5lKh/DeuNcP5yLWGsEbtLiOSp` +
			`a/TPg7JpJHxyendzWTBfX0cxOnKPjgBzi4diinWGdkF8kjdfnycQZXZeYGejmJl` +
			`ZeGl9i2icVqaNVailT6F5iJ90m6mvuTS4OK05M0vDk0Q4XUtwvKOzrcd3iq9uis` +
			`F81M1OIcR7lEewwcLp7tuNNkM3uNna3F2JQFo97Vriy/Xl4/f1cf5VWzXyym7PH` +
			`hhx4dbgYKAAA7`

	ns := []struct {
		photo string
		param map[string][]string
		valid bool
	}{
		{"http://www.example.com/pub/photos/jqpublic.gif", map[string][]string{"type": {"home"}}, true},
		{image, map[string][]string{"language": {"de"}, "pref": {"55"}}, false}, // no 'language' allowed
		{image, map[string][]string{"pref": {"55"}}, true},
		{"This is NoValidUri", map[string][]string{"pref": {"55"}}, false}, // invalid uri
		{"http://www.example.com/pub/photos/jqpublic.gif", nil, true},
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for _, tupel := range ns {
		t.Run(fmt.Sprintf("%v", tupel.photo), func(t *testing.T) {
			if err := j.AddPhoto(tupel.photo, tupel.param); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v. (%v)", !tupel.valid, tupel.valid, err)
			}
		})
	}
}

func TestAddBDay(t *testing.T) {
	ns := []struct {
		bday  string
		param map[string][]string
		valid bool
	}{
		{"--04-12", map[string][]string{"altid": {"47"}}, true},
		{"1970-04-12", nil, true},
		{"1970-04-12", map[string][]string{"calscale": {"gregorian"}, "altid": {"4711"}}, true},
		{"1970-04-12", map[string][]string{"calscale": {"invalid"}, "altid": {"4711"}}, false},
		{"1970-04-12", map[string][]string{"pref": {"55"}}, false},
		{"1970-04-12T22:50:07", nil, true},
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for _, tupel := range ns {
		t.Run(fmt.Sprintf("%v", tupel.bday), func(t *testing.T) {
			j.RemoveAll("bday")
			if err := j.AddBDay(tupel.bday, tupel.param); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v. (%v)", !tupel.valid, tupel.valid, err)
			}
		})
	}

	j.RemoveAll("bday")
	err := j.AddBDay("circa 1800", nil)
	if err != nil {
		t.Error("adding bday failed.", err)
	}
	props := j.Get("bday")
	if len(props) != 1 {
		t.Errorf("One bday expectec, got %v", len(props))
	}
	prop := props[0]
	if prop.Name != "bday" || prop.Type != "text" {
		t.Errorf("unexpected property: %v", prop)
	}
	svalue := fmt.Sprintf("%v", prop.Value)
	if svalue != "circa 1800" {
		t.Errorf("unexpected value: %v", svalue)
	}

}

func TestAddGender(t *testing.T) {
	ns := []struct {
		gender interface{}
		param  map[string][]string
		valid  bool
	}{
		{"", nil, true},
		{"M", nil, true},
		{"F", nil, true},
		{"O", nil, true},
		{"N", nil, true},
		{"U", nil, true},
		{"X", nil, false},
		{"M", map[string][]string{"pref": {"55"}}, false},
		{[]string{"", "it's complicated"}, nil, true},
		{[]string{"", "it's complicated", "max 2 params allowed"}, nil, false},
		{[]string{"abc", "it's complicated"}, nil, false},
		{[]string{"U"}, nil, true},
		{[]string{}, nil, false},
		{45, nil, false},
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for _, tupel := range ns {
		t.Run(fmt.Sprintf("%v", tupel.gender), func(t *testing.T) {
			j.RemoveAll("gender")
			if err := j.AddGender(tupel.gender, tupel.param); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v. (%v)", !tupel.valid, tupel.valid, err)
			}
		})
	}
}
