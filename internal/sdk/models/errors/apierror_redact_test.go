package errors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAPIErrorRedactsCredentialBody(t *testing.T) {
	t.Parallel()

	err := NewAPIError("unknown status code returned", http.StatusConflict, `{"code":"conflict","password":"from-the-body","message":"already exists"}`, nil)
	got := err.Error()

	require.NotContains(t, got, "from-the-body")
	require.Contains(t, got, "unknown status code returned: Status 409")
	require.Contains(t, got, `"password":"(sensitive)"`)
	require.Contains(t, got, `"message":"already exists"`)
}
