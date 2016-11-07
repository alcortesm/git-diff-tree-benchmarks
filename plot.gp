set terminal png font Verdana 9 size 1500,900
set output 'diff-tree-benchmark.png'

set key left top box
set xlabel "Number of files" font ",9"
set ylabel "Time (seconds)" font ",9"
set xtics mirror
set ytics mirror
set grid ytics xtics


f(x) = a*x + b
FIT_LIMIT = 1e-6
fit f(x) 'data.dat' using 3:($5/1000000) via a, b

set title "Time to calculate what files have changed between two commits\nversus the number of files in the commit with more files" font ",11"
plot \
    'data.dat' using ($3):($5/1000000) title "go-git v4" with points lt 1 pt 6 ps 2, \
    f(x) title "go-git v4 Marquardt-Levenberg fit" lt 1 lw 3
