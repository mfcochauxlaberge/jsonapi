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

	// ID and type
	id, t string

	// Structure
	attrs []Attr
	rels  map[string]Rel
}

// Wrap ...
func Wrap(v interface{}) *Wrapper {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Ptr {
		if val.Kind() != reflect.Struct {
			panic(errors.New("jsonapi: value has to be a pointer to a struct"))
		}

		val = reflect.New(val.Type())
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
	w.id, w.t = IDAndType(v)

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
	w.rels = map[string]Rel{}
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
			w.rels[jsonTag] = Rel{
				Name:        jsonTag,
				Type:        relTag[1],
				ToOne:       toOne,
				InverseName: invName,
				InverseType: w.t,
			}
		}
	}

	return w
}

// IDAndType ...
func (r *Wrapper) IDAndType() (string, string) {
	return r.id, r.t
}

// Attrs ...
func (r *Wrapper) Attrs() []Attr {
	return r.attrs
}

// Rels ...
func (r *Wrapper) Rels() map[string]Rel {
	return r.rels
}

// // AttrPtrs ...
// func (r *Wrapper) AttrPtrs(attrs []Attr) ([]string, []interface{}) {
// 	keys := []string{}
// 	dsts := []interface{}{}
//
// 	for i := range attrs {
// 		for j := range r.Attrs() {
// 			if attrs[i].Name == r.Attrs()[j].Name {
// 				dest[attrs[i].Name] = &sql.NullString{}
// 			}
// 		}
// 	}
//
// 	return dest
// }
//
// // ApplyAttrPtrs ...
// func (r *Wrapper) ApplyAttrPtrs(ptrs map[string]interface{}) {
// 	// TODO
// }

// New ...
func (r *Wrapper) New() Resource {
	newVal := reflect.New(r.val.Type())

	return Wrap(newVal.Interface())
}

// Get ...
func (r *Wrapper) Get(key string) interface{} {
	return r.getAttr(key, "")
}

// GetID ...
func (r *Wrapper) GetID() string {
	return r.getAttr("", "string").(string)
}

// GetString ...
func (r *Wrapper) GetString(key string) string {
	return r.getAttr(key, "string").(string)
}

// GetStringPtr ...
func (r *Wrapper) GetStringPtr(key string) *string {
	return r.getAttr(key, "*string").(*string)
}

// GetInt ...
func (r *Wrapper) GetInt(key string) int {
	return r.getAttr(key, "int").(int)
}

// GetIntPtr ...
func (r *Wrapper) GetIntPtr(key string) *int {
	return r.getAttr(key, "*int").(*int)
}

// GetInt8 ...
func (r *Wrapper) GetInt8(key string) int8 {
	return r.getAttr(key, "int8").(int8)
}

// GetInt8Ptr ...
func (r *Wrapper) GetInt8Ptr(key string) *int8 {
	return r.getAttr(key, "*int8").(*int8)
}

// GetInt16 ...
func (r *Wrapper) GetInt16(key string) int16 {
	return r.getAttr(key, "int16").(int16)
}

// GetInt16Ptr ...
func (r *Wrapper) GetInt16Ptr(key string) *int16 {
	return r.getAttr(key, "int16").(*int16)
}

// GetInt32 ...
func (r *Wrapper) GetInt32(key string) int32 {
	return r.getAttr(key, "int32").(int32)
}

// GetInt32Ptr ...
func (r *Wrapper) GetInt32Ptr(key string) *int32 {
	return r.getAttr(key, "*int32").(*int32)
}

// GetInt64 ...
func (r *Wrapper) GetInt64(key string) int64 {
	return r.getAttr(key, "int64").(int64)
}

// GetInt64Ptr ...
func (r *Wrapper) GetInt64Ptr(key string) *int64 {
	return r.getAttr(key, "*int64").(*int64)
}

// GetUint ...
func (r *Wrapper) GetUint(key string) uint {
	return r.getAttr(key, "uint").(uint)
}

// GetUintPtr ...
func (r *Wrapper) GetUintPtr(key string) *uint {
	return r.getAttr(key, "*uint").(*uint)
}

// GetUint8 ...
func (r *Wrapper) GetUint8(key string) uint8 {
	return r.getAttr(key, "uint8").(uint8)
}

// GetUint8Ptr ...
func (r *Wrapper) GetUint8Ptr(key string) *uint8 {
	return r.getAttr(key, "*uint8").(*uint8)
}

// GetUint16 ...
func (r *Wrapper) GetUint16(key string) uint16 {
	return r.getAttr(key, "uint16").(uint16)
}

// GetUint16Ptr ...
func (r *Wrapper) GetUint16Ptr(key string) *uint16 {
	return r.getAttr(key, "*uint16").(*uint16)
}

// GetUint32 ...
func (r *Wrapper) GetUint32(key string) uint32 {
	return r.getAttr(key, "uint32").(uint32)
}

// GetUint32Ptr ...
func (r *Wrapper) GetUint32Ptr(key string) *uint32 {
	return r.getAttr(key, "*uint32").(*uint32)
}

// GetBool ...
func (r *Wrapper) GetBool(key string) bool {
	return r.getAttr(key, "bool").(bool)
}

// GetBoolPtr ...
func (r *Wrapper) GetBoolPtr(key string) *bool {
	return r.getAttr(key, "*bool").(*bool)
}

// GetTime ...
func (r *Wrapper) GetTime(key string) time.Time {
	return r.getAttr(key, "time.Time").(time.Time)
}

// GetTimePtr ...
func (r *Wrapper) GetTimePtr(key string) *time.Time {
	return r.getAttr(key, "*time.Time").(*time.Time)
}

// Set ...
func (r *Wrapper) Set(key string, val interface{}) {
	r.setAttr(key, val)
}

// SetString ...
func (r *Wrapper) SetString(key string, val string) {
	r.setAttr(key, val)
}

// SetStringPtr ...
func (r *Wrapper) SetStringPtr(key string, val *string) {
	r.setAttr(key, val)
}

// SetInt ...
func (r *Wrapper) SetInt(key string, val int) {
	r.setAttr(key, val)
}

// SetIntPtr ...
func (r *Wrapper) SetIntPtr(key string, val *int) {
	r.setAttr(key, val)
}

// SetInt8 ...
func (r *Wrapper) SetInt8(key string, val int8) {
	r.setAttr(key, val)
}

// SetInt8Ptr ...
func (r *Wrapper) SetInt8Ptr(key string, val *int8) {
	r.setAttr(key, val)
}

// SetInt16 ...
func (r *Wrapper) SetInt16(key string, val int16) {
	r.setAttr(key, val)
}

// SetInt16Ptr ...
func (r *Wrapper) SetInt16Ptr(key string, val *int16) {
	r.setAttr(key, val)
}

// SetInt32 ...
func (r *Wrapper) SetInt32(key string, val int32) {
	r.setAttr(key, val)
}

// SetInt32Ptr ...
func (r *Wrapper) SetInt32Ptr(key string, val *int32) {
	r.setAttr(key, val)
}

// SetInt64 ...
func (r *Wrapper) SetInt64(key string, val int64) {
	r.setAttr(key, val)
}

// SetInt64Ptr ...
func (r *Wrapper) SetInt64Ptr(key string, val *int64) {
	r.setAttr(key, val)
}

// SetUint ...
func (r *Wrapper) SetUint(key string, val uint) {
	r.setAttr(key, val)
}

// SetUintPtr ...
func (r *Wrapper) SetUintPtr(key string, val *uint) {
	r.setAttr(key, val)
}

// SetUint8 ...
func (r *Wrapper) SetUint8(key string, val uint8) {
	r.setAttr(key, val)
}

// SetUint8Ptr ...
func (r *Wrapper) SetUint8Ptr(key string, val *uint8) {
	r.setAttr(key, val)
}

// SetUint16 ...
func (r *Wrapper) SetUint16(key string, val uint16) {
	r.setAttr(key, val)
}

// SetUint16Ptr ...
func (r *Wrapper) SetUint16Ptr(key string, val *uint16) {
	r.setAttr(key, val)
}

// SetUint32 ...
func (r *Wrapper) SetUint32(key string, val uint32) {
	r.setAttr(key, val)
}

// SetUint32Ptr ...
func (r *Wrapper) SetUint32Ptr(key string, val *uint32) {
	r.setAttr(key, val)
}

// SetFloat64 ...
func (r *Wrapper) SetFloat64(key string, val float64) {
	r.setAttr(key, val)
}

// SetBool ...
func (r *Wrapper) SetBool(key string, val bool) {
	r.setAttr(key, val)
}

// SetBoolPtr ...
func (r *Wrapper) SetBoolPtr(key string, val *bool) {
	r.setAttr(key, val)
}

// SetTime ...
func (r *Wrapper) SetTime(key string, val time.Time) {
	r.setAttr(key, val)
}

// SetTimePtr ...
func (r *Wrapper) SetTimePtr(key string, val *time.Time) {
	r.setAttr(key, val)
}

// GetToOne ...
func (r *Wrapper) GetToOne(key string) string {
	for i := 0; i < r.val.NumField(); i++ {
		field := r.val.Field(i)
		sf := r.val.Type().Field(i)

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
func (r *Wrapper) GetToMany(key string) []string {
	for i := 0; i < r.val.NumField(); i++ {
		field := r.val.Field(i)
		sf := r.val.Type().Field(i)

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
func (r *Wrapper) SetToOne(key string, rel string) {
	for i := 0; i < r.val.NumField(); i++ {
		field := r.val.Field(i)
		sf := r.val.Type().Field(i)

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
func (r *Wrapper) SetToMany(key string, rels []string) {
	for i := 0; i < r.val.NumField(); i++ {
		field := r.val.Field(i)
		sf := r.val.Type().Field(i)

		if key == sf.Tag.Get("json") {
			if strings.Split(sf.Tag.Get("api"), ",")[0] != "rel" {
				break
			}

			if field.Type().String() != "string" {
				panic(fmt.Sprintf("jsonapi: relationship %s is not 'to one'", key))
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
func (r *Wrapper) Validate(keys []string) []error {
	return nil
}

// MarshalJSON ...
// func (r *Wrapper) MarshalJSON() ([]byte, error) {
// 	return []byte{}, nil
// }

// MarshalJSONParams ...
func (r *Wrapper) MarshalJSONParams(params *Params) ([]byte, error) {
	mapPl := map[string]interface{}{}

	// ID and type
	mapPl["id"] = r.id
	mapPl["type"] = r.t

	// Attributes
	attrs := map[string]interface{}{}
	for _, attr := range r.Attrs() {
		if len(params.Fields[r.t]) == 0 {
			attrs[attr.Name] = r.Get(attr.Name)
		} else {
			for _, field := range params.Fields[r.t] {
				if field == attr.Name {
					attrs[attr.Name] = r.Get(attr.Name)
					break
				}
			}
		}
	}
	mapPl["attributes"] = attrs

	// Relationships
	rels := map[string]*json.RawMessage{}
	for _, rel := range r.Rels() {
		include := false
		if len(params.Fields[r.t]) == 0 {
			include = true
		} else {
			for _, field := range params.Fields[r.t] {
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
					"links": buildRelationshipLinks(r, "https://example.com", rel.Name),
				}

				for _, n := range params.RelData[r.t] {
					if n == rel.Name {
						id := r.GetToOne(rel.Name)
						if id != "" {
							s["data"] = map[string]string{
								"id":   r.GetToOne(rel.Name),
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
					"links": buildRelationshipLinks(r, "https://example.com", rel.Name),
				}

				for _, n := range params.RelData[r.t] {
					if n == rel.Name {
						data := []map[string]string{}

						for _, id := range r.GetToMany(rel.Name) {
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
		"self": buildSelfLink(r, "https://example.com/"), // TODO
	}

	return json.Marshal(mapPl)
}

// UnmarshalJSON ...
func (r *Wrapper) UnmarshalJSON(payload []byte) error {
	var err error

	// Resource
	ske := resourceSkeleton{}
	err = json.Unmarshal(payload, &ske)
	if err != nil {
		return err
	}

	// ID
	r.SetString("id", ske.ID)

	// Attributes
	attrs := map[string]interface{}{}
	err = json.Unmarshal(ske.Attributes, &attrs)
	if err != nil {
		return fmt.Errorf("jsonapi: the attributes could not be parsed: %s", err)
	}

	for _, attr := range r.Attrs() {
		k := attr.Name
		if v, ok := attrs[k]; ok {
			switch nv := v.(type) {
			case string:
				r.SetString(k, nv)
			case float64:
				r.SetFloat64(k, nv)
			case bool:
				r.SetBool(k, nv)
			default:
				panic(fmt.Errorf("jsonapi: attribute of unsupported type encountered"))
			}
		}
	}

	// Relationships
	for n, skeRel := range ske.Relationships {
		if rel, ok := r.Rels()[n]; ok {
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

	return nil
}

// Private methods

func (r *Wrapper) getAttr(key string, t string) interface{} {
	for i := 0; i < r.val.NumField(); i++ {
		field := r.val.Field(i)
		sf := r.val.Type().Field(i)

		if key == sf.Tag.Get("json") {
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

func (r *Wrapper) setAttr(key string, v interface{}) error {
	for i := 0; i < r.val.NumField(); i++ {
		field := r.val.Field(i)
		sf := r.val.Type().Field(i)

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
			case "*int", "*int8", "*int16", "*int32", "*int64":
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
			case "*uint", "*uint8", "*uint16", "*uint32":
				i, err := strconv.ParseUint(str, 10, 64)
				if err != nil {
					return err
				}
				field.Set(reflect.ValueOf(&i))
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
