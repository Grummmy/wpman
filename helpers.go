package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var varRegex *regexp.Regexp

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func expandVars(cmd string, env map[string]string, clear bool) string {
	if clear {
		os.Clearenv()
	}

	// add all env variable
	for k, v := range env {
		err := os.Setenv(k, v)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not set '"+k+"' to '"+v+"':", err)
		}
	}

	cmd = varRegex.ReplaceAllString(cmd, "\\$\x00")
	cmd = os.ExpandEnv(cmd)

	return strings.ReplaceAll(cmd, "\x00", "")
}
