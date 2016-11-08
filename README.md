# Introduction

This tool benchmarks how long does it takes to calculate what files has changed
between two git commits using different git implementations.  It generates
reports and comes with a gnuplot script for easier visualization.

The benchmark is performed as follows:

1. It generates a linear version of the history of a repository, going from the
   head up to the initial commit, choosing the first parent whenever it finds
   a merge.  Think of these step as a `git rev-list --first-parent HEAD`.

2. Measure the time to find out what files has changed between each commit and
   its parent.

3. Generate a report with the following format:

  ```
  # Fields:
  # 1. hash of the old commit
  # 2. hash of the new commit
  # 3. number of files changed between both commits
  # 4. number of files in the commit with more files
  # 5. duration of the diff tree operation in nanoseconds (time to find what
  files were added, deleted or modified
  #
  # repository URL = git@github.com:alcortesm/git-diff-tree-benchmarks.git
  # date = 2016-11-08 10:58:39.362789734 +0100 CET
  0000000000000000000000000000000000000000 c691e37fe800ff54e66b1f77b32f9f627112b91e         3         3           4514
  c691e37fe800ff54e66b1f77b32f9f627112b91e fd9403bf7aa36c9a8b46f0551ccb7f90dc98aca4         5         4          22474
  fd9403bf7aa36c9a8b46f0551ccb7f90dc98aca4 2017ba5cec817db28ff649e64bb7deae45de9528         5         3          49581
  ```

  And store it in a file named after the git implementation used for the
  measurements, in this case `go-git.dat`.

The git implementations supported by this command are:

- go-git, version 4.
- libgit2,  version 0.24.3, using git2go.

The `plot.gp` gnuplot script generates a graph to help visualizing the results.

# Example of use

Download and install

```bash
; go get git@github.com:alcortesm/git-diff-tree-benchmarks.git
; go install github.com/alcortesm/git-diff-tree-benchmarks
```

Benchmarks using data from the go-git repository (for example).

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
its the corresponding git2go (version 24).

When downloading repositories via SSH this programs relies on an already running
SSH-agent.

# Contact Information

Alberto Cortés <alberto@sourced.tech>

