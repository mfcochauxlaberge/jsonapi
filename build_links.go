package jsonapi

import "strings"

// buildSelfLink builds a URL that points to the resource
// represented by the value v.
//
// v has to be a struct or a pointer to a struct.
func buildSelfLink(res Resource, prepath string) string {
	if !strings.HasSuffix(prepath, "/") {
		prepath = prepath + "/"
	}

	if res.GetID() != "" && res.GetType() != "" {
		return prepath + res.GetType() + "/" + res.GetID()
	}

	return ""
}

func buildRelationshipLinks(res Resource, prepath, rel string) map[string]string {
	return map[string]string{
		"self":    buildSelfLink(res, prepath) + "/relationships/" + rel,
		"related": buildSelfLink(res, prepath) + "/" + rel,
	}
}
