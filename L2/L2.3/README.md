# L2.3

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

## Ответ

Интерфейс является `nil`, когда и динамический тип, и значение `nil`.
`fmt.Println(err)` - вывод значения `err` (`nil`).
`fmt.Println(err == nil)` - вывод результата сравнения `err` с `nil` (`false`).
