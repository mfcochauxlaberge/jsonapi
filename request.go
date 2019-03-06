package jsonapi

// Request ...
type Request struct {
	Method string
	URL    *URL
	Doc    Document
}
