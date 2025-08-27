# L2.4

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

## Ответ

Вывод от 0 до 9, затем deadlock, потому что нет закрытия канала `ch` писателем.
