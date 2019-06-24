package jsonapi

import (
	"encoding/json"
	"sort"
	"time"
)

// A Condition is used to define filters when querying collections.
type Condition struct {
	Field string      `json:"f"`
	Op    string      `json:"o"`
	Val   interface{} `json:"v"`
	Col   string      `json:"c"`
}

// cnd is an internal version of Condition.
type cnd struct {
	Field string          `json:"f"`
	Op    string          `json:"o"`
	Val   json.RawMessage `json:"v"`
	Col   string          `json:"c"`
}

// UnmarshalJSON parses the provided data and populates a Condition.
func (c *Condition) UnmarshalJSON(data []byte) error {
	tmpCnd := cnd{}
	err := json.Unmarshal(data, &tmpCnd)
	if err != nil {
		return err
	}

	c.Field = tmpCnd.Field
	c.Op = tmpCnd.Op
	c.Col = tmpCnd.Col

	if tmpCnd.Op == "and" || tmpCnd.Op == "or" {
		c.Field = ""

		cnds := []*Condition{}
		err := json.Unmarshal(tmpCnd.Val, &cnds)
		if err != nil {
			return err
		}
		c.Val = cnds
	} else {
		err := json.Unmarshal(tmpCnd.Val, &(c.Val)) // TODO parenthesis needed?
		if err != nil {
			return err
		}
	}

	return nil
}

// MarshalJSON marshals a Condition into JSON.
func (c *Condition) MarshalJSON() ([]byte, error) {
	payload := map[string]interface{}{}
	if c.Field != "" {
		payload["f"] = c.Field
	}
	if c.Op != "" {
		payload["o"] = c.Op
	}
	payload["v"] = c.Val
	if c.Col != "" {
		payload["c"] = c.Col
	}
	return json.Marshal(payload)
}

// FilterResource reports whether res is valid under the rules defined
// in cond.
func FilterResource(res Resource, cond *Condition) bool {
	var (
		val interface{}
		// typ string
	)
	if _, ok := res.Attrs()[cond.Field]; ok {
		val = res.Get(cond.Field)
	}
	if rel, ok := res.Rels()[cond.Field]; ok {
		if rel.ToOne {
			val = res.GetToOne(cond.Field)
		} else {
			val = res.GetToMany(cond.Field)
		}
	}

	switch cond.Op {
	case "and":
		conds := cond.Val.([]*Condition)
		for i := range conds {
			if !FilterResource(res, conds[i]) {
				return false
			}
		}
	case "or":
		conds := cond.Val.([]*Condition)
		for i := range conds {
			if FilterResource(res, conds[i]) {
				return true
			}
		}
	case "in":
		return checkIn(val.(string), cond.Val.([]string))
	case "=", "!=", "<", "<=", ">", ">=":
		return checkVal(cond.Op, val, cond.Val)
	}

	return false
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
