package ldap

import (
	"testing"
)

func TestMember_setDn(t *testing.T) {
	member := new(Member)
	member.SetDN("CN=Zahphod Beeblebrox,DC=lübthen,DC=hagenow,DC=parchim,DC=elkm,DC=de")
	if member.Name != "Zahphod Beeblebrox" {
		t.Fail()
		t.Logf("%+v \n", member)
	}
}

func TestMember_setDnWithComma(t *testing.T) {
	member := new(Member)
	member.SetDN("CN=Zahphod\\, Beeblebrox,DC=lübthen,DC=hagenow,DC=parchim,DC=elkm,DC=de")
	if member.Name != "Zahphod, Beeblebrox" {
		t.Fail()
		t.Logf("%+v \n", member)
	}
}
