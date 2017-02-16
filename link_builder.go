package jsonapi

import "strings"

// buildSelfLink builds a URL that points to the resource
// represented by the value v.
//
// v has to be a struct or a pointer to a struct.
func buildSelfLink(res Resource, pre string) string {
	id, typ := res.IDAndType()

	if !strings.HasSuffix(pre, "/") {
		pre = pre + "/"
	}

	if id != "" && typ != "" {
		return pre + typ + "/" + id
	}

	return ""
}

func buildRelationshipLinks(res Resource, pre string, rel string) map[string]string {
	return map[string]string{
		"self":    buildSelfLink(res, pre) + "/relationships/" + rel,
		"related": buildSelfLink(res, pre) + "/" + rel,
	}
}
