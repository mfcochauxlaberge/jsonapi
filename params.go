package jsonapi

// A Params object represents all the query parameters from the URL.
type Params struct {
	// Fields
	Fields  map[string][]string
	Attrs   map[string][]Attr
	Rels    map[string][]Rel
	RelData map[string][]string

	// Filters
	AttrFilters map[string]AttrFilter
	RelFilters  map[string]RelFilter

	// Sorting
	SortingRules []string

	// Pagination
	PageSize   int
	PageNumber int

	// Include
	Include [][]Rel
}

// NewParams ...
func NewParams() *Params {
	return &Params{
		Fields:       map[string][]string{},
		Attrs:        map[string][]Attr{},
		Rels:         map[string][]Rel{},
		RelData:      map[string][]string{},
		AttrFilters:  map[string]AttrFilter{},
		RelFilters:   map[string]RelFilter{},
		SortingRules: []string{},
		Include:      [][]Rel{},
	}
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
