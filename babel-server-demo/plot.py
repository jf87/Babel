import os
import sys
import time  
from datetime import datetime
import matplotlib
import matplotlib.pyplot as plt
import matplotlib.dates as mdates
import numpy as np
import time


filename = sys.argv[-1]

last_mtime = -1
try:
    while True:
        while True:
            try:
                cur_mtime = os.stat(filename).st_mtime
            except OSError:
                continue
            if cur_mtime == last_mtime:
                continue
            last_mtime = cur_mtime
            break
        print last_mtime
        with open(filename) as f:
            matched = False
            f.readline()
            data = [(int(l.split(",")[0]), int(l.split(",")[1])) for l in f.readlines()]
            data = zip(*data)
            print(
                datetime.fromtimestamp(
                    int("1284101485")
                ).strftime('%Y-%m-%d %H:%M:%S')
            )
            time = [datetime.fromtimestamp(t) for t in data[0]]
            time_mp = matplotlib.dates.date2num(time)
            cnt = data[1]
            for i in np.arange(len(cnt)):
                print cnt[i]
                if cnt[i] == 1:
                    matched = True
            #time = [t for t in data[0]]
            plt.clf()
            if matched:
                plt.plot(time, cnt, linewidth=2, marker='o', color ="g")
            else:
                plt.plot(time, cnt, linewidth=2, marker='o', color ="r")
            plt.gca().xaxis.set_major_formatter(mdates.DateFormatter('%H:%M:%S'))
            plt.gcf().autofmt_xdate()
            #plt.gca().xaxis.set_major_locator(mdates.DayLocator())
            #plt.show(block=False)
            plt.draw()
            plt.xticks(rotation=70)
            plt.xlabel("Time")
            plt.ylabel("Number of Possbile Points")
            plt.title('Babel Point Matching')
            #time.sleep(0.05)
            plt.pause(0.1)
except KeyboardInterrupt:
    plt.close('all')
    print "Done"