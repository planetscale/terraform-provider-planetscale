package hooks

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRedactedRolePasswordHookStripsPassword(t *testing.T) {
	t.Parallel()

	body := `{"id":"role1","username":"app","password":"super-secret"}`
	res, err := NewRedactedRolePasswordHook().AfterSuccess(AfterSuccessContext{
		HookContext: HookContext{OperationID: "create_redacted_role"},
	}, clientErrorResponse(http.StatusOK, body))

	require.NoError(t, err)
	got, readErr := io.ReadAll(res.Body)
	require.NoError(t, readErr)
	require.NotContains(t, string(got), "super-secret")
	require.Contains(t, string(got), `"username":"app"`)
	require.Contains(t, string(got), `"password":"(sensitive)"`)
	require.Equal(t, int64(len(got)), res.ContentLength)
}

func TestRedactedRolePasswordHookLeavesOtherOperations(t *testing.T) {
	t.Parallel()

	body := `{"password":"keep-for-state"}`
	original := clientErrorResponse(http.StatusOK, body)
	res, err := NewRedactedRolePasswordHook().AfterSuccess(AfterSuccessContext{
		HookContext: HookContext{OperationID: "create_role"},
	}, original)

	require.NoError(t, err)
	require.Same(t, original, res)
	got, readErr := io.ReadAll(res.Body)
	require.NoError(t, readErr)
	require.Equal(t, body, string(got))
}

func TestRedactedRolePasswordHookIgnoresEmptyBody(t *testing.T) {
	t.Parallel()

	res := &http.Response{StatusCode: http.StatusNoContent, Body: io.NopCloser(strings.NewReader(""))}
	got, err := NewRedactedRolePasswordHook().AfterSuccess(AfterSuccessContext{
		HookContext: HookContext{OperationID: "delete_redacted_role"},
	}, res)
	require.NoError(t, err)
	require.Same(t, res, got)
}
