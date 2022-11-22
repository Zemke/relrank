package main

import (
  "sort"
  "testing"
  "github.com/shopspring/decimal"
)

func TestByEffort(t *testing.T) {
  T := total{
    peru: map[int64]int64{ 1: 5, 2: 1, 3: 3, },
    mn: decimal.NewFromInt(1),
    mx: decimal.NewFromInt(5),
  }
  rett := []decimal.Decimal{
    byEffort(2, T),
    byEffort(3, T),
    byEffort(1, T),
  }
  isSorted := sort.SliceIsSorted(rett, func (i, j int) bool {
    return rett[i].Cmp(rett[j]) < 0
  })
  if !isSorted {
    t.Error(rett, "is not ascending")
  }
  if rett[0].Cmp(dmn) != 0 {
    t.Error(rett[0], "should be", dmn)
  }
  if rett[len(rett)-1].Cmp(dmx) != 0 {
    t.Error(rett[len(rett)-1], "should be", dmx)
  }
  two := decimal.NewFromInt(2)
  mid := rett[len(rett)-1].Div(two).Add(dmn.Div(two))
  if rett[1].Cmp(mid) != 0 {
    t.Error(rett[1], "should be", mid)
  }
}

func TestByQuality(t *testing.T) {
  tests := []struct {
    o map[int64]int64
    w int64
    up int64
    u int64
  }{
    { u: 1, w: 5, up: 1, o: map[int64]int64{ 2: 3, 4: 2 } },
    { u: 2, w: 7, up: 3, o: map[int64]int64{ 3: 4, 1: 2, 4: 1 } },
    { u: 3, w: 6, up: 2, o: map[int64]int64{ 2: 6 } },
    { u: 4, w: 4, up: 0, o: map[int64]int64{ 1: 1, 2: 3 } },
  }
  up := map[int64]int64{ 1: 1, 2: 3, 3: 2, 4: 0 }
  rrlen := decimal.NewFromInt(int64(len(tests)))
  rett := map[int64]decimal.Decimal{}
  for _, test := range tests {
    rett[test.u] = byQuality(test.o, test.w, up, rrlen)
  }
  exp := map[int]string{
    1: "0.60625",
    2: "0.2790178571428571296875",
    3: "1",
    4: "0.78125",
  }
  for u := 1; u <= len(exp); u++ {
    act := rett[int64(u)]
    if decimal.RequireFromString(exp[u]).Cmp(act) != 0 {
      t.Error("expecting", act, "to be", exp[u])
    }
  }
  if rett[3].Cmp(dmx) != 0 {
    t.Error("expecting", rett[3], "to be", dmx)
  }
}

func TestByFarming(t *testing.T) {
  tests := []struct{
    u int64
    w int64
    oo map[int64]int64
  }{
    { u: 1, w: 3, oo: map[int64]int64{ 1: 2, 3: 1 } },
    { u: 2, w: 7, oo: map[int64]int64{ 1: 4, 3: 3 } },
  }
  var mxWonOpp int64 = 4
  rett := map[int64]decimal.Decimal{}
  for _, test := range tests {
    rett[test.u] = byFarming(mxWonOpp, test.w, test.oo)
  }
  exp := map[int64]string{ 1: "0.835", 2: "0.4929838748980079" }
  for u, ex := range exp {
    if rett[u].Cmp(decimal.RequireFromString(ex)) != 0 {
      t.Error("Expecting", rett[u], "to be", ex)
    }
  }
}

func TestDistinctPositionsAsc(t *testing.T) {
  R := map[int64]decimal.Decimal{
    1: decimal.RequireFromString("1"),
    2: decimal.RequireFromString("2"),
    3: decimal.RequireFromString("3"),
    4: decimal.RequireFromString("3"),
    5: decimal.RequireFromString("10"),
    6: decimal.RequireFromString("10"),
    7: decimal.RequireFromString("9"),
    8: decimal.RequireFromString("8"),
    9: decimal.RequireFromString("11"),
  }
  act, actL := distinctPositionsAsc(R)
  if actL.IntPart() != 7 {
    t.Error("Expecting", actL, "to be", 7)
  }
  exp := map[int64]int64{
    1: 0,
    2: 1,
    3: 2,
    4: 2,
    5: 5,
    6: 5,
    7: 4,
    8: 3,
    9: 6,
  }
  for u, ex := range exp {
    if act[u] != ex {
      t.Error("Expecting", act[u], "to be", ex)
    }
  }
}

func TestCalcSteps(t *testing.T) {
  tests := []struct{
    a int
    exp int
  }{
    {   50, 15 },
    {  150, 21 },
    {   77, 19 },
    { 1000, 21 },
  }
  for _, test := range tests {
    G := []game{}
    for i := 1; i <= test.a; i++ {
      G = append(G, game{ 1, 1 + int64(i) ,1 ,0 })
    }
    if ret := calcSteps(G); ret != test.exp {
      t.Error("Expecting", ret, "to be", test.exp)
    }
  }
}

func TestPrepare(t *testing.T) {
  inp := []string {
    "1,2,3,0",
    "1,2,1,0",
    "3,1,1,3",
    "5,1,2,3",
  }
  ret := prepare(inp)

  // G
  expG := []game {
    game{ hi: 1, ai: 2, hs: 3, as: 0 },
    game{ hi: 1, ai: 2, hs: 1, as: 0 },
    game{ hi: 3, ai: 1, hs: 1, as: 3 },
    game{ hi: 5, ai: 1, hs: 2, as: 3 },
  }
  for i, g := range expG {
    if g != ret.G[i] {
      t.Error("Expected", ret.G[i], "to be", g)
    }
  }

  // R
  expR := map[int64]string{ 1: "10", 2: "0", 3: "1", 5: "2" }
  for u, r := range expR {
    if decimal.RequireFromString(r).Cmp(ret.R[u]) != 0 {
      t.Error("Expected", ret.R[u], "for user", u, "to be", r)
    }
  }

  // T
  if ret.T.mn.Cmp(decimal.NewFromInt(0)) != 0 {
    t.Error("Expected", ret.T.mn, "to be", 0)
  }
  if ret.T.mx.Cmp(decimal.NewFromInt(13)) != 0 {
    t.Error("Expected", ret.T.mx, "to be", 13)
  }
  expTperu := map[int64]int64{ 1: 13, 2: 4, 3: 4, 5: 5 }
  for u, v := range expTperu {
    if v != ret.T.peru[u] {
      t.Error("Expected", ret.T.peru[u], "for user", u, "to be", v)
    }
  }

  // WT
  expWT := map[int64]int64{ 1: 10, 2: 0, 3: 1, 5: 2 }
  for u, v := range expWT {
    if v != ret.WT[u] {
      t.Error("Expected", ret.WT[u], "for user", u, "to be", v)
    }
  }

  // OPP
  expOPP := map[int64]map[int64]int64{
    1: map[int64]int64{ 2: 4, 3: 3, 5: 3 },
    2: map[int64]int64{ 1: 0, 3: 0, 5: 0 },
    3: map[int64]int64{ 1: 1, 2: 0, 5: 0 },
    5: map[int64]int64{ 1: 2, 2: 0, 3: 0 },
  }
  for u, oo := range expOPP {
    if len(ret.OPP[u]) > 4 {
      t.Error("There should only be 4 opponents, is", len(ret.OPP[u]))
    }
    for o, w := range oo {
      if w != ret.OPP[u][o] {
        t.Error("Expected", ret.OPP[u][o], "for user", u, "against", o, "to be", w)
      }
    }
  }

  // mxWonOpp
  if ret.mxWonOpp != 4 {
    t.Error("Expected", ret.mxWonOpp, "to be", 4)
  }
}

func TestScale(t *testing.T) {
  R := map[int64]decimal.Decimal{
    1: decimal.NewFromInt(100),
    2: decimal.NewFromInt(53),
    3: decimal.NewFromInt(11),
  }
  ret := scale(R, decimal.NewFromInt(10))
  exp := map[int64]string{ 1: "10", 2: "5.3", 3: "1.1" }
  for u, ex := range exp {
    if decimal.RequireFromString(ex).Cmp(ret[u]) != 0 {
      t.Error("Expected", ret[u], "to equal", ex)
    }
  }
}

func TestRound(t *testing.T) {
  R := map[int64]decimal.Decimal{
    1: decimal.RequireFromString("1234.1234"),
    2: decimal.RequireFromString("12.35666"),
    3: decimal.RequireFromString("12"),
  }
  ret := round(R, 1)
  exp := map[int64]string{ 1: "1234.1", 2: "12.4", 3: "12" }
  for u, ex := range exp {
    if decimal.RequireFromString(ex).Cmp(ret[u]) != 0 {
      t.Error("Expected", ret[u], "to equal", ex)
    }
  }
}

func TestApply(t *testing.T) {
  decimal.DivisionPrecision = 5
  prep := prep{
    G: []game{game{1, 2, 3, 0}, game{2, 1, 3, 2}, game{3, 1, 2, 1}},
    R: map[int64]decimal.Decimal{
      1: decimal.NewFromInt(6),
      2: decimal.NewFromInt(3),
      3: decimal.NewFromInt(2),
    },
    T: total{
      peru: map[int64]int64{1: 11, 2: 8, 3: 3},
      mn: decimal.NewFromInt(0),
      mx: decimal.NewFromInt(11),
    },
    OPP: map[int64]map[int64]int64{
      1: map[int64]int64{2: 5, 3: 1},
      2: map[int64]int64{1: 3},
      3: map[int64]int64{1: 2},
    },
    WT: map[int64]int64{1: 6, 2: 3, 3: 2},
    mxWonOpp: 5,
  }
  ret := apply(prep, 3, decimal.NewFromInt(20), prep.R)
  expR := map[int64]string{
    1: "966.238766151319158",
    2: "524.216272629898488",
    3: "335.792088726969422",
  }
  for u, r := range expR {
    if decimal.RequireFromString(r).Cmp(ret[u]) != 0 {
      t.Error("Expected", ret[u], "to equal", r, "for user", u)
    }
  }
  decimal.DivisionPrecision = DEFAULT_PRECISION
}

