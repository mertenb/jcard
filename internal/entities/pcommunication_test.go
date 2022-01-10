package entities

import (
	"testing"

	"github.com/mertenb/jcard/internal/test"
)

func TestValidTelNumberRegex(t *testing.T) {
	numbers := []struct {
		value string
		valid bool
	}{
		{"tel:+49170-449", true},
		{"tel:0049170-449", false},
		{"tel:445943", false},
		{"tel:445943;phone-context=test.de", true},
		{"tel:123;phone-context=+1-914-555", true},
		{"tel:123;phone-context=+1 914 555", false},
		{"tel:+49(0)171-432345", true},
		{"tel:+49.171.432345", true},
		{"tel:+49.171+432345", false},
	}

	for _, n := range numbers {
		t.Run(n.value, func(t *testing.T) {
			if validTelNumber.MatchString(n.value) != n.valid {
				t.Errorf("'%v' should be  %v", n.value, n.valid)
			}
		})
	}
}

func TestValidEMailRegex(t *testing.T) {
	mails := []struct {
		value string
		valid bool
	}{
		{"t@test.de", true},
		{"ABC@abc.de", true},
		{"test@üiu.de", false},
		{"test@test@test.de", false},
		{"noAt.de", false},
		{"no@dot", false},
		{"", false},
	}

	for _, n := range mails {
		t.Run(n.value, func(t *testing.T) {
			if validMail.MatchString(n.value) != n.valid {
				t.Errorf("'%v' should be %v", n.value, n.valid)
			}
		})
	}
}

func TestNormalizeTel(t *testing.T) {
	numbers := []struct {
		src  string
		dest string
	}{
		{"0049 1754 4432", "tel:+49-1754-4432"},
		{"+49 1754 4432", "tel:+49-1754-4432"},
		{"0049        1754    4432", "tel:+49-1754-4432"},
		{"863-1234;phone-context=+1-914-555", "tel:863-1234;phone-context=+1-914-555"},
		{"863 1234;phone-context = +1 914  555", "tel:863-1234;phone-context=+1914555"},
		{"7042 ;phone-context = example.com", "tel:7042;phone-context=example.com"},
		{"tel:0049 175 4323456", "tel:+49-175-4323456"},
		{"0049 (0)175 4323456", "tel:+49-175-4323456"},
		{"00704 ;phone-context = example.com", "tel:00704;phone-context=example.com"},
	}
	for _, tupel := range numbers {
		t.Run(tupel.src, func(t *testing.T) {
			nt := normalizeTel(tupel.src)
			if nt != tupel.dest {
				t.Errorf("Got %q, want %q", nt, tupel.dest)
			}
		})
	}
}

func TestAddPhone(t *testing.T) {
	j, _ := NewVCard(test.LoadFile("jcard/ok_example.json"))
	if err := j.AddTel("0049 3455 4432", map[string][]string{"type": {"voice", "cell"}, "pref": {"55"}}); err != nil {
		t.Errorf("No error expected : all parameters are valid. %w", err)
	}
	if nil == j.AddTel("0049 3455 4432", map[string][]string{"type": {"voice", "cell"}, "pref": {"1", "2"}}) {
		t.Errorf("Error expected : too many prefs")
	}
	if nil == j.AddTel("0049 3455 4432", map[string][]string{"type": {"voice", "cell"}, "pref": {"200"}}) {
		t.Errorf("Error expected : max pref exceeded (200>100)")
	}
	if nil == j.AddTel("0049 3455 4432", map[string][]string{"type": {"voice", "cell", "badType"}}) {
		t.Errorf("Error expected : unknown phone type")
	}
}

func TestAddEmail(t *testing.T) {
	j, _ := NewVCard(test.LoadFile("jcard/ok_example.json"))
	if err := j.AddEmail("test@test.de", map[string][]string{"type": {"voice", "cell"}, "pref": {"55"}}); err == nil {
		t.Error("Error expected: type 'voice' and 'cell' are not allowed for emails.")
	}
	if err := j.AddEmail("test@test.de", map[string][]string{"type": {"home", "work"}, "pref": {"55"}}); err != nil {
		t.Errorf("No error expected. %w", err)
	}
	if err := j.AddEmail("test@test.de", map[string][]string{"type": {"home"}, "pref": {"55"}}); err != nil {
		t.Errorf("No error expected. %w", err)
	}
	if err := j.AddEmail("test@test.de", map[string][]string{}); err != nil {
		t.Errorf("No error expected. %w", err)
	}
}

func TestAddLang(t *testing.T) {
	j, _ := NewVCard(test.LoadFile("jcard/ok_example.json"))
	if err := j.AddLang("DE", map[string][]string{"type": {"uselessandfalse"}, "pref": {"66"}}); err == nil {
		t.Error("Error expected: type 'uselessandfalse' are not allowed for lang.")
	}
	if err := j.AddLang("UNKNOWN-Lang", map[string][]string{}); err == nil {
		t.Error("Error expected: type 'uselessandfalse' are not allowed for lang.")
	}
	if err := j.AddLang("FR", map[string][]string{"type": {"work"}, "pref": {"55"}}); err != nil {
		t.Errorf("No error expected. %w", err)
	}
	if err := j.AddLang("zh-Hant", map[string][]string{"type": {"home"}, "pref": {"55"}}); err != nil {
		t.Errorf("No error expected. %w", err)
	}
	if err := j.AddLang("de-Qaaa", map[string][]string{}); err != nil {
		t.Errorf("No error expected. %w", err)
	}
}

func TestPhoneNumber(t *testing.T) {
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	if j.Validate(); j.hasErrors() {
		t.Errorf("Got unexpected validation error.")
	}
}

func TestAddImpp(t *testing.T) {
	impps := []struct {
		uri   string
		param map[string][]string
		valid bool
	}{
		{"xmpp:alice@example.com", map[string][]string{"mediatype": {"application/xml-dtd", "application/zip"}, "pref": {"55"}}, true},
		{"noUri", nil, false},
		{"http://MVSXX.COMPANY.COM:44455/CICSPLEXSM//TOXTETH/VIEW/EYUSTARTPROGRAM.TABULAR?FILTERC=1", nil, true},
		{"sms:00493445456", nil, true},
		{"icq:message?uin=4711", map[string][]string{"mediatype": {"application/xml-dtd", "application/zip"}}, true},
		{"icq:message?uin=4711", map[string][]string{"mediatype": {"application/xml-dtd", "invalidMediaType"}}, false},
		{"skype:myusername?call", nil, true},
		{"ymsgr:addfriend?test@test.de", map[string][]string{"type": {"work"}, "pref": {"55"}}, true},
	}
	j, _ := NewVCard(test.LoadFile("jcard/ok_tel.json"))
	for _, tupel := range impps {
		t.Run(tupel.uri, func(t *testing.T) {
			if err := j.AddImpp(tupel.uri, tupel.param); (err == nil) != tupel.valid {
				t.Errorf("Got %v, want %v", !tupel.valid, tupel.valid)
			}
		})
	}
}
