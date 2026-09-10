package hooks

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNekiShardReassignmentHookTransformsRequest(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest(
		http.MethodPatch,
		"https://api.planetscale.com/v1/organizations/acme/databases/db/branches/main/configuration-profiles/analytics/shards",
		strings.NewReader(`{"id":"shard-1"}`),
	)
	require.NoError(t, err)

	hook := NewNekiShardReassignmentHook()
	req, err = hook.BeforeRequest(BeforeRequestContext{
		HookContext: HookContext{OperationID: reassignNekiShardOperationID},
	}, req)
	require.NoError(t, err)

	var body map[string][]string
	require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
	require.Equal(t, map[string][]string{"shard_ids": {"shard-1"}}, body)
	require.NotNil(t, req.GetBody)
	require.Positive(t, req.ContentLength)
}

func TestNekiShardReassignmentHookTransformsResponse(t *testing.T) {
	t.Parallel()

	res := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body: io.NopCloser(strings.NewReader(
			`[{"id":"shard-1","status":"assigned"}]`,
		)),
	}

	hook := NewNekiShardReassignmentHook()
	res, err := hook.AfterSuccess(AfterSuccessContext{
		HookContext: HookContext{OperationID: reassignNekiShardOperationID},
	}, res)
	require.NoError(t, err)

	var body map[string]string
	require.NoError(t, json.NewDecoder(res.Body).Decode(&body))
	require.Equal(t, map[string]string{"id": "shard-1"}, body)
	require.Positive(t, res.ContentLength)
}

func TestNekiShardReassignmentHookReportsAssignmentFailure(t *testing.T) {
	t.Parallel()

	res := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body: io.NopCloser(strings.NewReader(
			`[{"id":"shard-1","status":"failed","error":{"code":"incompatible","message":"profiles differ"}}]`,
		)),
	}

	hook := NewNekiShardReassignmentHook()
	_, err := hook.AfterSuccess(AfterSuccessContext{
		HookContext: HookContext{OperationID: reassignNekiShardOperationID},
	}, res)
	require.EqualError(
		t,
		err,
		`neki shard "shard-1" reassignment failed: incompatible: profiles differ`,
	)
}

func TestNekiShardReassignmentHookIgnoresOtherOperations(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest(
		http.MethodPatch,
		"https://api.planetscale.com/v1/other",
		strings.NewReader(`{"id":"shard-1"}`),
	)
	require.NoError(t, err)

	hook := NewNekiShardReassignmentHook()
	got, err := hook.BeforeRequest(BeforeRequestContext{
		HookContext: HookContext{OperationID: "other_operation"},
	}, req)
	require.NoError(t, err)
	require.Same(t, req, got)
}
