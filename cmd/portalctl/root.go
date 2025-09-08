package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/csc478-wcu/portalctl/portal"
)

type gflags struct {
	server, path, cert, key, cacert string
	port                             int
	verify                           bool
	timeout                          int
}

func (g *gflags) add(fs *flag.FlagSet) {
	fs.StringVar(&g.server, "server", "boss.emulab.net", "XML-RPC server")
	fs.IntVar(&g.port, "port", 3069, "XML-RPC port")
	fs.StringVar(&g.path, "path", "/usr/testbed", "Server path")
	fs.StringVar(&g.cert, "cert", "", "Client cert PEM (or combined PEM)")
	fs.StringVar(&g.key, "key", "", "Client key PEM (defaults to -cert)")
	fs.StringVar(&g.cacert, "cacert", "", "CA cert PEM (required if -verify)")
	fs.BoolVar(&g.verify, "verify", false, "Verify server cert using CA")
	fs.IntVar(&g.timeout, "timeout", 900, "HTTP timeout seconds")
}

func (g *gflags) clientOrExit() *portal.Client {
	cl, err := portal.New(portal.Options{
		Server:    g.server,
		Port:      g.port,
		Path:      g.path,
		CertPEM:   g.cert,
		KeyPEM:    g.key,
		CACertPEM: g.cacert,
		Verify:    g.verify,
		Timeout:   time.Duration(g.timeout) * time.Second,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	return cl
}

func printJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}
