package jsonapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
)

// URL ...
type URL struct {
	// URL
	Scheme    string
	Host      string
	Port      string
	Fragments []string // [users, u1, articles]
	Route     string   // /users/:id/articles

	// Data
	IsCol           bool
	ResType         string
	ResID           string
	RelKind         string
	Rel             Rel
	BelongsToFilter BelongsToFilter

	// Params
	Params *Params
}

// NewURL ...
func NewURL(reg *Registry, su SimpleURL) (*URL, error) {
	url := &URL{}

	// Scheme
	url.Scheme = su.Scheme
	if url.Scheme != "http" && url.Scheme != "https" {
		url.Scheme = ""
	}

	// Host
	url.Host = su.Host

	// Port
	url.Port = su.Port

	// Route
	url.Route = su.Route

	// Fragments
	url.Fragments = su.Fragments

	// IsCol, ResType, ResID, RelKind, Rel, BelongsToFilter
	var typ Type
	var ok bool
	if len(url.Fragments) == 0 {
		return nil, NewErrBadRequest("Empty path", "There is no path.")
	}
	if len(url.Fragments) >= 1 {
		if typ, ok = reg.Types[url.Fragments[0]]; !ok {
			return nil, NewErrUnknownTypeInURL(url.Fragments[0])
		}

		if len(url.Fragments) == 1 {
			url.IsCol = true
			url.ResType = typ.Name
		}

		if len(url.Fragments) == 2 {
			url.IsCol = false
			url.ResType = typ.Name
			url.ResID = url.Fragments[1]
		}
	}
	if len(url.Fragments) >= 3 {
		relName := url.Fragments[len(url.Fragments)-1]
		if url.Rel, ok = typ.Rels[relName]; !ok {
			return nil, NewErrUnknownRelationshipInPath(typ.Name, relName, su.Path())
		}

		url.IsCol = !url.Rel.ToOne
		url.ResType = url.Rel.Type
		url.BelongsToFilter = BelongsToFilter{
			Type:        url.Fragments[0],
			ID:          url.Fragments[1],
			Name:        url.Rel.Name,
			InverseName: url.Rel.InverseName,
		}

		if len(url.Fragments) == 3 {
			url.RelKind = "related"
		} else if len(url.Fragments) == 4 {
			url.RelKind = "self"
		}
	}

	// Params
	var err error
	url.Params, err = NewParams(reg, su, url.ResType)
	if err != nil {
		return nil, err
	}

	return url, nil
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

	// Filter
	if u.Params.Filter != nil {
		mf, err := json.Marshal(u.Params.Filter)
		if err != nil {
			// This should not happen since Filter should be validated
			// at this point.
			panic(fmt.Errorf("jsonapi: can't marshal filter: %s\n", err))
		}
		param := "filter=" + string(mf)
		urlParams = append(urlParams, param)
	}

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
			param += attr + "%2C"
		}
		param = param[:len(param)-3]

		urlParams = append(urlParams, param)
	}

	params := "?"
	for _, param := range urlParams {
		params += param + "&"
	}
	params = params[:len(params)-1]

	return path + params
}

// FullURL ...
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

	su, err := NewSimpleURL(url)
	if err != nil {
		return nil, err
	}

	return NewURL(reg, su)
}
