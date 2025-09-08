package main

import (
	"flag"
	"fmt"
	"os"
)

func cmdStatus(argv []string) {
	fs := flag.NewFlagSet("status", flag.ExitOnError)
	var g gflags; g.add(fs)

	asjson := fs.Bool("j", true, "server returns JSON payload")
	withcert := fs.Bool("k", false, "include instance cert/key (with -j)")
	refresh := fs.Bool("r", false, "refresh status on server")

	fs.Parse(argv)
	args := fs.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: portalctl status [flags] <experiment-uuid-or-pid,name>")
		os.Exit(2)
	}

	cl := g.clientOrExit()
	resp, err := cl.ExperimentStatus(args[0], *asjson, *withcert, *refresh)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(resp.Output)
}
