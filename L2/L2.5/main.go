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
	//lint:ignore S1021 this is expected behavior in this example
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
