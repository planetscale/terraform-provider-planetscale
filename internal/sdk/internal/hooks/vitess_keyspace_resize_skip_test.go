package hooks

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

const vitessTestKeyspacePath = vitessTestBranchPath + "/keyspaces/commerce"

func TestVitessKeyspaceResizeSkipHook(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		body           map[string]any
		keyspaceStatus int
		wantResize     bool
	}{
		"unchanged size and replicas skip the resize": {
			body: map[string]any{"cluster_size": "PS_10", "extra_replicas": 1},
		},
		"a new size resizes": {
			body:       map[string]any{"cluster_size": "PS_20", "extra_replicas": 1},
			wantResize: true,
		},
		"new replicas resize": {
			body:       map[string]any{"cluster_size": "PS_10", "extra_replicas": 2},
			wantResize: true,
		},
		"a failed keyspace lookup resizes": {
			body:           map[string]any{"cluster_size": "PS_10", "extra_replicas": 1},
			keyspaceStatus: http.StatusInternalServerError,
			wantResize:     true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var resized map[string]any
			_, client := NewVitessKeyspaceResizeSkipHook().SDKInit("https://api.planetscale.com", testHTTPClient(func(req *http.Request) (*http.Response, error) {
				require.Equal(t, "Bearer test", req.Header.Get("Authorization"))
				switch {
				case req.Method == http.MethodGet && req.URL.Path == vitessTestKeyspacePath:
					if tt.keyspaceStatus != 0 {
						return vitessTestResponse(req, tt.keyspaceStatus, `{"code":"error"}`), nil
					}
					return vitessTestResponse(req, http.StatusOK, `{"name":"commerce","cluster_name":"PS_10","extra_replicas":1}`), nil
				case req.Method == http.MethodPut && req.URL.Path == vitessTestKeyspacePath+"/resizes":
					require.NoError(t, json.NewDecoder(req.Body).Decode(&resized))
					return vitessTestResponse(req, http.StatusOK, `{"id":"resize-1","state":"pending"}`), nil
				}
				require.FailNow(t, "unexpected request", "%s %s", req.Method, req.URL.Path)
				return nil, nil
			}))

			res, err := client.Do(vitessTestRequest(t, http.MethodPut, vitessTestKeyspacePath+"/resizes", tt.body))
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, res.StatusCode)
			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			if !tt.wantResize {
				require.Nil(t, resized)
				require.JSONEq(t, `{"id":"`+skippedKeyspaceResizeRequestID+`","state":"completed"}`, string(body))
				return
			}
			require.Equal(t, map[string]any{"cluster_size": tt.body["cluster_size"], "extra_replicas": float64(tt.body["extra_replicas"].(int))}, resized)
			require.JSONEq(t, `{"id":"resize-1","state":"pending"}`, string(body))
		})
	}
}
