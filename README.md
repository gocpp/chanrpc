# Go chan like rpc

## Examples

```go
package main

import (
    "fmt"
    "github.com/gocpp/chanrpc"
)


func add(a int, b int) int {
	//fmt.Println(a + b)
	return a + b
}

func print(age int, name string) {
	fmt.Println(u.num, age, name)
}

func main() {
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
```

```go
package main

import (
    "fmt"
    "github.com/gocpp/chanrpc"
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

func main() {
	ut := NewUtil(20)

	s := NewServer(10)
	s.Register("add", ut.add)
	s.Register("print", ut.print)

	ch := make(chan int)
	go func() {
		req := <-s.ChanCall
		s.Exec(req)
		ch <- 1
	}()

	s.Call("print", 100, "cinder")
	<-ch
}

```
## Benchmark

![](https://tva1.sinaimg.cn/large/007S8ZIlgy1gfxeautx5oj30el0pfdiq.jpg)

