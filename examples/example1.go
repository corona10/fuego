// Copyright 2018 The gofire Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"

	"github.com/corona10/gofire"
)

func Add(a int, b int) (int, int) {
	fmt.Println(a, b)
	return a + b, 2*a + b
}

func main() {
	gofire.Fire(Add)
}
