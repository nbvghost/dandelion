package tag

import (
	"github.com/lib/pq"
	"github.com/nbvghost/dandelion/entity/extends"
	"regexp"
	"strings"
)

var kg = regexp.MustCompile(`\s+`)

func CreateUri(tags []extends.Tag) []extends.Tag {
	for i := 0; i < len(tags); i++ {
		tags[i].Uri = kg.ReplaceAllString(tags[i].Name, "-")
	}
	return tags
}
func ToTagsName(tags []extends.Tag) []extends.Tag {
	for i := 0; i < len(tags); i++ {
		tags[i].Name = strings.ReplaceAll(tags[i].Uri, "-", " ")
	}
	return tags
}
func ToTagsUri(arr pq.StringArray) []extends.Tag {
	var tags []extends.Tag
	for i := 0; i < len(arr); i++ {
		tags = append(tags, extends.Tag{
			Name:  arr[i],
			Count: 1,
			Uri:   "",
		})
	}
	return CreateUri(tags)
}
