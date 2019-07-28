package jsonapi

import "strings"

// buildSelfLink builds a URL that points to the resource represented by the
// value v.
//
// prepath is prepended to the path and usually represents a scheme and a
// domain name.
func buildSelfLink(res Resource, prepath string) string {
	if !strings.HasSuffix(prepath, "/") {
		prepath += "/"
	}

	if res.GetID() != "" && res.GetType().Name != "" {
		return prepath + res.GetType().Name + "/" + res.GetID()
	}

	return ""
}

// buildRelationshipLinks builds a links object (according to the JSON:API
// specification) that include both the self and related members.
func buildRelationshipLinks(res Resource, prepath, rel string) map[string]string {
	return map[string]string{
		"self":    buildSelfLink(res, prepath) + "/relationships/" + rel,
		"related": buildSelfLink(res, prepath) + "/" + rel,
	}
}
