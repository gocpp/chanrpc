package chanrpc

import (
	"fmt"
	"reflect"
	"sync"
)

type server struct {
	m         sync.RWMutex
	functions map[string]interface{}
	ChanCall  chan *CallInfo
}

type CallInfo struct {
	f    string
	args []interface{}
}

func NewServer(l int) *server {
	return &server{
		functions: make(map[string]interface{}),
		ChanCall:  make(chan *CallInfo, l)}
}

func (s *server) Register(name string, f interface{}) {
	s.m.Lock()
	defer s.m.Unlock()

	_, ok := s.functions[name]
	if ok {
		panic("chanrpc Register error")
	}

	s.functions[name] = f
}

func (s *server) Send(f string, args ...interface{}) {
	req := &CallInfo{f: f, args: args}
	s.ChanCall <- req
}

func (s *server) Call(f string, args ...interface{}) {
	req := &CallInfo{f: f, args: args}
	s.ChanCall <- req
}

func (s *server) Exec(r *CallInfo) error {
	var (
		f      interface{}
		ok     bool
		rType  reflect.Type
		rValue reflect.Value
	)

	f, ok = s.functions[r.f]
	if !ok {
		return fmt.Errorf("chanrpc Exec error, invalid function: %s", r.f)
	}

	rType = reflect.TypeOf(f)
	rValue = reflect.ValueOf(f)
	in := make([]reflect.Value, rType.NumIn())
	for i := 0; i < rType.NumIn(); i++ {
		in[i] = reflect.ValueOf(r.args[i])
	}

	rValue.Call(in)
	return nil
}
