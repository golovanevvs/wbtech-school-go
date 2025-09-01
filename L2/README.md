# L2

## [L2.1](L2.1)

Что выведет программа?

Объясните вывод программы.

```go
package main

import "fmt"

func main() {
  a := [5]int{76, 77, 78, 79, 80}
  var b []int = a[1:4]
  fmt.Println(b)
}
```

## [L2.2](L2.2)

Что выведет программа?

Объяснить порядок выполнения defer функций и итоговый вывод.

```go
package main

import "fmt"

func test() (x int) {
  defer func() {
    x++
  }()
  x = 1
  return
}

func anotherTest() int {
  var x int
  defer func() {
    x++
  }()
  x = 1
  return x
}

func main() {
  fmt.Println(test())
  fmt.Println(anotherTest())
}
```

## [L2.3](L2.3)

Что выведет программа?

Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.

```go
package main

import (
  "fmt"
  "os"
)

func Foo() error {
  var err *os.PathError = nil
  return err
}

func main() {
  err := Foo()
  fmt.Println(err)
  fmt.Println(err == nil)
}
```

## [L2.4](L2.4)

Что выведет программа?

Объяснить вывод программы.

```go
func main() {
  ch := make(chan int)
  go func() {
    for i := 0; i < 10; i++ {
    ch <- i
  }
}()
  for n := range ch {
    println(n)
  }
}
```

## [L2.5](L2.5)

Что выведет программа?

Объяснить вывод программы.

```go
package main

type customError struct {
  msg string
}

func (e *customError) Error() string {
  return e.msg
}

func test() *customError {
  // ... do something
  return nil
}

func main() {
  var err error
  err = test()
  if err != nil {
    println("error")
    return
  }
  println("ok")
}
```

## [L2.6](L2.6)

Что выведет программа?

Объяснить поведение срезов при передаче их в функцию.

```go
package main

import (
  "fmt"
)

func main() {
  var s = []string{"1", "2", "3"}
  modifySlice(s)
  fmt.Println(s)
}

func modifySlice(i []string) {
  i[0] = "3"
  i = append(i, "4")
  i[1] = "5"
  i = append(i, "6")
}
```

## [L2.7](L2.7)

Что выведет программа?

Объяснить работу конвейера с использованием select.

```go
package main

import (
  "fmt"
  "math/rand"
  "time"
)

func asChan(vs ...int) <-chan int {
  c := make(chan int)
  go func() {
    for _, v := range vs {
      c <- v
      time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
    }
  close(c)
}()
  return c
}

func merge(a, b <-chan int) <-chan int {
  c := make(chan int)
  go func() {
    for {
      select {
        case v, ok := <-a:
          if ok {
            c <- v
          } else {
            a = nil
          }
        case v, ok := <-b:
          if ok {
            c <- v
          } else {
            b = nil
          }
        }
        if a == nil && b == nil {
          close(c)
          return
        }
     }
   }()
  return c
}

  func main() {
    rand.Seed(time.Now().Unix())
    a := asChan(1, 3, 5, 7)
    b := asChan(2, 4, 6, 8)
    c := merge(a, b)
    for v := range c {
    fmt.Print(v)
  }
```
