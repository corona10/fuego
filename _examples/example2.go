// Copyright 2018 The gofire Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
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
