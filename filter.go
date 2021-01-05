package jsonapi

import (
	"encoding/json"
	"sort"
	"time"
)

// A Filter is used to define filters when querying collections.
type Filter struct {
	Field string      `json:"f"`
	Op    string      `json:"o"`
	Val   interface{} `json:"v"`
	Col   string      `json:"c"`
}

// filter is an internal version of Filter.
type filter struct {
	Field string          `json:"f,omitempty"`
	Op    string          `json:"o,omitempty"`
	Val   json.RawMessage `json:"v"`
	Col   string          `json:"c,omitempty"`
}

// UnmarshalJSON parses the provided data and populates a Filter.
func (f *Filter) UnmarshalJSON(data []byte) error {
	tmpFilter := filter{}

	err := json.Unmarshal(data, &tmpFilter)
	if err != nil {
		return err
	}

	f.Field = tmpFilter.Field
	f.Op = tmpFilter.Op
	f.Col = tmpFilter.Col

	switch tmpFilter.Op {
	case "and", "or":
		f.Field = ""

		filters := []*Filter{}

		err := json.Unmarshal(tmpFilter.Val, &filters)
		if err != nil {
			return err
		}

		f.Val = filters
	default:
		// Error checking ignored since it cannot fail at this
		// point. The first unmarshaling step of this function
		// has already checked the data.
		_ = json.Unmarshal(tmpFilter.Val, &f.Val)
	}

	return nil
}

// IsAllowed reports whether res is valid under the rules defined in the filter.
func (f *Filter) IsAllowed(res Resource) bool {
	var (
		val interface{}
		// typ string
	)

	if _, ok := res.Attrs()[f.Field]; ok {
		val = res.Get(f.Field)
	}

	if rel, ok := res.Rels()[f.Field]; ok {
		if rel.ToOne {
			val = res.GetToOne(f.Field)
		} else {
			val = res.GetToMany(f.Field)
		}
	}

	switch f.Op {
	case "and":
		filters := f.Val.([]*Filter)
		for i := range filters {
			if !filters[i].IsAllowed(res) {
				return false
			}
		}

		return true
	case "or":
		filters := f.Val.([]*Filter)
		for i := range filters {
			if filters[i].IsAllowed(res) {
				return true
			}
		}

		return false
	case "in":
		return checkIn(val.(string), f.Val.([]string))
	case "has":
		return checkIn(f.Val.(string), val.([]string))
	default:
		return checkVal(f.Op, val, f.Val)
	}
}

func checkVal(op string, rval, cval interface{}) bool {
	switch rval := rval.(type) {
	case string:
		return checkStr(op, rval, cval.(string))
	case int:
		return checkInt(op, int64(rval), int64(cval.(int)))
	case int8:
		return checkInt(op, int64(rval), int64(cval.(int8)))
	case int16:
		return checkInt(op, int64(rval), int64(cval.(int16)))
	case int32:
		return checkInt(op, int64(rval), int64(cval.(int32)))
	case int64:
		return checkInt(op, rval, cval.(int64))
	case uint:
		return checkUint(op, uint64(rval), uint64(cval.(uint)))
	case uint8:
		return checkUint(op, uint64(rval), uint64(cval.(uint8)))
	case uint16:
		return checkUint(op, uint64(rval), uint64(cval.(uint16)))
	case uint32:
		return checkUint(op, uint64(rval), uint64(cval.(uint32)))
	case uint64:
		return checkUint(op, rval, cval.(uint64))
	case bool:
		return checkBool(op, rval, cval.(bool))
	case time.Time:
		return checkTime(op, rval, cval.(time.Time))
	case []byte:
		return checkBytes(op, rval, cval.([]byte))
	case *string:
		if rval == nil || cval.(*string) == nil {
			switch op {
			case "=":
				return rval == cval.(*string)
			case "!=":
				return rval != cval.(*string)
			default:
				return false
			}
		}

		return checkStr(op, *rval, *cval.(*string))
	case *int:
		if rval == nil || cval.(*int) == nil {
			switch op {
			case "=":
				return rval == cval.(*int)
			case "!=":
				return rval != cval.(*int)
			default:
				return false
			}
		}

		return checkInt(op, int64(*rval), int64(*cval.(*int)))
	case *int8:
		if rval == nil || cval.(*int8) == nil {
			switch op {
			case "=":
				return rval == cval.(*int8)
			case "!=":
				return rval != cval.(*int8)
			default:
				return false
			}
		}

		return checkInt(op, int64(*rval), int64(*cval.(*int8)))
	case *int16:
		if rval == nil || cval.(*int16) == nil {
			switch op {
			case "=":
				return rval == cval.(*int16)
			case "!=":
				return rval != cval.(*int16)
			default:
				return false
			}
		}

		return checkInt(op, int64(*rval), int64(*cval.(*int16)))
	case *int32:
		if rval == nil || cval.(*int32) == nil {
			switch op {
			case "=":
				return rval == cval.(*int32)
			case "!=":
				return rval != cval.(*int32)
			default:
				return false
			}
		}

		return checkInt(op, int64(*rval), int64(*cval.(*int32)))
	case *int64:
		if rval == nil || cval.(*int64) == nil {
			switch op {
			case "=":
				return rval == cval.(*int64)
			case "!=":
				return rval != cval.(*int64)
			default:
				return false
			}
		}

		return checkInt(op, *rval, *cval.(*int64))
	case *uint:
		if rval == nil || cval.(*uint) == nil {
			switch op {
			case "=":
				return rval == cval.(*uint)
			case "!=":
				return rval != cval.(*uint)
			default:
				return false
			}
		}

		return checkUint(op, uint64(*rval), uint64(*cval.(*uint)))
	case *uint8:
		if rval == nil || cval.(*uint8) == nil {
			switch op {
			case "=":
				return rval == cval.(*uint8)
			case "!=":
				return rval != cval.(*uint8)
			default:
				return false
			}
		}

		return checkUint(op, uint64(*rval), uint64(*cval.(*uint8)))
	case *uint16:
		if rval == nil || cval.(*uint16) == nil {
			switch op {
			case "=":
				return rval == cval.(*uint16)
			case "!=":
				return rval != cval.(*uint16)
			default:
				return false
			}
		}

		return checkUint(op, uint64(*rval), uint64(*cval.(*uint16)))
	case *uint32:
		if rval == nil || cval.(*uint32) == nil {
			switch op {
			case "=":
				return rval == cval.(*uint32)
			case "!=":
				return rval != cval.(*uint32)
			default:
				return false
			}
		}

		return checkUint(op, uint64(*rval), uint64(*cval.(*uint32)))
	case *uint64:
		if rval == nil || cval.(*uint64) == nil {
			switch op {
			case "=":
				return rval == cval.(*uint64)
			case "!=":
				return rval != cval.(*uint64)
			default:
				return false
			}
		}

		return checkUint(op, *rval, *cval.(*uint64))
	case *bool:
		if rval == nil || cval.(*bool) == nil {
			switch op {
			case "=":
				return rval == cval.(*bool)
			case "!=":
				return rval != cval.(*bool)
			default:
				return false
			}
		}

		return checkBool(op, *rval, *cval.(*bool))
	case *time.Time:
		if rval == nil || cval.(*time.Time) == nil {
			switch op {
			case "=":
				return rval == cval.(*time.Time)
			case "!=":
				return rval != cval.(*time.Time)
			default:
				return false
			}
		}

		return checkTime(op, *rval, *cval.(*time.Time))
	case *[]byte:
		if rval == nil || cval.(*[]byte) == nil {
			switch op {
			case "=":
				return rval == cval.(*[]byte)
			case "!=":
				return rval != cval.(*[]byte)
			default:
				return false
			}
		}

		return checkBytes(op, *rval, *cval.(*[]byte))
	case []string:
		return checkSlice(op, rval, cval.([]string))
	default:
		return false
	}
}

func checkStr(op string, rval, cval string) bool {
	switch op {
	case "=":
		return rval == cval
	case "!=":
		return rval != cval
	case "<":
		return rval < cval
	case "<=":
		return rval <= cval
	case ">":
		return rval > cval
	case ">=":
		return rval >= cval
	default:
		return false
	}
}

func checkInt(op string, rval, cval int64) bool {
	switch op {
	case "=":
		return rval == cval
	case "!=":
		return rval != cval
	case "<":
		return rval < cval
	case "<=":
		return rval <= cval
	case ">":
		return rval > cval
	case ">=":
		return rval >= cval
	default:
		return false
	}
}

func checkUint(op string, rval, cval uint64) bool {
	switch op {
	case "=":
		return rval == cval
	case "!=":
		return rval != cval
	case "<":
		return rval < cval
	case "<=":
		return rval <= cval
	case ">":
		return rval > cval
	case ">=":
		return rval >= cval
	default:
		return false
	}
}

func checkBool(op string, rval, cval bool) bool {
	switch op {
	case "=":
		return rval == cval
	case "!=":
		return rval != cval
	default:
		return false
	}
}

func checkTime(op string, rval, cval time.Time) bool {
	switch op {
	case "=":
		return rval.Equal(cval)
	case "!=":
		return !rval.Equal(cval)
	case "<":
		return rval.Before(cval)
	case "<=":
		return rval.Before(cval) || rval.Equal(cval)
	case ">":
		return rval.After(cval)
	case ">=":
		return rval.After(cval) || rval.Equal(cval)
	default:
		return false
	}
}

func checkBytes(op string, rval, cval []byte) bool {
	switch op {
	case "=":
		for i := 0; i < len(rval) && i < len(cval); i++ {
			if rval[i] != cval[i] {
				return false
			}
		}

		return len(rval) == len(cval)
	case "!=":
		for i := 0; i < len(rval) && i < len(cval); i++ {
			if rval[i] != cval[i] {
				return true
			}
		}

		return len(rval) != len(cval)
	case "<":
		for i := 0; i < len(rval) && i < len(cval); i++ {
			if rval[i] < cval[i] {
				return true
			}
		}

		return len(rval) < len(cval)
	case "<=":
		for i := 0; i < len(rval) && i < len(cval); i++ {
			if rval[i] > cval[i] {
				return false
			}
		}

		return len(rval) <= len(cval)
	case ">":
		for i := 0; i < len(rval) && i < len(cval); i++ {
			if rval[i] > cval[i] {
				return true
			}
		}

		return len(rval) > len(cval)
	case ">=":
		for i := 0; i < len(rval) && i < len(cval); i++ {
			if rval[i] < cval[i] {
				return false
			}
		}

		return len(rval) >= len(cval)
	default:
		return false
	}
}

func checkSlice(op string, rval, cval []string) bool {
	equal := false

	if len(rval) == len(cval) {
		sort.Strings(rval)
		sort.Strings(cval)

		equal = true

		for i := 0; i < len(rval); i++ {
			if rval[i] != cval[i] {
				equal = false
				break
			}
		}
	}

	switch op {
	case "=":
		return equal
	case "!=":
		return !equal
	default:
		return false
	}
}

func checkIn(id string, ids []string) bool {
	for i := range ids {
		if id == ids[i] {
			return true
		}
	}

	return false
}
