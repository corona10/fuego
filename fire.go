// Package gofire is a library for automatically generating command line interfaces (CLIs)
// from function and struct.
package gofire

// Copyright 2018 The gofire Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

// Config is a configuration struct for gofire
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
		funcName := getFunctionName(target)
		f := reflect.ValueOf(target)
		numParams := typ.NumIn()
		in := make([]reflect.Value, numParams)
		info[funcName] = typ

		if len(args) != numParams+1 {
			printFunctionHelp(typ, args)
			return nil, errors.New("Invalid command")
		}

		params := args[1:]

		// TODO (corona10): Support more features.
		for idx, param := range params {
			t := f.Type().In(idx)
			switch t.Kind() {
			case reflect.Int:
				paramValue, _ := strconv.Atoi(param)
				arg := reflect.ValueOf(paramValue)
				in[idx] = arg
			case reflect.Float32:
				paramValue, _ := strconv.ParseFloat(param, 32)
				arg := reflect.ValueOf(paramValue)
				in[idx] = arg
			case reflect.Float64:
				paramValue, _ := strconv.ParseFloat(param, 64)
				arg := reflect.ValueOf(paramValue)
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
				paramValue, _ := strconv.Atoi(param)
				arg := reflect.ValueOf(paramValue)
				in[idx] = arg
			case reflect.Float32:
				paramValue, _ := strconv.ParseFloat(param, 32)
				arg := reflect.ValueOf(paramValue)
				in[idx] = arg
			case reflect.Float64:
				paramValue, _ := strconv.ParseFloat(param, 64)
				arg := reflect.ValueOf(paramValue)
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
	var retValues []string
	if len(rets) > 0 {
		// TODO (corona10): Support more features.
		for _, ret := range rets {
			switch ret.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				retValues = append(retValues, strconv.FormatInt(ret.Int(), 10))
			case reflect.String:
				retValues = append(retValues, ret.String())
			case reflect.Float32, reflect.Float64:
				retValues = append(retValues, strconv.FormatFloat(ret.Float(), 'f', -1, 64))
			}
		}
		fmt.Println(strings.Join(retValues, " "))
	}
}

func getFunctionName(fn interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	funcNames := strings.Split(name, ".")
	return funcNames[len(funcNames)-1]
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

	showCommands := strings.Join(commands, "\n")
	msg = fmt.Sprintf("%s%s\n", msg, showCommands)
	fmt.Printf(msg)
}
