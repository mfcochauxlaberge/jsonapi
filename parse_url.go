package jsonapi

import (
	"errors"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

// ParseURL ...
func ParseURL(reg *Registry, u *url.URL) (*URL, error) {
	url := &URL{}

	// Path
	path := strings.TrimRight(u.Path, "/")
	tempPaths := strings.Split(path, "/")
	paths := []string{}
	invalid := false

	for i, p := range tempPaths {
		if p != "" {
			paths = append(paths, p)
		} else {
			if i != 0 && i != len(tempPaths)-1 {
				invalid = true
			}
		}
	}

	if len(paths) == 0 || len(paths) > 4 {
		invalid = true
	}

	if invalid {
		return nil, errors.New("url is invalid")
	}

	url.Path = paths

	fromFilter := FromFilter{}

	// Router
	url.Route = deduceRoute(url.Path)

	// Resource
	const (
		meta = "meta"
	)

	rel := Rel{}

	if len(paths) >= 1 {
		url.ResType = paths[0]
	}

	if len(paths) >= 2 {
		if paths[1] != meta {
			url.ResID = paths[1]
		}
	}

	if len(paths) >= 3 {
		if paths[2] == "relationships" {
			url.RelKind = "self"
		} else if paths[2] != meta {
			url.RelKind = "related"
			if r, ok := reg.Types[paths[0]].Rels[paths[2]]; ok {
				url.ResType = r.Type
			}
			rel.Name = paths[2]
			fromFilter.Name = paths[2]
			fromFilter.Type = paths[0]
			fromFilter.ID = paths[1]
		}
	}

	if len(paths) >= 4 {
		if url.RelKind == "self" {
			rel.Name = paths[3]
		}
	}

	if len(paths) >= 5 {
		url.ResType = paths[0]
	}

	rel = reg.Types[url.ResType].Rels[rel.Name]

	// RelFiter
	if fromFilter.Type != "" {
		fromFilter.InverseName = reg.Types[fromFilter.Type].Rels[fromFilter.Name].InverseName
		url.ResType = reg.Types[fromFilter.Type].Rels[fromFilter.Name].Type
	}

	// Query parameters
	values := u.Query()

	// ctx.Types = map[string]struct{}{ctx.URL.ResType: struct{}{}}

	// attrs := map[string][]Attr{
	// 	ctx.URL.ResType: []Attr{},
	// }
	fields := map[string][]string{
		url.ResType: []string{},
	}
	attrFilters := map[string]AttrFilter{}
	relFilters := map[string]RelFilter{}
	sorting := []string{}
	pagination := map[string]uint64{
		"size":   1000,
		"number": 1,
	}
	includes := []string{}

	// Inclusions
	for _, vals := range values["include"] {
		// Remove duplicates and uncessary includes
		includes = strings.Split(vals, ",")
		sort.Strings(includes)
		for i := len(includes) - 1; i >= 0; i-- {
			if i > 0 {
				if strings.HasPrefix(includes[i], includes[i-1]) {
					includes = append(includes[:i-1], includes[i:]...)
				}
			}
		}

		for i, inc := range includes {
			words := strings.Split(inc, ".")

			incRel := Rel{Type: url.ResType}
			var ok bool
			for _, word := range words {
				if incRel, ok = reg.Types[incRel.Type].Rels[word]; ok {
					fields[incRel.Type] = []string{}
				} else {
					includes = append(includes[:i], includes[i+1:]...)
					break
				}
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

			if r, ok := reg.Types[url.ResType].Rels[field]; ok {
				relFilters[field] = RelFilter{
					Type:        r.Type,
					InverseName: r.InverseName,
					IDs:         targets,
				}
			} else if a, ok := reg.Types[url.ResType].Attrs[field]; ok {
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

				attrFilters[field] = AttrFilter{
					Type:    r.Type,
					Rules:   rules,
					Targets: targets,
				}
			} else {
				panic("relationship or attribute not found for filter")
			}
		} else if param == "sort" {
			// Sorting
			if vals[0] != "" {
				for _, v := range strings.Split(vals[0], ",") {
					if v != "" {
						attr := v
						if strings.HasPrefix(attr, "-") {
							attr = attr[1:]
						}
						if _, ok := reg.Types[url.ResType].Attrs[attr]; ok {
							sorting = append(sorting, v)
						}
					}
				}
			}
		} else if param == "page[size]" {

			// Page size
			if size, err := strconv.ParseUint(vals[0], 10, 64); err == nil {
				if size > 0 && size <= 100 {
					pagination["size"] = size
				}
			}
		} else if param == "page[number]" {
			// Page number
			if number, err := strconv.ParseUint(vals[0], 10, 64); err == nil {
				if number > 0 && number <= 10000 {
					pagination["number"] = number
				}
			}
		} else if strings.HasPrefix(param, "fields[") {
			// Fields
			resName := param[7 : len(param)-1]
			if _, ok := fields[resName]; ok {
				for _, v := range strings.Split(vals[0], ",") {
					if v != "" {
						for _, f := range reg.Types[resName].Fields {
							if f == v {
								fields[resName] = append(fields[resName], v)
							}
						}
					}
				}
			}
		}
	}

	attrs := map[string][]Attr{}
	// rels := map[string][]Rel{} // TODO
	for resName := range fields {
		if len(fields[resName]) == 0 {
			fields[resName] = reg.Types[resName].Fields
		}

		for _, field := range fields[resName] {
			if attr, ok := reg.Types[resName].Attrs[field]; ok {
				attrs[resName] = append(attrs[resName], attr)
			} else if _, ok := reg.Types[resName].Rels[field]; ok {
				// rels[resName] = append(rels[resName], rel) // TODO
			}
		}
	}

	// Set all params
	params := &Params{}
	params.Fields = fields
	params.Attrs = attrs
	// params.Rels = rels // TODO
	params.FromFilter = fromFilter
	params.AttrFilters = attrFilters
	params.RelFilters = relFilters
	params.SortingRules = sorting
	params.PageSize = uint(pagination["size"])
	params.PageNumber = uint(pagination["number"])
	params.Include = make([][]Rel, len(includes))
	for i := range includes {
		words := strings.Split(includes[i], ".")

		params.Include[i] = make([]Rel, len(words))

		var incRel Rel
		for w := range words {
			if w == 0 {
				incRel = reg.Types[url.ResType].Rels[words[0]]
			}

			params.Include[i][w] = incRel

			if w < len(words)-1 {
				incRel = reg.Types[incRel.Type].Rels[words[w+1]]
			}
		}
	}
	url.Params = params
	// ctx.IncludeRelIdentifiers = map[string]bool{}
	// for f := range fields {
	// 	ctx.Options.IncludeRelIdentifiers[f] = true
	// }

	// Normalize URL
	urlParams := []string{}

	for k := range fields {
		if len(fields[k]) == len(reg.Types[k].Fields) {
			delete(fields, k)
		}
	}
	urlParams = append(urlParams, stringifyParams(fields, "fields")...)

	filtersParams := map[string][]string{}
	for n, f := range params.RelFilters {
		filtersParams[n] = f.IDs
	}
	urlParams = append(urlParams, stringifyParams(filtersParams, "filters")...)

	urlParams = append(urlParams, stringifyParams(map[string][]string{
		"size":   []string{strconv.FormatUint(pagination["size"], 10)},
		"number": []string{strconv.FormatUint(pagination["number"], 10)},
	}, "page")...)

	if len(sorting) > 0 {
		sortParam := "sort="
		for i, v := range sorting {
			if i < len(sorting)-1 {
				v += ","
			}
			sortParam += v
		}
		urlParams = append(urlParams, sortParam)
	}

	if len(includes) > 0 {
		urlParams = append(urlParams, stringifyParams(map[string][]string{
			"include": includes,
		}, "")...)
	}

	sort.Strings(urlParams)
	normURL := path
	if len(urlParams) > 0 {
		normURL += "?"
	}
	for _, param := range urlParams {
		normURL += param + "&"
	}
	normURL = strings.TrimSuffix(normURL, "&")
	url.URL = u.String()
	url.URLNormalized = normURL

	url.Rel = rel

	return url, nil
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
			key = wrapper + "[" + key + "]"
		}
		param := key + "="
		sort.Strings(vals)
		for i, v := range vals {
			if i < len(vals)-1 {
				v += ","
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
