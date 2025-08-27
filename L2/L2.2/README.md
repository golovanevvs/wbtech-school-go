# L2.2

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

## Ответ

### `fmt.Println(test())`

`x` - является именованным возвращаемым значением. `defer` изменяет уже сохранённое возвращаемое значение. Вывод - `2`.

### `fmt.Println(anotherTest())`

`x` - не является именованным возвращаемым значением. `defer` инкрементирует `x` до `2`, но возвращаемое значение уже установлено в `1`.
