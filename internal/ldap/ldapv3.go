// Package ldap provides a simple ldap client to authenticate,
// and simple ldap search.
package ldap

import (
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

var membercache []Member

// Connectionparam contains values for connection to ldap server.
type Connectionparam struct {
	Host               string // e.g. "localhost"
	ServerName         string
	Port               int // e.g. 10389
	InsecureSkipVerify bool
	UseSSL             bool
	SkipTLS            bool
	ClientCertificates []tls.Certificate // Adding client certificates

}

// Authenticationparam contains values to authenticate on ldap
type Authenticationparam struct {
	BindDN       string // e.g. "uid=admin,ou=system"
	BindPassword string // e.g. "secret"
}

// Searchparam contains values for ldapsearch.
type Searchparam struct {
	Attributes []string // attributes to return (note: 'dn' is an extra field)
	BaseDN     string   // start search at BaseDN (e.g. "ou=kkv,o=elkm")
	Filter     string   // e.g. "(memberUid=%s)" or "(uid=%s)"
}

// Ldapconfig represent the handled ldap.
type Ldapconfig struct {
	ConParam    Connectionparam
	AuthParam   Authenticationparam
	SearchParam Searchparam
	conn        *ldap.Conn
}

//RefreshCache put all ldap entries dn's into a static cache.
func (lc *Ldapconfig) RefreshCache() {
	membercache = nil
	lc.Connect()
	defer func() { lc.Close() }()
	results, error := lc.searchEntries("*")
	if error != nil {
		panic(error)
	}
	membercache = make([]Member, len(results), len(results))
	for i, entry := range results {
		membercache[i] = toMember(entry)
		membercache[i].ID = i
	}
}

// Connect connects to the ldap backend.
func (lc *Ldapconfig) Connect() error {
	if lc.conn == nil {
		var l *ldap.Conn
		var err error
		address := fmt.Sprintf("%s:%d", lc.ConParam.Host, lc.ConParam.Port)
		if !lc.ConParam.UseSSL {
			l, err = ldap.Dial("tcp", address)
			if err != nil {
				return err
			}

			// Reconnect with TLS
			if !lc.ConParam.SkipTLS {
				err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
				if err != nil {
					return err
				}
			}
		} else {
			config := &tls.Config{
				InsecureSkipVerify: lc.ConParam.InsecureSkipVerify,
				ServerName:         lc.ConParam.ServerName,
			}
			if lc.ConParam.ClientCertificates != nil && len(lc.ConParam.ClientCertificates) > 0 {
				config.Certificates = lc.ConParam.ClientCertificates
			}
			l, err = ldap.DialTLS("tcp", address, config)
			if err != nil {
				return err
			}
		}

		lc.conn = l

	}
	return nil
}

// Close closes the ldap backend connection.
func (lc *Ldapconfig) Close() {
	if lc.conn != nil {
		lc.conn.Close()
		lc.conn = nil
	}
}

// Authenticate authenticates the user against the ldap backend.
// Return true if success.
func (lc *Ldapconfig) Authenticate() error {
	if err := lc.conn.Bind(lc.AuthParam.BindDN, lc.AuthParam.BindPassword); err != nil {
		return err
	}
	return nil
}

// SearchEntries deliver all memebers matching the given filter in 'RDN'.
func (lc *Ldapconfig) searchEntries(filter string) ([]*ldap.Entry, error) {
	searchRequest := ldap.NewSearchRequest(
		lc.SearchParam.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(lc.SearchParam.Filter, filter), // (make "cn=%s","*xyz*" to "cn=*xyz*")
		lc.SearchParam.Attributes,
		nil,
	)
	sr, err := lc.conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	return sr.Entries, nil
}

// findEntryByDN deliver exact one entry matching the given DN
func (lc *Ldapconfig) findEntryByDN(dn string) (*ldap.Entry, error) {
	lc.Connect()
	lc.SearchParam.Attributes = []string{"mail", "mobile"}
	searchRequest := ldap.NewSearchRequest(
		dn, //lc.SearchParam.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=*)",
		lc.SearchParam.Attributes,
		nil,
	)
	sr, err := lc.conn.Search(searchRequest)
	if err != nil || len(sr.Entries) != 1 {
		if len(sr.Entries) != 1 {
			err = fmt.Errorf("One result expected, but found %d. (%s)", len(sr.Entries), dn)
		}
		return nil, err
	}
	return sr.Entries[0], nil
}

// SearchMember from given ldap-connection matching given pattern.
func (lc *Ldapconfig) SearchMember(pattern string) []Member {
	if "*" == pattern {
		return membercache
	}
	var result []Member
	for _, member := range membercache {
		if strings.Contains(member.DN, pattern) {
			result = append(result, member)
		}
	}
	return result
}

// FindMember finds the member by given unique id
func (lc *Ldapconfig) FindMember(id int) Member {
	member := membercache[id]
	member = lc.findMemberByDN(member.DN)
	member.ID = id
	return member
}

// findMemberByDN finds the member by given unique DN
func (lc *Ldapconfig) findMemberByDN(dn string) Member {
	entry, _ := lc.findEntryByDN(dn)
	members := ToMembers([]*ldap.Entry{entry})
	member := members[0]
	member.Email = extract("mail", entry.Attributes)
	member.Mobile = extract("mobile", entry.Attributes)
	return member
}

func extract(attribName string, attribs []*ldap.EntryAttribute) string {
	for _, item := range attribs {
		if item.Name == attribName {
			return strings.Join(item.Values, " --- ")
		}
	}
	return ""
}

// ToMembers convert technical Entry to UI-friendly Member
func ToMembers(entries []*ldap.Entry) []Member {
	members := make([]Member, len(entries))
	for i, entry := range entries {
		members[i] = Member{ID: i}
		members[i].SetDN(entry.DN)
	}
	return members
}

// toMember convert technical Entry to UI-friendly Member
func toMember(entry *ldap.Entry) Member {
	var member Member
	member.SetDN(entry.DN)
	return member
}
