package jsonapi

import (
	"errors"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

// URL ...
type URL struct {
	Scheme string
	Host   string

	// URL
	Path      string   // /users/u1/articles?fields[users]=age,name
	Fragments []string // [users, u1, articles]
	Route     string   // /users/:id/articles
	Type      string   // col, res, related, self, meta

	// Resource
	BelongsToFilter BelongsToFilter
	ResType         string
	ResID           string
	RelKind         string
	Rel             Rel
	IsCol           bool

	// Params
	Params *Params
}

func NewURL() *URL {
	return &URL{
		Fragments: []string{},
		Params:    NewParams(),
	}
}

// BelongsToFilter ...
type BelongsToFilter struct {
	Type        string
	ID          string
	Name        string
	InverseName string
}

// NormalizePath ...
func (u *URL) NormalizePath() string {
	// Path
	path := "/"
	for _, p := range u.Fragments {
		path += p + "/"
	}
	path = path[:len(path)-1]

	// Params
	urlParams := []string{}

	// Fields
	for n := range u.Params.Fields {
		sort.Strings(u.Params.Fields[n])

		param := "fields%5B" + n + "%5D="
		for _, f := range u.Params.Fields[n] {
			param += f + "%2C"
		}
		param = param[:len(param)-3]

		urlParams = append(urlParams, param)
	}

	// Filters
	for n := range u.Params.RelFilters {
		sort.Strings(u.Params.RelFilters[n].IDs)

		param := "filter%5B" + n + "%5D="
		for _, id := range u.Params.RelFilters[n].IDs {
			param += id + "%2C"
		}
		param = param[:len(param)-3]

		urlParams = append(urlParams, param)
	}

	// TODO attribute filters

	// Pagination
	if u.IsCol {
		if u.Params.PageSize == 0 {
			u.Params.PageSize = 10
		}
		urlParams = append(urlParams, "page%5Bsize%5D="+strconv.FormatUint(uint64(u.Params.PageSize), 10))

		if u.Params.PageNumber == 0 {
			u.Params.PageNumber = 1
		}
		urlParams = append(urlParams, "page%5Bnumber%5D="+strconv.FormatUint(uint64(u.Params.PageNumber), 10))
	}

	// Sorting
	if len(u.Params.SortingRules) > 0 {
		param := "sort="
		for _, attr := range u.Params.SortingRules {
			param += attr + ","
		}
		param = param[:len(param)-1]

		urlParams = append(urlParams, param)
	}

	params := "?"
	for _, param := range urlParams {
		params += param + "&"
	}
	params = params[:len(params)-1]

	u.Path = path + params

	return u.Path
}

func (u *URL) FullURL() string {
	url := u.NormalizePath()

	if u.Scheme != "" && u.Host != "" {
		url = u.Scheme + "://" + u.Host + url
	}

	return url
}

// ParseRawURL ...
func ParseRawURL(reg *Registry, rawurl string) (*URL, error) {
	url, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	return ParseURL(reg, url)
}

// ParseURL ...
func ParseURL(reg *Registry, u *url.URL) (*URL, error) {
	url := &URL{
		Scheme: u.Scheme,
		Host:   u.Host,
	}

	// Path
	path := strings.TrimRight(u.Path, "/")
	path = strings.TrimLeft(path, "/")
	tempPaths := strings.Split(path, "/")
	paths := []string{}
	invalid := false

	for _, p := range tempPaths {
		if p != "" {
			paths = append(paths, p)
		} else {
			invalid = true
		}
	}

	if len(paths) == 0 || len(paths) > 4 {
		invalid = true
	}

	if invalid {
		return nil, errors.New("path is invalid")
	}

	url.Fragments = paths

	fromFilter := BelongsToFilter{}

	// Route
	url.Route = deduceRoute(url.Fragments)

	// Resource
	rel := Rel{}

	if len(paths) >= 1 {
		url.ResType = paths[0]
		url.IsCol = true
		url.Type = "col"
	}

	if len(paths) >= 2 {
		if paths[1] != "meta" {
			url.ResID = paths[1]
			url.IsCol = false
			url.Type = "res"
		} else {
			url.Type = "meta"
		}
	}

	if len(paths) >= 3 {
		fromFilter.Name = paths[2]
		fromFilter.Type = paths[0]
		fromFilter.ID = paths[1]

		if paths[2] == "relationships" {
			url.ResID = ""
			url.RelKind = "self"
			url.Type = "self"
		} else if paths[2] != "meta" {
			url.ResID = ""
			url.RelKind = "related"
			url.Type = "related"
			if r, ok := reg.Types[paths[0]].Rels[paths[2]]; ok {
				url.ResType = r.Type
				url.IsCol = !r.ToOne
			}
			rel.Name = paths[2]
		} else {
			url.Type = "meta"
		}
	}

	if len(paths) >= 4 {
		if url.RelKind == "self" {
			rel.Name = paths[3]
			fromFilter.Name = paths[3]
			if r, ok := reg.Types[paths[0]].Rels[paths[3]]; ok {
				url.ResType = r.Type
				url.IsCol = !r.ToOne
			}
		}
	}

	if len(paths) >= 5 {
		url.ResType = paths[0]
	}

	rel = reg.Types[paths[0]].Rels[rel.Name]

	// RelFiter
	if fromFilter.Type != "" {
		fromFilter.InverseName = reg.Types[fromFilter.Type].Rels[fromFilter.Name].InverseName
		url.ResType = reg.Types[fromFilter.Type].Rels[fromFilter.Name].Type
	}

	params, err := parseParams(reg, url.ResType, u)
	if err != nil {
		return url, err
	}
	url.Params = params

	url.BelongsToFilter = fromFilter
	url.Rel = rel

	url.NormalizePath()

	return url, nil
}

// parseParams ...
func parseParams(reg *Registry, resType string, u *url.URL) (*Params, error) {
	// Query parameters
	values := u.Query()

	params := &Params{
		Fields: map[string][]string{
			resType: []string{},
		},
		Attrs:        map[string][]Attr{},
		Rels:         map[string][]Rel{},
		RelData:      map[string][]string{},
		AttrFilters:  map[string]AttrFilter{},
		RelFilters:   map[string]RelFilter{},
		SortingRules: []string{},
		PageSize:     10,
		PageNumber:   1,
		Include:      [][]Rel{},
	}

	// Inclusions
	inclusions := []string{}
	// Get all values
	for _, vals := range values["include"] {
		incs := strings.Split(vals, ",")
		inclusions = append(inclusions, incs...)
	}
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
		var ok bool
		for _, word := range words {
			if incRel, ok = reg.Types[incRel.Type].Rels[word]; ok {
				params.Fields[incRel.Type] = []string{}
			} else {
				inclusions = append(inclusions[:i], inclusions[i+1:]...)
				break
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
				incRel = reg.Types[resType].Rels[words[0]]
			}

			params.Include[i][w] = incRel

			if w < len(words)-1 {
				incRel = reg.Types[incRel.Type].Rels[words[w+1]]
			}
		}
	}

	// Other params
	for param, vals := range values {
		if strings.HasPrefix(param, "filter[") {
			// Filters
			field := param[7 : len(param)-1]

			targets := []string{}
			for _, v := range strings.Split(vals[0], ",") {
				if v != "" {
					targets = append(targets, v)
				}
			}

			if r, ok := reg.Types[resType].Rels[field]; ok {
				params.RelFilters[field] = RelFilter{
					Type:        r.Type,
					InverseName: r.InverseName,
					IDs:         targets,
				}
			} else if a, ok := reg.Types[resType].Attrs[field]; ok {
				rules := []string{}

				if kind(a.Type) == "string" {
					for i := range targets {
						if len(targets[i]) > 2 {
							rules = append(rules, targets[i][:2])
							targets[i] = targets[i][2:]
						} else {
							panic("invalid url")
						}
					}
				} else if kind(a.Type) == "number" {

				} else if kind(a.Type) == "bool" {

				} else if kind(a.Type) == "time" {

				}

				params.AttrFilters[field] = AttrFilter{
					Type:    r.Type,
					Rules:   rules,
					Targets: targets,
				}
			} else {
				panic("relationship or attribute not found for filter")
			}
		} else if param == "sort" {
			// Sorting
			for _, val := range vals {
				if val != "" {
					for _, v := range strings.Split(val, ",") {
						if v != "" {
							attr := v
							if strings.HasPrefix(attr, "-") {
								attr = attr[1:]
							}
							if _, ok := reg.Types[resType].Attrs[attr]; ok {
								params.SortingRules = append(params.SortingRules, v)
							}
						}
					}
				}
			}
		} else if param == "page[size]" {
			// Page size
			if size, err := strconv.ParseInt(vals[0], 10, 64); err == nil {
				if size > 0 && size <= 10000 {
					params.PageSize = int(size)
				} else {
					params.PageSize = 10
				}
			}
		} else if param == "page[number]" {
			// Page number
			if number, err := strconv.ParseInt(vals[0], 10, 64); err == nil {
				if number > 0 && number <= 10000000 {
					params.PageNumber = int(number)
				} else {
					params.PageNumber = 1
				}
			}
		} else if strings.HasPrefix(param, "fields[") {
			// Fields
			resName := param[7 : len(param)-1]
			if _, ok := params.Fields[resName]; ok {
				for _, v := range strings.Split(vals[0], ",") {
					if v != "" {
						for _, f := range reg.Types[resName].Fields {
							if f == v {
								params.Fields[resName] = append(params.Fields[resName], v)
							}
						}
					}
				}
			}
		}
	}

	for resName := range params.Fields {
		if len(params.Fields[resName]) == 0 {
			params.Fields[resName] = reg.Types[resName].Fields
		}

		for _, field := range params.Fields[resName] {
			if attr, ok := reg.Types[resName].Attrs[field]; ok {
				params.Attrs[resName] = append(params.Attrs[resName], attr)
			} else if rel, ok := reg.Types[resName].Rels[field]; ok {
				params.Rels[resName] = append(params.Rels[resName], rel)
			}
		}
	}

	return params, nil
}

func deduceRoute(path []string) string {
	const (
		id   = "/:id"
		meta = "meta"
		rel  = "relationships"
	)

	route := ""

	if len(path) >= 1 {
		route = "/" + path[0]
	}

	if len(path) >= 2 {
		if path[1] == meta {
			route += "/" + meta
		} else {
			route += id
		}
	}

	if len(path) >= 3 {
		if path[2] == rel {
			route += "/" + rel
		} else if path[2] == meta {
			route += "/" + meta
		} else {
			route += "/" + path[2]
		}
	}

	if len(path) >= 4 {
		if path[3] == meta {
			route += "/" + meta
		} else {
			if path[2] == rel {
				route += "/" + path[3]
			} else {
				route += id
			}
		}
	}

	if len(path) >= 5 {
		if path[4] == meta {
			route += "/" + meta
		} else {
			route += id
		}
	}

	return route
}

func stringifyParams(params map[string][]string, wrapper string) []string {
	strParams := []string{}
	for key, vals := range params {
		if wrapper != "" {
			key = wrapper + "%5B" + key + "%5D"
		}
		param := key + "="
		sort.Strings(vals)
		for i, v := range vals {
			if i < len(vals)-1 {
				v += "%2C"
			}
			param += v
		}
		strParams = append(strParams, param)
	}

	return strParams
}

func kind(typ string) string {
	switch typ {
	case "string", "*string":
		return "string"

	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "*int", "*int8", "*int16", "*int32", "*int64", "*uint", "*uint8", "*uint16", "*uint32", "*uint64":
		return "int"

	case "bool", "*bool":
		return "bool"

	case "time.Time", "*time.Time":
		return "time"
	}

	panic("cannot find kind of the provided type")
}
