package main

import (
  "fmt"
  "strings"
  "bufio"
  "os"
  "math"
  "strconv"
  "sort"
  "github.com/shopspring/decimal"
)

var d99 = decimal.NewFromInt(99)
var d100 = decimal.NewFromInt(100)
var dmn = decimal.RequireFromString("0.01")
var d1 = decimal.NewFromInt(1)
var dn1 = decimal.NewFromInt(-1)
var d3 = decimal.NewFromInt(3)
var dmx = d1
var debug = getenv("DEBUG", "0")

type game struct {
  hi int64  // home user id
  ai int64  // away user id
  hs int64  // home user score
  as int64  // away user score
}

// total rounds played
type total struct {
  mn decimal.Decimal
  mx decimal.Decimal
  peru map[int64]int64
}

func getenv(env string, def string) string {
  if v, ok := os.LookupEnv(env); ok {
    return v
  }
  return def
}

func calcSteps(G []game) int {
  var relSteps, err = strconv.ParseFloat(getenv("RELRANK_RELTEPS", "15.9"), 64)
  dd("relSteps:", relSteps)
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
  dd("precision:", decimal.DivisionPrecision)
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
    dd(i, g)
  }
  relRel, err := decimal.NewFromString(getenv("RELRANK_RELREL", "20"));
  if err != nil {
    fmt.Println("RELRANK_RELREL is not a number")
  }
  dd("relRel:", relRel)
  R := map[int64]decimal.Decimal{}
  for _, g := range G {
    R[g.hi] = R[g.hi].Add(decimal.NewFromInt(g.hs))
    R[g.ai] = R[g.ai].Add(decimal.NewFromInt(g.as))
  }
  dd("R", R)
  T := total{ peru: map[int64]int64{}, }
  OPP := map[int64]map[int64]int64{}
  WT := map[int64]int64{}
  for u, _ := range R {
    OPP[u] = map[int64]int64{}
    for _, g := range G {
      if u == g.hi {
        OPP[u][g.ai] += g.hs
        WT[u] += g.hs
        T.peru[u] += g.hs + g.as
      } else if u == g.ai {
        OPP[u][g.hi] += g.as
        WT[u] += g.as
        T.peru[u] += g.hs + g.as
      }
    }
  }
  var mn int64 = 0
  var mx int64 = 0
  for _, w := range T.peru {
    if mn > w {
      mn = w
    }
    if mx < w {
      mx = w
    }
  }
  T.mn = decimal.NewFromInt(mn)
  T.mx = decimal.NewFromInt(mx)
  dd("T", T)
  dd("WT:", WT)
  dd("OPP:", OPP)
  var mxWonOpp int64 = 0
  for _, oo := range OPP {
    for _, w := range oo {
      if w > mxWonOpp {
        mxWonOpp = w
      }
    }
  }
  steps := calcSteps(G)
  dd("steps:", steps)
  for i := 1; i <= steps; i++ {
    up, L := distinctPositionsAsc(R)
    rels := map[int64]decimal.Decimal{}
    for u, _ := range R {
      relis := []decimal.Decimal{
        byQuality(OPP[u], WT[u], up, L),
        byFarming(mxWonOpp, T.peru[u], OPP[u]),
        byEffort(u, T),
      }
      sm := decimal.Sum(relRel, relis...)
      rels[u] = sm.Div(decimal.NewFromInt(int64(len(relis)+1)))
    }
    for u, rel := range rels {
      R[u] = R[u].Mul(rel)
    }
  }
  dd("output")
  for u, r := range R {
    fmt.Printf("%d,%s\n", u, r)
  }
}

func distinctPositionsAsc(R map[int64]decimal.Decimal) (map[int64]int64, decimal.Decimal) {
  rr := []decimal.Decimal{}
  for _, r := range R {
    rr = append(rr, r)
  }
  sort.Slice(rr, func (i, j int) bool {
    return rr[i].Cmp(rr[j]) < 0
  })
  var uniq int64 = 0
  var p int64 = 0
  upAsc := map[int64]int64{}
  done := map[int64]bool{}
  for i, r := range rr {
    var user int64
    for u, r1 := range R {
      if r1.Cmp(r) == 0 && !done[u] {
        user = u
        break
      }
    }
    if user == 0 {
      fmt.Printf("user with rating %s not found", r)
      os.Exit(1)
    }
    done[user] = true
    if i == 0 {
      upAsc[user] = p
      p += 1
      uniq += 1
    } else {
      if rr[i-1].Cmp(rr[i]) == 0 {
        upAsc[user] = p-1
      } else {
        upAsc[user] = p
        p += 1
        uniq += 1
      }
    }
  }
  return upAsc, decimal.NewFromInt(uniq)
}

func byQuality(o map[int64]int64, w int64, up map[int64]int64, L decimal.Decimal) decimal.Decimal {
  rel := decimal.Zero
  wt := decimal.NewFromInt(w)
  for u, w := range o {
    y := decimal.NewFromInt(up[u]+1).Div(L).Pow(d3)
    rel = rel.Add(y.Mul(decimal.NewFromInt(w).Div(wt)));
  }
  return rel
}

func byFarming(mxWonOpp int64, uw int64, oo map[int64]int64) decimal.Decimal {
  P := decimal.Zero
  if uw == 0 || mxWonOpp == 0 {
    return d1
  }
  for _, w := range oo {
    if w == 0 {
      continue
    }
    for i := 1; int64(i) <= w; i++ {
      sub := dn1.Mul(
          d99.Div(d100.Mul(decimal.NewFromFloat(math.Log(float64(mxWonOpp))))),
        ).Mul(decimal.NewFromFloat(math.Log(float64(i)))).Add(d1)
      P = P.Add(sub)
    }
  }
  return P.Div(decimal.NewFromInt(uw))
}

func byEffort(u int64, T total) decimal.Decimal {
  a, b, x := dmn, dmx, decimal.NewFromInt(T.peru[u])
  return a.Add(x.Sub(T.mn).Mul(b.Sub(a)).Div(T.mx.Sub(T.mn)))
}

func dd(ss ...any) {
  if debug != "0" {
    fmt.Println(ss...)
  }
}

func ddf(s string, ss ...any) {
  if debug != "0" {
    fmt.Printf(s, ss...)
  }
}

