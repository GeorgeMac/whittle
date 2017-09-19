# Generator For Table Driven Tests 

Generate the necessary boiler plate for constructing table driven tests for a given types functions.

## Usage

`whittle table -type TypeToGenerateOptionsFor`

## Example

**./thing.go**

```go
type Thing struct{}

func (t Thing) DoSomething() error {}

func (t *Thing) Close() error {}
```

Run the following command:

`whittle table -type Thing`

This will generate generates two outputs 

**./thing_test.go**

```go
package thing

import "testing"

func TestThingDoSomething(t *testing.T) {
	for _, testCase := range []thingDoSomethingCase{
		{name: "happy path"},
	} {
		t.Run(testCase.name, testCase.Run)
	}
}

func TestThingClose(t *testing.T) {
	for _, testCase := range []thingCloseCase{
		{name: "happy path"},
	} {
		t.Run(testCase.name, testCase.Run)
	}
}
```


**./thing_table_test.go**

```go
package thing

import "testing"

type thingDoSomethingCase struct {
	name string
}

func (c thingDoSomethingCase) Run(t *testing.T) {}

type thingCloseCase struct {
	name string
}

func (c thingCloseCase) Run(t *testing.T) {}
```
