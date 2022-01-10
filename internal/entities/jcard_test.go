package entities

// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.
// stolen from https://github.com/openrdap/rdap/blob/master/vcard_test.go

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/mertenb/jcard/internal/ldap"
	"github.com/mertenb/jcard/internal/test"
)

func TestVCardErrors(t *testing.T) {
	filenames := []string{
		"jcard/error_invalid_json.json",
		"jcard/error_bad_top_type.json",
		"jcard/error_bad_vcard_label.json",
		"jcard/error_bad_properties_array.json",
		"jcard/error_bad_property_size.json",
		"jcard/error_bad_property_name.json",
		"jcard/error_bad_property_type.json",
		"jcard/error_bad_property_parameters.json",
		"jcard/error_bad_property_parameters_2.json",
		"jcard/error_bad_property_nest_depth.json",
	}

	for _, filename := range filenames {
		j, err := NewVCard(test.LoadFile(filename))

		if j != nil || err == nil {
			t.Errorf("jCard with error unexpectedly parsed %s %v %s\n", filename, j, err)
		}
	}
}

func TestValidation(t *testing.T) {
	filenames := []string{
		"jcard/error_missing_version.json",
		"jcard/invalid_version.json",
		"jcard/invalid_member_kind_bad.json",
		"jcard/invalid_member_kind_nil.json",
		"jcard/error_pid_missing_clientpidmap.json",
		"jcard/error_pid_missing_map_2.json",
	}

	for _, filename := range filenames {
		j, err := NewVCard(test.LoadFile(filename))
		if err != nil {
			t.Errorf("VCard could not created from file %v. Error: %w", filename, err)
		} else {
			j.Validate()
			if !j.hasErrors() {
				t.Errorf("invalid jCard validated %s %v %s\n", filename, j, err)
			}
		}
	}
}

func TestVCardExample(t *testing.T) {
	j, err := NewVCard(test.LoadFile("jcard/ok_example.json"))
	if j == nil || err != nil {
		t.Errorf("jCard parse failed %v %s\n", j, err)
	}

	numProperties := 17
	if len(j.Properties) != numProperties {
		t.Errorf("Got %d properties expected %d", len(j.Properties), numProperties)
	}

	propArray := j.Get(Version)
	if len(propArray) != 1 {
		t.Errorf("Version is only allowed once.")
	}
	if err := propArray[0].Validate(); err != nil {
		t.Errorf("Version error: %w", err)
	}

	expectedN := &VCardProperty{
		Name:       "n",
		Parameters: make(map[string][]string),
		Type:       "text",
		Value:      []interface{}{"Perreault", "Simon", "", "", []interface{}{"ing. jr", "M.Sc."}},
		typeInfo:   nTypeInfo,
	}

	expectedFlatN := []string{
		"Perreault",
		"Simon",
		"",
		"",
		"ing. jr",
		"M.Sc.",
	}

	// deepequal does not work (function values are not compareable); use string representations
	s1 := fmt.Sprintf("%v", j.Get("n")[0])
	s2 := fmt.Sprintf("%v", expectedN)
	ti1 := fmt.Sprintf("%v", j.Get("n")[0].typeInfo)
	ti2 := fmt.Sprintf("%v", expectedN.typeInfo)
	if !((s1 == s2) && (ti1 == ti2)) {
		t.Errorf("n field incorrect. \n Got: %v \n Exp: %v \n \n Got: %v \n Exp: %v", j.Get("n")[0], expectedN, j.Get("n")[0].typeInfo, expectedN.typeInfo)
	}

	if !reflect.DeepEqual(j.Get("n")[0].Values(), expectedFlatN) {
		t.Errorf("n flat value incorrect")
	}

	expectedTel0 := &VCardProperty{
		Name:       "tel",
		Parameters: map[string][]string{"type": []string{"work", "voice"}, "pref": []string{"1"}},
		Type:       "uri",
		Value:      "tel:+1-418-656-9254;ext=102",
	}

	s1 = fmt.Sprintf("%v", j.Get("tel")[0])
	s2 = fmt.Sprintf("%v", expectedTel0)
	if !(s1 == s2) {
		t.Errorf("tel[0] field incorrect")
	}
}

func TestRemoveAll(t *testing.T) {
	j, _ := NewVCard(test.LoadFile("jcard/ok_example.json"))

	if len(j.Get("tel")) != 2 {
		t.Errorf("Expected 2 phone properties but got %v", len(j.Get("tel")))
	}

	b := j.RemoveAll("tel")

	if !b {
		t.Error("No error while removeAll(tel) expected.")
	}

	if len(j.Get("tel")) != 0 {
		t.Errorf("Expected 0 phone properties but got %v", len(j.Get("tel")))
	}
	b = j.RemoveAll("tel")

	if b {
		t.Error("removeAll(tel) should return false (no tel anymore).")
	}

	b = j.RemoveAll("notexists")

	if b {
		t.Error("removeAll(notexists) should return false.")
	}

}

func TestNewVcardFromLdap(t *testing.T) {
	person := new(ldap.InetOrgPerson)
	person.DisplayName = "Arthur Dent"
	person.Mail = "arthur@dent.de"
	person.Mobile = "+49 170 1895281"
	person.TelephoneNumber = "0049 3866 307"

	j, err := NewVCardFromLdap(person)
	if j.Email() != "arthur@dent.de" || j.Tel() != "tel:+49-170-1895281" || j.Fn() != "Arthur Dent" {
		t.Errorf("min. one of email,name,tel is incorrect: %v %s\n", j, err)
	}
}

func TestVCardMixedDatatypes(t *testing.T) {
	j, err := NewVCard(test.LoadFile("jcard/ok_mixed.json"))
	if j == nil || err != nil {
		t.Errorf("jCard parse failed %v %s\n", j, err)
	}

	expectedMixed := &VCardProperty{
		Name:       "mixed",
		Parameters: make(map[string][]string),
		Type:       "text",
		Value:      []interface{}{"abc", true, float64(42), nil, []interface{}{"def", false, float64(43)}},
	}

	expectedFlatMixed := []string{
		"abc",
		"true",
		"42",
		"",
		"def",
		"false",
		"43",
	}

	if !reflect.DeepEqual(j.Get("mixed")[0], expectedMixed) {
		t.Errorf("mixed field incorrect")
	}

	flattened := j.Get("mixed")[0].Values()
	if !reflect.DeepEqual(flattened, expectedFlatMixed) {
		t.Errorf("mixed flat value incorrect %v", flattened)
	}
}

func TestVCardQuickAccessors(t *testing.T) {
	j, err := NewVCard(test.LoadFile("jcard/ok_example.json"))
	if j == nil || err != nil {
		t.Errorf("jCard parse failed %v %s\n", j, err)
	}

	got := []string{
		j.Fn(),
		j.POBox(),
		j.ExtendedAddress(),
		j.StreetAddress(),
		j.Locality(),
		j.Region(),
		j.PostalCode(),
		j.Country(),
		j.Tel(),
		j.Fax(),
		j.Email(),
	}

	expected := []string{
		"Simon Perreault",
		"",
		"Suite D2-630",
		"2875 Laurier",
		"Quebec",
		"QC",
		"G1V 2M2",
		"Canada",
		"tel:+1-418-656-9254;ext=102",
		"",
		"simon.perreault@viagenie.ca",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Got %v expected %v\n", got, expected)
	}
}
func TestPid(t *testing.T) {
	vcard, err := NewVCard(test.LoadFile("jcard/ok_pid.json"))
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if vcard.Validate(); vcard.hasErrors() {
		t.Errorf("Error: %w", vcard.Errors)
	}

	if vcard, err = NewVCard(test.LoadFile("jcard/error_bad_pid.json")); err != nil {
		t.Errorf("Error: %v", err)
	}
}

type testRecord struct {
	value string
	valid bool
}

var dateStruct = []testRecord{
	{"1985-04-12", true},
	{"1985-04", true},
	{"1985", true},
	{"--04-12", true},
	{"--04", true},
	{"---12", true},
	{"198504", false},
	{"invalid", false},
	{"85-04-12", false},
	{"19850412", false},
	{"19850412", false},
}

// (https://stackoverflow.com/questions/522251/whats-the-difference-between-iso-8601-and-rfc-3339-date-formats)
func TestDate(t *testing.T) {
	for _, tupel := range dateStruct {
		t.Run(fmt.Sprintf("%v", tupel.value), func(t *testing.T) {
			ok := validateDate(&tupel.value)
			if ok != tupel.valid {
				t.Errorf("Got %v, want %v.", !tupel.valid, tupel.valid)
			}
		})
	}
}

var timeStruct = []testRecord{
	{"23:20:50", true},
	{"23:20:50.44", true},
	{"23:20", true},
	{"23", true},
	{"-20:50", true},
	{"-20:50.6543", true},
	{"-20", true},
	{"--50", true},
	{"23.5", false},
	{"23,5", false},
	{"-20:50,6543", false},
	{"invalid", false},
}

func TestTime(t *testing.T) {
	for _, tupel := range timeStruct {
		t.Run(fmt.Sprintf("%v", tupel.value), func(t *testing.T) {
			ok := validateTime(&tupel.value)
			if ok != tupel.valid {
				t.Errorf("Got %v, want %v.", !tupel.valid, tupel.valid)
			}
		})
	}
}

var dateTimeStruct = []testRecord{
	{"1985-04-12T23:20:50", true},
	{"1985-04-12T23:20:50Z", true},
	{"1985-04-12T23:20:50+04:00", true},
	{"1985-04-12T23:20:50+04", true},
	{"1985-04-12T23:20", true},
	{"1985-04-12T23", true},
	{"--04-12T23:20", true},
	{"--04T23:20", true},
	{"---12T23:20", true},
	{"--04-12T23:20", true},
	{"--04T23", true},
	{"1985-04-12T232050", false},
	{"1985-4-12T23:20:50T23:20:50", false},
	{"195-04-12T23:20:50T23:20:50", false},
	{"----12T23:20", false},
	{"1985-04-12", false},
	{"23:20", false},
	{"invalid", false},
}

func TestDateTime(t *testing.T) {
	for _, tupel := range dateTimeStruct {
		t.Run(fmt.Sprintf("%v", tupel.value), func(t *testing.T) {
			ok := validateDateTime(&tupel.value)
			if ok != tupel.valid {
				t.Errorf("Got %v, want %v.", !tupel.valid, tupel.valid)
			}
		})
	}
}

// TODO
func TestDateAndOrTime(t *testing.T) {
	testRecordsArray := [][]testRecord{dateStruct, timeStruct, dateTimeStruct}
	for _, testRecords := range testRecordsArray {
		for _, record := range testRecords {
			if record.valid { // invalid conditions from date, time and datetime can be valid here so check only the valid ones.
				t.Run(fmt.Sprintf("%v", record.value), func(t *testing.T) {
					ok := validateDateAndOrTime(&record.value)
					if !ok {
						t.Error("Should be valid.")
					}
				})
			}
		}
	}
}

// https://tools.ietf.org/html/rfc6350#section-4.7
func TestUtc(t *testing.T) {
	utcStruct := []testRecord{
		{"+1259", true},
		{"-1200", true},
		{"Z", true},
		{"-12:00", false},
		{"1200", false},
		{"", false},
		{"invalid", false},
		{"14:00", false},
		{"1301", false},
	}
	for _, tupel := range utcStruct {
		t.Run(fmt.Sprintf("%v", tupel.value), func(t *testing.T) {
			ok := validateUtcOffset(&tupel.value)
			if ok != tupel.valid {
				t.Errorf("Got %v, want %v.", !tupel.valid, tupel.valid)
			}
		})
	}
}

func TestTimestamp(t *testing.T) {
	utcStruct := []testRecord{
		{"19961022T140000", true},
		{"19961022T140000Z", true},
		{"19961022T140000-05", true},
		{"19961022T140000-0500", false},
		{"1200", false},
		{"", false},
		{"invalid", false},
		{"14:00", false},
		{"1301", false},
	}
	for _, tupel := range utcStruct {
		t.Run(fmt.Sprintf("%v", tupel.value), func(t *testing.T) {
			ok := validateTimestamp(&tupel.value)
			if ok != tupel.valid {
				t.Errorf("Got %v, want %v.", !tupel.valid, tupel.valid)
			}
		})
	}
}
