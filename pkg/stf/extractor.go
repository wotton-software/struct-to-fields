package stf

type Extractor struct {
	TagRequired bool // Will only extract fields which have a valid STF Tag, extracts all fields by default.
	ExcludeNils bool // Will not return fields which are nil, does by default.
}

func NewExtractor(opts ...Option) *Extractor {
	e := &Extractor{}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

type Option func(e *Extractor)

func TagRequiredOption(b bool) Option {
	return func(e *Extractor) {
		e.TagRequired = b
	}
}

func ExcludeNilsOption(b bool) Option {
	return func(e *Extractor) {
		e.ExcludeNils = b
	}
}
