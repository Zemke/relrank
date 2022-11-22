#!/usr/bin/env python3

import subprocess as sp
import math
import os
import sys
from collections import Counter

U = {i: n for i,n in [l.split(',') for l in open('csv/users.csv', 'r').read().splitlines()]}

debug = os.getenv('DEBUG', '0')
minw = int(os.getenv('MINW', '5'))
print('minw', minw)
print('env', env := f"DEBUG={debug} RELRANK_ROUND=2")
nn = [int(sys.argv[1])] if len(sys.argv) > 1 else [38,39,40,41]
maes, rmses = 0, 0
for n in nn:
  print('season', n)
  C = Counter()
  for l in open(f"csv/{n}games.csv", 'r').read().splitlines()[1:]:
    hi,ai,hsc,asc = l.split(',')
    C[hi] += int(hsc)
    C[ai] += int(asc)
  uws = [str(u) for u,w in C.items() if w >= minw]
  R = {u: r for u,r in [l.split(',') for l in open(f'csv/{n}ranking.csv', 'r').read().splitlines()]}
  mx = max([float(r) for _, r in list(R.items())[1:]])
  cmd = f'{env} RELRANK_SCALE_MAX={mx} go run . < csv/{n}games.csv'
  if debug == '1':
    print('cmd', cmd)
  out = sp.run(cmd, shell=True, capture_output=True, text=True)
  urr = out.stdout.splitlines()
  if debug == '1':
    print('stderr')
    print(out.stderr)
  if debug == '1':
    print('stdout')
    print(out.stdout)
  F = []
  pad = 0
  for l in urr:
    u,r = l.split(',')
    if u not in R or u not in uws:
      if debug == '1':
        print(f"skipping {u} ({U[u]})")
      continue
    pad = max(pad, len(r))
    F.append((float(r), float(R[u]), U[u]))
  F.sort()
  dd = []
  for rn,ro,u in F[::-1]:
    if debug == '1':
      print(f"{rn:6.2f},{ro:6.2f},{u}")
    dd.append(rn - ro)
  mae = sum([abs(dd[i-1] - dd[i]) for i in range(1, len(dd))]) / len(dd)
  print(n, 'MAE', mae)
  rmse = math.sqrt(sum([(dd[i-1] - dd[i])**2 for i in range(1, len(dd))]) / len(dd))
  print(n, 'RMSE', rmse)
  maes += mae
  rmses += rmse

print('average')
print('MAE', maes / len(nn))
print('RMSE', rmses / len(nn))

