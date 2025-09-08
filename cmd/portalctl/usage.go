package main

import "fmt"

func printUsage() {
	fmt.Print(`portalctl <command> [global flags] [command flags]

Global flags:
  -server string   XML-RPC server (default "boss.emulab.net")
  -port int        XML-RPC port (default 3069)
  -path string     Server path (default "/usr/testbed")
  -cert string     Client cert PEM (or combined PEM)
  -key string      Client key PEM (defaults to -cert)
  -cacert string   CA cert PEM (required if -verify)
  -verify          Verify server cert
  -timeout int     HTTP timeout seconds (default 900)

Commands:
  start       -project <pid[,gid]> -name <eid> [--bindings JSON] [--param k=v]... [--refspec url:ref] [--aggregate ...] [--site ...] [--duration H] [--start T] [--stop T] [--sshpubkey STR] <profile>
  status      [-j] [-k] [-r] <experiment>
  modify      [--bindings JSON] [--param k=v]... <experiment>
  terminate   <experiment>
  extend      [-m reason] <experiment> <hours>
  manifests   <experiment>
  reboot      [-f] <experiment> node [node ...]
  connect     <src-exp> <src-lan> <dst-exp> <dst-lan>
  disconnect  <experiment> <src-lan>

Notes:
- "bindings" must be a JSON object mapping parameter names to STRING values.
- Use --param to build that JSON quickly. For JSON-typed params, pass a stringified JSON with proper quoting.
`)
}
