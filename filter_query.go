package jsonapi

import (
	"encoding/json"
)

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
	case "=", "!=", "<", "<=", ">", ">=":
		if val == cond.Val {
			return true
		}
	}

	return false
}

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

		cnds := []Condition{}
		err := json.Unmarshal(tmpCnd.Val, &cnds)
		if err != nil {
			return err
		}
		c.Val = cnds
	} else if tmpCnd.Op == "=" ||
		tmpCnd.Op == "!=" ||
		tmpCnd.Op == "<" ||
		tmpCnd.Op == "<=" ||
		tmpCnd.Op == ">" ||
		tmpCnd.Op == ">=" {

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
