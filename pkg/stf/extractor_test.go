package stf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStfTagRequiredOption_true_hides_fields(t *testing.T) {
	d := struct {
		Data string `json:"data"`
	}{
		Data: "secret",
	}

	tags, err := NewExtractor(TagRequiredOption(true)).ExtractFields(d)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(tags))
}

func TestStfTagRequiredOption_false_returns_fields(t *testing.T) {
	d := struct {
		Data string `json:"data"`
	}{
		Data: "secret",
	}

	tags, err := NewExtractor(TagRequiredOption(false)).ExtractFields(d)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(tags))
	assert.Equal(t, "secret", tags["data"])
}

func TestExcludeNilsOption_true_removes_nils(t *testing.T) {
	d := struct {
		Data []string `json:"data"`
	}{
		Data: nil,
	}

	tags, err := NewExtractor(ExcludeNilsOption(true)).ExtractFields(d)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(tags))
}

func TestExcludeNilsOption_false_returns_nils(t *testing.T) {
	d := struct {
		Data []string `json:"data"`
	}{
		Data: nil,
	}

	tags, err := NewExtractor(ExcludeNilsOption(false)).ExtractFields(d)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(tags))
}
