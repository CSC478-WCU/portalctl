package main

import (
	"flag"
	"fmt"
	"os"
)

func cmdExtend(argv []string) {
	fs := flag.NewFlagSet("extend", flag.ExitOnError)
	var g gflags; g.add(fs)

	reason := fs.String("m", "", "reason for the extension (string)")
	fs.Parse(argv)
	args := fs.Args()
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: portalctl extend [-m reason] <experiment> <hours>")
		os.Exit(2)
	}

	cl := g.clientOrExit()
	resp, err := cl.ExtendExperiment(args[0], args[1], *reason)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	printJSON(resp)
}
