package errors

import "github.com/planetscale/terraform-provider-planetscale/internal/redact"

// redactBody removes generated credential values from API error details.
//
// script/patchsensitive rewrites the generated APIError.Error method to call
// this after Speakeasy regeneration.
func redactBody(body string) string {
	return redact.RedactCredentialFields(body)
}
