#!/usr/bin/env python3

import subprocess as sp

U = {i: n for i,n in [l.split(',') for l in open('csv/users.csv', 'r').read().splitlines()]}

for n in [38,39,40,41]:
  print('season', n)
  R = {u: r for u,r in [l.split(',') for l in open(f'csv/{n}ranking.csv', 'r').read().splitlines()]}
  mx = max([float(r) for _, r in list(R.items())[1:]])
  cmd = f'DEBUG=0 RELRANK_ROUND=2 RELRANK_SCALE_MAX={mx} go run . < csv/{n}games.csv'
  print('cmd', cmd)
  out = sp.run(cmd, shell=True, capture_output=True, text=True)
  print(out.stderr)
  urr = out.stdout.splitlines()
  F = []
  pad = 0
  for l in urr:
    u,r = l.split(',')
    pad = max(pad, len(r))
    F.append((float(r), U[u]))
  F.sort()
  for r,u in F[::-1]:
    print(f"{r:6.2f},{u}")

