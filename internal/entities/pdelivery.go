package entities

import (
	"fmt"
	"log"
)

// VCard methods related to addresses moved to own file to shrink jcard.go

// POBox returns the address's PO Box.
//
// Returns empty string if no address is present.
func (v *VCard) POBox() string {
	return v.getFirstAddressField(0)
}

// ExtendedAddress returns the "extended address", e.g. an apartment
// or suite number.
//
// Returns empty string if no address is present.
func (v *VCard) ExtendedAddress() string {
	return v.getFirstAddressField(1)
}

// StreetAddress returns the street address.
//
// Returns empty string if no address is present.
func (v *VCard) StreetAddress() string {
	return v.getFirstAddressField(2)
}

// Locality returns the address locality.
//
// Returns empty string if no address is present.
func (v *VCard) Locality() string {
	return v.getFirstAddressField(3)
}

// Region returns the address region (e.g. state or province).
//
// Returns empty string if no address is present.
func (v *VCard) Region() string {
	return v.getFirstAddressField(4)
}

// PostalCode returns the address postal code (e.g. zip code).
//
// Returns empty string if no address is present.
func (v *VCard) PostalCode() string {
	return v.getFirstAddressField(5)
}

// Country returns the address country name.
//
// This is the full country name.
//
// Returns empty string if no address is present.
func (v *VCard) Country() string {
	return v.getFirstAddressField(6)
}

func (v *VCard) getFirstAddressField(index int) string {
	adr := v.GetFirst("adr")
	if adr == nil {
		return ""
	}

	values := adr.Values()

	if index >= len(values) {
		return ""
	}

	return values[index]
}

var adrTypeInfo = VPropertyTypeInfo{
	Name:          Adr,
	Cardinal:      Many,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.3.1",
	Validate:      ValidateAdr,
}

// AddAdr adds an address to VCard. Return false if some validation fail.
// https://tools.ietf.org/html/rfc6350#section-6.3.1
// https://tools.ietf.org/html/rfc7095#section-3.3.1.3
func (v *VCard) AddAdr(value []string, params map[string][]string) error {
	property := &VCardProperty{
		Name:       Adr,
		Type:       Text,
		Parameters: params,
		Value:      value,
		typeInfo:   adrTypeInfo,
	}
	return v.append(property)
}

//ValidateAdr validate if property is a valid adr property.
func ValidateAdr(p *VCardProperty) error {

	for k, v := range p.Parameters {
		switch k {
		case "label", "tz":
			// TODO validate this 'text' fields if used!
		case "geo":
			if !validateGeo(&v[0]) {
				return vCardError(fmt.Sprintf("Invalid geo coordinates: %v", &v[0]))
			}
		case "lang":
			if err := validateLanguageParam(v); err != nil {
				return err
			}
		default:
			if err := validateDefaultParam(k, v); err != nil {
				return err
			}
		}
	}
	if p.Type != Text {
		return vCardError(fmt.Sprintf("type must be 'text' but is %v", p.Type))
	}
	var sa []string
	var ok bool
	if sa, ok = p.Value.([]string); !ok {
		return vCardError(fmt.Sprintf("value must be '[]string' but is %T", p.Value))
	}
	if len(sa) != 7 {
		return vCardError(` The adr must have 7 items:
		the post office box;
		the extended address (e.g., apartment or suite number);
		the street address;
		the locality (e.g., city);
		the region (e.g., state or province);
		the postal code;
		the country name; 
		but number of items is ` + fmt.Sprintf("%v", len(sa)))
	}
	if sa[0] != "" || sa[1] != "" {
		log.Printf("First two items of address array SHOULD be empty, but array is %v", sa)
	}
	return nil
}
