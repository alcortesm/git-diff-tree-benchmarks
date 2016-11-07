package main

import (
	"fmt"
	"io"
	"os"

	"github.com/alcortesm/git-diff-tree-benchmarks/gogit"
)

func main() {
	url, help, err := parseArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		usage(os.Stderr)
		os.Exit(1)
	}

	if help {
		usage(os.Stdout)
		os.Exit(0)
	}

	result, err := gogit.Benchmark(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := result.Report("go-git.dat"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func parseArgs() (string, bool, error) {
	if len(os.Args) != 2 {
		return "", false, fmt.Errorf("bad number of arguments")
	}

	if os.Args[1] == "help" || os.Args[1] == "--help" {
		return "", true, nil
	}

	return os.Args[1], false, nil
}

func usage(w io.Writer) {
	fmt.Fprintln(w, "usage:")
	fmt.Fprintf(w, "\t%s <git repository url> : to benchmark a git repository\n", os.Args[0])
	fmt.Fprintln(w, "or")
	fmt.Fprintf(w, "\t%s help : to get this help message\n", os.Args[0])
	fmt.Fprintln(w, "or")
	fmt.Fprintf(w, "\t%s --help : to get this help message\n", os.Args[0])
}
