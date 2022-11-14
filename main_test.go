package main

import (
  "fmt"
  "sort"
  "testing"
  "github.com/shopspring/decimal"
)

func TestByEffort(t *testing.T) {
  fmt.Println("form within the test")
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
  if rett[1].Cmp(decimal.RequireFromString(".505")) != 0 {
    t.Error(rett[1], "should be", dmx)
  }
}

