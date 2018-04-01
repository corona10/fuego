[![Build Status](https://travis-ci.org/corona10/gofire.svg?branch=master)](https://travis-ci.org/corona10/gofire)
[![GoDoc](https://godoc.org/github.com/corona10/goimghdr?status.svg)](https://godoc.org/github.com/corona10/gofire)
[![Go Report Card](https://goreportcard.com/badge/github.com/corona10/gofire)](https://goreportcard.com/report/github.com/corona10/gofire)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)

# gofire
> Inspired by Google [python-fire](https://github.com/google/python-fire)

gofire is a library for automatically generating command line interfaces (CLIs) from function and struct.

## Features
* gofire is a simple way to create a CLI in Go.
* gofire helps with exploring existing code or turning other people's code into a CLI.

## Installation
```
go get github.com/corona10/gofire
```

## TODO
- Support flag options
- More error handling
- Support more types

[![asciicast](https://asciinema.org/a/173759.png)](https://asciinema.org/a/173759)

## [Examples](/_examples)

```go
package main

import (
	"fmt"

	"github.com/corona10/gofire"
)

type Sample struct {
	Name string
}

func (s Sample) Add(a, b int) int {
	return a + b
}

func (s Sample) Minus(a, b int) int {
	return a - b
}

func (s Sample) HelloWorld() {
	fmt.Println(s.Name)
	fmt.Println("Hello world!")
}

func main() {
	var s Sample
	s.Name = "test"
	gofire.Fire(s)
}
```

## Special thanks to
* [Haeun Kim](https://github.com/haeungun/)
