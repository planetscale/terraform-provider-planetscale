// Package redact removes credential values from HTTP payloads before they are
// written to Terraform logs or error diagnostics.
package redact

import "regexp"

// credentialField matches JSON string fields that carry generated database
// credentials. The value may contain escaped characters.
var credentialField = regexp.MustCompile(`"(password|plain_text|plain_text_refresh_token)"\s*:\s*"(?:\\.|[^"\\])*"`)

// RedactCredentialFields replaces credential JSON string values with
// "(sensitive)". Non-credential text is left unchanged, including when the
// input is an HTTP dump rather than a bare JSON document.
func RedactCredentialFields(s string) string {
	if s == "" {
		return s
	}
	return credentialField.ReplaceAllString(s, `"$1":"(sensitive)"`)
}
