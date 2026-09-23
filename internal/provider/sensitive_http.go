package provider

import "github.com/planetscale/terraform-provider-planetscale/internal/redact"

// redactSensitiveHTTP removes generated credential values from HTTP bodies
// before they are written to debug logs or error diagnostics.
//
// script/patchsensitive rewrites the generated logging helpers to call this
// after Speakeasy regeneration.
func redactSensitiveHTTP(s string) string {
	return redact.RedactCredentialFields(s)
}
