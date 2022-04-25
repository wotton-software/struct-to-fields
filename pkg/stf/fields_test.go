package stf

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtractFields_errors_on_nil(t *testing.T) {
	_, err := NewExtractor().ExtractFields(nil)
	assert.Equal(t, errors.New("needs struct or pointer to struct type, got: invalid"), err)
}

func TestExtractFields_errors_on_non_struct(t *testing.T) {
	d := []string{"one"}

	_, err := NewExtractor().ExtractFields(d)
	assert.Equal(t, errors.New("needs struct or pointer to struct type, got: slice"), err)
}

func TestExtractFields_recursive_data(t *testing.T) {
	type Data struct {
		Name string `json:"name"`
		Main *Data  `json:"main"`
	}

	b := &Data{Name: "sub"}
	d := Data{Main: b, Name: "Root"}
	b.Main = &d

	tags, err := NewExtractor().ExtractFields(d)
	assert.Nil(t, err)

	assert.Equal(t, 3, len(tags))
	root := tags["name"]
	assert.Equal(t, "Root", root)
	sub := tags["main.name"]
	assert.Equal(t, "sub", sub)
	recursiveRoot := tags["main.main.name"]
	assert.Equal(t, "Root", recursiveRoot)
}

func TestExtractFields_handles_nils(t *testing.T) {
	type Four struct {
		Sub string `json:"s3"`
	}
	type Three struct {
		Sub *Four `json:"s2"`
	}
	type Two struct {
		Sub *Three `json:"s1"`
	}
	type One struct {
		Main Two `json:"main"`
	}

	d := One{Main: Two{Sub: &Three{Sub: nil}}}

	tags, err := NewExtractor().ExtractFields(d)
	assert.Nil(t, err)
	assert.Equal(t, len(tags), 1)

	var expected *Four
	assert.Equal(t, expected, tags["main.s1.s2"])
}

func TestExtractFields_handles_pointers(t *testing.T) {
	type Four struct {
		Sub string `json:"s3"`
	}
	type Three struct {
		Sub *Four `json:"s2"`
	}
	type Two struct {
		Sub *Three `json:"s1"`
	}
	type One struct {
		Main Two `json:"main"`
	}

	key := "12312"
	d := One{Main: Two{Sub: &Three{Sub: &Four{Sub: key}}}}

	tags, err := NewExtractor().ExtractFields(d)
	assert.Nil(t, err)
	assert.Equal(t, len(tags), 1)
	assert.Equal(t, key, tags["main.s1.s2.s3"])
}

func TestExtractFields_handles_structs(t *testing.T) {
	type Three struct {
		Sub []string `json:"sub_sub"`
	}
	type Two struct {
		Sub Three `json:"sub"`
	}
	type One struct {
		Main Two `json:"main"`
	}

	slice := []string{"1234", "6789"}
	d := One{Main: Two{Sub: Three{Sub: slice}}}

	tags, err := NewExtractor().ExtractFields(d)
	assert.Nil(t, err)
	assert.Equal(t, len(tags), 1)
	assert.Equal(t, slice, tags["main.sub.sub_sub"])
}

func TestExtractFields_defaults(t *testing.T) {
	data := struct {
		Tag string
	}{
		Tag: "boop",
	}

	tags, err := NewExtractor().ExtractFields(data)
	assert.Nil(t, err)
	assert.Equal(t, len(tags), 1)
	assert.Equal(t, "boop", tags["Tag"])
}

func TestExtractFields_prefers_json(t *testing.T) {
	data := struct {
		Tag string `json:"taggy",stf:"json"`
	}{
		Tag: "boop",
	}

	tags, err := NewExtractor().ExtractFields(data)
	assert.Nil(t, err)
	assert.Equal(t, len(tags), 1)
	assert.Equal(t, "boop", tags["taggy"])
}

func TestExtractFields_prefers_stf_tag(t *testing.T) {
	data := struct {
		Tag string `json:"tag",stf:"stf_tag"`
	}{
		Tag: "boop",
	}

	tags, err := NewExtractor().ExtractFields(data)
	assert.Nil(t, err)
	assert.Equal(t, len(tags), 1)
	assert.Equal(t, "boop", tags["stf_tag"])
}

func TestExtractFields_ignored(t *testing.T) {
	data := struct {
		Ignored string `json:"stew",stf:"-"`
	}{
		Ignored: "boop",
	}

	tags, err := NewExtractor().ExtractFields(data)
	assert.Nil(t, err)
	assert.Equal(t, len(tags), 0)
}
