package main

import (
	"flag"
	"fmt"
	"os"
)

func cmdConnect(argv []string) {
	fs := flag.NewFlagSet("connect", flag.ExitOnError)
	var g gflags; g.add(fs)
	fs.Parse(argv)
	args := fs.Args()
	if len(args) != 4 {
		fmt.Fprintln(os.Stderr, "usage: portalctl connect <src-exp> <src-lan> <dst-exp> <dst-lan>")
		os.Exit(2)
	}
	cl := g.clientOrExit()
	resp, err := cl.ConnectSharedLan(args[0], args[1], args[2], args[3])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	printJSON(resp)
}
