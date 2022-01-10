package entities

// Package prop contains all vCard properties based on version 4.
// Static properties.
const (
	VBegin = "BEGIN:VCARD"
	VEnd   = "END:VCARD"
)

// vcard property names.
// must be lowercase for jCard, should be uppercase for vCard.
const (
	Adr          = "adr"
	Anniversary  = "anniversary"
	BDay         = "bday"
	CalAdrURI    = "caladruri"
	CalURI       = "caluri"
	Categories   = "categories"
	Class        = "class"
	ClientPIDMap = "clientpidmap"
	Email        = "email"
	FbURL        = "fburl"
	Fn           = "fn"
	Gender       = "gender"
	Geo          = "geo"
	IMPP         = "impp"
	Key          = "key"
	Kind         = "kind"
	Lang         = "lang"
	Logo         = "logo"
	Mailer       = "mailer"
	Member       = "member"
	N            = "n"
	Name         = "name"
	Nickname     = "nickname"
	Note         = "note"
	Org          = "org"
	Photo        = "photo"
	ProdID       = "prodid"
	Profile      = "profile"
	Related      = "related"
	Rev          = "rev"
	Role         = "role"
	Sound        = "sound"
	Source       = "source"
	Tel          = "tel"
	Title        = "title"
	Tz           = "tz"
	UID          = "uid"
	URL          = "url"
	Version      = "version"
	XML          = "xml"

	// property types
	URI  = "uri"
	Text = "text"

	// Additional properties
	BirthPlace = "birthplace"
	DeathPlace = "deathplace"
	DeathDate  = "deathdate"
	Expertise  = "expertise"
	Hobby      = "hobby"

	Interest     = "interest"
	OrgDirectory = "org-directory"
)
