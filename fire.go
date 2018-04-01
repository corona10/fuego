// Package fuego is a library for automatically generating command line interfaces (CLIs)
// from function and struct.
package fuego

// Copyright 2018 The fuego Authors. All rights reserved.
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

// Config is a configuration struct for fuego
type Config struct {
	PrintReturnValuesOff bool
}

// Fire is a function for automatically generating command line interfaces (CLIs)
// from function and struct
func Fire(target interface{}, config ...Config) ([]reflect.Value, error) {
	var info map[string]Sym
	var conf Config
	info = make(map[string]Sym)
	typ := reflect.TypeOf(target)
	args := os.Args

	if len(config) > 0 {
		conf = config[0]
	}

	if typ.Kind() == reflect.Func {
		var sym Sym
		funcName := getFunctionName(target)
		sym.SetKind(Func)
		sym.SetName(funcName)
		sym.SetCall(reflect.ValueOf(target))
		sym.SetParams(args[1:])
		if len(args) != sym.GetNumOfNeededArgs() {
			printFunctionHelp(typ, args)
			return nil, errors.New("Invalid command")
		}
		ret := sym.Call()
		if !conf.PrintReturnValuesOff {
			printCallResult(ret)
		}
		return ret, nil
	} else if typ.Kind() == reflect.Struct {

		for i := 0; i < typ.NumMethod(); i++ {
			var sym Sym
			method := typ.Method(i)
			sym.SetKind(Method)
			sym.SetName(method.Name)
			sym.SetCall(reflect.ValueOf(target).MethodByName(method.Name))
			info[method.Name] = sym
		}

		if len(args) < 2 {
			printMethodHelp(info, args)
			return nil, errors.New("Invalid command")
		}

		sym, ok := info[args[1]]
		if !ok || len(args) != sym.GetNumOfNeededArgs() {
			printMethodHelp(info, args)
			return nil, errors.New("Invalid command")
		}

		params := args[2:]
		sym.SetParams(params)
		ret := sym.Call()
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

	msg := fmt.Sprintf("Usage:%s %s\n", file, strings.Join(params, " "))
	fmt.Printf(msg)
}

func printMethodHelp(info map[string]Sym, args []string) {
	var commands []string
	file := filepath.Base(args[0])
	msg := fmt.Sprintf("Usage: %s\n", file)
	for key, value := range info {
		var command []string
		command = append(command, file)
		command = append(command, key)

		for i := 0; i < value.GetNumIns(); i++ {
			command = append(command, value.GetIn(i).String())
		}
		commands = append(commands, strings.Join(command, " "))
	}

	showCommands := strings.Join(commands, "\n")
	msg = fmt.Sprintf("%s%s\n", msg, showCommands)
	fmt.Printf(msg)
}
