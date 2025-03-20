package main

import (
	"fmt"
	"strconv"
)

func Main2() {
	// ch - channel : type string
	// make is to init channel with type
	// Channel trong Golang là một cơ chế dùng để giao tiếp giữa các goroutine, giúp truyền dữ liệu một cách an toàn.
	ch := make(chan string)

	for i := 0; i < 10; i++ {
		go func(i int) {
			for j := 0; j < 10; j++ {
				ch <- "Goroutine : " + strconv.Itoa(i)
			}
			close(ch)
		}(i)
	}

	for msg := range ch {
		fmt.Println(msg)
	}

}
