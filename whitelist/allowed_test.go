package whitelist

import (
	"testing"
)

func TestTagAllowed(t *testing.T) {
	okTags := AllowedTags()
	tag := okTags[4]
	if !IsAllowedTag(tag) {
		t.Errorf("Tag %s was not allowed, but it's on the allowed list.", tag)
	}
}

func TestTagNotAllowed(t *testing.T) {
	tag := "badtag"
	if IsAllowedTag(tag) {
		t.Errorf("Tag %s not allowed, but it's not on the allowed list.", tag)
	}
}

func TestAttributeAllowed(t *testing.T) {
	okAttrs := AllowedAttributes()
	attr := okAttrs[4]
	if !IsAllowedAttribute(attr) {
		t.Errorf("Attribute %s was not allowed, but it's on the allowed list.", attr)
	}
}

func TestAttributeNotAllowed(t *testing.T) {
	attr := "badattr"
	if IsAllowedAttribute(attr) {
		t.Errorf("Attribute %s not allowed, but it's not on the allowed list.", attr)
	}
}
