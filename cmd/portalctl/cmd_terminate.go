package main

import (
	"flag"
	"fmt"
	"os"
)

func cmdTerminate(argv []string) {
	fs := flag.NewFlagSet("terminate", flag.ExitOnError)
	var g gflags; g.add(fs)
	fs.Parse(argv)
	args := fs.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: portalctl terminate <experiment-uuid-or-pid,name>")
		os.Exit(2)
	}

	cl := g.clientOrExit()
	resp, err := cl.TerminateExperiment(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	printJSON(resp)
}
