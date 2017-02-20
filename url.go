package jsonapi

// URL ...
type URL struct {
	Host string

	// URL
	URL           string   // /users/u1/articles?fields[users]=name,age
	URLNormalized string   // /users/u1/articles?fields[users]=age,name
	Path          []string // [users, u1, articles]
	Route         string   // /users/:id/articles

	// Resource
	ResType string
	ResID   string
	RelKind string
	Rel     Rel

	// Params
	Params *Params
}
