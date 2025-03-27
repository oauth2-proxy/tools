package generator

import (
	"fmt"

	"k8s.io/gengo/v2/parser"
	"k8s.io/gengo/v2/types"
)

// loadPackage loads and parses the given package.
func loadPackage(packageName string) (*types.Package, error) {
	b := parser.New()
	// the following may silently fail (turn on -v=4 to see logs)
	if err := b.LoadPackages(packageName); err != nil {
		return nil, err
	}

	universe, err := b.NewUniverse()
	if err != nil {
		return nil, fmt.Errorf("failed to find types for package: %v", err)
	}

	pkg := universe.Package(packageName)
	if pkg == nil {
		return nil, fmt.Errorf("package %q was not found by parser", packageName)
	}

	return pkg, nil
}

// isReferenceRequired determines if the type needs a reference generated.
func isReferenceRequired(t *types.Type, requiredTypes stringSet, allReferences map[*types.Type][]*types.Type) bool {
	if requiredTypes.has(t.Name.Name) {
		return true
	}
	for _, reference := range allReferences[t] {
		if isReferenceRequired(reference, requiredTypes, allReferences) {
			return true
		}
	}
	return false
}

// filterToRequestedTypes filters the type references given to only those that
// are requested or referenced by a requested type.
func filterToRequestedTypes(allTypes map[*types.Type][]*types.Type, requestedTypes stringSet) map[*types.Type][]*types.Type {
	filteredTypes := make(map[*types.Type][]*types.Type)
	for typ, ref := range allTypes {
		if isReferenceRequired(typ, requestedTypes, allTypes) {
			filteredTypes[typ] = ref
		}
	}

	return filteredTypes
}

// filterToPackageTypes filters the given references to only those present in the typeSet.
func filterToPackageTypes(allReferences map[*types.Type][]*types.Type, pkgTypes typeSet) map[*types.Type][]*types.Type {
	filteredTypes := make(map[*types.Type][]*types.Type)
	for typ, refs := range allReferences {
		if pkgTypes.has(typ) {
			filteredTypes[typ] = refs
		}
	}
	return filteredTypes
}

// findTypeReferences converts a list of types to a map of types and types that
// reference that type.
func findTypeReferences(allTypes map[string]*types.Type) map[*types.Type][]*types.Type {
	m := make(map[*types.Type]typeSet)
	for _, typ := range allTypes {
		// Ensure every type is initialised, if not already
		if _, ok := m[typ]; !ok {
			m[typ] = make(typeSet)
		}

		// add this type to other types that it references
		for _, member := range typ.Members {
			if hideMember(member) {
				// Don't include a reference if the member is private.
				continue
			}
			t := member.Type
			t = tryDereference(t)
			if _, ok := m[t]; !ok {
				m[t] = make(typeSet)
			}
			m[t].add(typ)
		}

		// Cater for aliases rather than structs
		if typ.Underlying != nil {
			t := tryDereference(typ.Underlying)
			if _, ok := m[t]; !ok {
				m[t] = make(typeSet)
			}
			m[t].add(typ)
		}
	}

	out := make(map[*types.Type][]*types.Type)
	for t, s := range m {
		out[t] = s.toList()
	}
	return out
}
