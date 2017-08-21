package important

type Option func(*Important)

type Options []Option

func (o Options) Apply(i *Important) {
	for _, opt := range o {
		opt(i)
	}
}

func WithField(field string) Option {
	return func(i *Important) {
		i.field = field
	}
}

func WithAttribute(attribute int) Option {
	return func(i *Important) {
		i.attribute = attribute
	}
}

func WithThings(mapOfThings map[string]string) Option {
	return func(i *Important) {
		i.mapOfThings = mapOfThings
	}
}
