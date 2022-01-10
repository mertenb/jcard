package ldap

// InetOrgPerson represents a simplified ldap inetorgperson. Only attributes which are used by elkm will be defined.
// see https://tools.ietf.org/html/rfc2798
type InetOrgPerson struct {
	// MANDATORY, DESCRIPTION
	DisplayName              string `json:"displayName"`              // Y, e.g. used in oneline summary list
	Sn                       string `json:"sn"`                       // Y, surename
	Cn                       string `json:"cn"`                       // Y, mandatory common name, persons full name;
	PersonalTitle            string `json:"personalTitle"`            // N, e.g. 'Prof. Dr'
	Ou                       string `json:"ou"`                       // N, elkm;wismar;wismar;dorf mecklenburg;groß stieten
	PostalAddress            string `json:"postalAddress"`            // max 6 lines a 30 character
	Mail                     string `json:"mail"`                     // 7bit IA5 Character Set email
	TelephoneNumber          string `json:"telephoneNumber"`          // printable String syntax
	FacsimileTelephoneNumber string `json:"facsimileTelephoneNumber"` // printable String syntax
	Mobile                   string `json:"mobile"`                   // printable String syntax
	UserPassword             string `json:"userPassword"`             // Octet String Syntax, 128 char
	UniqueIdentifier         string `json:"uniqueIdentifier"`         // ELKM00001 - UUID??
	EmployeeType             string `json:"employeeType"`             // KIDAT person Type (Pastor, etc)
}
