package important

import "os"

// Option is a functional option for the Important type
type Option func(*Important)

// Options is a slice of Option types
type Options []Option

// Apply calls each option in order to
// the supplied Important type
func (o Options) Apply(i *Important) {
	for _, opt := range o {
		opt(i)
	}
}

// WithField sets the field string on the
// Important type
func WithField(field string) Option {
	return func(i *Important) {
		i.field = field
	}
}

// WithAttribute sets the attribute int on the
// Important type
func WithAttribute(attribute int) Option {
	return func(i *Important) {
		i.attribute = attribute
	}
}

// WithThings sets the mapOfThings map[string]string on the
// Important type
func WithThings(mapOfThings map[string]string) Option {
	return func(i *Important) {
		i.mapOfThings = mapOfThings
	}
}

// WithPointerToThing sets the pointerToThing *string on the
// Important type
func WithPointerToThing(pointerToThing *string) Option {
	return func(i *Important) {
		i.pointerToThing = pointerToThing
	}
}

// WithPointerToStruct sets the pointerToStruct *os.File on the
// Important type
func WithPointerToStruct(pointerToStruct *os.File) Option {
	return func(i *Important) {
		i.pointerToStruct = pointerToStruct
	}
}
