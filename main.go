package main

import "fmt"

func HelloWorld(name string) string {
  return "Hello, " + name
}

func main() {
    fmt.Println(HelloWorld("Flori"))
}

