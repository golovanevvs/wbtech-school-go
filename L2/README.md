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
