package entities

// VCard methods related to phone moved to own file to shrink jcard.go

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/text/language"
)

var emailTypeInfo = VPropertyTypeInfo{
	Name:          Email,
	Cardinal:      Many,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.4.2",
	Validate:      ValidateEmail,
}

// AddEmail adds an email to VCard. Return false if some validation fail.
func (v *VCard) AddEmail(value string, params map[string][]string) error {
	property := &VCardProperty{
		Name:       Email,
		Type:       URI,
		Parameters: params,
		Value:      value,
		typeInfo:   emailTypeInfo,
	}
	return v.append(property)
}

var langTypeInfo = VPropertyTypeInfo{
	Name:          Lang,
	Cardinal:      Many,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.4.4",
	Validate:      ValidateLang,
}

// AddLang adds a language to VCard. Return false if some validation fail.
func (v *VCard) AddLang(value string, params map[string][]string) error {
	property := &VCardProperty{
		Name:       Lang,
		Type:       "language-tag",
		Parameters: params,
		Value:      value,
		typeInfo:   langTypeInfo,
	}
	return v.append(property)
}

var imppTypeInfo = VPropertyTypeInfo{
	Name:          IMPP,
	Cardinal:      Many,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.4.3",
	Validate:      ValidateImpp,
}

// AddImpp adds an URI for instant messaging. Return false if some validation fail.
func (v *VCard) AddImpp(value string, params map[string][]string) error {
	property := &VCardProperty{
		Name:       IMPP,
		Type:       URI,
		Parameters: params,
		Value:      value,
		typeInfo:   imppTypeInfo,
	}
	return v.append(property)
}

var telTypeInfo = VPropertyTypeInfo{
	Name:          Tel,
	Cardinal:      Many,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.4.1",
	Validate:      ValidateTel,
}

// AddTel adds a phone type to VCard. Return false if some validation fail.
func (v *VCard) AddTel(value string, params map[string][]string) error {
	property := &VCardProperty{
		Name:       Tel,
		Type:       URI,
		Parameters: params,
		Value:      normalizeTel(value),
		typeInfo:   telTypeInfo,
	}
	return v.append(property)
}

func normalizeTel(number string) string {
	number = strings.TrimSpace(number)
	number = strings.ReplaceAll(number, "tel:", "")
	s := strings.SplitN(number, ";", 2)
	if len(s) > 1 { // local number
		s[1] = strings.ReplaceAll(s[1], " ", "")
	} else { // global number
		if len(s[0]) > 1 && s[0][0:2] == "00" {
			s[0] = strings.Replace(s[0], "00", "+", 1)
			s[0] = strings.Replace(s[0], "(0)", "", 1)
		}
	}
	s[0] = strings.Join(strings.Fields(s[0]), " ") // replace multispaces with singlespace
	s[0] = strings.ReplaceAll(s[0], " ", "-")
	return fmt.Sprintf("tel:%v", strings.Join(s, ";"))
}

//https://tools.ietf.org/html/rfc3966#section-3
const regexTel = `^tel:((?:\+[\d().-]*\d[\d().-]*|[0-9A-F*#().-]*[0-9A-F*#][0-9A-F*#().-]*(?:;[a-z\d-]+(?:=(?:[a-z\d\[\]\/:&+$_!~*'().-]|%[\dA-F]{2})+)?)*;phone-context=(?:\+[\d().-]*\d[\d().-]*|(?:[a-z0-9]\.|[a-z0-9][a-z0-9-]*[a-z0-9]\.)*(?:[a-z]|[a-z][a-z0-9-]*[a-z0-9])))(?:;[a-z\d-]+(?:=(?:[a-z\d\[\]\/:&+$_!~*'().-]|%[\dA-F]{2})+)?)*(?:,(?:\+[\d().-]*\d[\d().-]*|[0-9A-F*#().-]*[0-9A-F*#][0-9A-F*#().-]*(?:;[a-z\d-]+(?:=(?:[a-z\d\[\]\/:&+$_!~*'().-]|%[\dA-F]{2})+)?)*;phone-context=\+[\d().-]*\d[\d().-]*)(?:;[a-z\d-]+(?:=(?:[a-z\d\[\]\/:&+$_!~*'().-]|%[\dA-F]{2})+)?)*)*)$`

var validTelNumber = regexp.MustCompile(regexTel)

//ValidateTel checks whether the syntax corresponds to the RFC specification.
func ValidateTel(p *VCardProperty) error {

	if b := validTelNumber.MatchString(fmt.Sprintf("%v", p.Value)); !b {
		return vCardError(fmt.Sprintf("Tel value is invalid (rfc3966#section3): %v", p.Value))
	}

	for k, v := range p.Parameters {
		switch k {
		case "type":
			allowed := map[string]struct{}{"home": {}, "work": {}, "text": {}, "voice": {}, "fax": {}, "cell": {}, "video": {}, "pager": {}, "textphone": {}}
			for _, s := range v {
				if _, contains := allowed[s]; !contains {
					return vCardError("unknown tel type: " + s)
				}
			}
		default:
			if err := validateDefaultParam(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

const regexpMail = `^[a-zA-Z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`

var validMail = regexp.MustCompile(regexpMail)

//ValidateEmail validate if property is a valid eMail property
func ValidateEmail(p *VCardProperty) error {
	if b := validMail.MatchString(fmt.Sprintf("%v", p.Value)); !b {
		return vCardError(fmt.Sprintf("Email value is invalid (https://tools.ietf.org/html/rfc6350#section-6.4.2): %v", p.Value))
	}
	return p.validateDefaultParam()
}

//ValidateLang validate if property is a valid language property
func ValidateLang(p *VCardProperty) error {

	lang := p.Value.(string)
	if _, err := language.Parse(lang); err != nil {
		return vCardError(fmt.Sprintf("Unknown language (https://tools.ietf.org/html/rfc5646): %v :%w", p.Value, err))
	}
	return p.validateDefaultParam()
}

//ValidateImpp validate if property is a valid impp uri property
func ValidateImpp(p *VCardProperty) error {

	impp := p.Value.(string)
	if !validateURI(&impp) {
		return vCardError(fmt.Sprintf("'%v' is not a valid uri", impp))
	}
	for k, v := range p.Parameters {
		switch k {
		case "mediatype":
			if err := validateMediatypeParam(v); err != nil {
				return err
			}
		default:
			if err := validateDefaultParam(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

// LANG-param = "VALUE=language-tag" / pid-param / pref-param / altid-param / type-param / any-param
// LANG-value = Language-Tag

// Tel returns the VCard's first (voice) telephone number.
// Returns empty string if the VCard contains no suitable telephone number.
func (v *VCard) Tel() string {
	properties := v.Get("tel")

	for _, p := range properties {
		isVoice := false

		if types, ok := p.Parameters["type"]; ok {
			for _, t := range types {
				if t == "voice" {
					isVoice = true
					break
				}
			}
		} else {
			isVoice = true
		}

		if isVoice && len(p.Values()) > 0 {
			return (p.Values())[0]
		}
	}
	return ""
}

// Fax returns the VCard's first fax number.
//
// Returns empty string if the VCard contains no fax number.
func (v *VCard) Fax() string {
	properties := v.Get("tel")

	for _, p := range properties {
		if types, ok := p.Parameters["type"]; ok {
			for _, t := range types {
				if t == "fax" {
					if len(p.Values()) > 0 {
						return (p.Values())[0]
					}
				}
			}
		}
	}

	return ""
}

// Email returns the VCard's first email address. Empty string if the VCard contains no email addresses.
func (v *VCard) Email() string {
	return v.getFirstPropertySingleString("email")
}

// Lang returns the first language. Empry if no lang attribute exists.
func (v *VCard) Lang() string {
	return v.getFirstPropertySingleString("lang")
}

// Impp returns the first instant messaging entry. Empty if no lang attribute exists.
func (v *VCard) Impp() string {
	return v.getFirstPropertySingleString("impp")
}
