package gogit

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/alcortesm/git-diff-tree-benchmarks/result"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/core"
)

func Benchmark(url string) (*result.Result, error) {
	r, err := downloadRepository(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ret := &result.Result{
		URL:  url,
		When: time.Now(),
	}

	ret.Data, err = benchmarkAllCommits(r)

	return ret, err
}

func downloadRepository(url string) (*git.Repository, error) {
	r := git.NewMemoryRepository()

	o := &git.CloneOptions{
		URL: url,
	}

	if err := r.Clone(o); err != nil {
		return nil, fmt.Errorf("cloning %q: %s", url, err)
	}

	return r, nil
}

func benchmarkAllCommits(r *git.Repository) ([]*result.Duration, error) {
	commits, err := flatHistory(r)
	if err != nil {
		return nil, err
	}

	ret := make([]*result.Duration, 0, len(commits)-1)

	var parent *git.Commit
	for _, c := range commits {
		d, err := benchmarkDiffTree(parent, c)
		if err != nil {
			if parent == nil {
				return nil, fmt.Errorf(
					"cannot benchmark diff tree between the empty repository and %s: %s",
					c.Hash, err)
			}
			return nil, fmt.Errorf(
				"cannot benchmark diff tree between %s and %s: %s",
				parent.Hash, c.Hash, err)
		}

		ret = append(ret, d)
		parent = c
	}

	return ret, nil
}

// Returns a flat version of the commit history of a repository: the
// first commit will be the initial commit, then the second commit, up
// until the head.  When it founds a merge, it chooses the first parent.
func flatHistory(r *git.Repository) ([]*git.Commit, error) {
	// the history is built from the head to the initial commit, going
	// from parent to parent, we will reverse the history at the end.
	headReference, err := r.Head()
	if err != nil {
		return nil, fmt.Errorf("cannot get the head: %s", err)
	}

	current, err := r.Commit(headReference.Hash())
	if err != nil {
		return nil, fmt.Errorf("cannot get the head commit: %s", err)
	}

	var ret []*git.Commit
	var found bool
	for {
		ret = append(ret, current)
		if current, found = getFirstParent(current); !found {
			break
		}
	}

	// reverse the history
	for i := 0; i < (len(ret) / 2); i++ {
		ret[i], ret[len(ret)-i-1] = ret[len(ret)-i-1], ret[i]
	}

	return ret, nil
}

func getFirstParent(c *git.Commit) (*git.Commit, bool) {
	if c.NumParents() == 0 {
		return nil, false
	}

	iter := c.Parents()
	defer iter.Close()

	p, err := iter.Next()
	if err != nil {
		return nil, false
	}

	return p, true
}

// measures the time to compare the trees of two commits.
// if the old commit is nil, the new is compared against an empty tree.
func benchmarkDiffTree(o, n *git.Commit) (*result.Duration, error) {
	var ot *git.Tree
	var err error
	if o != nil {
		ot, err = o.Tree()
		if err != nil {
			return nil, fmt.Errorf("cannot get tree from %s: %s", o.Hash, err)
		}
	}

	nt, err := n.Tree()
	if err != nil {
		return nil, fmt.Errorf("cannot get tree from %s: %s", n.Hash, err)
	}

	start := time.Now()
	changes, err := git.DiffTree(ot, nt)
	elapsed := time.Since(start)

	if err != nil {
		if o == nil {
			return nil, fmt.Errorf("cannot get changes between the empty repository and %s: %s",
				n.Hash, err)
		}
		return nil, fmt.Errorf("cannot get changes between %s and %s: %s",
			o.Hash, n.Hash, err)
	}

	hashOld := core.ZeroHash
	if o != nil {
		hashOld = o.Hash
	}

	nFiles, err := biggerNumberOfFiles(o, n)
	if err != nil {
		return nil, err
	}

	return &result.Duration{
		HashOld:  hashOld,
		HashNew:  n.Hash,
		NChanges: len(changes),
		NFiles:   nFiles,
		Duration: elapsed,
	}, nil
}

// returns the number of files in the commit with more files.
// nil commits are allowed and assumed to have 0 files.
func biggerNumberOfFiles(l ...*git.Commit) (int, error) {
	max := 0

	for _, c := range l {
		if c == nil {
			continue
		}

		n, err := numberOfFiles(c)
		if err != nil {
			return 0, fmt.Errorf("cannot get number of files: %s", err)
		}

		if n > max {
			max = n
		}
	}

	return max, nil
}

func numberOfFiles(c *git.Commit) (int, error) {
	iter, err := c.Files()
	if err != nil {
		return 0, fmt.Errorf("cannot get files: %s", err)
	}
	defer iter.Close()

	sum := 0
	for {
		_, err := iter.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, fmt.Errorf("counting files: %s", err)
		}
		sum++
	}

	return sum, nil
}
