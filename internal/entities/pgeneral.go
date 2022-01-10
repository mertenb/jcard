package entities

import (
	"fmt"
)

var versionTypeInfo = VPropertyTypeInfo{
	Name:          Version,
	Cardinal:      One,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.7.9",
	Validate:      ValidateVersion,
}

var kindTypeInfo = VPropertyTypeInfo{
	Name:          Kind,
	Cardinal:      OneOrZero,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.1.4",
	Validate:      ValidateKind,
}

// AddKind specify the kind of object the vCard represents.
// https://tools.ietf.org/html/rfc6350#section-6.1.4
func (v *VCard) AddKind(value string) error {

	property := &VCardProperty{
		Name:       Kind,
		Type:       Text,
		Parameters: nil,
		Value:      value,
		typeInfo:   kindTypeInfo,
	}
	return v.append(property)
}

// ValidateKind check if 'kind' fulfill the requirenments of https://tools.ietf.org/html/rfc6350#section-6.2.3
func ValidateKind(p *VCardProperty) error {
	valids := map[string][]struct{}{"individual": {}, "group": {}, "org": {}, "location": {}}

	_, ok := valids[fmt.Sprintf("%v", p.Value)]
	if !ok {
		return vCardError(fmt.Sprintf("invalid kind value. Expected one of %v but got %v", valids, p.Value))
	}

	if (p.Parameters != nil) && (len(p.Parameters) != 0) {
		return vCardError(fmt.Sprintf("No parameters allowed for 'kind', but got %v", p.Parameters))
	}

	return nil
}

// ValidateVersion check if version fulfill the requirenments.
func ValidateVersion(p *VCardProperty) error {

	if (p.Parameters != nil) && (len(p.Parameters) != 0) {
		return vCardError(fmt.Sprintf("Version must not have parameters, but has %v", p.Parameters))
	}

	sValue := p.Value.(string)
	if sValue != "4.0" {
		return vCardError(fmt.Sprintf("Only version 4.0 is supported, but version is %v", sValue))
	}

	return nil
}
