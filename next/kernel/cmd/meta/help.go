package main

import (
	_ "embed"
	"strings"
)

//go:embed help.txt
var helpText string

// helpFor returns the guidance section for a command ("## <name>" in help.txt), with its usage line first.
func helpFor(name string) (string, bool) {
	marker := "## " + name + "\n"
	i := strings.Index(helpText, marker)
	if i < 0 {
		return "", false
	}
	body := helpText[i+len(marker):]
	if j := strings.Index(body, "\n## "); j >= 0 {
		body = body[:j+1]
	}
	line := ""
	for _, u := range usage {
		if strings.HasPrefix(u, "meta "+name+" ") || u == "meta "+name {
			line = "usage: " + u + "\n\n"
			break
		}
	}
	return line + strings.TrimRight(body, "\n") + "\n" + exitContract, true
}

const exitContract = `
Exit codes: 0 success (result on stdout) | 1 verdict: invalid, gate failed or unchecked, merge refused
(full JSON report on stdout) | 2 could not run: {"error":{"kind","message"}} on stderr, stdout empty.
Error kinds: usage, io, config, precondition, not_a_run, invalid_json, bad_token, unknown_region,
ledger, ledger_corrupt, kernel_fault, exists, internal.
`

func wantsHelp(args []string) bool {
	for _, a := range args {
		if a == "--" {
			return false
		}
		if a == "--help" || a == "-h" {
			return true
		}
	}
	return false
}
