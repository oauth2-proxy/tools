package types

import (
	"text/template"
	"time"
)

// MyTestStruct contains a collection of fields all attempting to test various
// aspects of the code generation.
type MyTestStruct struct {
	// Name is the name of the MyTestStruct.
	Name string `yaml:"name"`

	// privateField is a private field and so should not be documented
	// in the generated docs.
	privateField string

	// LongMessageInt has a very long message, very very very very very very
	// very very very very very very very very very very very very very very
	// very very very very very very very very very very very very very very
	// very very very very very very very very very very very very very very
	// very very very very very very very very very very very very very very
	// long message attached to the top of it.
	// This should prove how the generator handles long doc strings.
	LongMessageInt int `yaml:"longMessageInt"`

	// SubStruct is a struct referenced from within the parent struct.
	// This should get its own section in the referenced docs.
	SubStruct SomeSubStruct `yaml:"subStruct"`

	// SubStructMap is a map of a known struct type.
	SubStructMap map[string]SomeSubStruct `yaml:"subStructMap"`

	// AnEmbeddedStruct is embedded here.
	AnEmbeddedStruct

	// AliasedDuration is a type alias to a duration.
	AliasedDuration MyDuration `yaml:"aliasedDuration"`

	// AliasDurationString is a type alias to a duration that should be documented
	// as a string type.
	AliasedDurationString MyDurationString `yaml:"aliasedDurationString"`

	// PointerString shows that the docs gen strips the pointer (*) from the beginning
	// of the type when documented.
	PointerString *string `yaml:"pointerString"`

	// Private should be included as a new struct, but without any documented members.
	Private PrivateMembers `yaml:"private"`

	// AliasedStruct is a type aliased struct
	AliasedStruct AliasSubStruct `yaml:"aliasedStruct"`

	// ExternalMap references and external map type outisde of the package.
	ExternalMap template.FuncMap `yaml:"externalMap"`

	// AliasExternalMap references an external map type outside of the package via an alias.
	AliasExternalMap AliasedExternalMap `yaml:"aliasExternalMap"`

	// Bytes is a slice of raw byte data.
	Bytes []byte `yaml:"bytes"`
}

// SomeSubStruct is a struct to go within another struct.
type SomeSubStruct struct {
	// NonTaggedField doesn't have a tag, so the name will be capitalised.
	NonTaggedField bool

	// privateStruct should not be included in the docs.
	privateStruct PrivateMembers
}

// AliasSubStruct is an aliased struct, it will be added to the documentation with an identical
// members table as the origin struct.
type AliasSubStruct SomeSubStruct

// AnEmbeddedStruct gets embedded within other structures.
type AnEmbeddedStruct struct {
	// EmbeddedDuration is a duration within an embedded struct.
	EmbeddedDuration time.Duration `yaml:"embeddedDuration"`
}

// PrivateMembers only has private members so when documented, should not have a members table printed.
type PrivateMembers struct {
	privateInt   int64
	privateBool  bool
	privateBytes []byte
}

// MyDuration is an alias to a duration.
type MyDuration time.Duration

// MyDuration is an alias to a duration with the type overridden as a string
// +reference-gen:alias-name=string
type MyDurationString time.Duration

// AliasedExternalMap is an alias type for a map type outside of the package.
type AliasedExternalMap template.FuncMap
