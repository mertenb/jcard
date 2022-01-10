package entities

import (
	"fmt"
)

// Fn returns the VCard's formatted name (fn). e.g. "John Smith". https://tools.ietf.org/html/rfc6350#section-6.2.1
func (v *VCard) Fn() string {
	return v.getFirstPropertySingleString("fn")
}

// N returns the components of the name of the object the vCard represents. https://tools.ietf.org/html/rfc6350#section-6.2.2
func (v *VCard) N() string {
	return v.getFirstPropertySingleString("n")
}

var fnTypeInfo = VPropertyTypeInfo{
	Name:          Fn,
	Cardinal:      OneOrMany,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.2.1",
	Validate:      ValidateFn,
}

// AddFn adds the name to VCard. Minimum one Fn is required. (= DN in LDAP). https://tools.ietf.org/html/rfc6350#section-6.2.1
func (v *VCard) AddFn(value string, params map[string][]string) error {
	property := &VCardProperty{
		Name:       Fn,
		Type:       Text,
		Parameters: params,
		Value:      value,
		typeInfo:   fnTypeInfo,
	}
	return v.append(property)
}

var nickNameTypeInfo = VPropertyTypeInfo{
	Name:          Nickname,
	Cardinal:      Many,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.2.3",
	Validate:      ValidateNickname,
}

// AddNickname The nickname is the descriptive name given instead of or in addition to the one belonging to the object the vCard
// represents.  It can also be used to specify a familiar form of a proper name specified by the FN or N properties.
// https://tools.ietf.org/html/rfc6350#section-6.2.3
func (v *VCard) AddNickname(value string, params map[string][]string) error {
	property := &VCardProperty{
		Name:       Nickname,
		Type:       Text,
		Parameters: params,
		Value:      value,
		typeInfo:   nickNameTypeInfo,
	}
	return v.append(property)
}

var nTypeInfo = VPropertyTypeInfo{
	Name:          N,
	Cardinal:      OneOrZero,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.2.2",
	Validate:      ValidateN,
}

// AddN to specify the components of the name of the object the vCard represents.
// https://tools.ietf.org/html/rfc6350#section-6.2.2
func (v *VCard) AddN(value interface{}, params map[string][]string) error {
	property := &VCardProperty{
		Name:       N,
		Type:       Text,
		Parameters: params,
		Value:      value,
		typeInfo:   nTypeInfo,
	}

	return v.append(property)
}

var photoTypeInfo = VPropertyTypeInfo{
	Name:          Photo,
	Cardinal:      Many,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.2.4",
	Validate:      ValidatePhoto,
}

// AddPhoto specify an image or photograph information that annotates some aspect of the object the vCard represents.
// https://tools.ietf.org/html/rfc6350#section-6.2.4
func (v *VCard) AddPhoto(value interface{}, params map[string][]string) error {
	property := &VCardProperty{
		Name:       Photo,
		Type:       URI,
		Parameters: params,
		Value:      value,
		typeInfo:   photoTypeInfo,
	}

	return v.append(property)
}

var bdayTypeInfo = VPropertyTypeInfo{
	Name:          BDay,
	Cardinal:      OneOrZero,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.2.5",
	Validate:      ValidateBDay,
}

// AddBDay specify the birth date of the object the vCard represents.
// (https://tools.ietf.org/html/rfc6350#section-6.2.5)
func (v *VCard) AddBDay(value interface{}, params map[string][]string) error {
	svalue := value.(string)
	var stype string
	if validateDateAndOrTime(&svalue) {
		stype = "date-and-or-time"
	} else {
		stype = Text
	}
	property := &VCardProperty{
		Name:       BDay,
		Type:       stype,
		Parameters: params,
		Value:      value,
		typeInfo:   bdayTypeInfo,
	}

	return v.append(property)
}

var anniversaryTypeInfo = VPropertyTypeInfo{
	Name:          Anniversary,
	Cardinal:      OneOrZero,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.2.6",
	Validate:      ValidateAnniversary,
}

// AddAnniversary the date of marriage, or equivalent, of the object the vCard represents.
// (https://tools.ietf.org/html/rfc6350#section-6.2.6)
func (v *VCard) AddAnniversary(value interface{}, params map[string][]string) error {
	svalue := value.(string)
	var stype string
	if validateDateAndOrTime(&svalue) {
		stype = "date-and-or-time"
	} else {
		stype = Text
	}
	property := &VCardProperty{
		Name:       Anniversary,
		Type:       stype,
		Parameters: params,
		Value:      value,
		typeInfo:   anniversaryTypeInfo,
	}

	return v.append(property)
}

var genderTypeInfo = VPropertyTypeInfo{
	Name:          Gender,
	Cardinal:      OneOrZero,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.2.7",
	Validate:      ValidateGender,
}

// AddGender https://tools.ietf.org/html/rfc6350#section-6.2.7
func (v *VCard) AddGender(value interface{}, params map[string][]string) error {
	property := &VCardProperty{
		Name:       Gender,
		Type:       Text,
		Parameters: params,
		Value:      value,
		typeInfo:   genderTypeInfo,
	}
	return v.append(property)
}

//ValidateFn check if fn fulfill the requirenments of https://tools.ietf.org/html/rfc6350#section-6.2.1
func ValidateFn(p *VCardProperty) error {

	for k, v := range p.Parameters {
		switch k {
		case "language":
			if err := validateLanguageParam(v); err != nil {
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

var allowedGender = map[string]struct{}{"": {}, "M": {}, "F": {}, "O": {}, "N": {}, "U": {}}

//ValidateGender check if fn fulfill the requirenments of https://tools.ietf.org/html/rfc6350#section-6.2.7
func ValidateGender(p *VCardProperty) error {

	if p.Parameters != nil && len(p.Parameters) > 0 {
		return vCardError(fmt.Sprintf("gender must not have parameters, but found %v", p.Parameters))
	}

	switch p.Value.(type) {
	case []string:
		values := p.Value.([]string)
		switch len(values) {
		case 1, 2:
			g := fmt.Sprintf("%v", values[0]) // second item is of type text and always valid.
			if _, ok := allowedGender[g]; !ok {
				return vCardError(fmt.Sprintf("Unknown gender: %v is none of %v", g, allowedGender))
			}
		default:
			return vCardError("Gender has too many values. 1 or 2 expected.")
		}
	case string:
		value := p.Value.(string)
		g := fmt.Sprintf("%v", value)
		if _, ok := allowedGender[g]; !ok {
			return vCardError(fmt.Sprintf("Unknown gender: %v is none of %v", g, allowedGender))
		}

	default:
		return vCardError(fmt.Sprintf("Unknown type of pValue: string or []string expected, but got %T", p.Value))
	}
	return nil
}

//ValidateN check if n fulfill the requirenments of https://tools.ietf.org/html/rfc6350#section-6.2.2
func ValidateN(p *VCardProperty) error {

	for k, v := range p.Parameters {
		switch k {
		case "language":
			if err := validateLanguageParam(v); err != nil {
				return err
			}
		case "sort-as":
			if err := validateSortAsParam(p.Value, v); err != nil {
				return err
			}

		case "altid":
			if err := validateAltidParam(v); err != nil {
				return err
			}
		default:
			return vCardError("unknown parameter type: " + k)
		}

	}
	return nil
}

//ValidateNickname check if nickname fulfill the requirenments of https://tools.ietf.org/html/rfc6350#section-6.2.3
func ValidateNickname(p *VCardProperty) error {

	for k, v := range p.Parameters {
		switch k {
		case "language":
			if err := validateLanguageParam(v); err != nil {
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

//ValidatePhoto https://tools.ietf.org/html/rfc6350#section-6.2.4
func ValidatePhoto(p *VCardProperty) error {

	photo := p.Value.(string)
	if !validateURI(&photo) {
		return vCardError(fmt.Sprintf("'%v' is not a valid uri", photo))
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

// summarize validation for BDate and Anniversary
func validateAnniversaryAndBDate(p *VCardProperty, name string) error {

	if p.Type == "date-and-or-time" {
		birthdate := p.Value.(string)
		if !validateDateAndOrTime(&birthdate) {
			return vCardError(fmt.Sprintf("'%v' is not a valid date-and-or-time.", birthdate))
		}
	} else {
		if p.Type != "text" {
			return vCardError(fmt.Sprintf("Only 'text' and 'date-and-or-time' allowed, but got %v", p.Type))
		}
	}
	for k, v := range p.Parameters {
		switch k {
		case "altid":
			if err := validateAltidParam(v); err != nil {
				return err
			}
		case "calscale":
			if err := validateCalscaleParam(v); err != nil {
				return err
			}
		default:
			return vCardError(fmt.Sprintf("Unknown parameter type: %v", v))
		}
	}
	return nil
}

// ValidateBDay https://tools.ietf.org/html/rfc6350#section-6.2.5
func ValidateBDay(p *VCardProperty) error {
	return validateAnniversaryAndBDate(p, BDay)

}

// ValidateAnniversary https://tools.ietf.org/html/rfc6350#section-6.2.6
func ValidateAnniversary(p *VCardProperty) error {
	return validateAnniversaryAndBDate(p, Anniversary)
}
