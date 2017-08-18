package jsonapi

import "strings"

// buildSelfLink builds a URL that points to the resource
// represented by the value v.
//
// v has to be a struct or a pointer to a struct.
func buildSelfLink(res Resource, scheme, host string) string {
	id, typ := res.IDAndType()

	pre := ""
	if scheme != "" && host != "" {
		pre = scheme + "://" + host
	}

	if !strings.HasSuffix(pre, "/") {
		pre = pre + "/"
	}

	if id != "" && typ != "" {
		return pre + typ + "/" + id
	}

	return ""
}

func buildRelationshipLinks(res Resource, scheme, host, rel string) map[string]string {
	return map[string]string{
		"self":    buildSelfLink(res, scheme, host) + "/relationships/" + rel,
		"related": buildSelfLink(res, scheme, host) + "/" + rel,
	}
}
