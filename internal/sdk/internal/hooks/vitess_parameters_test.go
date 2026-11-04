package hooks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	vitessTestBranchPath = "/v1/organizations/org/databases/db/branches/main"
	vitessTestCatalog    = `{
		"vtgate": [
			{"name": "query-timeout", "value": "30000", "override": true},
			{"name": "buffer_window", "value": "10s", "override": false},
			{"name": "max_memory_rows", "value": "300000", "override": false}
		],
		"keyspaces": {
			"commerce": {
				"vttablet": [
					{"name": "queryserver-config-pool-size", "value": "50", "override": true},
					{"name": "vreplication-parallel-insert-workers", "value": "1", "override": false}
				],
				"mysqld": [
					{"name": "wait_timeout", "value": "600", "override": true},
					{"name": "max_allowed_packet", "value": "67108864", "override": false},
					{"name": "innodb_adaptive_hash_index", "value": "false", "override": true}
				]
			}
		}
	}`
)

type vitessDraft struct {
	Path       string         `json:"-"`
	ChangeType string         `json:"change_type"`
	Options    map[string]any `json:"options"`
}

type vitessConfigChange struct {
	State        string
	ErrorMessage string
}

// vitessAPIFake serves the endpoints VitessParametersHook calls and records
// the drafts, submissions, and cancellations it makes.
type vitessAPIFake struct {
	catalog       string
	catalogStatus int
	failDraft     int
	submitStatus  int
	changes       map[string]vitessConfigChange

	requests  []string
	drafts    []vitessDraft
	submitted []string
	canceled  []string
}

func (f *vitessAPIFake) client(t *testing.T) HTTPClient {
	return testHTTPClient(func(req *http.Request) (*http.Response, error) {
		f.requests = append(f.requests, req.Method+" "+req.URL.Path)
		require.Equal(t, "Bearer test", req.Header.Get("Authorization"))
		require.Empty(t, req.URL.RawQuery)

		switch {
		case req.Method == http.MethodGet && req.URL.Path == vitessTestBranchPath+"/vitess-parameters":
			if f.catalogStatus != 0 {
				return vitessTestResponse(req, f.catalogStatus, `{"code":"error","message":"Request failed"}`), nil
			}
			if f.catalog != "" {
				return vitessTestResponse(req, http.StatusOK, f.catalog), nil
			}
			return vitessTestResponse(req, http.StatusOK, vitessTestCatalog), nil
		case req.Method == http.MethodPost && req.URL.Path == vitessTestBranchPath+"/config-changes/submit":
			var body struct {
				IDs []string `json:"ids"`
			}
			require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
			f.submitted = body.IDs
			if f.submitStatus != 0 {
				return vitessTestResponse(req, f.submitStatus, `{"code":"unprocessable","message":"Branch is not ready"}`), nil
			}
			return vitessTestResponse(req, http.StatusOK, `[]`), nil
		case req.Method == http.MethodPost && strings.HasSuffix(req.URL.Path, "/config-changes"):
			draft := vitessDraft{Path: req.URL.Path}
			require.NoError(t, json.NewDecoder(req.Body).Decode(&draft))
			f.drafts = append(f.drafts, draft)
			if len(f.drafts) == f.failDraft {
				return vitessTestResponse(req, http.StatusUnprocessableEntity, `{"code":"unprocessable","message":"A change is already in progress"}`), nil
			}
			return vitessTestResponse(req, http.StatusCreated, fmt.Sprintf(`{"id":"change-%d"}`, len(f.drafts))), nil
		case req.Method == http.MethodDelete && strings.Contains(req.URL.Path, "/config-changes/"):
			f.canceled = append(f.canceled, path.Base(req.URL.Path))
			return vitessTestResponse(req, http.StatusNoContent, ``), nil
		case req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/config-changes/"):
			id := path.Base(req.URL.Path)
			change, ok := f.changes[id]
			require.True(t, ok, "unexpected config change %s", id)
			payload, err := json.Marshal(map[string]any{
				"id":            id,
				"state":         change.State,
				"change_type":   "mysqld",
				"error_message": change.ErrorMessage,
			})
			require.NoError(t, err)
			return vitessTestResponse(req, http.StatusOK, string(payload)), nil
		}

		require.FailNow(t, "unexpected request", "%s %s", req.Method, req.URL.Path)
		return nil, nil
	})
}

func vitessTestResponse(req *http.Request, status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": {"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    req,
	}
}

func vitessTestRequest(t *testing.T, method string, path string, body any) *http.Request {
	t.Helper()

	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		require.NoError(t, err)
		reader = strings.NewReader(string(payload))
	}
	req, err := http.NewRequest(method, "https://api.planetscale.com"+path, reader)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer test")
	return req
}

func vitessTestDo(t *testing.T, fake *vitessAPIFake, req *http.Request) (*http.Response, error) {
	t.Helper()

	_, client := NewVitessParametersHook().SDKInit("https://api.planetscale.com", fake.client(t))
	return client.Do(req)
}

func vitessTestDecode(t *testing.T, res *http.Response) vitessParametersResponse {
	t.Helper()

	require.Equal(t, http.StatusOK, res.StatusCode)
	var payload vitessParametersResponse
	require.NoError(t, json.NewDecoder(res.Body).Decode(&payload))
	return payload
}

func TestVitessParametersHookReadsOverriddenVTGateParameters(t *testing.T) {
	t.Parallel()

	fake := &vitessAPIFake{}
	res, err := vitessTestDo(t, fake, vitessTestRequest(t, http.MethodGet, vitessTestBranchPath+"/vitess-parameters", nil))
	require.NoError(t, err)

	require.Equal(t, vitessParametersResponse{
		Parameters: map[string]map[string]string{
			"vtgate": {"query-timeout": "30000"},
		},
	}, vitessTestDecode(t, res))
	require.Equal(t, []string{"GET " + vitessTestBranchPath + "/vitess-parameters"}, fake.requests)
}

func TestVitessParametersHookReadsKeyspaceParameters(t *testing.T) {
	t.Parallel()

	fake := &vitessAPIFake{}
	res, err := vitessTestDo(t, fake, vitessTestRequest(t, http.MethodGet, vitessTestBranchPath+"/keyspaces/commerce/vitess-parameters", nil))
	require.NoError(t, err)

	require.Equal(t, vitessParametersResponse{
		Parameters: map[string]map[string]string{
			"vttablet": {"queryserver-config-pool-size": "50"},
			"mysqld": {
				"wait_timeout":               "600",
				"innodb_adaptive_hash_index": "false",
			},
		},
	}, vitessTestDecode(t, res))
}

func TestVitessParametersHookReadsNoOverridesAsAnEmptyMap(t *testing.T) {
	t.Parallel()

	fake := &vitessAPIFake{catalog: `{"vtgate": [{"name": "query-timeout", "value": "0", "override": false}], "keyspaces": {}}`}
	res, err := vitessTestDo(t, fake, vitessTestRequest(t, http.MethodGet, vitessTestBranchPath+"/vitess-parameters", nil))
	require.NoError(t, err)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.JSONEq(t, `{"parameters":{}}`, string(body))
}

func TestVitessParametersHookFailsReadsOfAMissingKeyspace(t *testing.T) {
	t.Parallel()

	_, err := vitessTestDo(t, &vitessAPIFake{}, vitessTestRequest(t, http.MethodGet, vitessTestBranchPath+"/keyspaces/missing/vitess-parameters", nil))
	require.ErrorContains(t, err, `did not list keyspace "missing"`)
}

func TestVitessParametersHookHandlesParameterListErrors(t *testing.T) {
	t.Parallel()

	_, err := vitessTestDo(t, &vitessAPIFake{catalogStatus: http.StatusNotFound}, vitessTestRequest(t, http.MethodGet, vitessTestBranchPath+"/vitess-parameters", nil))
	require.ErrorContains(t, err, "refusing to treat the branch as missing")

	res, err := vitessTestDo(t, &vitessAPIFake{catalogStatus: http.StatusUnprocessableEntity}, vitessTestRequest(t, http.MethodGet, vitessTestBranchPath+"/vitess-parameters", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)
}

func TestVitessParametersHookChangesNothingWithoutConfiguredParameters(t *testing.T) {
	t.Parallel()

	for _, method := range []string{http.MethodPatch, http.MethodPut} {
		for _, body := range []map[string]any{{}, {"parameters": nil}} {
			fake := &vitessAPIFake{}
			res, err := vitessTestDo(t, fake, vitessTestRequest(t, method, vitessTestBranchPath+"/keyspaces/commerce/vitess-parameters", body))
			require.NoError(t, err)

			require.Equal(t, vitessParametersResponse{
				Parameters: map[string]map[string]string{
					"vttablet": {"queryserver-config-pool-size": "50"},
					"mysqld": {
						"wait_timeout":               "600",
						"innodb_adaptive_hash_index": "false",
					},
				},
				ConfigChangeID: skippedVitessConfigChangeID,
			}, vitessTestDecode(t, res), "%s %v", method, body)
			require.Equal(t, []string{"GET " + vitessTestBranchPath + "/vitess-parameters"}, fake.requests)
		}
	}
}

func TestVitessParametersHookCreateOnlySetsConfiguredParameters(t *testing.T) {
	t.Parallel()

	desired := map[string]map[string]string{
		"vttablet": {"vreplication-parallel-insert-workers": "4"},
		"mysqld": {
			"wait_timeout":       "600",
			"max_allowed_packet": "67108864",
		},
	}
	fake := &vitessAPIFake{}
	res, err := vitessTestDo(t, fake, vitessTestRequest(t, http.MethodPatch, vitessTestBranchPath+"/keyspaces/commerce/vitess-parameters", map[string]any{
		"parameters": desired,
	}))
	require.NoError(t, err)

	require.Equal(t, vitessParametersResponse{
		Parameters:     desired,
		ConfigChangeID: "change-1,change-2",
	}, vitessTestDecode(t, res))
	require.Equal(t, []vitessDraft{
		{
			Path:       vitessTestBranchPath + "/keyspaces/commerce/config-changes",
			ChangeType: "vttablet",
			Options:    map[string]any{"vreplication-parallel-insert-workers": "4"},
		},
		{
			Path:       vitessTestBranchPath + "/keyspaces/commerce/config-changes",
			ChangeType: "mysqld",
			Options:    map[string]any{"max_allowed_packet": "67108864"},
		},
	}, fake.drafts)
	require.Equal(t, []string{"change-1", "change-2"}, fake.submitted)
	require.Empty(t, fake.canceled)
}

func TestVitessParametersHookUpdateResetsUnlistedOverrides(t *testing.T) {
	t.Parallel()

	fake := &vitessAPIFake{}
	res, err := vitessTestDo(t, fake, vitessTestRequest(t, http.MethodPut, vitessTestBranchPath+"/keyspaces/commerce/vitess-parameters", map[string]any{
		"parameters": map[string]map[string]string{
			"mysqld": {"wait_timeout": "300"},
		},
	}))
	require.NoError(t, err)
	require.Equal(t, "change-1,change-2", vitessTestDecode(t, res).ConfigChangeID)

	require.Equal(t, []vitessDraft{
		{
			Path:       vitessTestBranchPath + "/keyspaces/commerce/config-changes",
			ChangeType: "vttablet",
			Options:    map[string]any{"queryserver-config-pool-size": nil},
		},
		{
			Path:       vitessTestBranchPath + "/keyspaces/commerce/config-changes",
			ChangeType: "mysqld",
			Options: map[string]any{
				"wait_timeout":               "300",
				"innodb_adaptive_hash_index": nil,
			},
		},
	}, fake.drafts)
}

func TestVitessParametersHookSkipsUnchangedParameters(t *testing.T) {
	t.Parallel()

	fake := &vitessAPIFake{}
	_, client := NewVitessParametersHook().SDKInit("https://api.planetscale.com", fake.client(t))

	res, err := client.Do(vitessTestRequest(t, http.MethodPut, vitessTestBranchPath+"/vitess-parameters", map[string]any{
		"parameters": map[string]map[string]string{
			"vtgate": {"query-timeout": "30000"},
		},
	}))
	require.NoError(t, err)
	payload := vitessTestDecode(t, res)
	require.Equal(t, skippedVitessConfigChangeID, payload.ConfigChangeID)
	require.Empty(t, fake.drafts)

	res, err = client.Do(vitessTestRequest(t, http.MethodGet, vitessTestBranchPath+"/config-changes/"+payload.ConfigChangeID, nil))
	require.NoError(t, err)
	var status map[string]string
	require.NoError(t, json.NewDecoder(res.Body).Decode(&status))
	require.Equal(t, "completed", status["state"])
	require.Equal(t, []string{"GET " + vitessTestBranchPath + "/vitess-parameters"}, fake.requests)
}

func TestVitessParametersHookCancelsDraftsWhenSubmitFails(t *testing.T) {
	t.Parallel()

	fake := &vitessAPIFake{submitStatus: http.StatusUnprocessableEntity}
	res, err := vitessTestDo(t, fake, vitessTestRequest(t, http.MethodPut, vitessTestBranchPath+"/keyspaces/commerce/vitess-parameters", map[string]any{
		"parameters": map[string]map[string]string{
			"mysqld": {"wait_timeout": "300"},
		},
	}))
	require.NoError(t, err)
	require.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)
	require.Equal(t, []string{"change-1", "change-2"}, fake.canceled)
}

func TestVitessParametersHookCancelsEarlierDraftsWhenADraftFails(t *testing.T) {
	t.Parallel()

	fake := &vitessAPIFake{failDraft: 2}
	res, err := vitessTestDo(t, fake, vitessTestRequest(t, http.MethodPut, vitessTestBranchPath+"/keyspaces/commerce/vitess-parameters", map[string]any{
		"parameters": map[string]map[string]string{
			"mysqld": {"wait_timeout": "300"},
		},
	}))
	require.NoError(t, err)
	require.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Contains(t, string(body), "A change is already in progress")
	require.Equal(t, []string{"change-1"}, fake.canceled)
	require.Nil(t, fake.submitted)
}

func TestVitessParametersHookReportsCombinedConfigChangeState(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		changes   map[string]vitessConfigChange
		wantState string
		wantErr   string
	}{
		"completed and queued changes are done": {
			changes: map[string]vitessConfigChange{
				"change-1": {State: "completed"},
				"change-2": {State: "queued"},
			},
			wantState: "completed",
		},
		"a change still rolling out keeps waiting": {
			changes: map[string]vitessConfigChange{
				"change-1": {State: "completed"},
				"change-2": {State: "in_progress"},
			},
			wantState: "in_progress",
		},
		"a failed change reports the API error": {
			changes: map[string]vitessConfigChange{
				"change-1": {State: "completed"},
				"change-2": {State: "failed", ErrorMessage: "wait_timeout must be at most 31536000"},
			},
			wantErr: "mysqld config change change-2 failed: wait_timeout must be at most 31536000",
		},
		"a canceled change fails": {
			changes: map[string]vitessConfigChange{
				"change-1": {State: "canceled"},
				"change-2": {State: "completed"},
			},
			wantErr: "mysqld config change change-1 was canceled before it applied",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			fake := &vitessAPIFake{changes: tt.changes}
			res, err := vitessTestDo(t, fake, vitessTestRequest(t, http.MethodGet, vitessTestBranchPath+"/keyspaces/commerce/config-changes/change-1,change-2", nil))
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)

			var status map[string]string
			require.NoError(t, json.NewDecoder(res.Body).Decode(&status))
			require.Equal(t, map[string]string{"id": "change-1,change-2", "state": tt.wantState}, status)
		})
	}
}

func TestVitessParametersHookPassesThroughOtherRequests(t *testing.T) {
	t.Parallel()

	var calls int
	_, client := NewVitessParametersHook().SDKInit("https://api.planetscale.com", testHTTPClient(func(req *http.Request) (*http.Response, error) {
		calls++
		return vitessTestResponse(req, http.StatusOK, `{"passthrough":true}`), nil
	}))

	for _, req := range []*http.Request{
		vitessTestRequest(t, http.MethodGet, vitessTestBranchPath+"/keyspaces/commerce", nil),
		vitessTestRequest(t, http.MethodGet, vitessTestBranchPath+"/config-changes", nil),
		vitessTestRequest(t, http.MethodDelete, vitessTestBranchPath+"/config-changes/change-1", nil),
	} {
		res, err := client.Do(req)
		require.NoError(t, err)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		require.JSONEq(t, `{"passthrough":true}`, string(body))
	}
	require.Equal(t, 3, calls)
}
