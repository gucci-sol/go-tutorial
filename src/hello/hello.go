package main

import (
    "fmt"
    "example.com/greetings"
)

func main() {
    message := greetings.Hello("Taiga")
    fmt.Println(message)
}