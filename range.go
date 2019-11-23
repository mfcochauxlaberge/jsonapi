package jsonapi

import (
	"sort"
	"strings"
	"time"
)

// Range returns a subset of the collection arranged according to the given
// parameters.
//
// From collection c, only IDs from ids are considered. filter is applied if not
// nil. The resources are sorted in the order defined by sort, which may contain
// the names of some or all of the attributes. The result is split in pages of a
// certain size (defined by size). The page at index num is returned.
//
// A non-nil Collection is always returned, but it can be empty.
func Range(c Collection, ids []string, filter *Filter, sort []string, size uint, num uint) Collection {
	col := sortedResources{}

	// Filter IDs
	if len(ids) > 0 {
		for i := 0; i < c.Len(); i++ {
			for _, id := range ids {
				res := c.At(i)
				if res.GetID() == id {
					col.col = append(col.col, res)
				}
			}
		}
	} else {
		for i := 0; i < c.Len(); i++ {
			col.col = append(col.col, c.At(i))
		}
	}

	// Filter
	if filter != nil {
		i := 0
		for i < col.Len() {
			if !filter.IsAllowed(col.col[i]) {
				col.col = append(col.col[:i], col.col[i+1:]...)
			} else {
				i++
			}
		}
	}

	// Sort
	col.Sort(sort)

	// Pagination
	var page Resources

	skip := int(num * size)

	if skip >= len(col.col) {
		col = sortedResources{}
	} else {
		for i := skip; i < len(col.col) && i < skip+int(size); i++ {
			page = append(page, col.col[i])
		}
	}

	return &page
}

// sortedResources is an internal struct for sorting Collections with the Range
// function.
type sortedResources struct {
	rules []string
	col   Resources
}

// Sort rearranges the order of the collection according the rules.
func (s sortedResources) Sort(rules []string) {
	s.rules = rules
	if len(s.rules) == 0 {
		s.rules = []string{"id"}
	}

	sort.Sort(s)
}

// Len implements sort.Interface's Len method.
func (s sortedResources) Len() int {
	return len(s.col)
}

// Swap implements sort.Interface's Swap method.
func (s sortedResources) Swap(i, j int) {
	s.col[i], s.col[j] = s.col[j], s.col[i]
}

// Less implements sort.Interface's Less method.
func (s sortedResources) Less(i, j int) bool {
	for _, r := range s.rules {
		inverse := false

		if strings.HasPrefix(r, "-") {
			r = r[1:]
			inverse = true
		}

		if r == "id" {
			return s.col[i].GetID() < s.col[j].GetID() != inverse
		}

		v := s.col[i].Get(r)
		v2 := s.col[j].Get(r)

		// Here we return true if v < v2.
		// The "!= inverse" part acts as a XOR operation so that
		// the opposite boolean is returned when inverse sorting
		// is required.
		switch v := v.(type) {
		case string:
			v2 := v2.(string)
			if v == v2 {
				continue
			}

			return v < v2 != inverse
		case int:
			v2 := v2.(int)
			if v == v2 {
				continue
			}

			return v < v2 != inverse
		case int8:
			v2 := v2.(int8)
			if v == v2 {
				continue
			}

			return v < v2 != inverse
		case int16:
			v2 := v2.(int16)
			if v == v2 {
				continue
			}

			return v < v2 != inverse
		case int32:
			v2 := v2.(int32)
			if v == v2 {
				continue
			}

			return v < v2 != inverse
		case int64:
			v2 := v2.(int64)
			if v == v2 {
				continue
			}

			return v < v2 != inverse
		case uint:
			v2 := v2.(uint)
			if v == v2 {
				continue
			}

			return v < v2 != inverse
		case uint8:
			v2 := v2.(uint8)
			if v == v2 {
				continue
			}

			return v < v2 != inverse
		case uint16:
			v2 := v2.(uint16)
			if v == v2 {
				continue
			}

			return v < v2 != inverse
		case uint32:
			v2 := v2.(uint32)
			if v == v2 {
				continue
			}

			return v < v2 != inverse
		case bool:
			v2 := v2.(bool)
			if v == v2 {
				continue
			}

			return !v != inverse
		case time.Time:
			if v.Equal(v2.(time.Time)) {
				continue
			}

			return v.Before(v2.(time.Time)) != inverse
		case []byte:
			s2 := v2.([]byte)
			for i := 0; i < len(v) && i < len(s2); i++ {
				if v[i] == s2[i] {
					continue
				}

				return v[i] < s2[i] != inverse
			}

			if len(v) == len(s2) {
				continue
			}

			return len(v) < len(s2) != inverse
		case *string:
			v2 := v2.(*string)
			if v == v2 {
				continue
			}

			if v == nil {
				return !inverse
			}

			if v2 == nil {
				return inverse
			}

			if *v == *v2 {
				continue
			}

			return *v < *v2 != inverse
		case *int:
			v2 := v2.(*int)
			if v == v2 {
				continue
			}

			if v == nil {
				return !inverse
			}

			if v2 == nil {
				return inverse
			}

			if *v == *v2 {
				continue
			}

			return *v < *v2 != inverse
		case *int8:
			v2 := v2.(*int8)
			if v == v2 {
				continue
			}

			if v == nil {
				return !inverse
			}

			if v2 == nil {
				return inverse
			}

			if *v == *v2 {
				continue
			}

			return *v < *v2 != inverse
		case *int16:
			v2 := v2.(*int16)
			if v == v2 {
				continue
			}

			if v == nil {
				return !inverse
			}

			if v2 == nil {
				return inverse
			}

			if *v == *v2 {
				continue
			}

			return *v < *v2 != inverse
		case *int32:
			v2 := v2.(*int32)
			if v == v2 {
				continue
			}

			if v == nil {
				return !inverse
			}

			if v2 == nil {
				return inverse
			}

			if *v == *v2 {
				continue
			}

			return *v < *v2 != inverse
		case *int64:
			v2 := v2.(*int64)
			if v == v2 {
				continue
			}

			if v == nil {
				return !inverse
			}

			if v2 == nil {
				return inverse
			}

			if *v == *v2 {
				continue
			}

			return *v < *v2 != inverse
		case *uint:
			v2 := v2.(*uint)
			if v == v2 {
				continue
			}

			if v == nil {
				return !inverse
			}

			if v2 == nil {
				return inverse
			}

			if *v == *v2 {
				continue
			}

			return *v < *v2 != inverse
		case *uint8:
			v2 := v2.(*uint8)
			if v == v2 {
				continue
			}

			if v == nil {
				return !inverse
			}

			if v2 == nil {
				return inverse
			}

			if *v == *v2 {
				continue
			}

			return *v < *v2 != inverse
		case *uint16:
			v2 := v2.(*uint16)
			if v == v2 {
				continue
			}

			if v == nil {
				return !inverse
			}

			if v2 == nil {
				return inverse
			}

			if *v == *v2 {
				continue
			}

			return *v < *v2 != inverse
		case *uint32:
			v2 := v2.(*uint32)
			if v == v2 {
				continue
			}

			if v == nil {
				return !inverse
			}

			if v2 == nil {
				return inverse
			}

			if *v == *v2 {
				continue
			}

			return *v < *v2 != inverse
		case *bool:
			v2 := v2.(*bool)
			if v == v2 {
				continue
			}

			if v == nil {
				return !inverse
			}

			if v2 == nil {
				return inverse
			}

			if *v == *v2 {
				continue
			}

			return !*v != inverse
		case *time.Time:
			v2 := v2.(*time.Time)
			if v == v2 {
				continue
			}

			if v == nil {
				return !inverse
			}

			if v2 == nil {
				return inverse
			}

			if v.Equal(*v2) {
				continue
			}

			return v.Before(*v2) != inverse
		}
	}

	return false
}
