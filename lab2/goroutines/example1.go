package main

import (
	"fmt"
)

func Main1() {
	ch := make(chan string)

	go func() {
		ch <- "Hello golang"
	}()

	msg := <-ch
	fmt.Println(msg)
}
