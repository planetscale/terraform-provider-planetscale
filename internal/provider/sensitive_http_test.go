package provider

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDebugResponseRedactsCredentials(t *testing.T) {
	t.Parallel()

	body := `{"id":"pw1","username":"app","password":"role-secret","plain_text":"vitess-secret"}`
	req, err := http.NewRequest(http.MethodPost, "https://api.planetscale.com/v1/organizations/acme/databases/app/branches/main/roles", strings.NewReader(`{"name":"app"}`))
	require.NoError(t, err)
	req.Header.Set("Authorization", "token-id:super-token")

	res := &http.Response{
		StatusCode:    http.StatusOK,
		Status:        "200 OK",
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        make(http.Header),
		Body:          io.NopCloser(strings.NewReader(body)),
		ContentLength: int64(len(body)),
		Request:       req,
	}
	res.Header.Set("Content-Type", "application/json")

	got := debugResponse(res)
	require.NotContains(t, got, "role-secret")
	require.NotContains(t, got, "vitess-secret")
	require.NotContains(t, got, "super-token")
	require.Contains(t, got, `"password":"(sensitive)"`)
	require.Contains(t, got, `"plain_text":"(sensitive)"`)
	require.Contains(t, got, `"username":"app"`)
}

func TestDecomposeResponseForLoggingRedactsCredentials(t *testing.T) {
	t.Parallel()

	body := `{"password":"logged-secret","name":"app"}`
	res := &http.Response{
		StatusCode:    http.StatusOK,
		Status:        "200 OK",
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        make(http.Header),
		Body:          io.NopCloser(strings.NewReader(body)),
		ContentLength: int64(len(body)),
	}

	fields, err := decomposeResponseForLogging(res)
	require.NoError(t, err)

	logged, ok := fields[FieldHttpResponseBody].(string)
	require.True(t, ok)
	require.NotContains(t, logged, "logged-secret")
	require.Contains(t, logged, `"password":"(sensitive)"`)
	require.Contains(t, logged, `"name":"app"`)

	// The response body handed back to the SDK is unchanged.
	restored, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, body, string(restored))
}

func TestDecomposeRequestForLoggingRedactsCredentials(t *testing.T) {
	t.Parallel()

	body := `{"password":"request-secret","name":"app"}`
	req, err := http.NewRequest(http.MethodPost, "https://api.planetscale.com/v1/roles", strings.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "token-id:super-token")

	fields, err := decomposeRequestForLogging(req)
	require.NoError(t, err)

	logged, ok := fields[FieldHttpRequestBody].(string)
	require.True(t, ok)
	require.NotContains(t, logged, "request-secret")
	require.Contains(t, logged, `"password":"(sensitive)"`)
	require.Equal(t, "(sensitive)", fields["Authorization"])
}
