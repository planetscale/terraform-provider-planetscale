package hooks

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func clientErrorResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json; charset=utf-8"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestClientErrorHookReportsCodeAndMessage(t *testing.T) {
	res, err := NewClientErrorHook().AfterSuccess(AfterSuccessContext{}, clientErrorResponse(http.StatusUnprocessableEntity, `{"code":"unprocessable","message":"Minimum storage bytes must be greater than or equal to 10 GB"}`))

	require.EqualError(t, err, "PlanetScale API returned HTTP 422 (unprocessable): Minimum storage bytes must be greater than or equal to 10 GB")

	var clientErr *ClientError
	require.ErrorAs(t, err, &clientErr)
	require.Equal(t, http.StatusUnprocessableEntity, clientErr.StatusCode)
	require.Equal(t, "unprocessable", clientErr.Code)

	body, readErr := io.ReadAll(res.Body)
	require.NoError(t, readErr)
	require.JSONEq(t, `{"code":"unprocessable","message":"Minimum storage bytes must be greater than or equal to 10 GB"}`, string(body), "body should be restored for downstream consumers")
}

func TestClientErrorHookFallsBackToRawBody(t *testing.T) {
	_, err := NewClientErrorHook().AfterSuccess(AfterSuccessContext{}, clientErrorResponse(http.StatusForbidden, "<html>forbidden</html>"))
	require.EqualError(t, err, "PlanetScale API returned HTTP 403: <html>forbidden</html>")

	_, err = NewClientErrorHook().AfterSuccess(AfterSuccessContext{}, clientErrorResponse(http.StatusUnauthorized, ""))
	require.EqualError(t, err, "PlanetScale API returned HTTP 401")
}

func TestClientErrorHookIgnoresOtherStatuses(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusCreated, http.StatusNoContent, http.StatusNotFound, http.StatusInternalServerError} {
		res := clientErrorResponse(status, `{"code":"x","message":"y"}`)
		got, err := NewClientErrorHook().AfterSuccess(AfterSuccessContext{}, res)
		require.NoError(t, err, "status %d", status)
		require.Same(t, res, got, "status %d", status)
	}
}
