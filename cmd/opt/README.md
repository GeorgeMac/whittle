# Generator For Funtional `opt`ions 

Opinionated Code Generators For Functional Options In Golang

## Usage

`opt -type TypeToGenerateOptionsFor`

## Example

generate functional options for the current directories package

e.g. given package

```go
package worker

//go:generate opt -type Worker

type Worker struct {
	times     int    `opts`
	namespace string `opts:"SetNamespace"`
}

func New() *Worker { &Worker }
```

`go generate`

will produce:

```go
package worker

type Option func(*Worker)

type Options []Option

func (o Options) Apply(w *Worker) {
	for _, opt := range o {
		opt(w)
	}
}

func WithTimes(times int) Option {
	return func(w *Worker) {
		w.times = times
	}
}

func SetNamespace(namespace string) Option {
	return func(w *Worker) {
		w.namespace = namespace
	}
}
```

Now to wire in the new options, all you have to do is:

```go
func New(opts ...Option) *Worker {
    worker := &Worker{
        times: 5, // some sensible default
        namespace: "default_namespace" // another sensible default
    }

    Options(opts).Apply(worker)

    return worker
}
```

## Roadmap

- [x] `lib/options`: initial simple implementation of generating source for functional options
- [x] `lib/parse`: initial implementation of source code -> fun internals parsing
- [x] `cmd/fun`: initial implemenation of the executable
- [x] `lib/parse`: support parsing field tags on target structs
- [x] `lib/options`: support generation of provided options e.g. generate `func WithSomeFieldName(value string) Option`
- [x] `cmd/fun`: update binary to support wiring new options together
- [x] `cmd/opt`: rename `fun` to `opt` as it will now live as a sub-command of `whittle`
