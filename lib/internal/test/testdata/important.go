package important

// go:generate fun --include="Important"

// Important
type Important struct {
	Field       string            `opts`
	Attribute   int               `opts`
	MapOfThings map[string]string `opts:"WithThings"`
}
