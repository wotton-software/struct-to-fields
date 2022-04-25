package stf

import (
	"errors"
	"strings"
)

const (
	prefix    = "stf:"
	ignorable = prefix + `"-"`
)

func getStfTag(s string) (string, error) {
	tags := strings.Split(s, ",")

	for _, tag := range tags {
		if strings.HasPrefix(tag, prefix) {
			return trimTagToValue(tag, prefix), nil
		}
	}

	return "", errors.New("failed to find stf tag: " + s)
}

func hasStfTag(s string) bool {
	return strings.Contains(s, prefix)
}

func hasStfIgnoreTag(s string) bool {
	return strings.Contains(s, ignorable)
}
