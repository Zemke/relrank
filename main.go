package main

import (
  "fmt"
  "strings"
  "bufio"
  "os"
  "math"
  "strconv"
  "github.com/shopspring/decimal"
)

type game struct {
  hi int64  // home user id
  ai int64  // away user id
  hs int64  // home user score
  as int64  // away user score
}

type relParam struct {
  u int64
  r decimal.Decimal
  R map[int64]decimal.Decimal
  OPP map[int64]map[int64]int64
  G []game
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
  t := map[int64]int64{}
  for _, g := range G {
    t[g.hi] += g.hs + g.as
    t[g.ai] += g.hs + g.as
  }
  var mx int64 = -1
  for _, s := range t {
    if s > mx {
      mx = s
    }
  }
  steps := -1 + math.Log10(float64(mx) * float64(.13)) * relSteps
  return int(math.Max(math.Min(math.Round(steps), 21), 1))
}

func main() {
  prec, err := strconv.Atoi(getenv("RELRANK_PREC", "50"))
  if err != nil {
    fmt.Println("Precision from RELRANK_PREC is invalid - should be int")
    os.Exit(1)
  }
  decimal.DivisionPrecision = prec
  fmt.Println("precision:", decimal.DivisionPrecision)
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
    gvv := [4]int64{}
    for i, v := range vv  {
      if gvv[i], err = strconv.ParseInt(v.(string), 10, 64); err != nil {
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
  relRel, err := decimal.NewFromString(getenv("RELRANK_RELREL", "20"));
  if err != nil {
    fmt.Println("RELRANK_RELREL is not a number")
  }
  fmt.Println("relRel:", relRel)
  R := map[int64]decimal.Decimal{}
  for _, g := range G {
    R[g.hi] = R[g.hi].Add(decimal.NewFromInt(g.hs))
    R[g.ai] = R[g.ai].Add(decimal.NewFromInt(g.as))
  }
  OPP := map[int64]map[int64]int64{}
  for u, _ := range R {
    OPP[u] = map[int64]int64{}
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
    rels := map[int64]decimal.Decimal{}
    for u, r := range R {
      p := relParam{u: u, r: r, R: R, OPP: OPP, G: G}
      relis := []decimal.Decimal{
        byQuality(p),
        byFarming(p),
        byEffort(p),
      }
      sm := decimal.Sum(relRel, relis...)
      rels[u] = sm.Div(decimal.NewFromInt(int64(len(relis)+1)))
    }
    for u, rel := range rels {
      R[u] = R[u].Mul(rel)
    }
  }
  for u, r := range R {
    fmt.Printf("%d,%s\n", u, r)
  }
}

func byQuality(P relParam) decimal.Decimal {
  return decimal.NewFromInt(1)
}

func byFarming(P relParam) decimal.Decimal {
  return decimal.NewFromInt(1)
}

func byEffort(P relParam) decimal.Decimal {
  return decimal.NewFromInt(1)
}

