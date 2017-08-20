# `fun`ctional options

Opinionated Code Generators For Functional Options In Golang

## Roadmap

- [x] `lib/options`: initial simple implementation of generating source for functional options
- [x] `lib/parse`: initial implementation of source code -> fun internals parsing
- [x] `cmd/fun`: initial implemenation of the executable
- [ ] `lib/parse`: support parsing field tags on target structs
- [ ] `lib/options`: support generation of provided options e.g. generate `func WithSomeFieldName(value string) Option`
- [ ] `cmd/fun`: update binary to support wiring new options together
