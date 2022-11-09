#!/usr/bin/env python3 

import matplotlib.pyplot as plt
import math
import sys

plt.rcParams['text.usetex'] = True

xx = range(1, 100)
for a in [50, 30, 20, 5]:
  yy = [max(min(round(-1+math.log10(x*.13)*a),21),1) for x in xx]
  plt.plot(xx, yy, label=f'a={a}')

plt.ylabel('relativization steps S')
plt.xlabel('max number of games played of any user')
plt.legend()
plt.title(r'$S = \max(\min(\left\lfloor -1+\log(x*.13)*a \right\rceil,21),1)$')
if len(sys.argv) >= 2 and 'show' == sys.argv[1]:
  plt.show()
else:
  plt.savefig('images/relsteps.png')

