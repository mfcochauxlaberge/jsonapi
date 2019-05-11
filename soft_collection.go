package jsonapi

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

// SoftCollection ...
type SoftCollection struct {
	typ  *Type
	col  []*SoftResource
	sort []string

	m sync.Mutex
}

// Type ...
func (s *SoftCollection) Type() Type {
	return *s.typ
}

// SetType ...
func (s *SoftCollection) SetType(typ Type) {
	s.typ = &typ
}

// AddAttr ...
func (s *SoftCollection) AddAttr(attr Attr) error {
	return s.typ.AddAttr(attr)
}

// AddRel ...
func (s *SoftCollection) AddRel(rel Rel) error {
	return s.typ.AddRel(rel)

}

// Len ...
func (s *SoftCollection) Len() int {
	return len(s.col)
}

// Elem ...
func (s *SoftCollection) Elem(i int) Resource {
	if i >= 0 && i < len(s.col) {
		return s.col[i]
	}
	return nil
}

// Add ...
func (s *SoftCollection) Add(r Resource) {
	// A SoftResource is built from the Resource and
	// then it is added to the collection.
	sr := &SoftResource{}
	sr.id = r.GetID()
	sr.typ = s.typ

	for _, attr := range r.Attrs() {
		sr.AddAttr(attr)
		sr.Set(attr.Name, r.Get(attr.Name))
	}

	for _, rel := range r.Rels() {
		sr.AddRel(rel)
		if rel.ToOne {
			sr.SetToOne(rel.Name, r.GetToOne(rel.Name))
		} else {
			sr.SetToMany(rel.Name, r.GetToMany(rel.Name))
		}
	}

	s.col = append(s.col, sr)
}

// UnmarshalJSON ...
func (s *SoftCollection) UnmarshalJSON(payload []byte) error {
	return errors.New("jsonapi: SoftCollection.UnmarshalJSON unimplemented")
}

// Sort ...
func (s *SoftCollection) Sort(rules []string) {
	s.sort = rules

	if len(s.sort) == 0 {
		s.sort = []string{"id"}
	}

	fmt.Printf("Before: %+v\n", s.col)
	sort.Sort(s)
	fmt.Printf("After: %+v\n", s.col)
}

// Swap ...
func (s *SoftCollection) Swap(i, j int) {
	fmt.Printf("We are swapping %d and %d\n", i, j)
	s.col[i], s.col[j] = s.col[j], s.col[i]
}

// Less ...
func (s *SoftCollection) Less(i, j int) bool {
	for _, r := range s.sort {
		inverse := false
		if strings.HasPrefix(r, "-") {
			r = r[1:]
			inverse = true
		}

		switch v := s.col[i].data[r].(type) {
		case string:
			if v == s.col[j].data[r].(string) {
				continue
			}
			if inverse {
				return v > s.col[j].data[r].(string)
			}
			return v < s.col[j].data[r].(string)
		case int:
			if v == s.col[j].data[r].(int) {
				continue
			}
			if inverse {
				return v > s.col[j].data[r].(int)
			}
			return v < s.col[j].data[r].(int)
		case int8:
			if v == s.col[j].data[r].(int8) {
				continue
			}
			if inverse {
				return v > s.col[j].data[r].(int8)
			}
			return v < s.col[j].data[r].(int8)
		case int16:
			if v == s.col[j].data[r].(int16) {
				continue
			}
			if inverse {
				return v > s.col[j].data[r].(int16)
			}
			return v < s.col[j].data[r].(int16)
		case int32:
			if v == s.col[j].data[r].(int32) {
				continue
			}
			if inverse {
				return v > s.col[j].data[r].(int32)
			}
			return v < s.col[j].data[r].(int32)
		case int64:
			if v == s.col[j].data[r].(int64) {
				continue
			}
			if inverse {
				return v > s.col[j].data[r].(int64)
			}
			return v < s.col[j].data[r].(int64)
		case uint:
			if v == s.col[j].data[r].(uint) {
				continue
			}
			if inverse {
				return v > s.col[j].data[r].(uint)
			}
			return v < s.col[j].data[r].(uint)
		case uint8:
			if v == s.col[j].data[r].(uint8) {
				continue
			}
			if inverse {
				return v > s.col[j].data[r].(uint8)
			}
			return v < s.col[j].data[r].(uint8)
		case uint16:
			if v == s.col[j].data[r].(uint16) {
				continue
			}
			if inverse {
				return v > s.col[j].data[r].(uint16)
			}
			return v < s.col[j].data[r].(uint16)
		case uint32:
			if v == s.col[j].data[r].(uint32) {
				continue
			}
			if inverse {
				return v > s.col[j].data[r].(uint32)
			}
			return v < s.col[j].data[r].(uint32)
		case bool:
			if v == s.col[j].data[r].(bool) {
				continue
			}
			if inverse {
				return v
			}
			return !v
		case time.Time:
			if v.Equal(s.col[j].data[r].(time.Time)) {
				continue
			}
			if inverse {
				return v.After(s.col[j].data[r].(time.Time))
			}
			return v.Before(s.col[j].data[r].(time.Time))
		case *string:
			if *v == *(s.col[j].data[r].(*string)) {
				continue
			}
			if inverse {
				return *v > *(s.col[j].data[r].(*string))
			}
			return *v < *(s.col[j].data[r].(*string))
		case *int:
			if *v == *(s.col[j].data[r].(*int)) {
				continue
			}
			if inverse {
				return *v > *(s.col[j].data[r].(*int))
			}
			return *v < *(s.col[j].data[r].(*int))
		case *int8:
			if *v == *(s.col[j].data[r].(*int8)) {
				continue
			}
			if inverse {
				return *v > *(s.col[j].data[r].(*int8))
			}
			return *v < *(s.col[j].data[r].(*int8))
		case *int16:
			if *v == *(s.col[j].data[r].(*int16)) {
				continue
			}
			if inverse {
				return *v > *(s.col[j].data[r].(*int16))
			}
			return *v < *(s.col[j].data[r].(*int16))
		case *int32:
			if *v == *(s.col[j].data[r].(*int32)) {
				continue
			}
			if inverse {
				return *v > *(s.col[j].data[r].(*int32))
			}
			return *v < *(s.col[j].data[r].(*int32))
		case *int64:
			if *v == *(s.col[j].data[r].(*int64)) {
				continue
			}
			if inverse {
				return *v > *(s.col[j].data[r].(*int64))
			}
			return *v < *(s.col[j].data[r].(*int64))
		case *uint:
			if *v == *(s.col[j].data[r].(*uint)) {
				continue
			}
			if inverse {
				return *v > *(s.col[j].data[r].(*uint))
			}
			return *v < *(s.col[j].data[r].(*uint))
		case *uint8:
			if *v == *(s.col[j].data[r].(*uint8)) {
				continue
			}
			if inverse {
				return *v > *(s.col[j].data[r].(*uint8))
			}
			return *v < *(s.col[j].data[r].(*uint8))
		case *uint16:
			if *v == *(s.col[j].data[r].(*uint16)) {
				continue
			}
			if inverse {
				return *v > *(s.col[j].data[r].(*uint16))
			}
			return *v < *(s.col[j].data[r].(*uint16))
		case *uint32:
			if *v == *(s.col[j].data[r].(*uint32)) {
				continue
			}
			if inverse {
				return *v > *(s.col[j].data[r].(*uint32))
			}
			return *v < *(s.col[j].data[r].(*uint32))
		case *bool:
			if *v == *(s.col[j].data[r].(*bool)) {
				continue
			}
			if inverse {
				return *v
			}
			return !*v
		case *time.Time:
			if v.Equal(*(s.col[j].data[r].(*time.Time))) {
				continue
			}
			if inverse {
				return v.After(*(s.col[j].data[r].(*time.Time)))
			}
			return v.Before(*(s.col[j].data[r].(*time.Time)))
		}
	}

	fmt.Printf("nothing happened, it was a lie, a big lie\n")
	return false
}
