package main

import (
	"flag"
	"fmt"
	"os"
)

func cmdManifests(argv []string) {
	fs := flag.NewFlagSet("manifests", flag.ExitOnError)
	var g gflags; g.add(fs)

	fs.Parse(argv)
	args := fs.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: portalctl manifests <experiment-uuid-or-pid,name>")
		os.Exit(2)
	}

	cl := g.clientOrExit()
	resp, err := cl.ExperimentManifests(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	printJSON(resp)
}
