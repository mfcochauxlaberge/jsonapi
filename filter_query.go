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
	Field string          `json:"f"`
	Op    string          `json:"o"`
	Val   json.RawMessage `json:"v"`
	Col   string          `json:"c"`
}

// MarshalJSON marshals a filter into JSON.
func (f *Filter) MarshalJSON() ([]byte, error) {
	payload := map[string]interface{}{}
	if f.Field != "" {
		payload["f"] = f.Field
	}
	if f.Op != "" {
		payload["o"] = f.Op
	}
	payload["v"] = f.Val
	if f.Col != "" {
		payload["c"] = f.Col
	}
	return json.Marshal(payload)
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

	if tmpFilter.Op == "and" || tmpFilter.Op == "or" {
		f.Field = ""

		filters := []*Filter{}
		err := json.Unmarshal(tmpFilter.Val, &filters)
		if err != nil {
			return err
		}
		f.Val = filters
	} else {
		err := json.Unmarshal(tmpFilter.Val, &f.Val)
		if err != nil {
			return err
		}
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
	switch rval.(type) {
	case string:
		return checkStr(op, rval.(string), cval.(string))
	case int:
		return checkInt(op, int64(rval.(int)), int64(cval.(int)))
	case int8:
		return checkInt(op, int64(rval.(int8)), int64(cval.(int8)))
	case int16:
		return checkInt(op, int64(rval.(int16)), int64(cval.(int16)))
	case int32:
		return checkInt(op, int64(rval.(int32)), int64(cval.(int32)))
	case int64:
		return checkInt(op, rval.(int64), cval.(int64))
	case uint:
		return checkUint(op, uint64(rval.(uint)), uint64(cval.(uint)))
	case uint8:
		return checkUint(op, uint64(rval.(uint8)), uint64(cval.(uint8)))
	case uint16:
		return checkUint(op, uint64(rval.(uint16)), uint64(cval.(uint16)))
	case uint32:
		return checkUint(op, uint64(rval.(uint32)), uint64(cval.(uint32)))
	case uint64:
		return checkUint(op, rval.(uint64), cval.(uint64))
	case bool:
		return checkBool(op, rval.(bool), cval.(bool))
	case time.Time:
		return checkTime(op, rval.(time.Time), cval.(time.Time))
	case []byte:
		return checkBytes(op, rval.([]byte), cval.([]byte))
	case *string:
		if rval.(*string) == nil || cval.(*string) == nil {
			if op == "=" {
				return rval.(*string) == cval.(*string)
			} else if op == "!=" {
				return rval.(*string) != cval.(*string)
			} else {
				return false
			}
		}
		return checkStr(op, *rval.(*string), *cval.(*string))
	case *int:
		if rval.(*int) == nil || cval.(*int) == nil {
			if op == "=" {
				return rval.(*int) == cval.(*int)
			} else if op == "!=" {
				return rval.(*int) != cval.(*int)
			} else {
				return false
			}
		}
		return checkInt(op, int64(*rval.(*int)), int64(*cval.(*int)))
	case *int8:
		if rval.(*int8) == nil || cval.(*int8) == nil {
			if op == "=" {
				return rval.(*int8) == cval.(*int8)
			} else if op == "!=" {
				return rval.(*int8) != cval.(*int8)
			} else {
				return false
			}
		}
		return checkInt(op, int64(*rval.(*int8)), int64(*cval.(*int8)))
	case *int16:
		if rval.(*int16) == nil || cval.(*int16) == nil {
			if op == "=" {
				return rval.(*int16) == cval.(*int16)
			} else if op == "!=" {
				return rval.(*int16) != cval.(*int16)
			} else {
				return false
			}
		}
		return checkInt(op, int64(*rval.(*int16)), int64(*cval.(*int16)))
	case *int32:
		if rval.(*int32) == nil || cval.(*int32) == nil {
			if op == "=" {
				return rval.(*int32) == cval.(*int32)
			} else if op == "!=" {
				return rval.(*int32) != cval.(*int32)
			} else {
				return false
			}
		}
		return checkInt(op, int64(*rval.(*int32)), int64(*cval.(*int32)))
	case *int64:
		if rval.(*int64) == nil || cval.(*int64) == nil {
			if op == "=" {
				return rval.(*int64) == cval.(*int64)
			} else if op == "!=" {
				return rval.(*int64) != cval.(*int64)
			} else {
				return false
			}
		}
		return checkInt(op, *rval.(*int64), *cval.(*int64))
	case *uint:
		if rval.(*uint) == nil || cval.(*uint) == nil {
			if op == "=" {
				return rval.(*uint) == cval.(*uint)
			} else if op == "!=" {
				return rval.(*uint) != cval.(*uint)
			} else {
				return false
			}
		}
		return checkUint(op, uint64(*rval.(*uint)), uint64(*cval.(*uint)))
	case *uint8:
		if rval.(*uint8) == nil || cval.(*uint8) == nil {
			if op == "=" {
				return rval.(*uint8) == cval.(*uint8)
			} else if op == "!=" {
				return rval.(*uint8) != cval.(*uint8)
			} else {
				return false
			}
		}
		return checkUint(op, uint64(*rval.(*uint8)), uint64(*cval.(*uint8)))
	case *uint16:
		if rval.(*uint16) == nil || cval.(*uint16) == nil {
			if op == "=" {
				return rval.(*uint16) == cval.(*uint16)
			} else if op == "!=" {
				return rval.(*uint16) != cval.(*uint16)
			} else {
				return false
			}
		}
		return checkUint(op, uint64(*rval.(*uint16)), uint64(*cval.(*uint16)))
	case *uint32:
		if rval.(*uint32) == nil || cval.(*uint32) == nil {
			if op == "=" {
				return rval.(*uint32) == cval.(*uint32)
			} else if op == "!=" {
				return rval.(*uint32) != cval.(*uint32)
			} else {
				return false
			}
		}
		return checkUint(op, uint64(*rval.(*uint32)), uint64(*cval.(*uint32)))
	case *uint64:
		if rval.(*uint64) == nil || cval.(*uint64) == nil {
			if op == "=" {
				return rval.(*uint64) == cval.(*uint64)
			} else if op == "!=" {
				return rval.(*uint64) != cval.(*uint64)
			} else {
				return false
			}
		}
		return checkUint(op, *rval.(*uint64), *cval.(*uint64))
	case *bool:
		if rval.(*bool) == nil || cval.(*bool) == nil {
			if op == "=" {
				return rval.(*bool) == cval.(*bool)
			} else if op == "!=" {
				return rval.(*bool) != cval.(*bool)
			}
		}
		return checkBool(op, *rval.(*bool), *cval.(*bool))
	case *time.Time:
		if rval.(*time.Time) == nil || cval.(*time.Time) == nil {
			if op == "=" {
				return rval.(*time.Time) == cval.(*time.Time)
			} else if op == "!=" {
				return rval.(*time.Time) != cval.(*time.Time)
			} else {
				return false
			}
		}
		return checkTime(op, *rval.(*time.Time), *cval.(*time.Time))
	case *[]byte:
		if rval.(*[]byte) == nil || cval.(*[]byte) == nil {
			if op == "=" {
				return rval.(*[]byte) == cval.(*[]byte)
			} else if op == "!=" {
				return rval.(*[]byte) != cval.(*[]byte)
			} else {
				return false
			}
		}
		return checkBytes(op, *rval.(*[]byte), *cval.(*[]byte))
	case []string:
		return checkSlice(op, rval.([]string), cval.([]string))
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
