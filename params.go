package jsonapi

import (
	"sort"
	"strconv"
	"strings"
)

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
	PageSize   int
	PageNumber int

	// Include
	Include [][]Rel
}

// NewParams creates and returns a Params object built from a SimpleURL
// and a given resource type. A schema is used for validation.
//
// If validation is not expected, it is recommended to simply build a
// SimpleURL object with NewSimpleURL.
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
	inclusions := make([]string, len(su.Include))
	copy(inclusions, su.Include)
	sort.Strings(inclusions)

	// Remove duplicates and uncessary includes
	for i := len(inclusions) - 1; i >= 0; i-- {
		if i > 0 {
			if strings.HasPrefix(inclusions[i], inclusions[i-1]) {
				inclusions = append(inclusions[:i-1], inclusions[i:]...)
			}
		}
	}

	// Check inclusions
	for i := 0; i < len(inclusions); i++ {
		words := strings.Split(inclusions[i], ".")

		incRel := Rel{Type: resType}
		for _, word := range words {
			if typ, ok := schema.GetType(incRel.Type); ok {
				if incRel, ok = typ.Rels[word]; ok {
					params.Fields[incRel.Type] = []string{}
				} else {
					inclusions = append(inclusions[:i], inclusions[i+1:]...)
					break
				}
			}
		}
	}

	// Build params.Include
	params.Include = make([][]Rel, len(inclusions))
	for i := range inclusions {
		words := strings.Split(inclusions[i], ".")

		params.Include[i] = make([]Rel, len(words))

		var incRel Rel
		for w := range words {
			if w == 0 {
				typ, _ := schema.GetType(resType)
				incRel = typ.Rels[words[0]]
			}

			params.Include[i][w] = incRel

			if w < len(words)-1 {
				typ, _ := schema.GetType(incRel.Type)
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
			if _, ok := schema.GetType(t); !ok {
				return nil, NewErrUnknownTypeInURL(t)
			}
		}
		if typ, ok := schema.GetType(t); ok {
			params.Fields[t] = []string{}
			for _, f := range fields {
				for _, ff := range typ.Fields() {
					if f == ff {
						params.Fields[t] = append(params.Fields[t], f)
					}
				}
			}
		}
	}
	for t := range params.Fields {
		if len(params.Fields[t]) == 0 {
			typ, _ := schema.GetType(t)
			params.Fields[t] = make([]string, len(typ.Fields()))
			copy(params.Fields[t], typ.Fields())
		}
	}

	// Attrs and Rels
	for typeName, fields := range params.Fields {
		var (
			typ  Type
			attr Attr
			rel  Rel
			ok   bool
		)
		if typ, ok = schema.GetType(typeName); !ok {
			return nil, NewErrUnknownTypeInURL(typeName)
		}

		params.Attrs[typeName] = []Attr{}
		params.Rels[typeName] = []Rel{}

		for _, field := range typ.Fields() {
			for _, field2 := range fields {
				if field == field2 {
					// Append to list of fields
					// params.Fields[typeName] = append(params.Fields[typeName], field)

					if typ, ok = schema.GetType(typeName); ok {
						if attr, ok = typ.Attrs[field]; ok {
							// Append to list of attributes
							params.Attrs[typeName] = append(params.Attrs[typeName], attr)
						}
					}

					if typ, ok = schema.GetType(typeName); ok {
						if rel, ok = typ.Rels[field]; ok {
							// Append to list of relationships
							params.Rels[typeName] = append(params.Rels[typeName], rel)
						}
					}
				}
			}
		}
	}

	// Filter
	params.FilterLabel = su.FilterLabel
	params.Filter = su.Filter
	// TODO

	// Sorting
	// TODO All of the following is just to figure out
	// if the URL represents a single resource or a
	// collection. It should be done in a better way.
	isCol := false
	if len(su.Fragments) == 1 {
		isCol = true
	} else if len(su.Fragments) >= 3 {
		relName := su.Fragments[len(su.Fragments)-1]
		typ, _ := schema.GetType(su.Fragments[0])
		var (
			rel Rel
			ok  bool
		)
		if rel, ok = typ.Rels[relName]; !ok {
			return nil, NewErrUnknownRelationshipInPath(typ.Name, relName, su.Path())
		}
		isCol = !rel.ToOne
	}
	if isCol {
		typ, _ := schema.GetType(resType)
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
	if params.PageSize < 0 || params.PageSize > 100 {
		return nil, NewErrInvalidPageSizeParameter(strconv.FormatInt(int64(params.PageSize), 10))
	}
	params.PageNumber = su.PageNumber
	if params.PageNumber < 0 {
		return nil, NewErrInvalidPageNumberParameter(strconv.FormatInt(int64(params.PageNumber), 10))
	}

	return params, nil
}
