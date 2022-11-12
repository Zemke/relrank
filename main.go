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
    t[g.hi] += g.hs + g.as
    t[g.ai] += g.hs + g.as
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
  steps := calcSteps(G)
  fmt.Printf("steps: %d\n", steps)
  relRel, err := strconv.ParseFloat(getenv("RELRANK_RELREL", "20"), 64);
  if err != nil {
    fmt.Println("RELRANK_RELREL is not a number")
  }
  fmt.Println("relRel:", relRel)
  R := map[int]float64{}
  for _, g := range G {
    R[g.hi] += float64(g.hs)
    R[g.ai] += float64(g.as)
  }
  OPP := map[int]map[int]int{}
  for u, _ := range R {
    OPP[u] = map[int]int{}
    for _, g := range G {
      if u == g.hi {
        OPP[u][g.ai] += g.hs
      } else if u == g.ai {
        OPP[u][g.hi] += g.as
      }
    }
  }
  fmt.Println(OPP)
  for i := 1; i <= steps; i++ {
    for u, r := range R {
      relis := [4]float64{
        relRel,
        byQuality(u, r, R, G),
        byFarming(u, r, R, G),
        byEffort(u, r, R, G),
      }
      sm := 0.
      for _, reli := range relis {
        sm += reli
      }
      rel := sm / float64(len(relis))
      R[u] = R[u] * rel
    }
  }
  for u, r := range R {
    fmt.Printf("%d,%f\n", u, r)
  }
}

func byQuality(u int, r float64, R map[int]float64, G []game) float64 {
  return 1.
}

func byFarming(u int, r float64, R map[int]float64, G []game) float64 {
  return 1.
}

func byEffort(u int, r float64, R map[int]float64, G []game) float64 {
  return 1.
}

