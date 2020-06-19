# Go chan like rpc

## Examples

Send:

```go
package main

import (
    "fmt"
    "github.com/gocpp/chanrpc"
)

func main() {
	s := NewServer(10)

	s.Register("add", func (a int, b int) int {
                      	return a + b
                      })
	s.Register("print", func (age int, name string) {
                        	fmt.Println(age, name)
                        })

	ch := make(chan int)
	go func() {
		req := <-s.R()
		s.Exec(req)
		ch <- 1
	}()

	s.Send("print", 100, "cinder")
	<-ch
}
```

Call: returns need assert

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
		req := <-s.R()
		s.Exec(req)
		ch <- 1
	}()

	r, err := s.Call("add", 100, 200)
    fmt.Println(r[0].(int), err)
	<-ch
}

```
## Benchmark

![](https://tva1.sinaimg.cn/large/007S8ZIlgy1gfxp2v6bu0j30el0dpjsw.jpg)
