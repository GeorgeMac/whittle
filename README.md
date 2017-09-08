# Whittle - Go Code Generation Tools + Libraries 

Opinionated Code Generators For Golang 

[![Go Report Card](https://goreportcard.com/badge/github.com/georgemac/whittle)](https://goreportcard.com/report/github.com/georgemac/whittle)

## Commands

1. [options](./cmd/whittle/options) for generating functional options for struct definitions 

e.g.

```go
//go:generate whittle options -type TypeToGeneratorFor
```

## Ideas Roadmap

- [ ] generation of table drive test cases
- [ ] explore generation of implementations for interfaces, with quick mock generation in mind.
- [ ] quick "constructor" style function insertion, with support for functional options.
