// Package ldap provides a simple ldap client to authenticate,
// retrieve basic information and groups for a user.
package ldap

import (
	"testing"
)

func Test_SearchMembers_All(t *testing.T) {
	var lc = Ldapconfig{
		ConParam: Connectionparam{
			Host: "localhost",
			Port: 10389,
		},
		SearchParam: Searchparam{
			BaseDN: "o=elkm",
			Filter: "(cn=%s)",
		},
	}
	lc.Connect()
	lc.RefreshCache()
	defer func() { lc.Close() }()
	members := lc.SearchMember("*")
	if len(members) < 2 {
		t.Errorf("Members: %v", members)
	}
}

func Test_SearchMembers_One(t *testing.T) {
	var lc = Ldapconfig{
		ConParam: Connectionparam{
			Host: "localhost",
			Port: 10389,
		},
		SearchParam: Searchparam{
			BaseDN: "o=elkm",
			Filter: "(cn=%s)",
		},
	}
	lc.Connect()
	lc.RefreshCache()
	defer func() { lc.Close() }()
	members := lc.SearchMember("arvin")
	if len(members) < 1 {
		t.Errorf("arvin not found. Members: %v", members)
	}
}

func Test_FindEntry(t *testing.T) {
	var lc = Ldapconfig{
		ConParam: Connectionparam{
			Host: "localhost",
			Port: 10389,
		},
	}
	lc.SearchParam.Attributes = []string{"mail"}
	lc.Connect()
	defer func() { lc.Close() }()
	entry, _ := lc.findEntryByDN("cn=Marvin,ou=neubrandenburg,ou=neustrelitz,ou=gemeinden,o=elkm")
	if "test@marvin.de" != entry.GetAttributeValue("mail") {
		t.Errorf("mail for marwin not found")
	}
	t.Logf("%v", entry)
}

func Test_searchEntry(t *testing.T) {
	var lc = Ldapconfig{
		ConParam: Connectionparam{
			Host: "localhost",
			Port: 10389,
		},
		SearchParam: Searchparam{
			BaseDN: "o=elkm",
			Filter: "(cn=%s)",
		},
	}
	lc.SearchParam.Attributes = []string{"mail"}
	lc.Connect()
	defer func() { lc.Close() }()
	members, _ := lc.searchEntries("*arvi*")
	if len(members) < 1 {
		t.Errorf("*arvi* not found.")
		t.Logf("%v", members)
	}

}

func Test_FindMemberByDn(t *testing.T) {
	var lc = Ldapconfig{
		ConParam: Connectionparam{
			Host: "localhost",
			Port: 10389,
		},
		SearchParam: Searchparam{
			BaseDN: "o=elkm",
			Filter: "(cn=%s)",
		},
	}
	lc.SearchParam.Attributes = []string{"mail"}
	lc.Connect()
	defer func() { lc.Close() }()
	member := lc.findMemberByDN("cn=Marvin,ou=neubrandenburg,ou=neustrelitz,ou=gemeinden,o=elkm")
	if "test@marvin.de" != member.Email {
		t.Errorf("mail for marwin not found")
	}
	t.Logf("%v", member)
}

func Test_Authenticate(t *testing.T) {
	var lc = Ldapconfig{
		AuthParam: Authenticationparam{
			BindDN:       "uid=admin,ou=system",
			BindPassword: "secret",
		},
		ConParam: Connectionparam{
			Host: "localhost",
			Port: 10389,
		},
	}

	lc.Connect()
	defer func() { lc.Close() }()
	err := lc.Authenticate()
	if err != nil {
		panic("error:")
	}
}
