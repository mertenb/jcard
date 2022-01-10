package ldap

import (
	"regexp"
	"strings"
)

// Member represents a ldap entry
type Member struct {
	ID     int    `json:"id"`               // internal id
	Name   string `json:"name"`             // rdn human readable
	Path   string `json:"path"`             // path to rdn human readable
	DN     string `json:"dn"`               // distinguished name: original from LDAP
	Mobile string `json:"mobile,omitempty"` // optional mobilephone number
	Email  string `json:"email,omitempty"`  // optional email
}

var trenner = regexp.MustCompile("[,]?[A-Za-z]+=")         // remove 'CN=' and ',OU=' etc.
var rdnPattern = regexp.MustCompile("[A-Za-z]=.*[^\\\\],") // extract RDN

// SetDN set the distinguished name and fills other fields like readable name and path.
func (m *Member) SetDN(dn string) {
	m.DN = dn
	var dns []string
	if m.DN == "" {
		panic("DN is empty. Can not supplement other values.")
	}
	dns = trenner.Split(m.DN, -1)
	m.setPath(dns)
	m.setName(dns)
}

func (m *Member) setPath(dns []string) {
	if len(dns) < 2 {
		return
	}
	var path string
	if len(dns) > 2 {
		path = strings.Join(dns[2:len(dns)], " --- ")
	} else {
		path = dns[2]
	}
	m.Path = decode(path)
}

func (m *Member) setName(dns []string) {
	m.Name = decode(dns[1])
}

// remove comma escaping backslashes
func decode(s string) string {
	return strings.ReplaceAll(s, "\\,", ",")
}
