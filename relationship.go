package jsonapi

// RelTemp represents a relationship between two types.
type RelTemp struct {
	Type1  string
	Name1  string
	ToOne1 bool
	Type2  string
	Name2  string
	ToOne2 bool
}

// Names builds and returns the name of the relationship.
func (r RelTemp) Name() string {
	var name string
	if r.Type1 < r.Type2 {
		name = r.Type1 + "_" + r.Name1 + "_" + r.Type2 + "_" + r.Name2
	} else {
		name = r.Type2 + "_" + r.Name2 + "_" + r.Type1 + "_" + r.Name1
	}
	return name
}
