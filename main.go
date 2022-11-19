package main

import (
  "fmt"
  "log"
  "strings"
  "io"
  "bufio"
  "os"
  "math"
  "strconv"
  "sort"
  "github.com/shopspring/decimal"
)

const DEFAULT_PRECISION int = 20

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

type prep struct {
  G []game
  R map[int64]decimal.Decimal
  T total
  OPP map[int64]map[int64]int64
  WT map[int64]int64
  mxWonOpp int64
}

func getenv(env string, def string) string {
  if v, ok := os.LookupEnv(env); ok {
    return v
  }
  return def
}

func setPrecision(prec int) {
  decimal.DivisionPrecision = prec
}

func calcSteps(G []game) int {
  var relSteps, err = strconv.ParseFloat(getenv("RELRANK_RELTEPS", "15.9"), 64)
  dd("relSteps:", relSteps)
  if err != nil {
    log.Fatalf("%f is invalid for RELRANK_RELSTEPS\n", relSteps)
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

func prepare(inp []string) prep {
  var G []game
  for k, l := range inp {
    sp := strings.Split(l, ",")
    if _, err := strconv.ParseInt(sp[2] + sp[3], 10, 64); err != nil && k == 0 {
      dd("skipping csv header")
      continue
    }
    var vv, err = [...]interface{}{sp[0], sp[1], sp[2], sp[3]}, error(nil)
    gvv := [4]int64{}
    for i, v := range vv  {
      if gvv[i], err = strconv.ParseInt(v.(string), 10, 64); err != nil {
        log.Fatalf("%s is not an integer", vv[i])
      }
    }
    G = append(G, game{gvv[0], gvv[1], gvv[2], gvv[3]})
  }

  R := map[int64]decimal.Decimal{}
  for _, g := range G {
    R[g.hi] = R[g.hi].Add(decimal.NewFromInt(g.hs))
    R[g.ai] = R[g.ai].Add(decimal.NewFromInt(g.as))
  }

  T := total{ peru: map[int64]int64{}, }
  OPP := map[int64]map[int64]int64{}
  WT := map[int64]int64{}
  for u := range R {
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

  var mxWonOpp int64 = 0
  for _, oo := range OPP {
    for _, w := range oo {
      if w > mxWonOpp {
        mxWonOpp = w
      }
    }
  }
  return prep{
    G: G,
    R: R,
    T: T,
    OPP: OPP,
    WT: WT,
    mxWonOpp: mxWonOpp,
  }
}

func main() {
  log.SetFlags(0)
  prec, err := strconv.Atoi(getenv("RELRANK_PREC", strconv.Itoa(DEFAULT_PRECISION)))
  if err != nil {
    log.Fatalln("Precision from RELRANK_PREC is invalid - should be int")
  }
  setPrecision(prec)
  dd("precision:", decimal.DivisionPrecision)

  stat, _ := os.Stdin.Stat()
  var inp []string
  if (stat.Mode() & os.ModeCharDevice) == 0 {
    reader := bufio.NewReader(os.Stdin)
    var ll string
    for {
      l, err := reader.ReadString('\n')
      if err == io.EOF {
        break
      }
      ll += l
    }
    inp = strings.Split(ll[:len(ll)-1], "\n")
  } else {
    log.Fatalln("Pass input into stdin")
  }

  prep := prepare(inp)
  dd("G:", prep.G)
  dd("R:", prep.R)
  dd("T.peru:", prep.T.peru)
  dd("T.mn, T.mx:", prep.T.mn, prep.T.mx)
  dd("WT:", prep.WT)
  dd("OPP:", prep.OPP)
  dd("mxWonOpp:", prep.mxWonOpp)

  relRel, err := decimal.NewFromString(getenv("RELRANK_RELREL", "20"));
  if err != nil {
    log.Fatalln("RELRANK_RELREL is not a number")
  }
  dd("relRel:", relRel)
  steps := calcSteps(prep.G)
  dd("steps:", steps)

  R := apply(prep, steps, relRel, prep.R)

  if v, ok := os.LookupEnv("RELRANK_SCALE_MAX"); ok {
    scaleMx, err := decimal.NewFromString(v)
    if err != nil {
      log.Fatalln("RELRANK_SCALE_MAX is not a number")
    }
    if scaleMx.Cmp(decimal.Zero) < 0 {
      log.Fatalln("RELRANK_SCALE_MAX must be positive")
    }
    R = scale(R, scaleMx)
  }

  if v, ok := os.LookupEnv("RELRANK_ROUND"); ok {
    rnd, err := strconv.Atoi(v)
    if err != nil {
      log.Fatalln("RELRANK_ROUND must be an integer")
    }
    R = round(R, int32(rnd))
  }

  for u, r := range R {
    fmt.Printf("%d,%s\n", u, r)
  }
}

func apply(prep prep,
           steps int,
           relRel decimal.Decimal,
           R map[int64]decimal.Decimal) map[int64]decimal.Decimal {
  up, L := distinctPositionsAsc(R)
  rels := map[int64]decimal.Decimal{}
  for u := range prep.R {
    relis := []decimal.Decimal{
      byQuality(prep.OPP[u], prep.WT[u], up, L),
      byFarming(prep.mxWonOpp, prep.WT[u], prep.OPP[u]),
      byEffort(u, prep.T),
    }
    sm := decimal.Sum(relRel, relis...)
    rels[u] = sm.Div(decimal.NewFromInt(int64(len(relis)+1)))
  }
  R2 := map[int64]decimal.Decimal{}
  for u, rel := range rels {
    R2[u] = R[u].Mul(rel)
  }
  if steps > 1 {
    R2 = apply(prep, steps-1, relRel, R2)
  }
  return R2
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
      log.Fatalf("user with rating %s not found", r)
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
  if w == 0 {
    return decimal.Zero
  }
  rel := decimal.Zero
  wt := decimal.NewFromInt(w)
  for u, w1 := range o {
    y := decimal.NewFromInt(up[u]+1).Div(L).Pow(d3)
    rel = rel.Add(y.Mul(decimal.NewFromInt(w1).Div(wt)));
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

func scale(R map[int64]decimal.Decimal,
           scaleMx decimal.Decimal) map[int64]decimal.Decimal {
  mx := decimal.Zero
  for _, r := range R {
    if r.Cmp(mx) > 0 {
      mx = r
    }
  }
  z := decimal.Zero
  R2 := map[int64]decimal.Decimal{}
  for u, r := range R {
    R2[u] = z.Add(r.Sub(z).Mul(scaleMx.Sub(z)).Div(mx.Sub(z)))
  }
  return R2
}

func round(R map[int64]decimal.Decimal,
           rnd int32) map[int64]decimal.Decimal {
  R2 := map[int64]decimal.Decimal{}
  for u, r := range R {
    R2[u] = r.Round(rnd)
  }
  return R2
}

func dd(ss ...any) {
  if debug != "0" {
    log.Println(ss...)
  }
}

