package jsonapi

// Request ...
type Request struct {
	Method string
	URL    *URL
	User   string
	Doc    Document
}
