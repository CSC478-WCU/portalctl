package main

import (
	"flag"
	"fmt"
	"os"
)

func cmdDisconnect(argv []string) {
	fs := flag.NewFlagSet("disconnect", flag.ExitOnError)
	var g gflags; g.add(fs)
	fs.Parse(argv)
	args := fs.Args()
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: portalctl disconnect <experiment> <src-lan>")
		os.Exit(2)
	}
	cl := g.clientOrExit()
	resp, err := cl.DisconnectSharedLan(args[0], args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	printJSON(resp)
}
