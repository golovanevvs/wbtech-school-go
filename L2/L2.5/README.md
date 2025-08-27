# L2.5

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

## Ответ

Выведет `error`. Интерфейс является `nil`, когда и динамический тип, и значение `nil`.
