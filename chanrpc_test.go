package chanrpc

import (
	"fmt"
	"testing"
)

type util struct {
	num int
}

func NewUtil(n int) *util {
	return &util{
		num: n,
	}
}

func (u *util) add(a int, b int) int {
	//fmt.Println(a + b)
	return a + b
}

func (u *util) print(age int, name string) {
	fmt.Println(u.num, age, name)
}

func TestNewServer(t *testing.T) {
	ut := NewUtil(20)

	s := NewServer(10)
	s.Register("add", ut.add)
	s.Register("print", ut.print)

	ch := make(chan int)
	go func() {
		req := <-s.R()
		s.Exec(req)
		ch <- 1
	}()

	s.Call("print", 100, "cinder")
	<-ch
}

func BenchmarkServer_Exec(b *testing.B) {
	ut := NewUtil(100)

	s := NewServer(1)
	s.Register("add", ut.add)
	s.Register("print", ut.print)

	var args []interface{}
	args = append(args, 100)
	args = append(args, 100)
	//args = append(args, "cinder")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Exec(&CallInfo{f: "add", args: args})
	}
}
