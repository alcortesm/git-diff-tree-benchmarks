package libgit2

import (
	"fmt"

	"github.com/alcortesm/git-diff-tree-benchmarks/gogit"
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

	r, err := git.Clone(url, "web", cloneOptions)
	if err != nil {
		return nil, err
	}

	head, err := r.Head()
	if err != nil {
		return nil, err
	}

	commitHead := head.Target()

	fmt.Println(commitHead)

	return gogit.Benchmark(url) // TODO implement this
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
