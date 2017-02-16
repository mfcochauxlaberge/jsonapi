# jsonapi

jsonapi is a complete library to marshal and unmarshal JSON API payloads.

The official specification can be found at [jsonapi.org/format](http://jsonapi.org/format).

Here's a list of features supported by the library:

 * Marshal and Unmarshal JSON API payloads
 * The generated payload is guaranteed to be valid according to the specification
 * All top-level members are supported (data, included, meta, errors...)
 * Structs can be used or interfaces can also be implemented
 * A bunch of utilities are provided to help build tools on top of this library

## State

The library is in **alpha** and its API is subject to change until v1 is released. Also, tests are not complete.

## Quick start

The simplest way to start using jsonapi is to use the Marshal and Unmarshal functions.

```
func Marshal(v interface{}, params *Params, extra Extra) ([]byte, error)
func Unmarshal(payload []byte, v interface{}) error
```

A struct has to follow certain rules in order to be understood by the library, but interfaces are also provided which let the library avoid the reflect package and be more efficient.

```
type User struct {
  ID string `json:"id" api:"users"` // ID is mandatory and the api tag sets the type

  // Attributes
  Name string `json:"name" api:"attr"` // attr means it is an attribute
  BornAt time.Time `json:"born-at" api:"attr"`

  // Relationships
  Articles []string `json:"articles" api:"rel,articles"`
}
```

It is recommended to use the json tag which sets the name of the field. Relationships can be a bit tricky. To-one relationships are defined with a string and to-many relationships are defined with a []slice. They contain the IDs of the related resources. The api tag has to take the form of "rel,xxx[,yyy]" where yyy is optional. xxx is the type of the relationship and yyy is the name of the inverse relationship when dealing with a two-way relationship. In this case, our Article struct would have an author relationship of type users:

```
Author string `json:"author" api:"rel,users,articles"`
```

## Utilities

Some structs and functions are provided to help build tools on top of this library.

For example:

```
// CheckType checks v to make sure it can be used as a type.
func CheckType(v interface{}) error {

// IDAndType returns the ID and type of v.
func IDAndType(v interface{}) (string, string)
```

See schema.go, document.go, url.go, and params.go for more.
