package entities

import (
	"fmt"
)

var tzTypeInfo = VPropertyTypeInfo{
	Name:          Tz,
	Cardinal:      Many,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.5.1",
	Validate:      ValidateTz,
}

// AddTz add information related to the time zone of the object the jCard represents. https://tools.ietf.org/html/rfc6350#section-6.5.1
func (v *VCard) AddTz(value string, params map[string][]string) error {
	stype := Text
	if validateURI(&value) {
		stype = URI
	} else {
		if validateUtcOffset(&value) {
			stype = "utc-offset"
		}
	}

	property := &VCardProperty{
		Name:       Tz,
		Type:       stype,
		Parameters: params,
		Value:      value,
		typeInfo:   tzTypeInfo,
	}
	return v.append(property)
}

var geoTypeInfo = VPropertyTypeInfo{
	Name:          Geo,
	Cardinal:      Many,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.5.2",
	Validate:      ValidateGeo,
}

// AddGeo add specify information related to the global positioning of the object the jCard represents.
// https://tools.ietf.org/html/rfc6350#section-6.5.2
func (v *VCard) AddGeo(value string, params map[string][]string) error {
	property := &VCardProperty{
		Name:       Geo,
		Type:       URI,
		Parameters: params,
		Value:      value,
		typeInfo:   geoTypeInfo,
	}
	return v.append(property)
}

// ValidateTz check if th fulfill the requirenments of
// https://tools.ietf.org/html/rfc6350#section-6.5.1
func ValidateTz(p *VCardProperty) error {

	for k, v := range p.Parameters {
		switch k {
		case "mediatype":
			return validateMediatypeParam(v)
		default:
			if err := validateDefaultParam(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

//ValidateGeo check if geo fulfill the requirenments of https://tools.ietf.org/html/rfc6350#section-6.5.2
func ValidateGeo(p *VCardProperty) error {

	s := p.Value.(string)
	if !(validateGeo(&s)) {
		return vCardError(fmt.Sprintf("Valid geo uri expected, but got %v", p.Value))
	}

	for k, v := range p.Parameters {
		switch k {
		case "mediatype":
			return validateMediatypeParam(v)
		default:
			if err := validateDefaultParam(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}
