package jsonapi

// A Params object represents all the query parameters from the URL.
type Params struct {
	// Fields
	Fields  map[string][]string
	Attrs   map[string][]Attr
	Rels    map[string][]Rel
	RelData map[string][]string

	// Filters
	FromFilter  FromFilter
	AttrFilters map[string]AttrFilter
	RelFilters  map[string]RelFilter

	// Sorting
	SortingRules []string

	// Pagination
	PageSize   uint
	PageNumber uint

	// Include
	Include [][]Rel
}

// FromFilter ...
type FromFilter struct {
	Name        string
	Type        string
	InverseName string
	ID          string
}

// AttrFilter ...
type AttrFilter struct {
	Type    string
	Rules   []string
	Targets []string
}

// RelFilter ...
type RelFilter struct {
	Type        string
	InverseName string
	IDs         []string
}
