package main

import (
	"fmt"
	"time"

	"gopkg.in/tylerb/graceful.v1"
)

func main() {
	graceful.Run(":8888", time.Second, &HTTPProxyHandler{})

	fmt.Println("test conn")
	time.Sleep(time.Second * 10)
	fmt.Println("ok conn")
}
