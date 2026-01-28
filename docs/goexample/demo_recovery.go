package main

import (
	"errors"
	"fmt"
	"reflect"
)

/*
panic本质:
*/

func main() {
	m1()
}

func m1() {
	defer func() {
		if err := recover(); err != nil {
			t := reflect.TypeOf(err).String()
			fmt.Printf("recover a panic, t:%v v:%v\n", t, err)
		}
	}()
	fmt.Printf("m1 before\n")
	m2()
	fmt.Printf("m1 after\n")
}

func m2() {
	fmt.Printf("m2 before\n")
	m3()
	fmt.Printf("m3 after\n")
	if r := recover(); r != nil {
		// 这里根本捕获不到 panic, 需要在defer中捕获
		fmt.Printf("recovered: %v\n", r)
	}
}

func m3() {
	fmt.Printf("has panic\n")
	panic(errors.New("this is a panic error"))
}
