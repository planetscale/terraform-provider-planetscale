// Command patchsensitive reapplies credential redaction to Speakeasy-generated
// files. Speakeasy overwrites those files on `make generate`; this program
// puts the redaction calls back.
package main

import (
	"fmt"
	"go/format"
	"os"
	"strings"
)

type replacement struct {
	path string
	old  string
	new  string
}

func main() {
	replacements := []replacement{
		{
			path: "internal/provider/utils.go",
			old:  "\tfields[FieldHttpRequestBody] = bodyFromRestOfRequestReader(reqReader)\n",
			new:  "\t// Credential fields are redacted by script/patchsensitive after generation.\n\tfields[FieldHttpRequestBody] = redactSensitiveHTTP(bodyFromRestOfRequestReader(reqReader))\n",
		},
		{
			path: "internal/provider/utils.go",
			old:  "\tfields[FieldHttpResponseBody] = string(resBody)\n",
			new:  "\t// Credential fields are redacted by script/patchsensitive after generation.\n\tfields[FieldHttpResponseBody] = redactSensitiveHTTP(string(resBody))\n",
		},
		{
			path: "internal/provider/utils.go",
			old:  "\treturn fmt.Sprintf(\"**Request**:\\n%s\\n**Response**:\\n%s\", string(dumpReq), string(dumpRes))\n",
			new:  "\t// Credential fields are redacted by script/patchsensitive after generation.\n\treturn redactSensitiveHTTP(fmt.Sprintf(\"**Request**:\\n%s\\n**Response**:\\n%s\", string(dumpReq), string(dumpRes)))\n",
		},
		{
			path: "internal/sdk/models/errors/apierror.go",
			old:  "\t\tbody = fmt.Sprintf(\"\\n%s\", e.Body)\n",
			new:  "\t\t// Credential fields are redacted by script/patchsensitive after generation.\n\t\tbody = fmt.Sprintf(\"\\n%s\", redactBody(e.Body))\n",
		},
	}

	seen := map[string]struct{}{}
	var failed bool
	for _, repl := range replacements {
		if err := apply(repl); err != nil {
			fmt.Fprintf(os.Stderr, "patchsensitive: %s: %v\n", repl.path, err)
			failed = true
			continue
		}
		seen[repl.path] = struct{}{}
	}
	if failed {
		os.Exit(1)
	}

	for path := range seen {
		if err := gofmtFile(path); err != nil {
			fmt.Fprintf(os.Stderr, "patchsensitive: format %s: %v\n", path, err)
			os.Exit(1)
		}
	}
}

func apply(repl replacement) error {
	raw, err := os.ReadFile(repl.path)
	if err != nil {
		return err
	}
	src := string(raw)
	if strings.Contains(src, repl.new) {
		return nil
	}
	if !strings.Contains(src, repl.old) {
		return fmt.Errorf("expected snippet not found:\n%s", repl.old)
	}
	updated := strings.Replace(src, repl.old, repl.new, 1)
	return os.WriteFile(repl.path, []byte(updated), 0o644)
}

func gofmtFile(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	formatted, err := format.Source(raw)
	if err != nil {
		return err
	}
	if string(formatted) == string(raw) {
		return nil
	}
	return os.WriteFile(path, formatted, 0o644)
}
