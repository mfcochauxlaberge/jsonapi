package jsonapi

type Meta map[string]interface{}

func (m *Meta) GetString(key string) string {
	v, _ := (*m)[key].(string)
	return v
}
