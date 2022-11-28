#!/usr/bin/env python3

import subprocess as sp
import math
import os
import sys
from collections import Counter
import itertools
from tqdm.auto import tqdm
import time

class MinMax:
  def __init__(self, a, b, mn, mx):
    self.a = a
    self.b = b
    self.mn = mn
    self.mx = mx
  def __call__(self, x):
    numer = (x-self.mn)*(self.b-self.a)
    denom = self.mx-self.mn
    return self.a + (numer/denom)

U = {i: n for i,n in [l.split(',') for l in open('csv/users.csv', 'r').read().splitlines()]}

permenv = os.getenv('PERM', '0')
if permenv == '1':
  perms = list(itertools.product(
    [15.9+(ch/10) for ch in range(-10, 11, 1)],
    [20+ch for ch in range(-10, 11, 3)],
  ))
elif ',' in permenv:
  perms = [tuple([float(n) for n in permenv.split(',')])]
else:
  perms = [(None, None)]

for relRel,relSteps in tqdm(perms):
  debug = os.getenv('DEBUG', '0')
  minw = int(os.getenv('MINW', '5'))
  print('minw', minw)
  env = f"DEBUG={debug} RELRANK_ROUND=2"
  if relRel is not None:
    env += f" RELRANK_RELREL={relRel:.2f}"
  if relSteps is not None:
    env += f" RELRANK_RELSTEPS={relSteps}"
  print('env', env)
  nn = [int(sys.argv[1])] if len(sys.argv) > 1 else [39,40,41,42]
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
    cmd = f'{env} go run . < csv/{n}games.csv'
    if debug == '1':
      print('cmd', cmd)
    out = sp.run(cmd, shell=True, capture_output=True, text=True)
    if debug == '1':
      print('stderr')
      print(out.stderr)
    if debug == '1':
      print('stdout')
      print(out.stdout)
    if out.returncode != 0:
      print("There's an error in the gateway script", file=sys.stderr)
      print(out.stderr, file=sys.stderr)
      exit(out.returncode)
    urr = out.stdout.splitlines()
    F = []
    pad = 0
    a, b = +math.inf, -math.inf
    mn, mx = +math.inf, -math.inf
    a = max([float(r) for _, r in list(R.items())[1:]])
    mn = +math.inf
    for l in urr:
      u,r = l.split(',')
      if u not in R or u not in uws:
        if debug == '1':
          print(f"skipping {u} ({U[u]})")
        continue
      pad = max(pad, len(r))
      F.append((float(r), float(R[u]), U[u]))
      a = min(a, float(R[u]))
      b = max(b, float(R[u]))
      mn = min(mn, float(r))
      mx = max(mx, float(r))
    minmax = MinMax(a, b, mn, mx)
    F.sort()
    dd = []
    for rn,ro,u in F[::-1]:
      rn_norm = minmax(rn)
      if debug == '1':
        print(f"{rn_norm:6.2f},{ro:6.2f},{u}")
      dd.append(rn_norm - ro)
    mae = sum([abs(dd[i-1] - dd[i]) for i in range(1, len(dd))]) / len(dd)
    print(n, 'MAE', mae)
    rmse = math.sqrt(sum([(dd[i-1] - dd[i])**2 for i in range(1, len(dd))]) / len(dd))
    print(n, 'RMSE', rmse)
    maes += mae
    rmses += rmse

  print('average')
  print('MAE', maes / len(nn))
  print('RMSE', rmses / len(nn))

