package redact

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRedactCredentialFields(t *testing.T) {
	t.Parallel()

	input := `{"id":"role1","name":"app","password":"s3cret\"value","plain_text":"hunter2","plain_text_refresh_token":"refresh-secret","username":"pscale"}`
	got := RedactCredentialFields(input)

	require.NotContains(t, got, `s3cret\"value`)
	require.NotContains(t, got, "hunter2")
	require.NotContains(t, got, "refresh-secret")
	require.Contains(t, got, `"id":"role1"`)
	require.Contains(t, got, `"name":"app"`)
	require.Contains(t, got, `"username":"pscale"`)
	require.Contains(t, got, `"password":"(sensitive)"`)
	require.Contains(t, got, `"plain_text":"(sensitive)"`)
	require.Contains(t, got, `"plain_text_refresh_token":"(sensitive)"`)
	require.Equal(t, got, RedactCredentialFields(got))
}

func TestRedactCredentialFieldsIgnoresUnrelatedText(t *testing.T) {
	t.Parallel()

	input := "HTTP/1.1 422 Unprocessable Entity\r\n\r\n{\"code\":\"unprocessable\",\"message\":\"invalid password policy\"}"
	require.Equal(t, input, RedactCredentialFields(input))
	require.Equal(t, "", RedactCredentialFields(""))
}

func TestRedactCredentialFieldsNestedAndWhitespace(t *testing.T) {
	t.Parallel()

	input := "{\n  \"data\": [{\"password\": \"nested-secret\"}]\n}"
	got := RedactCredentialFields(input)
	require.NotContains(t, got, "nested-secret")
	require.True(t, strings.Contains(got, `"password":"(sensitive)"`))
}
