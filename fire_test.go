// Copyright 2018 The gofire Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package gofire

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func Func1(a, b int) (int, int) {
	return a + b, a - b
}

func Func2(a, b float64) (float64, float64) {
	return a + b, a - b
}

type SampleStruct struct {
	Name string
}

func (s SampleStruct) Add(a, b int) int {
	return a + b
}

func (s SampleStruct) Minus(a, b int) int {
	return a - b
}

func (s SampleStruct) String(str string) string {
	return strings.ToUpper(str)
}

func TestFunc1(t *testing.T) {
	os.Args = []string{"TestFunc1", "3", "5"}
	ret, err := Fire(Func1)
	expectedNumOut := reflect.ValueOf(Func1).Type().NumOut()
	if len(ret) != expectedNumOut {
		t.Errorf("%d return value expected, got %d", expectedNumOut, len(ret))
	}

	expected_ret1, expected_ret2 := Func1(3, 5)
	got_ret1 := int(ret[0].Int())
	got_ret2 := int(ret[1].Int())
	if expected_ret1 != got_ret1 || expected_ret2 != got_ret2 {
		t.Errorf("(%v, %v) is expected but got (%v, %v)", expected_ret1,
			expected_ret2, got_ret1, got_ret2)
	}
	if err != nil {
		t.Errorf("Error is not expected but got %v", err)
	}
}

func TestFunc2(t *testing.T) {
	os.Args = []string{"TestFunc2", "3.5", "5.4"}
	ret, err := Fire(Func2)
	expectedNumOut := reflect.ValueOf(Func1).Type().NumOut()
	if len(ret) != expectedNumOut {
		t.Errorf("%d return value expected, got %d", expectedNumOut, len(ret))
	}

	expected_ret1, expected_ret2 := Func2(3.5, 5.4)
	got_ret1 := ret[0].Float()
	got_ret2 := ret[1].Float()
	if expected_ret1 != got_ret1 || expected_ret2 != got_ret2 {
		t.Errorf("(%v, %v) is expected but got (%v, %v)", expected_ret1,
			expected_ret2, got_ret1, got_ret2)
	}
	if err != nil {
		t.Errorf("Error is not expected but got %v", err)
	}
}

func TestSampleStruct1(t *testing.T) {
	var s SampleStruct
	os.Args = []string{"TestSampleStruct", "Add", "3", "5"}
	ret, err := Fire(s)
	expectedNumOut := reflect.ValueOf(s.Add).Type().NumOut()
	if len(ret) != expectedNumOut {
		t.Errorf("%d return value expected, got %d", expectedNumOut, len(ret))
	}

	expected_ret := s.Add(3, 5)
	got_ret := int(ret[0].Int())
	if expected_ret != got_ret {
		t.Errorf("(%v) is expected but got (%v)", expected_ret, got_ret)
	}
	if err != nil {
		t.Errorf("Error is not expected but got %v", err)
	}
}

func TestSampleStruct2(t *testing.T) {
	var s SampleStruct
	os.Args = []string{"TestSampleStruct", "Minus", "3", "5"}
	ret, err := Fire(s)
	expectedNumOut := reflect.ValueOf(s.Add).Type().NumOut()
	if len(ret) != expectedNumOut {
		t.Errorf("%d return value expected, got %d", expectedNumOut, len(ret))
	}

	expected_ret := s.Minus(3, 5)
	got_ret := int(ret[0].Int())
	if expected_ret != got_ret {
		t.Errorf("(%v) is expected but got (%v)", expected_ret, got_ret)
	}
	if err != nil {
		t.Errorf("Error is not expected but got %v", err)
	}
}

func TestSampleStruct3(t *testing.T) {
	var s SampleStruct
	os.Args = []string{"TestSampleStruct", "String", "hello, world"}
	ret, err := Fire(s)
	expectedNumOut := reflect.ValueOf(s.String).Type().NumOut()
	if len(ret) != expectedNumOut {
		t.Errorf("%d return value expected, got %d", expectedNumOut, len(ret))
	}

	expected_ret := s.String("hello, world")
	got_ret := ret[0].String()
	if expected_ret != got_ret {
		t.Errorf("(%v) is expected but got (%v)", expected_ret, got_ret)
	}
	if err != nil {
		t.Errorf("Error is not expected but got %v", err)
	}
}

func TestSampleStructWrongCommand(t *testing.T) {
	var s SampleStruct
	os.Args = []string{"TestSampleStruct", "Wrong", "hello, world"}
	ret, err := Fire(s)
	if len(ret) != 0 {
		t.Errorf("%d return value expected, got %d", 0, len(ret))
	}

	if ret != nil {
		t.Errorf("(%v) is expected but got (%v)", nil, ret)
	}

	if err == nil || err.Error() != "Invalid command" {
		t.Errorf("Error is expected but got %v", err)
	}
}
