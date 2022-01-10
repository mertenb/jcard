package entities

import (
	"fmt"
)

var fburlTypeInfo = VPropertyTypeInfo{
	Name:          FbURL,
	Cardinal:      Many,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.9.1",
	Validate:      ValidateFburl,
}

// AddFburl To specify the URI for the busy time associated with the object that the vCard represents.
// https://tools.ietf.org/html/rfc6350#section-6.9.1
func (v *VCard) AddFburl(url string, params map[string][]string) error {
	property := &VCardProperty{
		Name:       FbURL,
		Type:       URI,
		Parameters: params,
		Value:      url,
		typeInfo:   fburlTypeInfo,
	}
	return v.append(property)
}

var caladruriTypeInfo = VPropertyTypeInfo{
	Name:          CalAdrURI,
	Cardinal:      Many,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.9.2",
	Validate:      ValidateCaladruri,
}

// AddCaladruri To specify the calendar user address [RFC5545] to which a scheduling request [RFC5546] should be sent for the object represented by the vCard.
// https://tools.ietf.org/html/rfc6350#section-6.9.2
func (v *VCard) AddCaladruri(url string, params map[string][]string) error {
	property := &VCardProperty{
		Name:       CalAdrURI,
		Type:       URI,
		Parameters: params,
		Value:      url,
		typeInfo:   caladruriTypeInfo,
	}
	return v.append(property)
}

var caluriTypeInfo = VPropertyTypeInfo{
	Name:          CalURI,
	Cardinal:      Many,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.9.3",
	Validate:      ValidateCaluri,
}

// AddCaluri To specify the URI for a calendar associated with the object represented by the vCard.
// https://tools.ietf.org/html/rfc6350#section-6.9.3
func (v *VCard) AddCaluri(url string, params map[string][]string) error {
	property := &VCardProperty{
		Name:       CalURI,
		Type:       URI,
		Parameters: params,
		Value:      url,
		typeInfo:   caluriTypeInfo,
	}
	return v.append(property)
}

// ValidateFburl check if th fulfill the requirenments of
// https://tools.ietf.org/html/rfc6350#section-6.9.1
func ValidateFburl(p *VCardProperty) error {
	url := p.Value.(string)
	if !validateURI(&url) {
		return vCardError(fmt.Sprintf("'%v' is not a valid uri", url))
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

// ValidateCaladruri check if th fulfill the requirenments of
// https://tools.ietf.org/html/rfc6350#section-6.9.2
func ValidateCaladruri(p *VCardProperty) error {

	url := p.Value.(string)
	if !validateURI(&url) {
		return vCardError(fmt.Sprintf("'%v' is not a valid uri", url))
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

// ValidateCaluri check if th fulfill the requirenments of
// https://tools.ietf.org/html/rfc6350#section-6.9.3
func ValidateCaluri(p *VCardProperty) error {

	url := p.Value.(string)
	if !validateURI(&url) {
		return vCardError(fmt.Sprintf("'%v' is not a valid uri", url))
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
