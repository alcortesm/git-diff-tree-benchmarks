set terminal png font Verdana 9 size 1500,900
set output 'diff-tree-benchmark.png'

set key left top box
set xlabel "Number of files" font ",9"
set ylabel "Time (milliseconds)" font ",9"
set xtics mirror
set ytics mirror
set grid ytics xtics

N = floor(system("cat go-git.dat | egrep -v '^#.*$' | wc -l go-git.dat"))
if (N < 2) {
    print "cannot calculate fit for go-git.dat, need at least two data points"
    exit
}

N = floor(system("cat libgit2.dat | egrep -v '^#.*$' | wc -l libgit2.dat"))
if (N < 2) {
    print "cannot calculate fit for libgit2.dat, need at least two data points"
    exit
}

N = floor(system("cat go-git-dev.dat | egrep -v '^#.*$' | wc -l go-git-dev.dat"))
if (N < 2) {
    print "cannot calculate fit for go-git-dev.dat, need at least two data points"
    exit
}

url = system("cat go-git.dat | egrep '^# repository URL = .*' | cut -d'=' -f2")
date = system("cat go-git.dat | egrep '^# date = .*' | cut -d'=' -f2")

f1(x) = a1*x + b1
FIT_LIMIT = 1e-6
fit f1(x) 'go-git.dat' using 3:($5/1000000) via a1, b1

f2(x) = a2*x + b2
fit f2(x) 'libgit2.dat' using 3:($5/1000000) via a2, b2

f3(x) = a3*x + b3
fit f3(x) 'go-git-dev.dat' using 3:($5/1000000) via a3, b3

set title sprintf("Time to calculate what files have changed between two commits\nversus the number of files in the commit with more files.\n\nRepository URL =%s\nDate =%s\n\ngo-git slope = %f (milliseconds/file)\nlibgit2 slope = %f (milliseconds/file)\ngo-git-dev slope = %f (milliseconds/file)\n\n(The linear regression was made using the nonlinear least-squares (NLLS) Marquardt-Levenberg algorithm)\n", url, date, a1, a2, a3) font ",11"

plot \
    'go-git.dat' using ($3):($5/1000000) title "go-git v4" with points lt 1 pt 6 ps 2, \
    f1(x) notitle lt 1 lw 3, \
    'libgit2.dat' using ($3):($5/1000000) title "libgit2 0.24" with points lt 2 pt 6 ps 2, \
    f2(x) notitle lt 2 lw 3, \
    'go-git-dev.dat' using ($3):($5/1000000) title "go-git developement version" with points lt 3 pt 6 ps 2, \
    f3(x) notitle lt 3 lw 3
