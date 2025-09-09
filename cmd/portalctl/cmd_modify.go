package main

import (
	"flag"
	"fmt"
	"os"
)

func cmdModify(argv []string) {
	fs := flag.NewFlagSet("modify", flag.ExitOnError)
	var g gflags; g.add(fs)

	bindings := fs.String("bindings", "", "JSON object of parameter->value strings")
	bindingsFile := fs.String("bindings-file", "", "Path to JSON file with parameter->value strings")
	spec     := fs.String("spec", "", "Experiment JSON for parameter 'spec_json'")
	specFile := fs.String("spec-file", "", "Path to experiment JSON file for parameter 'spec_json'")

	var paramsKV strSlice
	fs.Var(&paramsKV, "param", "Parameter as name=value (repeatable)")

	fs.Parse(argv)
	args := fs.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: portalctl modify [flags] <experiment-uuid-or-pid,name>")
		os.Exit(2)
	}

	if (*bindings != "" && *bindingsFile != "") ||
		((*bindings != "" || *bindingsFile != "") && len(paramsKV) > 0) {
		fmt.Fprintln(os.Stderr, "error: choose one of --bindings-file, --bindings, or --param ...")
		os.Exit(2)
	}

	bindingsStr, err := BuildBindings(*spec, *specFile, *bindings, *bindingsFile, paramsKV)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}

	p := map[string]any{"experiment": args[0]}
	if bindingsStr != "" {
		p["bindings"] = bindingsStr
	}

	cl := g.clientOrExit()
	resp, err := cl.ModifyExperiment(p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	printJSON(resp)
}
