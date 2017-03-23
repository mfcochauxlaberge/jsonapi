package jsonapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Wrapper ...
type Wrapper struct {
	val reflect.Value // Actual value (with content)

	// Structure
	typ   string
	attrs []Attr
	rels  []Rel
}

// Wrap ...
func Wrap(v interface{}) *Wrapper {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Ptr {
		panic(errors.New("jsonapi: value has to be a pointer to a struct"))
	}

	if val.Elem().Kind() != reflect.Struct {
		panic(errors.New("jsonapi: value has to be a pointer to a struct"))
	}

	val = val.Elem()

	err := CheckType(val.Interface())
	if err != nil {
		panic(fmt.Sprintf("jsonapi: invalid type: %s", err))
	}

	w := &Wrapper{
		val: val,
	}

	// ID and type
	_, w.typ = IDAndType(v)

	// Attributes
	w.attrs = []Attr{}
	for i := 0; i < w.val.NumField(); i++ {
		fs := w.val.Type().Field(i)
		jsonTag := fs.Tag.Get("json")
		apiTag := fs.Tag.Get("api")

		if apiTag == "attr" {
			w.attrs = append(w.attrs, Attr{
				Name: jsonTag,
				Type: fs.Type.String(),
			})
		}
	}

	// Relationships
	w.rels = []Rel{}
	for i := 0; i < w.val.NumField(); i++ {
		fs := w.val.Type().Field(i)
		jsonTag := fs.Tag.Get("json")
		relTag := strings.Split(fs.Tag.Get("api"), ",")
		invName := ""
		if len(relTag) == 3 {
			invName = relTag[2]
		}

		toOne := true
		if fs.Type.String() == "[]string" {
			toOne = false
		}

		if relTag[0] == "rel" {
			w.rels = append(w.rels, Rel{
				Name:        jsonTag,
				Type:        relTag[1],
				ToOne:       toOne,
				InverseName: invName,
				InverseType: w.typ,
			})
		}
	}

	return w
}

// IDAndType ...
func (w *Wrapper) IDAndType() (string, string) {
	return IDAndType(w.val.Interface())
}

// Attrs ...
func (w *Wrapper) Attrs() []Attr {
	return w.attrs
}

// Rels ...
func (w *Wrapper) Rels() []Rel {
	return w.rels
}

// Attr ...
func (w *Wrapper) Attr(key string) Attr {
	for _, attr := range w.attrs {
		if attr.Name == key {
			return attr
		}
	}

	panic(fmt.Sprintf("jsonapi: attribute %s does not exist", key))
}

// Rel ...
func (w *Wrapper) Rel(key string) Rel {
	for _, rel := range w.rels {
		if rel.Name == key {
			return rel
		}
	}

	panic(fmt.Sprintf("jsonapi: relationship %s does not exist", key))
}

// New ...
func (w *Wrapper) New() Resource {
	newVal := reflect.New(w.val.Type())

	return Wrap(newVal.Interface())
}

// Get ...
func (w *Wrapper) Get(key string) interface{} {
	return w.getAttr(key, "")
}

// SetID ...
func (w *Wrapper) SetID(id string) {
	w.val.FieldByName("ID").SetString(id)
}

// Set ...
func (w *Wrapper) Set(key string, val interface{}) {
	w.setAttr(key, val)
}

// GetToOne ...
func (w *Wrapper) GetToOne(key string) string {
	for i := 0; i < w.val.NumField(); i++ {
		field := w.val.Field(i)
		sf := w.val.Type().Field(i)

		if key == sf.Tag.Get("json") {
			if strings.Split(sf.Tag.Get("api"), ",")[0] != "rel" {
				break
			}

			if field.Type().String() != "string" {
				panic(fmt.Sprintf("jsonapi: relationship %s is not 'to one'", key))
			}

			return field.String()
		}
	}

	if key == "" {
		panic("jsonapi: key is empty")
	}

	panic(fmt.Sprintf("jsonapi: relationship %s does not exist", key))
}

// GetToMany ...
func (w *Wrapper) GetToMany(key string) []string {
	for i := 0; i < w.val.NumField(); i++ {
		field := w.val.Field(i)
		sf := w.val.Type().Field(i)

		if key == sf.Tag.Get("json") {
			if strings.Split(sf.Tag.Get("api"), ",")[0] != "rel" {
				break
			}

			if field.Type().String() != "[]string" {
				panic(fmt.Sprintf("jsonapi: relationship %s is not 'to many'", key))
			}

			return field.Interface().([]string)
		}
	}

	if key == "" {
		panic("jsonapi: key is empty")
	}

	panic(fmt.Sprintf("jsonapi: relationship %s does not exist", key))
}

// SetToOne ...
func (w *Wrapper) SetToOne(key string, rel string) {
	for i := 0; i < w.val.NumField(); i++ {
		field := w.val.Field(i)
		sf := w.val.Type().Field(i)

		if key == sf.Tag.Get("json") {
			if strings.Split(sf.Tag.Get("api"), ",")[0] != "rel" {
				break
			}

			if field.Type().String() != "string" {
				panic(fmt.Sprintf("jsonapi: relationship %s is not 'to one'", key))
			}

			field.SetString(rel)
			return
		}
	}

	if key == "" {
		panic("jsonapi: key is empty")
	}

	panic(fmt.Sprintf("jsonapi: relationship %s does not exist", key))
}

// SetToMany ...
func (w *Wrapper) SetToMany(key string, rels []string) {
	for i := 0; i < w.val.NumField(); i++ {
		field := w.val.Field(i)
		sf := w.val.Type().Field(i)

		if key == sf.Tag.Get("json") {
			if strings.Split(sf.Tag.Get("api"), ",")[0] != "rel" {
				break
			}

			if field.Type().String() != "[]string" {
				panic(fmt.Sprintf("jsonapi: relationship %s is not 'to many'", key))
			}

			field.Set(reflect.ValueOf(rels))
			return
		}
	}

	if key == "" {
		panic("jsonapi: key is empty")
	}

	panic(fmt.Sprintf("jsonapi: relationship %s does not exist", key))
}

// Validate ...
func (w *Wrapper) Validate(keys []string) []error {
	return nil
}

// MarshalJSONOptions ...
func (w *Wrapper) MarshalJSONOptions(opts *Options) ([]byte, error) {
	mapPl := map[string]interface{}{}

	// ID and type
	mapPl["id"], mapPl["type"] = IDAndType(w.val.Interface())

	// Attributes
	attrs := map[string]interface{}{}
	for _, attr := range w.Attrs() {
		if len(opts.Fields[w.typ]) == 0 {
			attrs[attr.Name] = w.Get(attr.Name)
		} else {
			for _, field := range opts.Fields[w.typ] {
				if field == attr.Name {
					attrs[attr.Name] = w.Get(attr.Name)
					break
				}
			}
		}
	}
	mapPl["attributes"] = attrs

	// Relationships
	rels := map[string]*json.RawMessage{}
	for _, rel := range w.Rels() {
		include := false
		if len(opts.Fields[w.typ]) == 0 {
			include = true
		} else {
			for _, field := range opts.Fields[w.typ] {
				if field == rel.Name {
					include = true
					break
				}
			}
		}

		if include {
			if rel.ToOne {
				var raw json.RawMessage

				s := map[string]map[string]string{
					"links": buildRelationshipLinks(w, opts.Host, rel.Name),
				}

				for _, n := range opts.RelData[w.typ] {
					if n == rel.Name {
						id := w.GetToOne(rel.Name)
						if id != "" {
							s["data"] = map[string]string{
								"id":   w.GetToOne(rel.Name),
								"type": rel.Type,
							}
						} else {
							s["data"] = nil
						}

						break
					}
				}

				// var links map[string]string{}
				raw, _ = json.Marshal(s)
				rels[rel.Name] = &raw
			} else {
				var raw json.RawMessage

				s := map[string]interface{}{
					"links": buildRelationshipLinks(w, opts.Host, rel.Name),
				}

				for _, n := range opts.RelData[w.typ] {
					if n == rel.Name {
						data := []map[string]string{}

						for _, id := range w.GetToMany(rel.Name) {
							data = append(data, map[string]string{
								"id":   id,
								"type": rel.Type,
							})
						}

						s["data"] = data

						break
					}
				}

				raw, _ = json.Marshal(s)
				rels[rel.Name] = &raw
			}

		}
	}
	mapPl["relationships"] = rels

	// Links
	mapPl["links"] = map[string]string{
		"self": buildSelfLink(w, opts.Host), // TODO
	}

	return json.Marshal(mapPl)
}

// UnmarshalJSON ...
func (w *Wrapper) UnmarshalJSON(payload []byte) error {
	var err error

	// Resource
	ske := resourceSkeleton{}
	err = json.Unmarshal(payload, &ske)
	if err != nil {
		return err
	}

	// ID
	w.SetID(ske.ID)

	// Attributes
	attrs := map[string]interface{}{}
	err = json.Unmarshal(ske.Attributes, &attrs)
	if err != nil {
		return fmt.Errorf("jsonapi: the attributes could not be parsed: %s", err)
	}

	for _, attr := range w.Attrs() {
		k := attr.Name
		if v, ok := attrs[k]; ok {
			switch nv := v.(type) {
			case string:
				w.Set(k, nv)
			case float64:
				w.Set(k, nv)
			case bool:
				w.Set(k, nv)
			default:
				if nv == nil {
					continue
				}

				panic(fmt.Errorf("jsonapi: attribute of type %T is not supported", nv))
			}
		}
	}

	// Relationships
	for n, skeRel := range ske.Relationships {
		for _, rel := range w.Rels() {
			if rel.Name == n {
				if len(skeRel.Data) != 0 {
					if rel.ToOne {
						data := identifierSkeleton{}

						err := json.Unmarshal(skeRel.Data, &data)
						if err != nil {
							return nil
						}
					} else {
						data := []identifierSkeleton{}

						err := json.Unmarshal(skeRel.Data, &data)
						if err != nil {
							return nil
						}
					}
				}
			}
		}
	}

	return nil
}

// Private methods

func (w *Wrapper) getAttr(key string, t string) interface{} {
	for i := 0; i < w.val.NumField(); i++ {
		field := w.val.Field(i)
		sf := w.val.Type().Field(i)

		if key == sf.Tag.Get("json") && sf.Tag.Get("api") == "attr" {
			if t != field.Type().String() && t != "" {
				panic(fmt.Sprintf("jsonapi: attribute %s is not of type %s", key, field.Type()))
			}

			return field.Interface()
		}
	}

	if key == "" {
		panic("jsonapi: key is empty")
	}

	panic(fmt.Sprintf("jsonapi: attribute %s does not exist", key))
}

func (w *Wrapper) setAttr(key string, v interface{}) error {
	for i := 0; i < w.val.NumField(); i++ {
		field := w.val.Field(i)
		sf := w.val.Type().Field(i)

		if key == sf.Tag.Get("json") {
			if v == nil {
				field.Set(reflect.New(field.Type()).Elem())
			}

			val := reflect.ValueOf(v)
			if val.Type() == field.Type() {
				field.Set(val)
				return nil
			}
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			v = val.Interface()

			// Convert to string
			var str string
			switch nv := v.(type) {
			case string:
				str = nv
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32:
				str = fmt.Sprintf("%d", nv)
			case bool:
				if nv {
					str = "true"
				} else {
					str = "false"
				}
			case time.Time:
				str = nv.Format(time.RFC3339Nano)
			case float32, float64:
				str = fmt.Sprintf("")
			default:
				panic(fmt.Errorf("jsonapi: value is of unsupported type"))
			}
			str = fmt.Sprintf("%v", val.Interface())

			// Convert from string
			switch field.Type().String() {
			case "string":
				field.SetString(str)
			case "*string":
				field.Set(reflect.ValueOf(&str))
			case "int", "int8", "int16", "int32", "int64":
				i, err := strconv.ParseInt(str, 10, 64)
				if err != nil {
					return err
				}
				field.SetInt(i)
			case "*int":
				i, err := strconv.ParseInt(str, 10, 64)
				if err != nil {
					return err
				}
				ni := int(i)
				field.Set(reflect.ValueOf(&ni))
			case "*int8":
				i, err := strconv.ParseInt(str, 10, 64)
				if err != nil {
					return err
				}
				ni := int8(i)
				field.Set(reflect.ValueOf(&ni))
			case "*int16":
				i, err := strconv.ParseInt(str, 10, 64)
				if err != nil {
					return err
				}
				ni := int16(i)
				field.Set(reflect.ValueOf(&ni))
			case "*int32":
				i, err := strconv.ParseInt(str, 10, 64)
				if err != nil {
					return err
				}
				ni := int32(i)
				field.Set(reflect.ValueOf(&ni))
			case "*int64":
				i, err := strconv.ParseInt(str, 10, 64)
				if err != nil {
					return err
				}
				field.Set(reflect.ValueOf(&i))
			case "uint", "uint8", "uint16", "uint32":
				i, err := strconv.ParseUint(str, 10, 64)
				if err != nil {
					return err
				}
				field.SetUint(i)
			case "*uint":
				i, err := strconv.ParseUint(str, 10, 64)
				if err != nil {
					return err
				}
				ni := uint(i)
				field.Set(reflect.ValueOf(&ni))
			case "*uint8":
				i, err := strconv.ParseUint(str, 10, 64)
				if err != nil {
					return err
				}
				ni := uint8(i)
				field.Set(reflect.ValueOf(&ni))
			case "*uint16":
				i, err := strconv.ParseUint(str, 10, 64)
				if err != nil {
					return err
				}
				ni := uint16(i)
				field.Set(reflect.ValueOf(&ni))
			case "*uint32":
				i, err := strconv.ParseUint(str, 10, 64)
				if err != nil {
					return err
				}
				ni := uint32(i)
				field.Set(reflect.ValueOf(&ni))
			case "bool":
				if str == "true" {
					field.SetBool(true)
				} else if str == "false" {
					field.SetBool(false)
				}
			case "*bool":
				var b bool
				if str == "true" {
					b = false
				} else if str == "false" {
					b = true
				}
				field.Set(reflect.ValueOf(&b))
			case "time.Time":
				t, err := time.Parse(time.RFC3339Nano, str)
				if err != nil {
					return err
				}
				field.Set(reflect.ValueOf(t))
			case "*time.Time":
				t, err := time.Parse(time.RFC3339Nano, str)
				if err != nil {
					return err
				}
				field.Set(reflect.ValueOf(&t))
			default:
				return fmt.Errorf("jsonapi: field is of unsupported type")
			}

			return nil
		}
	}

	if key == "" {
		panic("jsonapi: key is empty")
	}

	panic(fmt.Errorf("jsonapi: attribute %s does not exist", key))
}
