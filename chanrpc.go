package chanrpc

import (
	"fmt"
	"reflect"
)

type server struct {
	functions map[string]interface{}
	chanReq   chan *Request
}

type Request struct {
	f        string
	args     []interface{}
	resp     bool
	chanResp chan *Response
}

type Response struct {
	rets []interface{}
	err  error
}

func NewServer(l int) *server {
	return &server{
		functions: make(map[string]interface{}),
		chanReq:   make(chan *Request, l)}
}

func (s *server) R() chan *Request {
	return s.chanReq
}

func (s *server) Register(name string, f interface{}) {
	_, ok := s.functions[name]
	if ok {
		panic("chanrpc Register error")
	}

	s.functions[name] = f
}

func (s *server) Send(f string, args ...interface{}) {
	req := &Request{f: f, args: args, resp: false}
	s.chanReq <- req
}

func (s *server) Call(f string, args ...interface{}) (rets []interface{}, err error) {
	req := &Request{f: f, args: args, resp: true, chanResp: make(chan *Response)}
	s.chanReq <- req
	resp := <-req.chanResp
	rets = resp.rets
	err = resp.err
	close(req.chanResp)
	return
}

func (s *server) Exec(r *Request) (err error) {
	var (
		f             interface{}
		ok            bool
		rType         reflect.Type
		rValue        reflect.Value
		retValues     []reflect.Value
		retInterfaces []interface{}
	)

	f, ok = s.functions[r.f]
	if !ok {
		err = fmt.Errorf("chanrpc Exec error, invalid function: %s", r.f)
		if r.resp {
			r.chanResp <- &Response{rets: retInterfaces, err: err}
		}
		return
	}

	rType = reflect.TypeOf(f)
	rValue = reflect.ValueOf(f)

	in := make([]reflect.Value, rType.NumIn())
	for i := 0; i < rType.NumIn(); i++ {
		in[i] = reflect.ValueOf(r.args[i])
	}
	retValues = rValue.Call(in)

	if r.resp {
		retInterfaces = make([]interface{}, len(retValues))
		for i, rv := range retValues {
			retInterfaces[i] = rv.Interface()
		}
		r.chanResp <- &Response{rets: retInterfaces, err: nil}
	}
	return
}
