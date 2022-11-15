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
