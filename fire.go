// Copyright 2018 The gofire Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package gofire

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

type Config struct {
	PrintReturnValuesOff bool
}

// Fire is a function for automatically generating command line interfaces (CLIs)
// from function and struct
func Fire(target interface{}, config ...Config) ([]reflect.Value, error) {
	var info map[string]reflect.Type
	var conf Config
	info = make(map[string]reflect.Type)
	typ := reflect.TypeOf(target)
	args := os.Args

	if len(config) > 0 {
		conf = config[0]
	}

	if typ.Kind() == reflect.Func {
		func_name := getFunctionName(target)
		f := reflect.ValueOf(target)
		num_params := typ.NumIn()
		in := make([]reflect.Value, num_params)
		info[func_name] = typ

		if len(args) != num_params+1 {
			printFunctionHelp(typ, args)
			return nil, errors.New("Invalid command")
		}

		params := args[1:]

		// TODO (corona10): Support more features.
		for idx, param := range params {
			t := f.Type().In(idx)
			switch t.Kind() {
			case reflect.Int:
				param_value, _ := strconv.Atoi(param)
				arg := reflect.ValueOf(param_value)
				in[idx] = arg
			case reflect.Float32:
				param_value, _ := strconv.ParseFloat(param, 32)
				arg := reflect.ValueOf(param_value)
				in[idx] = arg
			case reflect.Float64:
				param_value, _ := strconv.ParseFloat(param, 64)
				arg := reflect.ValueOf(param_value)
				in[idx] = arg
			case reflect.String:
				in[idx] = reflect.ValueOf(param)
			}
		}

		ret := f.Call(in)

		if !conf.PrintReturnValuesOff {
			printCallResult(ret)
		}

		return ret, nil
	} else if typ.Kind() == reflect.Struct {
		for i := 0; i < typ.NumMethod(); i++ {
			method := typ.Method(i)
			info[method.Name] = method.Type
		}

		if len(args) < 2 || info[args[1]] == nil || len(args[2:]) != info[args[1]].NumIn()-1 {
			printMethodHelp(info, args)
			return nil, errors.New("Invalid command")
		}

		method := info[args[1]]
		params := args[2:]
		in := make([]reflect.Value, len(params))
		// TODO (corona10): Support more features.
		for idx, param := range params {
			t := method.In(idx + 1)
			switch t.Kind() {
			case reflect.Int:
				param_value, _ := strconv.Atoi(param)
				arg := reflect.ValueOf(param_value)
				in[idx] = arg
			case reflect.Float32:
				param_value, _ := strconv.ParseFloat(param, 32)
				arg := reflect.ValueOf(param_value)
				in[idx] = arg
			case reflect.Float64:
				param_value, _ := strconv.ParseFloat(param, 64)
				arg := reflect.ValueOf(param_value)
				in[idx] = arg
			case reflect.String:
				in[idx] = reflect.ValueOf(param)
			}
		}
		ret := reflect.ValueOf(target).MethodByName(args[1]).Call(in)
		if !conf.PrintReturnValuesOff {
			printCallResult(ret)
		}
		return ret, nil
	} else {
		panic("Not supported yet")
	}
}

func printCallResult(rets []reflect.Value) {
	var rets_values []string
	if len(rets) > 0 {
		// TODO (corona10): Support more features.
		for _, ret := range rets {
			switch ret.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				rets_values = append(rets_values, strconv.FormatInt(ret.Int(), 10))
			case reflect.String:
				rets_values = append(rets_values, ret.String())
			case reflect.Float32, reflect.Float64:
				rets_values = append(rets_values, strconv.FormatFloat(ret.Float(), 'f', -1, 64))
			}
		}
		fmt.Println(strings.Join(rets_values, " "))
	}
}

func getFunctionName(fn interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	func_names := strings.Split(name, ".")
	return func_names[len(func_names)-1]
}

func printFunctionHelp(info reflect.Type, args []string) {
	var params []string
	file := filepath.Base(args[0])
	for i := 0; i < info.NumIn(); i++ {
		params = append(params, info.In(i).String())
	}

	msg := fmt.Sprintf("Usage:%2s%s %s\n", "", file, strings.Join(params, " "))
	fmt.Printf(msg)
}

func printMethodHelp(info map[string]reflect.Type, args []string) {
	var commands []string
	file := filepath.Base(args[0])
	msg := fmt.Sprintf("Usage:%2s%s\n", "", file)
	for key, value := range info {
		var command []string
		command = append(command, file)
		command = append(command, key)

		for i := 1; i < value.NumIn(); i++ {
			command = append(command, value.In(i).String())
		}
		commands = append(commands, strings.Join(command, " "))
	}

	show_commands := strings.Join(commands, "\n")
	msg = fmt.Sprintf("%s%s\n", msg, show_commands)
	fmt.Printf(msg)
}
