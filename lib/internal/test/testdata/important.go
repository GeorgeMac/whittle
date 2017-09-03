package important

import "os"

// go:generate fun --include="Important"

// Important
type Important struct {
	field           string            `opts`
	attribute       int               `opts`
	mapOfThings     map[string]string `opts:"WithThings"`
	pointerToThing  *string           `opts`
	pointerToStruct *os.File          `opts`
}
