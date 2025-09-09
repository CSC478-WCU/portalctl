package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func parseKV(kv string) (string, string, error) {
	i := strings.IndexRune(kv, '=')
	if i <= 0 || i == len(kv)-1 {
		return "", "", fmt.Errorf("bad --param %q, expected name=value", kv)
	}
	return kv[:i], kv[i+1:], nil
}

func BuildBindings(
	spec string, specFile string,
	bindings string, bindingsFile string,
	paramsKV []string,
) (string, error) {

	// bindings-file vs bindings
	if bindingsFile != "" && bindings != "" {
		return "", fmt.Errorf("use either --bindings-file or --bindings, not both")
	}
	// spec-file vs spec
	if specFile != "" && spec != "" {
		return "", fmt.Errorf("use either --spec-file or --spec, not both")
	}

	// Base map (parameter name -> string value)
	bmap := map[string]string{}

	// Seed from bindings-file / bindings / params
	if bindingsFile != "" {
		b, err := os.ReadFile(bindingsFile)
		if err != nil {
			return "", fmt.Errorf("read --bindings-file: %v", err)
		}
		if err := json.Unmarshal(b, &bmap); err != nil {
			return "", fmt.Errorf("invalid --bindings-file JSON (expecting object of string values): %v", err)
		}
	} else if bindings != "" {
		if err := json.Unmarshal([]byte(bindings), &bmap); err != nil {
			return "", fmt.Errorf("invalid --bindings JSON (expecting object of string values): %v", err)
		}
	} else if len(paramsKV) > 0 {
		// Build from repeated --param name=value
		for _, kv := range paramsKV {
			k, v, err := parseKV(kv)
			if err != nil {
				return "", err
			}
			bmap[k] = v
		}
	}

	// Inject spec_json (from spec-file/spec) if provided
	if specFile != "" || spec != "" {
		if _, exists := bmap["spec_json"]; exists {
			return "", fmt.Errorf("conflict: spec_json provided both via --spec/--spec-file and --bindings/--param")
		}

		var specText string
		if specFile != "" {
			b, err := os.ReadFile(specFile)
			if err != nil {
				return "", fmt.Errorf("read --spec-file: %v", err)
			}
			specText = string(b)
		} else {
			specText = spec
		}

		// Validate and minify the experiment JSON before embedding as string
		var tmp interface{}
		if err := json.Unmarshal([]byte(specText), &tmp); err != nil {
			return "", fmt.Errorf("spec JSON is invalid: %v", err)
		}
		min, _ := json.Marshal(tmp)
		bmap["spec_json"] = string(min)
	}

	// If no bindings at all, return empty (caller should omit the bindings field)
	if len(bmap) == 0 {
		return "", nil
	}

	out, _ := json.Marshal(bmap)
	return string(out), nil
}