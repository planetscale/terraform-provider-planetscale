package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net"
	"net/http"
	"net/http/httptest"
	"slices"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/planetscale/terraform-provider-planetscale/internal/sdk"
	"github.com/planetscale/terraform-provider-planetscale/internal/sdk/models/shared"
	"github.com/stretchr/testify/require"
)

// fakeVitessAPI serves one Vitess branch, named "main" with the ID
// "branch-id", and its "commerce" keyspace: enough of the API for
// planetscale_vitess_branch and planetscale_vitess_keyspace to read and
// update, applying submitted parameter drafts to its overrides.
type fakeVitessAPI struct {
	mu              sync.Mutex
	safeMigrations  bool
	keyspaceResizes int
	droppedPolls    int
	defaults        map[string]map[string]string
	overrides       map[string]map[string]string
	drafts          map[string]fakeVitessDraft
	states          map[string]string
}

type fakeVitessDraft struct {
	component  string
	changeType string
	options    map[string]*string
}

func newFakeVitessAPI(t *testing.T) (*fakeVitessAPI, *sdk.PlanetScale) {
	api := &fakeVitessAPI{
		defaults: map[string]map[string]string{
			"vtgate":            {"query-timeout": "0", "max_memory_rows": "300000"},
			"commerce/vttablet": {"vreplication-parallel-insert-workers": "1"},
			"commerce/mysqld":   {"wait_timeout": "28800"},
		},
		overrides: map[string]map[string]string{
			"vtgate":            {"max_memory_rows": "500000"},
			"commerce/vttablet": {},
			"commerce/mysqld":   {"wait_timeout": "600"},
		},
		drafts: map[string]fakeVitessDraft{},
		states: map[string]string{},
	}

	const branch = "/v1/organizations/org/databases/db/branches/{branch}"
	const keyspace = branch + "/keyspaces/commerce"
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+branch, api.getBranch)
	mux.HandleFunc("PATCH "+branch, api.getBranch)
	mux.HandleFunc("PUT "+branch+"/safe-migrations", api.updateSafeMigrations(t))
	mux.HandleFunc("PUT "+branch+"/resizes", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("GET "+branch+"/vitess-parameters", api.listParameters)
	mux.HandleFunc("POST "+branch+"/config-changes", api.createDraft(t, ""))
	mux.HandleFunc("POST "+branch+"/config-changes/submit", api.submit(t))
	mux.HandleFunc("GET "+branch+"/config-changes/{id}", api.getChange)
	mux.HandleFunc("GET "+keyspace, api.getKeyspace)
	mux.HandleFunc("PUT "+keyspace+"/resizes", func(w http.ResponseWriter, _ *http.Request) {
		api.mu.Lock()
		defer api.mu.Unlock()
		api.keyspaceResizes++
		writeFakeJSON(w, http.StatusOK, map[string]string{"id": "resize-1", "state": "completed"})
	})
	mux.HandleFunc("POST "+keyspace+"/config-changes", api.createDraft(t, "commerce"))
	mux.HandleFunc("GET "+keyspace+"/config-changes/{id}", api.getChange)
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		t.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// The branch resource addresses the branch by ID and the keyspace
		// resource by name.
		if name := req.PathValue("branch"); name != "" && name != "branch-id" && name != "main" {
			t.Errorf("unexpected branch %q in %s", name, req.URL.Path)
		}
		mux.ServeHTTP(w, req)
	}))
	t.Cleanup(server.Close)

	// A new connection per request keeps the transport from silently retrying a
	// dropped connection, so the SDK sees it as the provider would.
	transport := server.Client().Transport.(*http.Transport).Clone()
	transport.DisableKeepAlives = true

	return api, sdk.New(
		sdk.WithServerURL(server.URL+"/v1"),
		sdk.WithClient(&http.Client{Transport: transport}),
		sdk.WithSecurity(shared.Security{ServiceToken: "token", ServiceTokenID: "id"}),
	)
}

func (a *fakeVitessAPI) getBranch(w http.ResponseWriter, _ *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()
	writeFakeJSON(w, http.StatusOK, map[string]any{
		"id":                 "branch-id",
		"name":               "main",
		"state":              "ready",
		"ready":              true,
		"cluster_name":       "PS_10",
		"safe_migrations":    a.safeMigrations,
		"deletion_protected": false,
	})
}

func (a *fakeVitessAPI) updateSafeMigrations(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var body struct {
			SafeMigrations bool `json:"safe_migrations"`
		}
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))

		a.mu.Lock()
		a.safeMigrations = body.SafeMigrations
		a.mu.Unlock()
		writeFakeJSON(w, http.StatusOK, map[string]any{})
	}
}

func (a *fakeVitessAPI) getKeyspace(w http.ResponseWriter, _ *http.Request) {
	writeFakeJSON(w, http.StatusOK, map[string]any{
		"id":                        "keyspace-id",
		"name":                      "commerce",
		"shards":                    1,
		"extra_replicas":            0,
		"cluster_name":              "PS_10",
		"ready":                     true,
		"default":                   true,
		"resizing":                  false,
		"resize_pending":            false,
		"config_change_in_progress": false,
	})
}

func (a *fakeVitessAPI) currentOverrides() map[string]map[string]string {
	a.mu.Lock()
	defer a.mu.Unlock()
	current := map[string]map[string]string{}
	for component, overrides := range a.overrides {
		current[component] = maps.Clone(overrides)
	}
	return current
}

func (a *fakeVitessAPI) draftCount() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.drafts)
}

func (a *fakeVitessAPI) parameterList(component string) []map[string]any {
	names := slices.Sorted(maps.Keys(a.defaults[component]))
	list := make([]map[string]any, 0, len(names))
	for _, name := range names {
		value, override := a.overrides[component][name]
		if !override {
			value = a.defaults[component][name]
		}
		list = append(list, map[string]any{
			"name":          name,
			"value":         value,
			"default_value": a.defaults[component][name],
			"override":      override,
		})
	}
	return list
}

func (a *fakeVitessAPI) listParameters(w http.ResponseWriter, _ *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()
	writeFakeJSON(w, http.StatusOK, map[string]any{
		"vtgate": a.parameterList("vtgate"),
		"keyspaces": map[string]any{
			"commerce": map[string]any{
				"vttablet": a.parameterList("commerce/vttablet"),
				"mysqld":   a.parameterList("commerce/mysqld"),
			},
		},
	})
}

func (a *fakeVitessAPI) createDraft(t *testing.T, keyspace string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var body struct {
			ChangeType string             `json:"change_type"`
			Options    map[string]*string `json:"options"`
		}
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))

		component := body.ChangeType
		if keyspace != "" {
			component = keyspace + "/" + body.ChangeType
		}

		a.mu.Lock()
		defer a.mu.Unlock()
		_, ok := a.defaults[component]
		require.True(t, ok, "unexpected change_type %q", body.ChangeType)
		id := fmt.Sprintf("change-%d", len(a.drafts)+1)
		a.drafts[id] = fakeVitessDraft{component: component, changeType: body.ChangeType, options: body.Options}
		a.states[id] = "draft"
		writeFakeJSON(w, http.StatusCreated, map[string]string{"id": id})
	}
}

func (a *fakeVitessAPI) submit(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var body struct {
			IDs []string `json:"ids"`
		}
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))

		a.mu.Lock()
		defer a.mu.Unlock()
		for _, id := range body.IDs {
			draft, ok := a.drafts[id]
			require.True(t, ok, "unknown draft %q", id)
			for name, value := range draft.options {
				if value == nil {
					delete(a.overrides[draft.component], name)
				} else {
					a.overrides[draft.component][name] = *value
				}
			}
			a.states[id] = "completed"
		}
		writeFakeJSON(w, http.StatusOK, []any{})
	}
}

func (a *fakeVitessAPI) getChange(w http.ResponseWriter, req *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.droppedPolls > 0 {
		a.droppedPolls--
		resetFakeConnection(w)
		return
	}
	id := req.PathValue("id")
	writeFakeJSON(w, http.StatusOK, map[string]string{
		"id":          id,
		"state":       a.states[id],
		"change_type": a.drafts[id].changeType,
	})
}

// resetFakeConnection closes the connection with a TCP reset, so the client
// reads "connection reset by peer" instead of a response.
func resetFakeConnection(w http.ResponseWriter) {
	conn, _, err := http.NewResponseController(w).Hijack()
	if err != nil {
		panic(err)
	}
	if tcp, ok := conn.(*net.TCPConn); ok {
		_ = tcp.SetLinger(0)
	}
	_ = conn.Close()
}

func writeFakeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func terraformValue(t *testing.T, ctx context.Context, s schema.Schema, model any) tftypes.Value {
	t.Helper()

	object, diags := types.ObjectValueFrom(ctx, s.Type().(types.ObjectType).AttrTypes, model)
	require.False(t, diags.HasError(), diags.Errors())
	raw, err := object.ToTerraformValue(ctx)
	require.NoError(t, err)
	return raw
}

func stateParameters(t *testing.T, ctx context.Context, state tfsdk.State) map[string]map[string]string {
	t.Helper()

	var parameters map[string]map[string]string
	diags := state.GetAttribute(ctx, path.Root("parameters"), &parameters)
	require.False(t, diags.HasError(), diags.Errors())
	return parameters
}

// readVitessBranch reads the fake branch as Terraform does after import.
func readVitessBranch(t *testing.T, ctx context.Context, r *VitessBranchResource, s schema.Schema) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: s, Raw: terraformValue(t, ctx, s, &VitessBranchResourceModel{
		Organization: types.StringValue("org"),
		Database:     types.StringValue("db"),
		ID:           types.StringValue("branch-id"),
	})}
	response := resource.ReadResponse{State: state}
	r.Read(ctx, resource.ReadRequest{State: state}, &response)
	require.False(t, response.Diagnostics.HasError(), response.Diagnostics.Errors())
	return response.State
}

func updateVitessBranch(t *testing.T, ctx context.Context, r *VitessBranchResource, state tfsdk.State, changes map[string]any) tfsdk.State {
	t.Helper()

	plan := tfsdk.Plan{Schema: state.Schema, Raw: state.Raw.Copy()}
	for name, value := range changes {
		diags := plan.SetAttribute(ctx, path.Root(name), value)
		require.False(t, diags.HasError(), diags.Errors())
	}
	response := resource.UpdateResponse{State: tfsdk.State{Schema: state.Schema, Raw: plan.Raw.Copy()}}
	r.Update(ctx, resource.UpdateRequest{Plan: plan, State: state}, &response)
	require.False(t, response.Diagnostics.HasError(), response.Diagnostics.Errors())
	return response.State
}

func TestVitessBranchResourceReadsParameterOverrides(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	_, client := newFakeVitessAPI(t)
	r := &VitessBranchResource{client: client}
	var schemaResponse resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)

	state := readVitessBranch(t, ctx, r, schemaResponse.Schema)
	require.Equal(t, map[string]map[string]string{
		"vtgate": {"max_memory_rows": "500000"},
	}, stateParameters(t, ctx, state))
}

func TestVitessBranchResourceUpdateLeavesUnsetParametersUnchanged(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	api, client := newFakeVitessAPI(t)
	r := &VitessBranchResource{client: client}
	var schemaResponse resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
	state := readVitessBranch(t, ctx, r, schemaResponse.Schema)

	// Terraform plans an unset Optional and Computed attribute as unknown
	// when something else about the resource changes.
	state = updateVitessBranch(t, ctx, r, state, map[string]any{
		"safe_migrations": types.BoolValue(true),
		"parameters":      types.MapUnknown(types.MapType{ElemType: types.StringType}),
	})

	require.Zero(t, api.draftCount())
	require.Equal(t, map[string]string{"max_memory_rows": "500000"}, api.currentOverrides()["vtgate"])
	require.Equal(t, map[string]map[string]string{
		"vtgate": {"max_memory_rows": "500000"},
	}, stateParameters(t, ctx, state))
}

func TestVitessBranchResourceUpdateReplacesConfiguredParameters(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	api, client := newFakeVitessAPI(t)
	r := &VitessBranchResource{client: client}
	var schemaResponse resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
	state := readVitessBranch(t, ctx, r, schemaResponse.Schema)

	configured := map[string]map[string]string{"vtgate": {"query-timeout": "30000"}}
	parameters, diags := types.MapValueFrom(ctx, types.MapType{ElemType: types.StringType}, configured)
	require.False(t, diags.HasError(), diags.Errors())
	state = updateVitessBranch(t, ctx, r, state, map[string]any{"parameters": parameters})

	require.Equal(t, map[string]string{"query-timeout": "30000"}, api.currentOverrides()["vtgate"])
	require.Equal(t, configured, stateParameters(t, ctx, state))
}

func TestVitessBranchResourceUpdateRetriesDroppedConfigChangePolls(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	api, client := newFakeVitessAPI(t)
	r := &VitessBranchResource{client: client}
	var schemaResponse resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
	state := readVitessBranch(t, ctx, r, schemaResponse.Schema)

	api.mu.Lock()
	api.droppedPolls = 1
	api.mu.Unlock()

	configured := map[string]map[string]string{"vtgate": {"query-timeout": "30000"}}
	parameters, diags := types.MapValueFrom(ctx, types.MapType{ElemType: types.StringType}, configured)
	require.False(t, diags.HasError(), diags.Errors())
	state = updateVitessBranch(t, ctx, r, state, map[string]any{"parameters": parameters})

	require.Zero(t, api.droppedPolls)
	require.Equal(t, configured, stateParameters(t, ctx, state))
}
