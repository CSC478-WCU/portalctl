package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}
	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "start":
		cmdStart(args)
	case "status":
		cmdStatus(args)
	case "modify":
		cmdModify(args)
	case "terminate":
		cmdTerminate(args)
	case "extend":
		cmdExtend(args)
	case "manifests":
		cmdManifests(args)
	case "reboot":
		cmdReboot(args)
	case "connect":
		cmdConnect(args)
	case "disconnect":
		cmdDisconnect(args)
	case "-h", "--help", "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
		printUsage()
		os.Exit(2)
	}
}
