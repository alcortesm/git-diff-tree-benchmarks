package result

import (
	"fmt"
	"io"
	"os"
	"time"
)

// A duration value tell you how much did it took to calculate the
// diff-tree of two commits, their hashes, and the number of files of
// the one with more files.
type Sample struct {
	HashOld  string
	HashNew  string
	NFiles   int
	NChanges int
	Duration time.Duration
}

type Result struct {
	URL  string
	When time.Time
	Data []*Sample
}

func (l Result) Report(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func() {
		errClose := f.Close()
		if err == nil {
			err = errClose
		}
	}()

	l.report(f)

	return nil
}

func (r Result) report(w io.Writer) {
	fmt.Fprintln(w, "# Fields (separated by one or more spaces):")
	fmt.Fprintln(w, "# 1. hash of the old commit")
	fmt.Fprintln(w, "# 2. hash of the new commit")
	fmt.Fprintln(w, "# 3. number of files changed between both commits")
	fmt.Fprintln(w, "# 4. number of files in the commit with more files")
	fmt.Fprintln(w, "# 5. duration of the diff tree operation in nanoseconds (time to find what files were added, deleted or modified)")
	fmt.Fprintln(w, "#")
	fmt.Fprintln(w, "# repository URL =", r.URL)
	fmt.Fprintln(w, "# date =", r.When)

	for _, d := range r.Data {
		fmt.Fprintf(w, "%s %s %9d %9d %14d\n",
			d.HashOld,
			d.HashNew,
			d.NFiles,
			d.NChanges,
			d.Duration.Nanoseconds())
	}
}
