package jsonapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"
)

// Condition ...
type Condition struct {
	Kind string // string, number, time
	Op   string
	Val  interface{}
}

// UnmarshalJSON ...
func (c *Condition) UnmarshalJSON(data []byte) error {
	if c.Kind == "" {
		return errors.New("jsonapi: can't unmarshal Condition with empty Kind")
	}

	s := map[string]json.RawMessage{}

	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	if len(s) == 0 || len(s) > 1 {
		return errors.New("jsonapi: a json object may contain one and only one property")
	}

	for op, val := range s {
		c.Op = op

		if string(val) == "null" {
			c.Val = nil
		} else {
			if op == "and" || op == "or" {
				raws := []json.RawMessage{}

				err := json.Unmarshal(val, &raws)
				if err != nil {
					return err
				}

				cgs := make([]Condition, len(raws))
				for i := range raws {
					cgs[i].Kind = c.Kind
					err := json.Unmarshal(raws[i], &(cgs[i]))
					if err != nil {
						return err
					}
				}

				c.Val = cgs
			} else if c.Kind == "string" {
				if op == "=" || op == "!=" || op == "<" || op == ">" || op == "<=" || op == ">=" || op == "~" || op == "!~" {
					var str string

					err := json.Unmarshal(val, &str)
					if err != nil {
						return err
					}

					c.Val = str
				} else if op == "in" || op == "notin" {
					var sl []string

					err := json.Unmarshal(val, &sl)
					if err != nil {
						return err
					}

					sort.Strings(sl)
					c.Val = sl
				} else {
					return fmt.Errorf("jsonapi: operation '%s' unknown", op)
				}
			} else if c.Kind == "number" {
				if op == "=" || op == "!=" || op == "<" || op == ">" || op == "<=" || op == ">=" {
					var num int

					err := json.Unmarshal(val, &num)
					if err != nil {
						return err
					}

					c.Val = num
				} else if op == "in" || op == "notin" {
					var sl []int

					err := json.Unmarshal(val, &sl)
					if err != nil {
						sort.Ints(sl)
						return err
					}

					c.Val = sl
				} else {
					return fmt.Errorf("jsonapi: operation '%s' unknown", op)
				}
			} else if c.Kind == "time" {
				// Time types
				// year, month, day, hour, minute, second, ms
				if op == "=" || op == "!=" || op == "<" || op == ">" || op == "<=" || op == ">=" {
					var t time.Time

					err := json.Unmarshal(val, &t)
					if err != nil {
						return err
					}

					c.Val = t
				} else if op == "in" || op == "notin" {
					var sl []time.Time

					err := json.Unmarshal(val, &sl)
					if err != nil {
						return err
					}

					c.Val = sl
				} else {
					return fmt.Errorf("jsonapi: operation '%s' unknown", op)
				}
			} else {
				return fmt.Errorf("jsonapi: kind '%s' unknown", c.Kind)
			}
		}
	}

	return nil
}

// MarshalJSON ...
func (c *Condition) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`{"`)

	if len(c.Op) == 0 {
		return nil, errors.New("jsonapi: can't marshal Condition with empty Op")
	}

	buf.WriteString(c.Op)
	buf.WriteString(`":`)

	if c.Val == nil {
		buf.WriteString("null")
	} else {
		val, err := json.Marshal(c.Val)
		if err != nil {
			return nil, err
		}
		buf.Write(val)
	}

	buf.WriteString(`}`)

	return buf.Bytes(), nil
}
