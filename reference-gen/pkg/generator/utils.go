package generator

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
