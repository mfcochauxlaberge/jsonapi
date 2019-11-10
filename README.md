# jsonapi

<div align="center" style="text-align: center;">
  <img src="assets/logo.png" height="120">
  <br>
  <a href="https://travis-ci.com/mfcochauxlaberge/jsonapi">
    <img src="https://travis-ci.com/mfcochauxlaberge/jsonapi.svg?branch=master">
  </a>
  <a href="https://goreportcard.com/report/github.com/mfcochauxlaberge/jsonapi">
    <img src="https://goreportcard.com/badge/github.com/mfcochauxlaberge/jsonapi">
  </a>
  <a href="https://codecov.io/gh/mfcochauxlaberge/jsonapi">
    <img src="https://img.shields.io/codecov/c/github/mfcochauxlaberge/jsonapi">
  </a>
  <br>
  <a href="https://github.com/mfcochauxlaberge/jsonapi/blob/master/go.mod">
    <img src="https://img.shields.io/badge/go%20version-go1.11%2B-red">
  </a>
  <a href="https://github.com/mfcochauxlaberge/jsonapi/blob/master/go.mod">
    <img src="https://img.shields.io/github/v/release/mfcochauxlaberge/jsonapi?include_prereleases&sort=semver">
  </a>
  <a href="https://github.com/mfcochauxlaberge/jsonapi/blob/master/LICENSE">
    <img src="https://img.shields.io/github/license/mfcochauxlaberge/jsonapi?color=a33">
  </a>
  <a href="https://godoc.org/github.com/mfcochauxlaberge/jsonapi">
    <img src="https://godoc.org/github.com/golang/gddo?status.svg">
  </a>
</div>

jsonapi is a complete library to marshal and unmarshal JSON:API payloads.

It also offers many tools for parsing URLs, handling resources, managing schemas, and building other tools on top of this library.

The official specification can be found at [jsonapi.org/format](http://jsonapi.org/format).

## State

The library is in **beta** and its API is subject to change until v1 is released.

## Requirements

The supported versions of Go are the latest patch releases of every minor release starting with Go 1.11.

## Quick start

The simplest way to start using jsonapi is to use the MarshalDocument and UnmarshalDocument functions.

```go
func MarshalDocument(doc *Document, url *URL) ([]byte, error)
func UnmarshalDocument(payload []byte, schema *Schema) (*Document, error)
```

A struct has to follow certain rules in order to be understood by the library, but interfaces are also provided which let the library avoid the reflect package and be more efficient.

See the following section for more information about how to define structs for this library.

## Concepts

Here are some of the main concepts covered by the library.

### Request

A `Request` represents an HTTP request structured in a format easily readable from a JSON:API point of view.

If you are familiar to the specification, reading the `Request` struct and its fields (`URL`, `Document`, etc) should be straightforward.

### Type

A JSON:API type is generally defined with a struct.

There needs to be an ID field of type string. The `api` tag represents the name of the type.

```go
type User struct {
  ID string `json:"id" api:"users"` // ID is mandatory and the api tag sets the type

  // Attributes
  Name string `json:"name" api:"attr"` // attr means it is an attribute
  BornAt time.Time `json:"born-at" api:"attr"`

  // Relationships
  Articles []string `json:"articles" api:"rel,articles"`
}
```

Other fields with the `api` tag (`attr` or `rel`) can be added as attributes or relationships.

#### Attribute

Attributes can be of the following types:

```go
string
int, int8, int16, int32, int64
uint, uint8, uint16, uint32, uint64
bool
time.Time
[]byte
*string
*int, *int8, *int16, *int32, *int64
*uint, *uint8, *uint16, *uint32, *uint64
*bool
*time.Time
*[]byte
```

Using a pointer allows the field to be nil.

#### Relationship

Relationships can be a bit tricky. To-one relationships are defined with a string and to-many relationships are defined with a slice of strings. They contain the IDs of the related resources. The api tag has to take the form of "rel,xxx[,yyy]" where yyy is optional. xxx is the type of the relationship and yyy is the name of the inverse relationship when dealing with a two-way relationship. In the following example, our Article struct defines a relationship named author of type users:

```go
Author string `json:"author" api:"rel,users,articles"`
```

### Wrapper

A struct can be wrapped using the `Wrap` function which returns a pointer to a `Wrapper`. A `Wrapper` implements the `Resource` interface and can be used with this library. Modifying a Wrapper will modify the underlying struct. The resource's type is defined from reflecting on the struct.

```go
user := User{}
wrap := Wrap(&user)
wrap.Set("name", "Mike")
fmt.Printf(wrap.Get("name")) // Output: Mike
fmt.Printf(user.Name) // Output: Mike
```

### SoftResource

A SoftResource is a struct whose type (name, attributes, and relationships) can be modified indefinitely just like its values. When an attribute or a relationship is added, the new value is the zero value of the field type. For example, if you add an attribute named `my-attribute` of type string, then `softresource.Get("my-attribute")` will return an empty string.

```go
sr := SoftResource{}
sr.AddAttr(Attr{
  Name: "attr",
  Type: AttrTypeInt,
  Null: false,
})
fmt.Println(sr.Get("attr")) // Output: 0
```

Take a look at the `SoftCollection` struct for a similar concept applied to an entire collection of resources.

### URLs

From a raw string that represents a URL, it is possible that create a `SimpleURL` which contains the information stored in the URL in a structure that is easier to handle.

It is also possible to build a `URL` from a `Schema` and a `SimpleURL` which contains additional information taken from the schema. `NewURL` returns an error if the URL does not respect the schema.

## Documentation

Check out the [incomplete documentation](https://godoc.org/github.com/mfcochauxlaberge/jsonapi).

The best way to learn how to use it is to look at the source code and its comments.
