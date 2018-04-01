// Copyright 2018 The fuego Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"

	"github.com/corona10/fuego"
)

func Add(a int, b int) (int, int) {
	fmt.Println(a, b)
	return a + b, a - b
}

func main() {
	config := fuego.Config{PrintReturnValuesOff: true}
	fuego.Fire(Add, config)
}
