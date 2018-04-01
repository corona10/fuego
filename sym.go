package fuego

import (
	"reflect"
	"strconv"
)

type FuegoKind int

const (
	// Unknown is a enum value of the unknown FuegoKind.
	Unknown FuegoKind = iota
	// Func is a enum value of the Func.
	Func
	// Method is a enum value of the Go's method.
	Method
)

// Sym is a symbol of method/functions
type Sym struct {
	name   string
	kind   FuegoKind
	in     []reflect.Value
	params []string
	call   reflect.Value
}

func (s *Sym) Call() []reflect.Value {
	for idx, param := range s.params {
		var t reflect.Type
		t = s.call.Type().In(idx)

		switch t.Kind() {
		case reflect.Int:
			paramValue, _ := strconv.Atoi(param)
			arg := reflect.ValueOf(paramValue)
			s.in = append(s.in, arg)
		case reflect.Float32:
			paramValue, _ := strconv.ParseFloat(param, 32)
			arg := reflect.ValueOf(paramValue)
			s.in = append(s.in, arg)
		case reflect.Float64:
			paramValue, _ := strconv.ParseFloat(param, 64)
			arg := reflect.ValueOf(paramValue)
			s.in = append(s.in, arg)
		case reflect.String:
			s.in = append(s.in, reflect.ValueOf(param))
		}
	}
	return s.call.Call(s.in)
}

func (s Sym) GetKind() FuegoKind {
	return s.kind
}

func (s *Sym) GetIn(idx int) reflect.Type {
	return s.call.Type().In(idx)
}
func (s *Sym) GetNumIns() int {
	return s.call.Type().NumIn()
}

func (s *Sym) SetName(name string) {
	s.name = name
}

func (s *Sym) SetKind(kind FuegoKind) {
	s.kind = kind
}

func (s *Sym) SetParams(params []string) {
	s.params = params
}

func (s *Sym) SetCall(call reflect.Value) {
	s.call = call
}

func (s Sym) GetNumOfNeededArgs() int {
	if s.kind == Method {
		return s.call.Type().NumIn() + 2
	} else if s.kind == Func {
		return s.call.Type().NumIn() + 1
	} else {
		panic("Not supported yet")
	}
}
