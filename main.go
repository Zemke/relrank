package main

import (
  "fmt"
  "strings"
  "bufio"
  "os"
)

func main() {
  stat, _ := os.Stdin.Stat()
  var G []string
  if (stat.Mode() & os.ModeCharDevice) == 0 {
    reader := bufio.NewReader(os.Stdin)
    inp, _ := reader.ReadString('\n')
    for true {
      text, _ := reader.ReadString('\n')
      if text == "EOF" {
        break
      }
      inp = inp + text
    }
    G = strings.Split(inp[:len(inp)-1], "\n")
  } else {
    if len(os.Args) > 1 {
      file := os.Args[1]
      fmt.Printf("File is %s but it's not yet supported\n", file)
      os.Exit(1)
    }
  }
  for i, g := range G {
    fmt.Println(i, g)
  }
}

