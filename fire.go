// Package fuego is a library for automatically generating command line interfaces (CLIs)
// from function and struct.
package fuego

// Copyright 2018 The fuego Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
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
	var conf Config
	typ := reflect.TypeOf(target)
	args := os.Args
	_, callerFile, _, _ := runtime.Caller(1)
	info, err := getASTInfo(callerFile)

	if err != nil {
		return nil, err
	}

	if len(config) > 0 {
		conf = config[0]
	}

	if typ.Kind() == reflect.Func {
		funcName := getFunctionName(target)
		sym, ok := info[funcName]
		if !ok {
			msg := fmt.Sprintf("Documentation of %s should be found from AST information.\n", funcName)
			sym.SetDoc(msg)
		}
		sym.SetKind(Func)
		sym.SetCall(reflect.ValueOf(target))
		sym.SetParams(args[1:])
		// TODO(corona10): Fix this, not to check isvalid.
		sym.SetValid(true)
		if len(args) != sym.GetNumOfNeededArgs() {
			printFunctionHelp(sym, args)
			return nil, errors.New("Invalid command")
		}
		ret := sym.Call()
		if !conf.PrintReturnValuesOff {
			printCallResult(ret)
		}
		return ret, nil
	} else if typ.Kind() == reflect.Struct {

		for i := 0; i < typ.NumMethod(); i++ {
			method := typ.Method(i)
			sym, ok := info[method.Name]
			if !ok {
				msg := fmt.Sprintf("Documentation of %s should be found from AST information.\n", method.Name)
				sym.SetDoc(msg)
			}
			sym.SetKind(Method)
			sym.SetCall(reflect.ValueOf(target).MethodByName(method.Name))
			// TODO(corona10): Fix this, not to check isvalid.
			sym.SetValid(true)
			info[method.Name] = sym
		}

		if len(args) < 2 {
			printMethodHelp(info, args)
			return nil, errors.New("Invalid command")
		}

		// add support for lowcase methodName
		methodTitleName := strings.Title(args[1])
		sym, ok := info[methodTitleName]
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

func printFunctionHelp(info Sym, args []string) {
	var params []string
	file := filepath.Base(args[0])
	for i := 0; i < info.GetNumIns(); i++ {
		param := fmt.Sprintf("<:%s>", info.GetIn(i).String())
		params = append(params, param)
	}
	msg := fmt.Sprintf("Usage of %s:\n%s %s -> %s", file, file, strings.Join(params, " "), info.GetDoc())
	fmt.Printf(msg)
}

func printMethodHelp(info map[string]Sym, args []string) {
	var commands []string
	file := filepath.Base(args[0])
	msg := fmt.Sprintf("Usage of %s:\n", file)
	for key, value := range info {
		if !value.IsValid() {
			continue
		}
		var command []string
		command = append(command, file)
		command = append(command, key)
		for i := 0; i < value.GetNumIns(); i++ {
			param := fmt.Sprintf("<:%s>", value.GetIn(i).String())
			command = append(command, param)
		}
		command = append(command, "->")
		command = append(command, value.GetDoc())
		commands = append(commands, strings.Join(command, " "))
	}

	showCommands := strings.Join(commands, "")
	msg = fmt.Sprintf("%s%s", msg, showCommands)
	fmt.Printf(msg)
}

func getASTInfo(filePath string) (map[string]Sym, error) {
	info := make(map[string]Sym)
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	for _, f := range node.Decls {
		fn, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		}
		var sym Sym
		sym.SetName(fn.Name.Name)
		if fn.Doc.Text() != "" {
			sym.SetDoc(fn.Doc.Text())
		} else {
			sym.SetDoc("No document existed.\n")
		}
		info[fn.Name.Name] = sym
	}
	return info, err
}
