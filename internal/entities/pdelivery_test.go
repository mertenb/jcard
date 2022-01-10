package entities

import (
	"fmt"
	"testing"

	"github.com/mertenb/jcard/internal/test"
)

func TestAddAdr(t *testing.T) {
	adrs := []struct {
		adr   []string
		param map[string][]string
		valid bool
	}{
		{[]string{"", "", "123 Main Street", "Any Town", "CA", "91921-1234", "U.S.A."}, nil, true},
		{[]string{"", "", "", "", "", "", ""}, map[string][]string{"label": {"123 Maple Ave\nSuite 901\nVancouver BC\nA1B 2C9\nCanada"}}, true},
		{[]string{"", "", "", "", "", "", ""}, map[string][]string{"label": {"123 Maple Ave\nSuite 901\nVancouver BC\nA1B 2C9\nCanada"}, "pref": {"55"}}, true},
		{[]string{"", "", "", "", "", "", ""}, map[string][]string{"label": {"123 Maple Ave\nSuite 901\nVancouver BC\nA1B 2C9\nCanada"}, "unknown": {"55"}}, false},
		{[]string{"", "", "", ""}, nil, false},
		{[]string{"", "", "", "", "", "", "", "", "", "", "", ""}, nil, false},
		{[]string{"SHOULD BE", "EMPTY, BUT NEED NOT", "123 Main Street", "Any Town", "CA", "91921-1234", "U.S.A."}, nil, true},
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for _, tupel := range adrs {
		t.Run(fmt.Sprintf("%v", tupel.adr), func(t *testing.T) {
			if err := j.AddAdr(tupel.adr, tupel.param); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v", !tupel.valid, tupel.valid)
			}
		})
	}
}
