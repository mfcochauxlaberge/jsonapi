package jsonapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
)

// A URL stores all the information from a URL formatted for a JSON:API
// request.
//
// The data structure allows to have more information than what the URL
// itself stores.
type URL struct {
	// URL
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

// NewURL builds a URL from a SimpleURL and a schema for validating and
// supplementing the object with extra information.
func NewURL(schema *Schema, su SimpleURL) (*URL, error) {
	url := &URL{}

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
		if typ = schema.GetType(url.Fragments[0]); typ.Name == "" {
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
	url.Params, err = NewParams(schema, su, url.ResType)
	if err != nil {
		return nil, err
	}

	return url, nil
}

// NewURLFromRaw parses rawurl to make a *url.URL before making and returning
// a *URL.
func NewURLFromRaw(schema *Schema, rawurl string) (*URL, error) {
	url, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	su, err := NewSimpleURL(url)
	if err != nil {
		return nil, err
	}

	return NewURL(schema, su)
}

// A BelongsToFilter represents a parent resource, used to filter out
// resources that are not children of the parent.
//
// For example, in /articles/abc123/comments, the parent is the article
// with the ID abc123.
type BelongsToFilter struct {
	Type        string
	ID          string
	Name        string
	InverseName string
}

// String returns a string representation of the URL where special characters
// are escaped.
//
// The URL is normalized, so it always returns exactly the same string given
// the same URL.
func (u *URL) String() string {
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
			panic(fmt.Errorf("jsonapi: can't marshal filter: %s", err))
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
