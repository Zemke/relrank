#!/usr/bin/env python3 

import matplotlib.pyplot as plt
import math
import sys
import os

plt.rcParams['text.usetex'] = True
if os.getenv('DARK', '0') == '1':
  plt.style.use("dark_background")

def disp(name):
  if os.getenv('SHOW', '0') == '1':
    plt.show()
  else:
    plt.savefig(name)


plt.figure()
xx = range(1, 100+1)
for a in [50, 30, 20, 5]:
  yy = [max(min(round(-1+math.log10(x*.13)*a),21),1) for x in xx]
  plt.plot(xx, yy, label=f'a={a}')
plt.ylabel('relativization steps S')
plt.xlabel('max number of games played of any user')
plt.legend()
plt.title(r'$S = \max(\min(\left\lfloor -1+\log(x*.13)*a \right\rceil,21),1)$')
disp('images/relsteps.png')

plt.figure()
mx=15
xx = range(1, mx+1)
yy = [((x+1)/mx)**3 for x in xx]
plt.plot(xx, yy, label=f"n={mx}")
plt.legend()
plt.ylabel('relativizer')
plt.xlabel('opponent position')
plt.title(f'$((x+1)/n)^3$')
disp('images/quality.png')

plt.figure()
mx=100
xx = range(1, mx+1)
yy = [-math.log(x)+1 for x in xx]
plt.plot(xx, yy)
plt.ylabel('relativizer')
plt.xlabel('number of round against same opponent')
plt.title('y=-ln(x)+1')
disp('images/farming_fundamental.png')

plt.figure()
mx=100
a=mx
xx = range(1, mx+1)
yy = [-(99/(100*math.log(a)))*math.log(x)+1 for x in xx]
plt.plot(xx, yy, label=f'a={a}')
plt.legend()
plt.ylabel('relativizer')
plt.xlabel('number of round against same opponent')
plt.title('y=-(99/(100ln(a)))ln(x)+1')
disp('images/farming.png')

