package main

import (
  "fmt"
  "strings"
  "bufio"
  "os"
  "math"
  "strconv"
)

type game struct {
  hi int  // home user id
  ai int  // away user id
  hs int  // home user score
  as int  // away user score
}

func getenv(env string, def string) string {
  if v, ok := os.LookupEnv(env); ok {
    return v
  }
  return def
}


func calcSteps(G []game) int {
  var relSteps, err = strconv.ParseFloat(getenv("RELRANK_RELTEPS", "15.9"), 64)
  fmt.Println("relSteps:", relSteps)
  if err != nil {
    fmt.Printf("%f is invalid for RELRANK_RELSTEPS\n", relSteps)
    os.Exit(1)
  }
  t := map[int]int{}
  for _, g := range G {
    t[g.hi] += g.hs
    t[g.ai] += g.as
  }
  mx := -1
  for _, s := range t {
    if s > mx {
      mx = s
    }
  }
  steps := -1 + math.Log10(float64(mx) * float64(.13)) * relSteps
  return int(math.Max(math.Min(math.Round(steps), 21), 1))
}

func main() {
  stat, _ := os.Stdin.Stat()
  var inp []string
  if (stat.Mode() & os.ModeCharDevice) == 0 {
    reader := bufio.NewReader(os.Stdin)
    ll, _ := reader.ReadString('\n')
    for {
      l, _ := reader.ReadString('\n')
      if l == "EOF" {
        break
      }
      ll += l
    }
    inp = strings.Split(ll[:len(ll)-1], "\n")
  } else {
    if len(os.Args) > 1 {
      file := os.Args[1]
      fmt.Printf("File is %s but it's not yet supported\n", file)
      os.Exit(1)
    }
  }
  var G []game
  for _, l := range inp {
    sp := strings.Split(l, ",")
    var vv, err = [...]interface{}{sp[0], sp[1], sp[2], sp[3]}, error(nil)
    gvv := [4]int{}
    for i, v := range vv  {
      if gvv[i], err = strconv.Atoi(v.(string)); err != nil {
        fmt.Printf("%s is not an integer", vv[i])
        os.Exit(1)
      }
    }
    G = append(G, game{gvv[0], gvv[1], gvv[2], gvv[3]})
  }
  for i, g := range G {
    fmt.Println(i, g)
  }
  fmt.Printf("steps: %d\n", calcSteps(G))
  relRel, err := strconv.Atoi(getenv("RELRANK_RELREL", "20"));
  if err != nil {
    fmt.Println("RELRANK_RELREL is not a number")
  }
  fmt.Println("relRel:", relRel)
}

