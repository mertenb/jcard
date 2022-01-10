package entities

import (
	"fmt"
	"strconv"
)

var categoriesTypeInfo = VPropertyTypeInfo{
	Name:          Categories,
	Cardinal:      Many,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.7.1",
	Validate:      ValidateCategories,
}

// AddCategories to specify application category information about the vCard, also known as "tags".
// https://tools.ietf.org/html/rfc6350#section-6.7.1
func (v *VCard) AddCategories(value []string, params map[string][]string) error {
	property := &VCardProperty{
		Name:       Categories,
		Type:       Text,
		Parameters: params,
		Value:      value,
		typeInfo:   categoriesTypeInfo,
	}
	return v.append(property)
}

var noteTypeInfo = VPropertyTypeInfo{
	Name:          Note,
	Cardinal:      Many,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.7.2",
	Validate:      ValidateNote,
}

// AddNote To specify supplemental information or a comment that is associated with the vCard.
// https://tools.ietf.org/html/rfc6350#section-6.7.2
func (v *VCard) AddNote(value string, params map[string][]string) error {
	property := &VCardProperty{
		Name:       Note,
		Type:       Text,
		Parameters: params,
		Value:      value,
		typeInfo:   noteTypeInfo,
	}
	return v.append(property)
}

var prodIDTypeInfo = VPropertyTypeInfo{
	Name:          ProdID,
	Cardinal:      OneOrZero,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.7.3",
	Validate:      ValidateProdid,
}

// AddProdid to specify the identifier for the product that created the vCard object.
// https://tools.ietf.org/html/rfc6350#section-6.7.3
func (v *VCard) AddProdid(value string, params map[string][]string) error {
	property := &VCardProperty{
		Name:       ProdID,
		Type:       Text,
		Parameters: params,
		Value:      value,
		typeInfo:   prodIDTypeInfo,
	}
	return v.append(property)
}

var revTypeInfo = VPropertyTypeInfo{
	Name:          Rev,
	Cardinal:      OneOrZero,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.7.4",
	Validate:      ValidateRev,
}

// AddRev To specify revision information about the current vCard.
// https://tools.ietf.org/html/rfc6350#section-6.7.4
func (v *VCard) AddRev(value string, params map[string][]string) error {
	property := &VCardProperty{
		Name:       Rev,
		Type:       "timestamp",
		Parameters: params,
		Value:      value,
		typeInfo:   revTypeInfo,
	}
	return v.append(property)
}

var soundTypeInfo = VPropertyTypeInfo{
	Name:          Sound,
	Cardinal:      Many,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.7.5",
	Validate:      ValidateSound,
}

// AddSound To specify a digital sound content information that annotates some aspect of the vCard.
// This property is often used to specify the proper pronunciation of the name property value of the vCard.
// https://tools.ietf.org/html/rfc6350#section-6.7.5
func (v *VCard) AddSound(value interface{}, params map[string][]string) error {
	property := &VCardProperty{
		Name:       Sound,
		Type:       URI,
		Parameters: params,
		Value:      value,
		typeInfo:   soundTypeInfo,
	}

	return v.append(property)
}

var uidTypeInfo = VPropertyTypeInfo{
	Name:          UID,
	Cardinal:      OneOrZero,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.7.6",
	Validate:      ValidateUID,
}

// AddUID To specify To specify a value that represents a globally unique identifier corresponding to the entity associated with the vCard.
// https://tools.ietf.org/html/rfc6350#section-6.7.6
func (v *VCard) AddUID(value interface{}, params map[string][]string) error {
	svalue := value.(string)
	stype := Text
	if validateURI(&svalue) {
		stype = URI
	}
	property := &VCardProperty{
		Name:       UID,
		Type:       stype,
		Parameters: params,
		Value:      value,
		typeInfo:   uidTypeInfo,
	}
	return v.append(property)
}

var clientPIDMapTypeInfo = VPropertyTypeInfo{
	Name:          ClientPIDMap,
	Cardinal:      Many,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.7.7",
	Validate:      ValidateClientPIDMap,
}

// AddClientPIDMap To give a global meaning to a local PID source identifier.
// https://tools.ietf.org/html/rfc6350#section-6.7.7
func (v *VCard) AddClientPIDMap(value interface{}, params map[string][]string) error {

	property := &VCardProperty{
		Name:       ClientPIDMap,
		Type:       "1*DIGIT; URI",
		Parameters: params,
		Value:      value,
		typeInfo:   clientPIDMapTypeInfo,
	}
	return v.append(property)
}

var urlTypeInfo = VPropertyTypeInfo{
	Name:          URL,
	Cardinal:      Many,
	Specification: "https://tools.ietf.org/html/rfc6350#section-6.7.7",
	Validate:      ValidateURL,
}

// AddURL To specify a uniform resource locator associated with the object to which the vCard refers.
// https://tools.ietf.org/html/rfc6350#section-6.7.8
func (v *VCard) AddURL(value interface{}, params map[string][]string) error {

	property := &VCardProperty{
		Name:       URL,
		Type:       URI,
		Parameters: params,
		Value:      value,
		typeInfo:   urlTypeInfo,
	}
	return v.append(property)
}

// ValidateCategories check if categories fulfill the requirenments of
// https://tools.ietf.org/html/rfc6350#section-6.7.1
func ValidateCategories(p *VCardProperty) error {

	for k, v := range p.Parameters {
		if err := validateDefaultParam(k, v); err != nil {
			return err
		}
	}
	return nil
}

// ValidateNote check if note fulfill the requirenments of
// https://tools.ietf.org/html/rfc6350#section-6.7.2
func ValidateNote(p *VCardProperty) error {

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

// ValidateProdid check if prodid fulfill the requirenments of
// https://tools.ietf.org/html/rfc6350#section-6.7.3
func ValidateProdid(p *VCardProperty) error {

	if (p.Parameters != nil) && (len(p.Parameters) > 0) {
		return vCardError("prodid must not have parameter.")
	}
	return nil
}

// ValidateRev check if rev fulfill the requirenments of
// https://tools.ietf.org/html/rfc6350#section-6.7.4
func ValidateRev(p *VCardProperty) error {

	if (p.Parameters != nil) && (len(p.Parameters) > 0) {
		return vCardError("rev must not have parameter.")
	}
	sValue := p.Value.(string)
	if !validateTimestamp(&sValue) {
		return vCardError(fmt.Sprintf("Expected valid timestamp, but got %v", sValue))
	}
	return nil
}

// ValidateSound  check if sound fulfill the requirenments of
// https://tools.ietf.org/html/rfc6350#section-6.7.5
func ValidateSound(p *VCardProperty) error {

	sound := p.Value.(string)
	if !validateURI(&sound) {
		return vCardError(fmt.Sprintf("'%v' is not a valid uri", sound))
	}
	for k, v := range p.Parameters {
		switch k {
		case "mediatype":
			if err := validateMediatypeParam(v); err != nil {
				return err
			}
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

// ValidateUID check if UID fulfill the requirenments of
// https://tools.ietf.org/html/rfc6350#section-6.7.6
func ValidateUID(p *VCardProperty) error {

	if len(p.Parameters) != 0 {
		return vCardError(fmt.Sprintf("UID must not have parameter, but got %v", p.Parameters))
	}

	val := p.Value.(string)
	if p.Type == URI {
		if !validateURI(&val) {
			return vCardError(fmt.Sprintf("Valid uri expected, but got %v", val))
		}

	} else if p.Type != Text {
		return vCardError(fmt.Sprintf("type 'uri' or 'text' expected, but is %v", p.Type))
	}
	return nil
}

// ValidateClientPIDMap check if clientPIDMap fulfill the requirenments of
// https://tools.ietf.org/html/rfc6350#section-6.7.7
func ValidateClientPIDMap(p *VCardProperty) error {

	if len(p.Parameters) != 0 {
		return vCardError(fmt.Sprintf("clientPIDMap must not have parameter, but got %v", p.Parameters))
	}

	aval := p.Value.([]interface{})
	if len(aval) != 2 {
		return vCardError(fmt.Sprintf("a pair of values epxected (digit,uri), but got %v", aval))
	}
	if i, err := strconv.Atoi(aval[0].(string)); (err != nil) || (i < 1) {
		return vCardError(fmt.Sprintf("The first field is a small integer (>0) corresponding to the second field of a PID parameter instance, but got %v. (%w)", aval[0], err))

	}
	uri := aval[1].(string)
	if !validateURI(&uri) {
		return vCardError(fmt.Sprintf("The second field must be a uri, but is %v.", p.Type))
	}
	return nil
}

// ValidateURL check if url fulfill the requirenments of
// https://tools.ietf.org/html/rfc6350#section-6.7.8
func ValidateURL(p *VCardProperty) error {

	url := p.Value.(string)
	if !validateURI(&url) {
		return vCardError(fmt.Sprintf("'%v' is not a valid uri", url))
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
