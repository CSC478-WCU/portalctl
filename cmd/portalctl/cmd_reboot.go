package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func cmdReboot(argv []string) {
	fs := flag.NewFlagSet("reboot", flag.ExitOnError)
	var g gflags; g.add(fs)

	power := fs.Bool("f", false, "power cycle instead of reboot")
	fs.Parse(argv)
	args := fs.Args()
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: portalctl reboot [-f] <experiment> node [node ...]")
		os.Exit(2)
	}

	exp := args[0]
	nodes := args[1:]

	cl := g.clientOrExit()
	resp, err := cl.RebootNodes(exp, nodes, *power)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("requested reboot (power=%v) for nodes: %s\n", *power, strings.Join(nodes, ","))
	printJSON(resp)
}
