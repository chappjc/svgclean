package whitelist

func AllowedTags() []string {
	all := append(HTMLTags, SVGTags...)
	all = append(all, MathMLTags...)
	return append(all, HTMLTags...)
}

func AllowedAttributes() []string {
	all := append(HTMLTags, SVGAttributes...)
	all = append(all, MathMLAttributes...)
	all = append(all, XMLAttributes...)
	return append(all, HTMLAttributes...)
}

var (
	allowedTags       map[string]struct{}
	allowedAttributes map[string]struct{}
)

func init() {
	tags := AllowedTags()
	allowedTags = make(map[string]struct{}, len(tags))
	for i := range tags {
		allowedTags[tags[i]] = struct{}{}
	}

	attrs := AllowedAttributes()
	allowedAttributes = make(map[string]struct{}, len(attrs))
	for i := range attrs {
		allowedAttributes[attrs[i]] = struct{}{}
	}
}

func IsAllowedTag(tag string) bool {
	_, ok := allowedTags[tag]
	return ok
}

func IsAllowedAttribute(attr string) bool {
	_, ok := allowedAttributes[attr]
	return ok
}
