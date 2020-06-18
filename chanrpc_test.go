package chanrpc

import (
	"testing"
)

func add(a int, b int) int {
	//fmt.Println(a + b)
	return a + b
}

func print(age int, name string) {
	//fmt.Println(age, name)
}

func TestNewServer(t *testing.T) {
	s := NewServer(10)
	s.Register("add", add)
	s.Register("print", print)

	ch := make(chan int)
	go func() {
		req := <-s.ChanCall
		s.Exec(req)
		ch <- 1
	}()

	s.Call("print", 100, "cinder")
	<-ch
}

func BenchmarkServer_Exec(b *testing.B) {
	s := NewServer(1)
	s.Register("add", add)
	s.Register("print", print)

	var args []interface{}
	args = append(args, 100)
	args = append(args, 100)
	//args = append(args, "cinder")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Exec(&CallInfo{f: "add", args: args})
	}
}
