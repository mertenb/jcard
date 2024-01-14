package entities

// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.
// stolen from https://github.com/openrdap/rdap/blob/master/vcard.go
// see https://tools.ietf.org/html/rfc7095

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mertenb/jcard/internal/ldap"
	"golang.org/x/text/language"
)

// VCard represents a vCard.
//
// A vCard represents information about an individual or entity. It can include
// a name, telephone number, e-mail, delivery address, and other information.
//
// There are several vCard text formats. This implementation encodes/decodes the
// jCard format used by RDAP, as defined in https://tools.ietf.org/html/rfc7095.
//
// A jCard consists of an array of properties (e.g. "fn", "tel") describing the
// individual or entity. Properties may be repeated, e.g. to represent multiple
// telephone numbers. RFC6350 documents a set of standard properties.
//
// RFC7095 describes the JSON document format, which looks like:
//   ["vcard", [
//     [
//       ["version", {}, "text", "4.0"],
//       ["fn", {}, "text", "Joe Appleseed"],
//       ["tel", {
//             "type":["work", "voice"],
//           },
//           "uri",
//           "tel:+1-555-555-1234;ext=555"
//       ],
//       ...
//     ]
//   ]
type VCard struct {
	Properties []*VCardProperty
	Errors     map[string]error `json:"-"`
}

// VCardProperty represents a single vCard property.
//
// Each vCard property has four fields, these are:
//    Name   Parameters                  Type   Value
//    -----  --------------------------  -----  -----------------------------
//   ["tel", {"type":["work", "voice"]}, "uri", "tel:+1-555-555-1234;ext=555"]
type VCardProperty struct {
	Name string

	// vCard parameters can be a string, or array of strings.
	//
	// To simplify our usage, single strings are represented as an array of
	// length one.
	Parameters map[string][]string
	Type       string

	// A property value can be a simple type (string/float64/bool/nil), or be
	// an array. Arrays can be nested, and can contain a mixture of types.
	//
	// Value is one of the following:
	//   * string
	//   * float64
	//   * bool
	//   * nil
	//   * []interface{}. Can contain a mixture of these five types.
	//
	// To retrieve the property value flattened into a []string, use Values().
	Value interface{}

	typeInfo VPropertyTypeInfo
}

// Cardinality defines the permitted number of property
type Cardinality int

// Define possible cardinalities
const (
	Many Cardinality = iota + 1
	OneOrMany
	OneOrZero
	One
)

// VPropertyTypeInfo is wired by name to a concrete VCardProperty
// Each VCardProperty has a concrete corresponding VPropertyTypeInfo.
type VPropertyTypeInfo struct {
	Name          string
	Cardinal      Cardinality
	Specification string
	Validate      func(v *VCardProperty) error
	// todo: validate type (=text, =uri, =... or check if this is already done in property validation.)
}

// Values returns a simplified representation of the VCardProperty value.
//
// This is convenient for accessing simple unstructured data (e.g. "fn", "tel").
//
// The simplified []string representation is created by flattening the
// (potentially nested) VCardProperty value, and converting all values to strings.
func (p *VCardProperty) Values() []string {
	strings := make([]string, 0, 1)

	p.appendValueStrings(p.Value, &strings)

	return strings
}

func (p *VCardProperty) appendValueStrings(v interface{}, strings *[]string) {
	switch v := v.(type) {
	case nil:
		*strings = append(*strings, "")
	case bool:
		*strings = append(*strings, strconv.FormatBool(v))
	case float64:
		*strings = append(*strings, strconv.FormatFloat(v, 'f', -1, 64))
	case string:
		*strings = append(*strings, v)
	case []interface{}:
		for _, v2 := range v {
			p.appendValueStrings(v2, strings)
		}
	default:
		panic("Unknown type")
	}

}

// Validate all properties and VCard constraints of this vcard and put errors to vcard's error map.
func (v *VCard) Validate() {
	v.Errors = make(map[string]error)
	v.validateMandatoryProperties()
	v.validateProperties()
	v.validateMemberGroup()
	v.validateCardinality()

	// FIXME: ok_pid.json is invalid because
	// pid 5.2 has no corresponding entry in clientpidmap. Only '1' is covered there!
	// testPid should fail!
	v.validateClientPIDMap()
}

// Check if mandatory properties are present (like Version).
// The cardinality will be checked while property validation.
func (v *VCard) validateMandatoryProperties() {
	if v.GetFirst(Version) == nil {
		v.Errors[fmt.Sprintf("missing_property_version")] = vCardError("Version expected")
	}
}

// validate all properties of vcard.
func (v *VCard) validateProperties() {
	for i, prop := range v.Properties {
		err := prop.Validate()
		if err != nil {
			v.Errors[fmt.Sprintf("%v_%v validationError", i, prop.Name)] = err
		}
	}
}

// https://tools.ietf.org/html/rfc6350#section-6.6.5
// If Property 'member' exist, then property 'kind' must have the value 'group'
func (v *VCard) validateMemberGroup() {
	if len(v.Get(Member)) != 0 {
		prop := v.GetFirst(Kind)
		if prop == nil {
			v.Errors["kindIsNil"] = vCardError("Property 'kind' with value 'group' must exists if 'member' exists but is 'nil' ")
		} else {
			if prop.Value != "group" {
				v.Errors["kindValueNotGroup"] = vCardError(fmt.Sprintf("Kind value must be 'group' but is '%v' if member exist.", prop.Value))
			}
		}
	}
}

func (v *VCard) validateCardinality() {
	for _, prop := range v.Properties {
		switch prop.typeInfo.Cardinal {
		case Many:
			{
				// do nothing
			}
		case One:
			{
				l := len(v.Get(prop.Name))
				if l != 1 {
					v.Errors[fmt.Sprintf("%v_cardinality", prop.Name)] = vCardError(fmt.Sprintf("%v cardinalitry is 1 but found %v occurence", prop.Name, l))
				}
			}
		case OneOrMany:
			{
				l := len(v.Get(prop.Name))
				if l == 0 {
					v.Errors[fmt.Sprintf("%v_cardinality", prop.Name)] = vCardError(fmt.Sprintf("%v cardinalitry is 1* but found no occurence", prop.Name))
				}
			}
		case OneOrZero:
			{
				l := len(v.Get(prop.Name))
				if l > 1 {
					v.Errors[fmt.Sprintf("%v_cardinality", prop.Name)] = vCardError(fmt.Sprintf("%v cardinalitry is 01 but found %v occurence", prop.Name, l))
				}
			}
		default:
			{
				v.Errors[fmt.Sprintf("%v_cardinality", prop.Name)] = vCardError(fmt.Sprintf("%v cardinalitry is unknown: %v", prop.Name, prop.typeInfo.Cardinal))
			}
		}
	}
}

// https://tools.ietf.org/html/rfc6350#section-6.7.7
// Special notes:  PID source identifiers (the source identifier is the
// second field in a PID parameter instance) are small integers that
// only have significance within the scope of a single vCard
// instance.  Each distinct source identifier present in a vCard MUST
// have an associated CLIENTPIDMAP.  See Section 7 for more details
// on the usage of CLIENTPIDMAP.
func (v *VCard) validateClientPIDMap() {
	pids := v.getPIDs()
	cpms := v.getCpms()
	for k := range pids {

		if _, exist := cpms[k]; !exist {
			v.Errors[fmt.Sprintf("Pid %v has no ClientPIDMapping", k)] = vCardError(pids[k].String())
		}
	}
}

func (v *VCard) getCpms() map[int]*VCardProperty {
	result := make(map[int]*VCardProperty)
	cpms := v.Get(ClientPIDMap)
	for _, cpm := range cpms {
		aval := cpm.Value.([]interface{})
		if len(aval) != 2 {
			continue
		}
		if i, err := strconv.Atoi(aval[0].(string)); err == nil {
			result[i] = cpm
		}
	}
	return result
}

// get from all properties with param 'pid'  the pid's second fields.
// e.g. if 'pid=4.2' the '2' will be part of the result.
func (v *VCard) getPIDs() map[int]*VCardProperty {

	result := make(map[int]*VCardProperty)
	for _, p := range v.Properties {
		pid := p.Parameters["pid"]
		if pid != nil && len(pid) > 0 {
			for _, j := range pid {
				ints := strings.Split(j, ".")
				if len(ints) == 2 {
					v, err := strconv.Atoi(ints[1])
					if err == nil {
						result[v] = p
					}
				}
			}
		}
	}
	return result
}

func (v *VCard) hasErrors() bool {
	return v.Errors != nil && len(v.Errors) > 0
}

// String returns the vCard as a multiline human readable string. For example:
//
//   vCard[
//     version (type=text, parameters=map[]): [4.0]
//     mixed (type=text, parameters=map[]): [abc true 42 <nil> [def false 43]]
//   ]
//
// This is intended for debugging only, and is not machine parsable.
func (v *VCard) String() string {
	s := make([]string, 0, len(v.Properties))

	for _, s2 := range v.Properties {
		s = append(s, s2.String())
	}

	return "vCard[\n" + strings.Join(s, "\n") + "\n]"
}

// String returns the VCardProperty as a human readable string. For example:
//
//     mixed (type=text, parameters=map[]): [abc true 42 <nil> [def false 43]]
//
// This is intended for debugging only, and is not machine parsable.
func (p *VCardProperty) String() string {
	return fmt.Sprintf("  %s (type=%s, parameters=%v): %v", p.Name, p.Type, p.Parameters, p.Value)
}

// NewVCard creates a VCard from jsonBlob.
func NewVCard(jsonBlob []byte) (*VCard, error) {
	var top []interface{}
	err := json.Unmarshal(jsonBlob, &top)

	if err != nil {
		return nil, err
	}

	var vcard *VCard
	vcard, err = newVCardImpl(top)

	if (vcard != nil) && (err == nil) {
		for _, p := range vcard.Properties {
			p.addTypeInfo()
		}
	}

	return vcard, err
}

func newVCardImpl(src interface{}) (*VCard, error) {
	top, ok := src.([]interface{}) // type assertion: ok = true if 'src' is of type '[]interface{}' (top = src)

	if !ok || len(top) != 2 { // src(=jCard) = array with 2 elements. 1st: "vcard", 2nd: Array Of JCardProperties
		return nil, vCardError("structure is not a jCard (expected len=2 top level array)")
	} else if s, ok := top[0].(string); !(ok && s == "vcard") { // type assertion: first element is of type String and value 'vcard'
		return nil, vCardError("structure is not a jCard (missing 'vcard')")
	}

	var properties []interface{}

	properties, ok = top[1].([]interface{})
	if !ok {
		return nil, vCardError("structure is not a jCard (bad properties array)")
	}

	v := &VCard{
		Properties: make([]*VCardProperty, 0, len(properties)),
	}

	var p interface{}
	for _, p = range top[1].([]interface{}) {
		var a []interface{}
		var ok bool
		a, ok = p.([]interface{})

		if !ok {
			return nil, vCardError("jCard property was not an array")
		} else if len(a) < 4 {
			return nil, vCardError("jCard property too short (>=4 array elements required)")
		}

		name, ok := a[0].(string)

		if !ok {
			return nil, vCardError("jCard property name invalid")
		}

		var parameters map[string][]string
		var err error
		parameters, err = readParameters(a[1])

		if err != nil {
			return nil, err
		}

		propertyType, ok := a[2].(string)

		if !ok {
			return nil, vCardError("jCard property type invalid")
		}

		var value interface{}
		if len(a) == 4 {
			value, err = readValue(a[3], 0)
		} else {
			value, err = readValue(a[3:], 0)
		}

		if err != nil {
			return nil, err
		}

		property := &VCardProperty{
			Name:       name,
			Type:       propertyType,
			Parameters: parameters,
			Value:      value,
		}

		v.Properties = append(v.Properties, property)
	}

	return v, nil
}


// NewVCardFromLdap creates a VCard from given ldap person.
func NewVCardFromLdap(person *ldap.InetOrgPerson) (*VCard, error) {
	var vcard VCard

	property := &VCardProperty{
		Name:       "version",
		Type:       "text",
		Parameters: nil,
		Value:      "4.0",
	}
	vcard.Properties = append(vcard.Properties, property)

	property = &VCardProperty{
		Name:       Fn,
		Type:       Text,
		Parameters: nil,
		Value:      person.DisplayName,
	}
	vcard.Properties = append(vcard.Properties, property)

	property = &VCardProperty{
		Name:       Email,
		Type:       Text,
		Parameters: map[string][]string{"type": {"work"}, "pref": {"1"}},
		Value:      person.Mail,
	}
	vcard.Properties = append(vcard.Properties, property)

	vcard.AddTel(person.Mobile, map[string][]string{"type": {"voice", "cell"}, "pref": {"1"}})
	vcard.AddTel(person.TelephoneNumber, map[string][]string{"type": {"work", "voice"}, "pref": {"2"}})
	vcard.AddTel(person.FacsimileTelephoneNumber, map[string][]string{"type": {"work", "fax"}, "pref": {"3"}})

	return &vcard, nil
}


// Get returns a list of the vCard Properties with VCardProperty name |name|.
func (v *VCard) Get(name string) []*VCardProperty {
	var properties []*VCardProperty

	for _, p := range v.Properties {
		if p.Name == name {
			properties = append(properties, p)
		}
	}

	return properties
}

// GetFirst returns the first vCard Property with name |name|.
//
// TODO(tfh): Implement "pref" ordering, instead of taking the first listed property?
func (v *VCard) GetFirst(name string) *VCardProperty {
	properties := v.Get(name)

	if len(properties) == 0 {
		return nil
	}

	return properties[0]
}

func vCardError(e string) error {
	return fmt.Errorf("jCard error: %s", e)
}

func readParameters(p interface{}) (map[string][]string, error) {
	params := map[string][]string{}

	if _, ok := p.(map[string]interface{}); !ok {
		return nil, vCardError("jCard parameters invalid")
	}

	for k, v := range p.(map[string]interface{}) {
		if s, ok := v.(string); ok {
			params[k] = append(params[k], s)
		} else if arr, ok := v.([]interface{}); ok {
			for _, value := range arr {
				if s, ok := value.(string); ok {
					params[k] = append(params[k], s)
				}
			}
		}
	}

	return params, nil
}

func readValue(value interface{}, depth int) (interface{}, error) {
	switch value := value.(type) {
	case nil:
		return nil, nil
	case string:
		return value, nil
	case bool:
		return value, nil
	case float64:
		return value, nil
	case []interface{}:
		if depth == 3 {
			return "", vCardError("Structured value too deep")
		}

		result := make([]interface{}, 0, len(value))

		for _, v2 := range value {
			v3, err := readValue(v2, depth+1)

			if err != nil {
				return nil, err
			}

			result = append(result, v3)
		}

		return result, nil
	default:
		return nil, vCardError("Unknown JSON datatype in jCard value")
	}
}

func (v *VCard) getFirstPropertySingleString(name string) string {
	property := v.GetFirst(name)

	if property == nil {
		return ""
	}

	return strings.Join(property.Values(), " ")
}

// RemoveAll items with given name. Returns true if all properties are removed, false otherwise.
func (v *VCard) RemoveAll(name string) bool {
	result := false
	ok := true
	for _, p := range v.Get(name) {
		result = v.Remove(p)
		if !result {
			ok = false
		}
	}
	return result && ok
}

// Remove the given property if it exists. If removed, return true, false otherwise.
func (v *VCard) Remove(property *VCardProperty) bool {
	i := v.index(property)
	if i == -1 {
		return false
	}
	v.Properties[i] = v.Properties[len(v.Properties)-1]
	v.Properties = v.Properties[:len(v.Properties)-1]
	return true
}

func (v *VCard) index(value *VCardProperty) int {
	for p, v := range v.Properties {
		if v == value {
			return p
		}
	}
	return -1
}

func (v *VCard) append(prop *VCardProperty) error {
	if error := prop.Validate(); error != nil {
		log.Printf("Can't add %v to VCard. \n(\nValue: %v \nParameter: %v\n)", prop.Name, prop.Value, prop.Parameters)
		return fmt.Errorf("Adding %v failed: %v", prop.Name, error)
	}
	v.Properties = append(v.Properties, prop)
	return nil
}

// Validate the property
func (p *VCardProperty) Validate() error {
	if !(p.Name == p.typeInfo.Name) {
		return vCardError(fmt.Sprintf("Property '%v' has bad typeInfo '%v'", p.Name, p.typeInfo.Name))
	}
	return p.typeInfo.Validate(p)
}

// Validate the property
func (p *VCardProperty) addTypeInfo() error {
	switch p.Name {
	case Adr:
		p.typeInfo = adrTypeInfo
	case Anniversary:
		p.typeInfo = anniversaryTypeInfo
	case BDay:
		p.typeInfo = bdayTypeInfo
	case CalAdrURI:
		p.typeInfo = caladruriTypeInfo
	case CalURI:
		p.typeInfo = caluriTypeInfo
	case Categories:
		p.typeInfo = categoriesTypeInfo
	case ClientPIDMap:
		p.typeInfo = clientPIDMapTypeInfo
	case Email:
		p.typeInfo = emailTypeInfo
	case FbURL:
		p.typeInfo = fburlTypeInfo
	case Fn:
		p.typeInfo = fnTypeInfo
	case Gender:
		p.typeInfo = genderTypeInfo
	case Geo:
		p.typeInfo = geoTypeInfo
	case IMPP:
		p.typeInfo = imppTypeInfo
	case Key:
		p.typeInfo = keyTypeInfo
	case Kind:
		p.typeInfo = kindTypeInfo
	case Lang:
		p.typeInfo = langTypeInfo
	case Logo:
		p.typeInfo = logoTypeInfo
	case Member:
		p.typeInfo = memberTypeInfo
	case N:
		p.typeInfo = nTypeInfo
	case Nickname:
		p.typeInfo = nickNameTypeInfo
	case Note:
		p.typeInfo = noteTypeInfo
	case Org:
		p.typeInfo = orgTypeInfo
	case Photo:
		p.typeInfo = photoTypeInfo
	case ProdID:
		p.typeInfo = prodIDTypeInfo
	case Related:
		p.typeInfo = relatedTypeInfo
	case Rev:
		p.typeInfo = revTypeInfo
	case Role:
		p.typeInfo = roleTypeInfo
	case Sound:
		p.typeInfo = soundTypeInfo
	case Tel:
		p.typeInfo = telTypeInfo
	case Title:
		p.typeInfo = titleTypeInfo
	case Tz:
		p.typeInfo = tzTypeInfo
	case UID:
		p.typeInfo = uidTypeInfo
	case URL:
		p.typeInfo = urlTypeInfo
	case Version:
		p.typeInfo = versionTypeInfo
	default:
		return vCardError(fmt.Sprintf("unknown property name %v", p.Name))
	}
	return nil
}

func (p *VCardProperty) validateDefaultParam() error {
	for k, v := range p.Parameters {
		if err := validateDefaultParam(k, v); err != nil {
			return err
		}
	}
	return nil
}

// many properties have this common default parameter
func validateDefaultParam(key string, value []string) error {
	switch key {
	case "type":
		return validateTypeParam(value)
	case "pref":
		return validatePrefParam(value)
	case "pid":
		return validatePidParam(value)
	case "altid":
		return validateAltidParam(value)
	default:
		return vCardError("unknown parameter type: " + key)
	}
}

// https://tools.ietf.org/html/rfc6350#page-17
func validatePrefParam(params []string) error {
	if len(params) != 1 {
		return vCardError(fmt.Sprintf("only one value (1-100) allowed for pref, but is %v", params))
	}
	if value, err := strconv.Atoi(params[0]); (err != nil) || (value < 0) || (value > 100) {
		return vCardError(fmt.Sprintf("pref must be an integer between 0 and 100 but is %v ", params[0]))
	}
	return nil
}

// https://tools.ietf.org/html/rfc6350#section-5.5 (e.g. 5.4 OR 2 OR 4.1.4)
var validPid = regexp.MustCompile(`^\d+(.\d+)*$`)

func validatePidParam(params []string) error {
	for _, v := range params {
		if !validPid.MatchString(v) {
			return vCardError(fmt.Sprintf("PID value invalid. Positive <int>.<int> but is %v", v))
		}
	}
	return nil
}

// https://tools.ietf.org/html/rfc6350#section-5.6
func validateTypeParam(params []string) error {
	for _, v := range params {
		if !(v == "work" || v == "home") {
			return vCardError(fmt.Sprintf("Type must be one of 'work|home', but is %v", v))
		}
	}
	return nil
}

// https://tools.ietf.org/html/rfc6350#page-18
func validateAltidParam(params []string) error {
	for _, v := range params {
		if v != "" {
			return nil
		}
	}
	return nil
}

const regexpMedia = "(application|audio|font|example|image|message|model|multipart|text|video|x-(?:[0-9A-Za-z!#$%&'*+.^_`|~-]+))/([0-9A-Za-z!#$%&'*+.^_`|~-]+)"

var validMedia = regexp.MustCompile(regexpMedia)

func validateMediatypeParam(params []string) error {
	for _, v := range params {
		if !validMedia.MatchString(v) {
			return vCardError(fmt.Sprintf("invalid media type: %v", v))
		}
	}
	return nil
}

func validateLanguageParam(params []string) error {
	if len(params) > 1 {
		return vCardError(fmt.Sprintf("Only one language tag expected, but got %d", len(params)))
	}
	for _, v := range params {
		if _, err := language.Parse(v); err != nil {
			return vCardError(fmt.Sprintf("Unknown language (https://tools.ietf.org/html/rfc5646): %v :%v", v, err))
		}
	}
	return nil
}

func validateSortAsParam(values interface{}, paramvalues []string) error {
	var vArray []string
	var ok bool
	if vArray, ok = values.([]string); !ok {
		return vCardError(fmt.Sprintf("[]strings expected, but got '%v'", values))
	}

	if len(paramvalues) > len(vArray) {
		return vCardError(` The parameters values  MUST have as
		many or fewer elements as the corresponding property value has components.`)
	}
	return nil
}

// It is used to define the calendar system in which a date or date-time value is expressed
func validateCalscaleParam(params []string) error {
	for _, v := range params {
		if v != "gregorian" {
			return vCardError(fmt.Sprintf("Can only handle gregorian calender, but got %v", v))
		}
	}
	return nil
}

const regexpURI = `(?i)^([a-z0-9+.-]+):(?://(?:((?:[a-z0-9-._~!$&'()*+,;=:]|%[0-9A-F]{2})*)@)?((?:[a-z0-9-._~!$&'()*+,;=]|%[0-9A-F]{2})*)(?::(\d*))?(/(?:[a-z0-9-._~!$&'()*+,;=:@/]|%[0-9A-F]{2})*)?|(/?(?:[a-z0-9-._~!$&'()*+,;=:@]|%[0-9A-F]{2})+(?:[a-z0-9-._~!$&'()*+,;=:@/]|%[0-9A-F]{2})*)?)(?:\?((?:[a-z0-9-._~!$&'()*+,;=:/?@]|%[0-9A-F]{2})*))?(?:#((?:[a-z0-9-._~!$&'()*+,;=:/?@]|%[0-9A-F]{2})*))?$`

var validURI = regexp.MustCompile(regexpURI)

func validateURI(sURI *string) bool {
	return validURI.MatchString(*sURI)
}

const regexpUtcOffset = `^Z|[+-][01]\d[0-5]\d$`

var validUtcOffset = regexp.MustCompile(regexpUtcOffset)

func validateUtcOffset(s *string) bool {
	return validUtcOffset.MatchString(*s)
}

// https://tools.ietf.org/html/rfc7095#section-3.5.3
func validateDate(date *string) bool {
	validLayouts := []string{
		time.RFC3339,
		"2006-01-02",
		"2006-01",
		"2006",
		"--01-02",
		"---02",
		"--01",
	}
	return validateTemporal(&validLayouts, date)
}

// https://tools.ietf.org/html/rfc7095#section-3.5.4
func validateTime(date *string) bool {
	validLayouts := []string{
		"15:04:05",
		"15:04",
		"15",
		"-04",
		"-04:05",
		"--05",
	}
	return validateTemporal(&validLayouts, date)
}

// https://tools.ietf.org/html/rfc7095#section-3.5.5
// https://stackoverflow.com/questions/522251/whats-the-difference-between-iso-8601-and-rfc-3339-date-formats
// NOTE that 'T' is a mandatory separator between Date and Time in this implementation.
func validateDateTime(dt *string) bool {
	validLayouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05-07",
	}
	if validateTemporal(&validLayouts, dt) {
		return true
	}
	sa := strings.Split(*dt, "T")
	if len(sa) != 2 {
		return false
	}
	return (validateDate(&sa[0]) && validateTime(&sa[1]))
}

// https://tools.ietf.org/html/rfc7095#section-3.5.6
func validateDateAndOrTime(date *string) bool {
	return validateDate(date) || validateTime(date) || validateDateTime(date)
}

// NOTE that the golang time layouts allow fractional seconds which are NOT supported in https://tools.ietf.org/html/rfc7095#section-3.5.4
func validateTemporal(layouts *[]string, t *string) bool {
	for _, s := range *layouts {
		if _, err := time.Parse(s, *t); err == nil {
			return true
		}
	}
	return false
}

const regexpGeo = `^(geo:)?\-?\d+\.\d+?,\s*\-?\d+\.\d+?$`

var preferedGeo = regexp.MustCompile(regexpGeo)

// TODO validate against all geo types
//geo:46.772673,-71.282945
func validateGeo(g *string) bool {
	if !preferedGeo.MatchString(*g) {
		log.Printf("unknown geo schema: %v", *g)
		return false
	}
	return true
}

//
func validateTimestamp(t *string) bool {
	if _, err := time.Parse("20060102T150405Z07", *t); err == nil {
		return true
	}
	if _, err := time.Parse("20060102T150405", *t); err == nil {
		return true
	}
	return false
}
