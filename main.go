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
  w map[int64]int64
  G []game
}

type total struct {
  mn decimal.Decimal
  mx decimal.Decimal
  peru map[int64]int64
}

type sortedRanking struct {
  uu []int64
  r decimal.Decimal
  p int64
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
  fmt.Println("R", R)
  T := total{ peru: map[int64]int64{}, }
  OPP := map[int64]map[int64]int64{}
  W := map[int64]int64{}
  for u, _ := range R {
    OPP[u] = map[int64]int64{}
    for _, g := range G {
      if u == g.hi {
        OPP[u][g.ai] += g.hs
        W[u] += g.hs
        T.peru[u] += g.hs + g.as
      } else if u == g.ai {
        OPP[u][g.hi] += g.as
        W[u] += g.as
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
  fmt.Println("W:", W)
  fmt.Println("OPP:", OPP)
  for i := 1; i <= steps; i++ {
    rr := sortRankings(R)
    rels := map[int64]decimal.Decimal{}
    for u, r := range R {
      p := relParam{u: u, r: r, R: R, OPP: OPP, G: G, w: W}
      relis := []decimal.Decimal{
        byQuality(OPP[u], W[u], rr, u),
        byFarming(p),
        byEffort(u, T),
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

func sortRankings(R map[int64]decimal.Decimal) []sortedRanking {
    rankings := []decimal.Decimal{}
    for _, r := range R {
      rankings = append(rankings, r)
    }
    sort.Slice(rankings, func (i, j int) bool {
      return rankings[i].Cmp(rankings[j]) < 0
    })
    srr := []sortedRanking{}
    for i, r := range rankings {
      sr := sortedRanking{ r: r, p: int64(i), uu: []int64{} }
      for u, r1 := range R {
        if r == r1 {
          sr.uu = append(sr.uu, u)
        }
      }
      fmt.Println("sr r p uu", sr.r, sr.p, sr.uu)
      srr = append(srr, sr)
    }
    return srr
}

func byQuality(o map[int64]int64, w int64, rr []sortedRanking, u int64) decimal.Decimal {
  var rru int64 = -1
  for _, r := range rr {
    for _, u1 := range r.uu {
      if u1 == u {
        rru = r.p
        break
      }
    }
    if rru != -1 {
      break
    }
  }
  if rru == -1 {
    fmt.Println("user p not found")
    os.Exit(1)
  }
  rel := decimal.Zero
  wt := decimal.NewFromInt(w)
  L := decimal.NewFromInt(int64(len(rr)))
  for _, w := range o {
    y := decimal.NewFromInt(rru+1).Div(L).Pow(decimal.NewFromInt(3))
    rel = rel.Add(y.Mul(decimal.NewFromInt(w).Div(wt)));
  }
  return decimal.NewFromInt(1)
}

func byFarming(P relParam) decimal.Decimal {
  return decimal.NewFromInt(1)
}

func byEffort(u int64, T total) decimal.Decimal {
  a := decimal.RequireFromString("0.01")
  b := decimal.RequireFromString("1")
  x := decimal.NewFromInt(T.peru[u])
  return a.Add(x.Sub(T.mn).Mul(b.Sub(a)).Div(T.mx.Sub(T.mn)))
}

