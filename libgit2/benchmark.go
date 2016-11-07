package libgit2

import (
	"github.com/alcortesm/git-diff-tree-benchmarks/gogit"
	"github.com/alcortesm/git-diff-tree-benchmarks/result"
)

func Benchmark(url string) (*result.Result, error) {
	return gogit.Benchmark(url) // TODO implement this
}
