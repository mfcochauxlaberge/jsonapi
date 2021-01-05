package jsonapi

import (
	"sort"
	"strings"
)

// NewParams creates and returns a Params object built from a SimpleURL and a
// given resource type. A schema is used for validation.
//
// If validation is not expected, it is recommended to simply build a SimpleURL
// object with NewSimpleURL.
func NewParams(schema *Schema, su SimpleURL, resType string) (*Params, error) {
	params := &Params{
		Fields:       map[string][]string{},
		Attrs:        map[string][]Attr{},
		Rels:         map[string][]Rel{},
		RelData:      map[string][]string{},
		Filter:       nil,
		SortingRules: []string{},
		Include:      [][]Rel{},
	}

	// Include
	incs := make([]string, len(su.Include))
	copy(incs, su.Include)
	sort.Strings(incs)

	// Remove duplicates and uncessary includes
	for i := len(incs) - 1; i >= 0; i-- {
		if i > 0 {
			if strings.HasPrefix(incs[i], incs[i-1]) {
				incs = append(incs[:i-1], incs[i:]...)
			}
		}
	}

	// Check inclusions
	for i := 0; i < len(incs); i++ {
		words := strings.Split(incs[i], ".")

		incRel := Rel{ToType: resType}

		for _, word := range words {
			if typ := schema.GetType(incRel.ToType); typ.Name != "" {
				var ok bool
				if incRel, ok = typ.Rels[word]; ok {
					params.Fields[incRel.ToType] = []string{}
				} else {
					incs = append(incs[:i], incs[i+1:]...)
					break
				}
			}
		}
	}

	// Build params.Include
	params.Include = make([][]Rel, len(incs))

	for i := range incs {
		words := strings.Split(incs[i], ".")

		params.Include[i] = make([]Rel, len(words))

		var incRel Rel

		for w := range words {
			if w == 0 {
				typ := schema.GetType(resType)
				incRel = typ.Rels[words[0]]
			}

			params.Include[i][w] = incRel

			if w < len(words)-1 {
				typ := schema.GetType(incRel.ToType)
				incRel = typ.Rels[words[w+1]]
			}
		}
	}

	if resType != "" {
		params.Fields[resType] = []string{}
	}

	// Fields
	for t, fields := range su.Fields {
		if t != resType {
			if typ := schema.GetType(t); typ.Name == "" {
				return nil, NewErrUnknownTypeInURL(t)
			}
		}

		if typ := schema.GetType(t); typ.Name != "" {
			params.Fields[t] = []string{}

			for _, f := range fields {
				if f == "id" {
					params.Fields[t] = append(params.Fields[t], "id")
				} else {
					for _, ff := range typ.Fields() {
						if f == ff {
							params.Fields[t] = append(params.Fields[t], f)
						}
					}
				}
			}
			// Check for duplicates
			for i := range params.Fields[t] {
				for j := i + 1; j < len(params.Fields[t]); j++ {
					if params.Fields[t][i] == params.Fields[t][j] {
						return nil, NewErrDuplicateFieldInFieldsParameter(
							typ.Name,
							params.Fields[t][i],
						)
					}
				}
			}
		}
	}

	for t := range params.Fields {
		if len(params.Fields[t]) == 0 {
			typ := schema.GetType(t)
			params.Fields[t] = make([]string, len(typ.Fields()))
			copy(params.Fields[t], typ.Fields())
		}
	}

	// Attrs and Rels
	for typeName, fields := range params.Fields {
		// This should always return a type since
		// it is checked earlier.
		typ := schema.GetType(typeName)

		params.Attrs[typeName] = make([]Attr, 0, len(typ.Attrs))
		params.Rels[typeName] = make([]Rel, 0, len(typ.Attrs))

		for _, field := range typ.Fields() {
			for _, field2 := range fields {
				if field == field2 {
					if typ = schema.GetType(typeName); typ.Name != "" {
						if attr, ok := typ.Attrs[field]; ok {
							// Append to list of attributes
							params.Attrs[typeName] = append(
								params.Attrs[typeName],
								attr,
							)
						} else if rel, ok := typ.Rels[field]; ok {
							// Append to list of relationships
							params.Rels[typeName] = append(
								params.Rels[typeName],
								rel,
							)
						}
					}
				}
			}
		}
	}

	// Filter
	params.FilterLabel = su.FilterLabel
	params.Filter = su.Filter
	// TODO Check whether the filter is valid

	// Sorting
	// TODO All of the following is just to figure out
	// if the URL represents a single resource or a
	// collection. It should be done in a better way.
	isCol := false
	if len(su.Fragments) == 1 {
		isCol = true
	} else if len(su.Fragments) >= 3 {
		relName := su.Fragments[len(su.Fragments)-1]
		typ := schema.GetType(su.Fragments[0])
		// Checked earlier, assuming should be safe
		rel := typ.Rels[relName]
		isCol = !rel.ToOne
	}

	if isCol {
		typ := schema.GetType(resType)
		sortingRules := make([]string, 0, len(typ.Attrs))
		idFound := false

		for _, rule := range su.SortingRules {
			urule := rule
			if urule[0] == '-' {
				urule = urule[1:]
			}

			if urule == "id" {
				idFound = true

				sortingRules = append(sortingRules, rule)

				break
			}

			for _, attr := range typ.Attrs {
				if urule == attr.Name {
					sortingRules = append(sortingRules, rule)
					break
				}
			}
		}

		// Add 1 because of id
		restOfRules := make([]string, 0, len(typ.Attrs)+1-len(sortingRules))

		for _, attr := range typ.Attrs {
			found := false

			for _, rule := range sortingRules {
				urule := rule
				if urule[0] == '-' {
					urule = urule[1:]
				}

				if urule == attr.Name {
					found = true
					break
				}
			}

			if !found {
				restOfRules = append(restOfRules, attr.Name)
			}
		}

		sort.Strings(restOfRules)
		sortingRules = append(sortingRules, restOfRules...)

		if !idFound {
			sortingRules = append(sortingRules, "id")
		}

		params.SortingRules = sortingRules
	}

	// Pagination
	params.PageSize = su.PageSize
	params.PageNumber = su.PageNumber

	return params, nil
}

// A Params object represents all the query parameters from the URL.
type Params struct {
	// Fields
	Fields  map[string][]string
	Attrs   map[string][]Attr
	Rels    map[string][]Rel
	RelData map[string][]string

	// Filter
	FilterLabel string
	Filter      *Filter

	// Sorting
	SortingRules []string

	// Pagination
	PageSize   uint
	PageNumber uint

	// Include
	Include [][]Rel
}
