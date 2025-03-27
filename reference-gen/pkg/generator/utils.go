package generator

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"unicode"

	"k8s.io/gengo/v2/types"
	"k8s.io/klog/v2"
)

var (
	commonTypes = map[string]string{
		"time.Duration": "duration",
	}
)

// BEGIN: stringSet

// stringSet is a utility for checking a set of strings
type stringSet map[string]struct{}

func (s stringSet) add(new string) {
	s[new] = struct{}{}
}

func (s stringSet) has(required string) bool {
	_, ok := s[required]
	return ok
}

func (s stringSet) isEmpty() bool {
	return len(s) == 0
}

func newStringSet(list []string) stringSet {
	set := make(stringSet)
	for _, s := range list {
		set.add(s)
	}
	return set
}

// END: stringSet

// BEGIN: typeSet

// typeSet is a utility for checking a set of types
type typeSet map[*types.Type]struct{}

func (s typeSet) add(new *types.Type) {
	s[new] = struct{}{}
}

func (s typeSet) has(required *types.Type) bool {
	_, ok := s[required]
	return ok
}

func (s typeSet) toList() []*types.Type {
	list := []*types.Type{}
	for typ := range s {
		list = append(list, typ)
	}
	return list
}

func newTypeSetFromStringMap(typeMap map[string]*types.Type) typeSet {
	set := make(typeSet)
	for _, t := range typeMap {
		set.add(t)
	}
	return set
}

func newTypeSetFromList(typeList []*types.Type) typeSet {
	set := make(typeSet)
	for _, t := range typeList {
		set.add(t)
	}
	return set
}

// END: typeSet

// tryDereference returns the underlying type when t is a pointer, map, or slice.
func tryDereference(t *types.Type) *types.Type {
	if t.Elem != nil {
		return t.Elem
	}
	return t
}

// createTypeList converts a map of types into a list of types
func createTypeList(typesForList map[*types.Type][]*types.Type) []*types.Type {
	out := []*types.Type{}
	for typ := range typesForList {
		out = append(out, typ)
	}
	return out
}

// BEGIN: template functions

// aliasDisplayNameFunc constructs a aliasDisplayName function for the template
func aliasDisplayNameFunc(knownTypes typeSet) func(t *types.Type) string {
	return func(t *types.Type) string {
		return aliasDisplayName(t, knownTypes)
	}
}

// extractCommentTags parses comments for lines of the form:
//
//	'marker' + "key=value".
//
// Values are optional; "" is the default.  A tag can be specified more than
// one time and all values are returned.  If the resulting map has an entry for
// a key, the value (a slice) is guaranteed to have at least 1 element.
//
// Example: if you pass "+" for 'marker', and the following lines are in
// the comments:
//
//	+foo=value1
//	+bar
//	+foo=value2
//	+baz="qux"
//
// Then this function will return:
//
//	map[string][]string{"foo":{"value1, "value2"}, "bar": {""}, "baz": {"qux"}}
//
// Source: https://github.com/kubernetes/gengo/blob/1244d31929d7/types/comments.go#L42
func extractCommentTags(marker string, lines []string) map[string][]string {
	out := map[string][]string{}
	for _, line := range lines {
		line = strings.Trim(line, " ")
		if len(line) == 0 {
			continue
		}
		if !strings.HasPrefix(line, marker) {
			continue
		}
		// TODO: we could support multiple values per key if we split on spaces
		kv := strings.SplitN(line[len(marker):], "=", 2)
		if len(kv) == 2 {
			out[kv[0]] = append(out[kv[0]], kv[1])
		} else if len(kv) == 1 {
			out[kv[0]] = append(out[kv[0]], "")
		}
	}
	return out
}

// aliasDisplayName allows types to replace their alias with an alternate
// alias display name.
// This can be useful when a type has a custom marshalling rule.
func aliasDisplayName(t *types.Type, knownTypes typeSet) string {
	tags := extractCommentTags("+", t.CommentLines)
	if alias, ok := tags["reference-gen:alias-name"]; ok {
		// There should only be one entry
		return alias[0]
	}

	if t.Underlying != nil {
		return typeDisplayName(t.Underlying, knownTypes)
	}

	return ""
}

// anchorIDForLocalType returns the #anchor string for the local type
func anchorIDForLocalType(t *types.Type) string {
	return strings.ToLower(t.Name.Name)
}

// backtick wraps the text in backticks
func backtick(s string) string {
	return "`" + s + "`"
}

// fieldEmbedded detemines if the field is embedded or not
func fieldEmbedded(m types.Member) bool {
	return m.Embedded
}

// fieldName extracts the field name from the json tag
func fieldName(m types.Member) string {
	vj := reflect.StructTag(m.Tags).Get("json")
	vj = strings.Split(vj, ",")[0]

	if vj != "" {
		return vj
	}

	vy := reflect.StructTag(m.Tags).Get("yaml")
	vy = strings.Split(vy, ",")[0]
	if vy != "" {
		return vy
	}

	return m.Name
}

// filterCommentTags removes comment hints before they are rendered
func filterCommentTags(comments []string) []string {
	var out []string
	for _, v := range comments {
		if !strings.HasPrefix(strings.TrimSpace(v), "+") {
			out = append(out, v)
		}
	}
	return out
}

// hideMember determines if a member is to private
func hideMember(m types.Member) bool {
	return unicode.IsLower(rune(m.Name[0]))
}

// hideType determines if a type is to private
func hideType(t *types.Type) bool {
	return unicode.IsLower(rune(t.Name.Name[0]))
}

// isOptionalMember determines if a member is marked optional
func isOptionalMember(m types.Member) bool {
	tags := extractCommentTags("+", m.CommentLines)
	_, ok := tags["optional"]
	return ok
}

// linkForTypeFunc constructs a linkForType function for the template
func linkForTypeFunc(knownTypes typeSet) func(t *types.Type) string {
	return func(t *types.Type) string {
		return linkForType(t, knownTypes)
	}
}

// linkForType returns an anchor to the type if it can be generated. returns
// empty string if it is not a local type or unrecognized external type.
func linkForType(t *types.Type, knownTypes typeSet) string {
	if t == nil {
		return ""
	}

	t = tryDereference(t) // dereference kind=Pointer

	if knownTypes.has(t) {
		return "#" + anchorIDForLocalType(t)
	}

	return ""
}

// renderComments filters comments and joins them to a single string using the
// join sequence provided.
func renderComments(s []string, join string) string {
	s = filterCommentTags(s)
	doc := strings.Join(s, join)
	return doc
}

func renderCommentsBR(s []string) string {
	return renderComments(s, "<br/>")
}

func renderCommentsLF(s []string) string {
	return renderComments(s, "\n")
}

// sortTypes sorts types alphabetically
func sortTypes(typs []*types.Type) []*types.Type {
	sort.Slice(typs, func(i, j int) bool {
		t1, t2 := typs[i], typs[j]
		return t1.Name.Name < t2.Name.Name
	})
	return typs
}

// typeDisplayNameFunc constructs a typeDisplayName function for the template
func typeDisplayNameFunc(knownTypes typeSet) func(t *types.Type) string {
	return func(t *types.Type) string {
		return typeDisplayName(t, knownTypes)
	}
}

// typeDisplayName works out the display name to be printed for a type
func typeDisplayName(t *types.Type, knownTypes typeSet) string {
	s := typeIdentifier(t)
	dt := tryDereference(t)
	if knownTypes.has(dt) {
		s = dt.Name.Name
	}
	if t.Kind == types.Pointer {
		s = strings.TrimLeft(s, "*")
	}

	switch t.Kind {
	case types.Struct,
		types.Interface,
		types.Alias,
		types.Pointer,
		types.Slice,
		types.Builtin:
		// noop
	case types.Map:
		// construct map based on element name
		return fmt.Sprintf("map[%s]%s", t.Key.Name.Name, s)
	default:
		klog.Fatalf("type %s has kind=%v which is unhandled", t.Name, t.Kind)
	}

	if alias, ok := commonTypes[s]; ok {
		s = alias
	}

	if t.Kind == types.Slice {
		s = "[]" + s
	}

	return s
}

// typeIdentifier produces the type ID in the form of <pkg>.<type>.
func typeIdentifier(t *types.Type) string {
	t = tryDereference(t)
	return t.Name.String() // {PackagePath.Name}
}

// typeReferencesFunc constructs a typeReferences function for the template
func typeReferencesFunc(references map[*types.Type][]*types.Type, knownTypes typeSet) func(t *types.Type) []*types.Type {
	return func(t *types.Type) []*types.Type {
		return typeReferences(t, references, knownTypes)
	}
}

// typeReferences creates a list of types that reference the given type
func typeReferences(t *types.Type, references map[*types.Type][]*types.Type, knownTypes typeSet) []*types.Type {
	out := []*types.Type{}
	for _, typ := range references[t] {
		if knownTypes.has(typ) {
			out = append(out, typ)
		}
	}
	sortTypes(out)
	return out
}

// visibleMembers filters the members to only those that are exported
func visibleMembers(in []types.Member) []types.Member {
	var out []types.Member
	for _, t := range in {
		if !hideMember(t) {
			out = append(out, t)
		}
	}
	return out
}

// visibleTypes filters the types to only those that are exported
func visibleTypes(in []*types.Type) []*types.Type {
	var out []*types.Type
	for _, t := range in {
		if !hideType(t) {
			out = append(out, t)
		}
	}
	return out
}

// END: template functions
