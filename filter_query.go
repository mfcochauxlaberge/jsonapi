package jsonapi

import (
	"encoding/json"
)

// Condition ...
type Condition struct {
	Field string      `json:"f"`
	Op    string      `json:"o"`
	Val   interface{} `json:"v"`
	Col   string      `json:"c"`
}

// cnd ...
type cnd struct {
	Field string          `json:"f"`
	Op    string          `json:"o"`
	Val   json.RawMessage `json:"v"`
	Col   string          `json:"c"`
}

// UnmarshalJSON ...
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

// MarshalJSON ...
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
