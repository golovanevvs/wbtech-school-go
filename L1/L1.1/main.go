package main

import "fmt"

type Human struct {
	name     string
	age      int
	heightMM int
}

type Action struct {
	Human
}

func (hm Human) getAge() int {
	return hm.age
}

func (hm Human) getHeight() int {
	return hm.heightMM
}

func (hm Human) introduce() string {
	return fmt.Sprintf("my name is %s", hm.name)
}

func main() {
	hm := Human{
		name:     "Bion",
		age:      19,
		heightMM: 1900,
	}
	ac := Action{
		Human: hm,
	}

	fmt.Printf("Hi, %s\n", ac.introduce())
	fmt.Printf("My age is %d\n", ac.getAge())
	fmt.Printf("My height is %d mm\n", ac.getHeight())
}
