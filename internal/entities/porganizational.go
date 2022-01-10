package entities

import (
	"fmt"
)

var titleTypeInfo = VPropertyTypeInfo{
	Name:          Title,
	Cardinal:      Many,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.6.1",
	Validate:      ValidateTitle,
}

// AddTitle specify the position or job of the object the vCard represents. https://tools.ietf.org/html/rfc6350#section-6.6.1
func (v *VCard) AddTitle(value string, params map[string][]string) error {
	property := &VCardProperty{
		Name:       Title,
		Type:       Text,
		Parameters: params,
		Value:      value,
		typeInfo:   titleTypeInfo,
	}
	return v.append(property)
}

var roleTypeInfo = VPropertyTypeInfo{
	Name:          Role,
	Cardinal:      Many,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.6.2",
	Validate:      ValidateRole,
}

// AddRole specify the function or part played in a particular situation by the object the vCard represents. https://tools.ietf.org/html/rfc6350#section-6.6.2
func (v *VCard) AddRole(value string, params map[string][]string) error {
	property := &VCardProperty{
		Name:       Role,
		Type:       Text,
		Parameters: params,
		Value:      value,
		typeInfo:   roleTypeInfo,
	}
	return v.append(property)
}

var logoTypeInfo = VPropertyTypeInfo{
	Name:          Logo,
	Cardinal:      Many,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.6.3",
	Validate:      ValidateLogo,
}

// AddLogo specify a graphic image of a logo associated with the object the vCard represents.
// https://tools.ietf.org/html/rfc6350#section-6.6.3
func (v *VCard) AddLogo(value interface{}, params map[string][]string) error {
	property := &VCardProperty{
		Name:       Logo,
		Type:       URI,
		Parameters: params,
		Value:      value,
		typeInfo:   logoTypeInfo,
	}
	return v.append(property)
}

var orgTypeInfo = VPropertyTypeInfo{
	Name:          Org,
	Cardinal:      Many,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.6.4",
	Validate:      ValidateOrg,
}

// AddOrg To specify the organizational name and units associated with the vCard.
// https://tools.ietf.org/html/rfc6350#section-6.6.4
func (v *VCard) AddOrg(value interface{}, params map[string][]string) error {
	property := &VCardProperty{
		Name:       Org,
		Type:       Text,
		Parameters: params,
		Value:      value,
		typeInfo:   orgTypeInfo,
	}
	return v.append(property)
}

var memberTypeInfo = VPropertyTypeInfo{
	Name:          Member,
	Cardinal:      Many,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.6.5",
	Validate:      ValidateMember,
}

// AddMember To specify the organizational name and units associated with the vCard.
// https://tools.ietf.org/html/rfc6350#section-6.6.5
func (v *VCard) AddMember(value interface{}, params map[string][]string) error {
	property := &VCardProperty{
		Name:       Member,
		Type:       URI,
		Parameters: params,
		Value:      value,
		typeInfo:   memberTypeInfo,
	}
	return v.append(property)
}

var relatedTypeInfo = VPropertyTypeInfo{
	Name:          Related,
	Cardinal:      Many,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.6.6",
	Validate:      ValidateRelated,
}

// AddRelated To specify a relationship between another entity and the entity represented by this vCard.
// https://tools.ietf.org/html/rfc6350#section-6.6.6
func (v *VCard) AddRelated(value interface{}, params map[string][]string) error {
	svalue := value.(string)
	stype := Text
	if validateURI(&svalue) {
		stype = URI
	}
	property := &VCardProperty{
		Name:       Related,
		Type:       stype,
		Parameters: params,
		Value:      value,
		typeInfo:   relatedTypeInfo,
	}
	return v.append(property)
}

// ValidateTitle check if title fulfill the requirenments of
// https://tools.ietf.org/html/rfc6350#section-6.6.1
func ValidateTitle(p *VCardProperty) error {

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

// ValidateRole check if role fulfill the requirenments of
// https://tools.ietf.org/html/rfc6350#section-6.6.2
// TODO maybe make a map with key="language" value = 'validation-func' to reduce code of switch-statements.
func ValidateRole(p *VCardProperty) error {

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

// ValidateLogo check if logo fulfill the requirenments of https://tools.ietf.org/html/rfc6350#section-6.6.3
// TODO maybe make a map with key="language" value = 'validation-func' to reduce code of switch-statements.
func ValidateLogo(p *VCardProperty) error {

	value := p.Value.(string)
	if !validateURI(&value) {
		return vCardError(fmt.Sprintf("%v is not a valid uri.", p.Name))
	}

	for k, v := range p.Parameters {
		switch k {
		case "language":
			if err := validateLanguageParam(v); err != nil {
				return err
			}
		case "media-type":
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

// ValidateOrg check if logo fulfill the requirenments of
// https://tools.ietf.org/html/rfc6350#section-6.6.4
func ValidateOrg(p *VCardProperty) error {

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
		default:
			if err := validateDefaultParam(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

// ValidateMember check if logo fulfill the requirenments of
// https://tools.ietf.org/html/rfc6350#section-6.6.5
// NOTE: parameter 'type' (=general parameter) must be 'group'
func ValidateMember(p *VCardProperty) error {

	val := p.Value.(string)
	if !validateURI(&val) {
		return vCardError(fmt.Sprintf("Valid uri expected, but got %v", val))
	}

	for k, v := range p.Parameters {
		switch k {
		case "mediatype":
			if err := validateMediatypeParam(v); err != nil {
				return err
			}
		case "pref":
			if err := validatePrefParam(v); err != nil {
				return err
			}
		case "pid":
			if err := validatePidParam(v); err != nil {
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

// ValidateRelated check if logo fulfill the requirenments of
// https://tools.ietf.org/html/rfc6350#section-6.6.6
func ValidateRelated(p *VCardProperty) error {

	val := p.Value.(string)
	if p.Type == URI {
		if !validateURI(&val) {
			return vCardError(fmt.Sprintf("Valid uri expected, but got %v", val))
		}
		for k, v := range p.Parameters {
			switch k {
			case "mediatype":
				if err := validateMediatypeParam(v); err != nil {
					return err
				}
			default:
				if err := validateRelatedParams(k, v); err != nil {
					return err
				}
			}
		}
	} else {
		if p.Type != Text {
			return vCardError(fmt.Sprintf("type 'uri' or 'text' expected, but is %v", p.Type))
		}
		for k, v := range p.Parameters {
			switch k {
			case "language":
				if err := validateLanguageParam(v); err != nil {
					return err
				}
			default:
				if err := validateRelatedParams(k, v); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func validateRelatedParams(key string, value []string) error {
	validTypes := map[string][]struct{}{"contact": {}, "acquaintance": {}, "friend": {}, "": {}, "met": {},
		"co-worker": {}, "colleague": {}, "co-resident": {}, "neighbor": {}, "child": {}, "parent": {}, "sibling": {},
		"spouse": {}, "kin": {}, "muse": {}, "crush": {}, "date": {}, "sweetheart": {}, "me": {}, "agent": {}, "emergency": {}}

	switch key {

	case "type":
		for _, v := range value {
			if _, b := validTypes[v]; !b {
				return vCardError(fmt.Sprintf("expected one of '%v', but got '%v'", validTypes, value))
			}
		}
	case "pref":
		return validatePrefParam(value)
	case "pid":
		return validatePidParam(value)
	case "altid":
		return validateAltidParam(value)
	default:
		return vCardError("unknown parameter type: " + key)
	}
	return nil
}
