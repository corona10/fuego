// Copyright 2018 The fuego Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"

	"github.com/corona10/fuego"
)

type Sample struct {
	Name string
}

// Add is a method for Add.
func (s Sample) Add(a, b int) int {
	return a + b
}

// Minus is a method for Minus.
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
	fuego.Fire(s)
}
