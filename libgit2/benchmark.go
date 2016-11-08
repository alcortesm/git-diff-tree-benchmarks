package libgit2

import (
	"fmt"
	"time"

	"github.com/alcortesm/git-diff-tree-benchmarks/result"

	git "gopkg.in/libgit2/git2go.v24"
)

func Benchmark(url string) (*result.Result, error) {
	cloneOptions := &git.CloneOptions{}
	cloneOptions.FetchOptions = &git.FetchOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback:      credentialsCallback,
			CertificateCheckCallback: certificateCheckCallback,
		},
	}

	r, err := git.Clone(url, "libgit2-cloned", cloneOptions)
	if err != nil {
		return nil, err
	}

	ret := &result.Result{
		URL:  url,
		When: time.Now(),
	}

	ret.Data, err = benchmarkAllCommits(r)

	return ret, err
}

func credentialsCallback(
	url string, username string, allowedTypes git.CredType) (git.ErrorCode, *git.Cred) {

	ret, cred := git.NewCredSshKeyFromAgent("git")
	return git.ErrorCode(ret), &cred
}

// Made this one just return 0 during troubleshooting...
func certificateCheckCallback(
	cert *git.Certificate, valid bool, hostname string) git.ErrorCode {

	return 0
}

func benchmarkAllCommits(r *git.Repository) ([]*result.Sample, error) {
	history, err := flatHistory(r)
	if err != nil {
		return nil, err
	}

	ret := make([]*result.Sample, 0, len(history)-1)

	for i := 0; i < len(history)-1; i++ {
		new := history[i]
		old := history[i+1]

		sample, err := benchmarkDiffTree(r, old, new)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot benchmark diff tree between %s and %s: %s",
				old.Id(), new.Id(), err)
		}

		ret = append(ret, sample)
	}

	return ret, nil
}

// returns a flat version of the history of the repository, from the
// head (first element) to the initial commit (last element), using the
// first parent whenever we find a merge.
func flatHistory(r *git.Repository) ([]*git.Commit, error) {
	head, err := r.Head()
	if err != nil {
		return nil, err
	}

	current, err := r.LookupCommit(head.Target())
	if err != nil {
		panic(err)
	}

	ret := make([]*git.Commit, 0)
	for {
		ret = append(ret, current)

		if current.ParentCount() == 0 {
			break
		}

		current = current.Parent(0)
	}

	return ret, nil
}

func benchmarkDiffTree(r *git.Repository, o, n *git.Commit) (*result.Sample, error) {
	ot, err := o.Tree()
	if err != nil {
		return nil, fmt.Errorf("cannot get tree from %s: %s", o.Id(), err)
	}

	nt, err := n.Tree()
	if err != nil {
		return nil, fmt.Errorf("cannot get tree from %s: %s", n.Id(), err)
	}

	opts, err := git.DefaultDiffOptions()
	if err != nil {
		return nil, fmt.Errorf("cannot get default diff options")
	}

	start := time.Now()
	changes, err := r.DiffTreeToTree(ot, nt, &opts)
	elapsed := time.Since(start)

	if err != nil {
		return nil, fmt.Errorf("cannot get changes between %s and %s: %s",
			o.Id(), n.Id(), err)
	}

	stats, err := changes.Stats()
	if err != nil {
		return nil, fmt.Errorf("cannot get changes stats %s and %s: %s",
			o.Id(), n.Id(), err)
	}

	nFiles, err := biggerNumberOfFiles(ot, nt)
	if err != nil {
		return nil, err
	}

	return &result.Sample{
		HashOld:  o.Id().String(),
		HashNew:  n.Id().String(),
		NChanges: stats.FilesChanged(),
		NFiles:   nFiles,
		Duration: elapsed,
	}, nil
}

// returns the number of files in the commit with more files.
// nil commits are allowed and assumed to have 0 files.
func biggerNumberOfFiles(l ...*git.Tree) (int, error) {
	max := 0

	for _, t := range l {
		n := int(t.EntryCount())
		if n > max {
			max = n
		}
	}

	return max, nil
}
