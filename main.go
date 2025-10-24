package main

import (
	"fmt"
	"time"

	"github.com/ipoluianov/mic/mic"
)

func main() {
	fmt.Println("Started")
	mic := mic.NewMicapController()
	mic.Start()
	time.Sleep(1000 * time.Second)
	return
}
