package stf

import (
	"errors"
	"strings"
)

const (
	json       = "json"
	jsonPrefix = json + ":"
)

func hasJsonTag(s string) bool {
	return strings.Contains(s, jsonPrefix)
}

func getJsonTag(s string) (string, error) {
	tags := strings.Split(s, ",")

	for _, tag := range tags {
		if strings.HasPrefix(tag, jsonPrefix) {
			return trimTagToValue(tag, jsonPrefix), nil
		}
	}

	return "", errors.New("failed to find json tag: " + s)
}
