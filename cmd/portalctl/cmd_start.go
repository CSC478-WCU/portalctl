package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type strSlice []string
func (s *strSlice) String() string { return fmt.Sprint([]string(*s)) }
func (s *strSlice) Set(v string) error { *s = append(*s, v); return nil }

func cmdStart(argv []string) {
	fs := flag.NewFlagSet("start", flag.ExitOnError)
	var g gflags; g.add(fs)

	project := fs.String("project", "", "pid[,gid] (required)")
	name    := fs.String("name", "", "experiment EID (required)")

	// Optional pass-throughs
	bindings := fs.String("bindings", "", "JSON object of parameter->value STRINGS")
	bindingsFile := fs.String("bindings-file", "", "Path to JSON file with parameter->value strings")
	refspec   := fs.String("refspec", "", "repo_url:ref (for git-backed profile)")
	aggregate := fs.String("aggregate", "", "aggregate URN override")
	site      := fs.String("site", "", "site binding string (e.g., \"site:1=urn:...\")")
	duration  := fs.String("duration", "", "initial expiration hours")
	startAt   := fs.String("start", "", "schedule start (unix time)")
	stopAt    := fs.String("stop", "", "schedule stop (unix time)")
	sshpub    := fs.String("sshpubkey", "", "SSH public key string")
	noemail   := fs.Bool("noemail", false, "Suppress portal emails")
	nopending := fs.Bool("nopending", false, "No pending flag")

	var paramsKV strSlice
	fs.Var(&paramsKV, "param", "Parameter as name=value (repeatable)")

	fs.Parse(argv)
	args := fs.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: portalctl start [flags] <profile-uuid-or-pid,name>")
		os.Exit(2)
	}
	if *project == "" || *name == "" {
		fmt.Fprintln(os.Stderr, "-project and -name are required")
		os.Exit(2)
	}
	profile := args[0]

	// Build "bindings": precedence bindings-file > bindings > param list
	var bindingsStr string
	if *bindingsFile != "" && *bindings != "" {
		fmt.Fprintln(os.Stderr, "error: use either --bindings-file or --bindings, not both")
		os.Exit(2)
	}
	if *bindingsFile != "" {
		b, err := os.ReadFile(*bindingsFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read --bindings-file: %v\n", err)
			os.Exit(2)
		}
		bindingsStr = string(b)
	} else if *bindings != "" {
		bindingsStr = *bindings
	} else if len(paramsKV) > 0 {
		mp := map[string]string{}
		for _, kv := range paramsKV {
			var k, v string
			n, _ := fmt.Sscanf(kv, "%[^=]=%s", &k, &v)
			if n != 2 || k == "" {
				fmt.Fprintf(os.Stderr, "bad --param %q, expected name=value\n", kv)
				os.Exit(2)
			}
			mp[k] = v
		}
		b, _ := json.Marshal(mp)
		bindingsStr = string(b)
	}

	p := map[string]any{
		"proj":    *project,
		"profile": profile,
		"name":    *name,
	}
	if bindingsStr != "" { p["bindings"] = bindingsStr }
	if *refspec   != "" { p["refspec"]   = *refspec }
	if *aggregate != "" { p["aggregate"] = *aggregate }
	if *site      != "" { p["site"]      = *site }
	if *duration  != "" { p["duration"]  = *duration }
	if *startAt   != "" { p["start"]     = *startAt }
	if *stopAt    != "" { p["stop"]      = *stopAt }
	if *sshpub    != "" { p["sshpubkey"] = *sshpub }
	if *noemail          { p["noemail"]   = 1 }
	if *nopending        { p["nopending"] = 1 }

	cl := g.clientOrExit()
	resp, err := cl.StartExperiment(p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	printJSON(resp)
}
