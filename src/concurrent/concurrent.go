package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	// concurrent()
	sequential()
	fmt.Printf("経過：%vms\n", time.Since(now).Milliseconds())
}

func sequential() {
	fmt.Println("Start")
	f("A")
	f("B")
	fmt.Println("End")
}

func concurrent() {
	ch1 := make(chan bool)
	ch2 := make(chan bool)

	fmt.Println("Start")
	go func() {
		f("A")
		ch1 <- true
	}()
	go func() {
		f("B")
		ch2 <- true
	}()

	<-ch1
	<-ch2
	fmt.Println("End")
}

func f(str string) {
	for i := 0; i <= 2; i++ {
		time.Sleep(1 * time.Second)
		fmt.Println(i, str)
	}
}
