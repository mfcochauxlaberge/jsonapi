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

	sort.Sort(s)
}

// Swap ...
func (s *SoftCollection) Swap(i, j int) {
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
			fmt.Printf("why?\n")
			if v.Equal(s.col[j].data[r].(time.Time)) {
				fmt.Printf("time is equal!\n")
				continue
			}
			if inverse {
				fmt.Printf("inverse: %v\n", v.After(s.col[j].data[r].(time.Time)))
				return v.After(s.col[j].data[r].(time.Time))
			}
			fmt.Printf("inverse: %v\n", v.Before(s.col[j].data[r].(time.Time)))
			return v.Before(s.col[j].data[r].(time.Time))
		case *string:
			p := s.col[j].data[r].(*string)
			if v == p {
				return false
			}
			if v == nil {
				return false
			}
			if p == nil {
				return true
			}
			if *v == *p {
				return false
			}
			if inverse {
				return *v > *(s.col[j].data[r].(*string))
			}
			return *v < *(s.col[j].data[r].(*string))
		case *int:
			p := s.col[j].data[r].(*int)
			if v == p {
				return false
			}
			if v == nil {
				return false
			}
			if p == nil {
				return true
			}
			if *v == *p {
				return false
			}
			if inverse {
				return *v > *(s.col[j].data[r].(*int))
			}
			return *v < *(s.col[j].data[r].(*int))
		case *int8:
			p := s.col[j].data[r].(*int8)
			if v == p {
				return false
			}
			if v == nil {
				return false
			}
			if p == nil {
				return true
			}
			if *v == *p {
				return false
			}
			if inverse {
				return *v > *(s.col[j].data[r].(*int8))
			}
			return *v < *(s.col[j].data[r].(*int8))
		case *int16:
			p := s.col[j].data[r].(*int16)
			if v == p {
				return false
			}
			if v == nil {
				return false
			}
			if p == nil {
				return true
			}
			if *v == *p {
				return false
			}
			if inverse {
				return *v > *(s.col[j].data[r].(*int16))
			}
			return *v < *(s.col[j].data[r].(*int16))
		case *int32:
			p := s.col[j].data[r].(*int32)
			if v == p {
				return false
			}
			if v == nil {
				return false
			}
			if p == nil {
				return true
			}
			if *v == *p {
				return false
			}
			if inverse {
				return *v > *(s.col[j].data[r].(*int32))
			}
			return *v < *(s.col[j].data[r].(*int32))
		case *int64:
			p := s.col[j].data[r].(*int64)
			if v == p {
				return false
			}
			if v == nil {
				return false
			}
			if p == nil {
				return true
			}
			if *v == *p {
				return false
			}
			if inverse {
				return *v > *(s.col[j].data[r].(*int64))
			}
			return *v < *(s.col[j].data[r].(*int64))
		case *uint:
			p := s.col[j].data[r].(*uint)
			if v == p {
				return false
			}
			if v == nil {
				return false
			}
			if p == nil {
				return true
			}
			if *v == *p {
				return false
			}
			if inverse {
				return *v > *(s.col[j].data[r].(*uint))
			}
			return *v < *(s.col[j].data[r].(*uint))
		case *uint8:
			p := s.col[j].data[r].(*uint8)
			if v == p {
				return false
			}
			if v == nil {
				return false
			}
			if p == nil {
				return true
			}
			if *v == *p {
				return false
			}
			if inverse {
				return *v > *(s.col[j].data[r].(*uint8))
			}
			return *v < *(s.col[j].data[r].(*uint8))
		case *uint16:
			p := s.col[j].data[r].(*uint16)
			if v == p {
				return false
			}
			if v == nil {
				return false
			}
			if p == nil {
				return true
			}
			if *v == *p {
				return false
			}
			if inverse {
				return *v > *(s.col[j].data[r].(*uint16))
			}
			return *v < *(s.col[j].data[r].(*uint16))
		case *uint32:
			p := s.col[j].data[r].(*uint32)
			if v == p {
				return false
			}
			if v == nil {
				return false
			}
			if p == nil {
				return true
			}
			if *v == *p {
				return false
			}
			if inverse {
				return *v > *(s.col[j].data[r].(*uint32))
			}
			return *v < *(s.col[j].data[r].(*uint32))
		case *bool:
			p := s.col[j].data[r].(*bool)
			if v == p {
				return false
			}
			if v == nil {
				return false
			}
			if p == nil {
				return true
			}
			if *v == *p {
				return false
			}
			if inverse {
				return *v
			}
			return !*v
		case *time.Time:
			p := s.col[j].data[r].(*time.Time)
			if v == p {
				return false
			}
			if v == nil {
				return false
			}
			if p == nil {
				return true
			}
			if v.Equal(*p) {
				continue
			}
			if inverse {
				return v.After(*p)
			}
			return v.Before(*p)
		}
	}

	fmt.Printf("hey!\n")
	return false
}
