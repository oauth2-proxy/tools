package generator

import "k8s.io/gengo/types"

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

func newTypeSetFromStringMap(typeMap map[string]*types.Type) typeSet {
	set := make(typeSet)
	for _, t := range typeMap {
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
