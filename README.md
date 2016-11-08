# Introduction

This tool benchmarks how long does it takes to calculate what files has changed
between two git commits using different git implementations.  It generates
reports and comes with a gnuplot script for easier visualization.

The goal is get an idea about the algorithmic complexity of the different
implementations.

The benchmark is performed as follows:

1. It generates a linear version of the history of a repository, going from the
   head up to the initial commit, choosing the first parent whenever it finds
   a merge.  Think of this step as a `git rev-list --first-parent HEAD`.

2. Measure the time to find out what files has changed between each commit and
   its parent.

3. Generate a report with the following format:

  ```
  # Fields (separated by one or more spaces):
  # 1. hash of the old commit
  # 2. hash of the new commit
  # 3. number of files in the commit with more files
  # 4. number of files changed between both commits
  # 5. duration of the diff tree operation in nanoseconds (time to find what
  files were added, deleted or modified)
  #
  # repository URL = git@github.com:src-d/go-git.git
  # date = 2016-11-08 15:59:22.52953147 +0100 CET
  8298b32c5e42204ead8605727ef282621e82f9be 4043678d84215b929ad81cf7a6764d16c9a8fb43       207         1        7600604
  5a72b12c035ef65a3f2335474fe7e66a3c135ace 8298b32c5e42204ead8605727ef282621e82f9be       207        13        5661423
  aba9ec84aab55bd6b22c7c5eac22ca343297061b 5a72b12c035ef65a3f2335474fe7e66a3c135ace       207        37        5883455
  ff6f760e2f8887918e057b398e8c9c5005ffcaae aba9ec84aab55bd6b22c7c5eac22ca343297061b       207         1        5662017
  66c06bbe438c71ce0f854ade4b917a9a51a53488 ff6f760e2f8887918e057b398e8c9c5005ffcaae       207         2        5635016
  ...
  ```

  And store it in a file named after the git implementation used for the
  measurements, in this case `go-git.dat`.

The git implementations supported by this command are:

- go-git, version 4.
- libgit2,  version 0.24.3, using git2go.

The `plot.gp` gnuplot script generates a graph to help visualizing the results.

# Example of use

Download and install:

```bash
; go get git@github.com:alcortesm/git-diff-tree-benchmarks.git
; go install github.com/alcortesm/git-diff-tree-benchmarks
```

Benchmarks using data from (for example) the go-git repository:

```bash
; mkdir /tmp/benchmarks
; cd /tmp/benchmarks
;
; git-diff-tree-benchmarks git@github.com:src-d/go-git.git
; ls
go-git.dat  libgit2.dat
;
; gnuplot $GOPATH/src/alcortesm/git-diff-tree-benchmarks/plot.gp
; file diff-tree-benchmark.png
diff-tree-benchmark.png: PNG image data, 1500 x 900, 8-bit colormap, non-interlaced
```

# Gotchas

This program depends on a working libgit2 installation, version 0.24.3 and
its corresponding git2go (version 24).

When downloading repositories via SSH this program expects an already running
SSH-agent.

# Contact Information

Alberto Cortés <alberto@sourced.tech>

