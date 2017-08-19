package important

type Option func(*Important)

type Options []Option

func (o Options) Apply(i *Important) {
	for _, opt := range o {
		opt(i)
	}
}
