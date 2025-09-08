package main

import (
	"fmt"
	"strings"
)

func parseKV(kv string) (string, string, error) {
	i := strings.IndexRune(kv, '=')
	if i <= 0 || i == len(kv)-1 {
		return "", "", fmt.Errorf("bad --param %q, expected name=value", kv)
	}
	return kv[:i], kv[i+1:], nil
}