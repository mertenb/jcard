package entities

var keyTypeInfo = VPropertyTypeInfo{
	Name:          Key,
	Cardinal:      Many,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.8.1",
	Validate:      ValidateKey,
}

// AddKey To specify a public key or authentication certificate associated with the object that the vCard represents.
// https://tools.ietf.org/html/rfc6350#section-6.8.1
func (v *VCard) AddKey(value string, params map[string][]string) error {
	stype := Text
	if validateURI(&value) {
		stype = URI
	}

	property := &VCardProperty{
		Name:       Key,
		Type:       stype,
		Parameters: params,
		Value:      value,
		typeInfo:   keyTypeInfo,
	}
	return v.append(property)
}

// ValidateKey check if th fulfill the requirenments of
// https://tools.ietf.org/html/rfc6350#section-6.8.1
func ValidateKey(p *VCardProperty) error {

	for k, v := range p.Parameters {
		if err := validateDefaultParam(k, v); err != nil {
			return err
		}
	}
	return nil
}
