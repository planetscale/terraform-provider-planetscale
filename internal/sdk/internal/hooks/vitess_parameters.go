package hooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

// skippedVitessConfigChangeID stands in for a config change when the
// configured parameters already match, so polling completes without calling
// the API.
const skippedVitessConfigChangeID = "__speakeasy_skipped_config_change__"

var (
	vitessParametersPathPattern   = regexp.MustCompile(`^(/v1/organizations/[^/]+/databases/[^/]+/branches/[^/]+)(?:/keyspaces/([^/]+))?/vitess-parameters$`)
	vitessConfigChangePathPattern = regexp.MustCompile(`^(/v1/organizations/[^/]+/databases/[^/]+/branches/[^/]+(?:/keyspaces/[^/]+)?/config-changes)/([^/]+)$`)
)

// VitessParametersHook implements the VTGate, VTTablet, and MySQL parameter
// operations that the Terraform overlays define on top of the API's draft
// config change requests:
//   - GET .../vitess-parameters (VTGate) and GET .../keyspaces/{keyspace}/vitess-parameters
//     reduce the branch's parameter list to the overridden parameters of one
//     scope.
//   - PATCH (merge) and PUT (replace) on the same paths create a draft config
//     change per changed component, submit them together, and return their IDs.
//     A request without parameters changes nothing.
//   - GET .../config-changes/{ids} reports the combined state of those changes.
type VitessParametersHook struct{}

var _ sdkInitHook = (*VitessParametersHook)(nil)

func NewVitessParametersHook() *VitessParametersHook {
	return &VitessParametersHook{}
}

func (h *VitessParametersHook) SDKInit(baseURL string, client HTTPClient) (string, HTTPClient) {
	if client == nil {
		return baseURL, client
	}

	return baseURL, &vitessParametersClient{client: client}
}

type vitessParametersClient struct {
	client HTTPClient
}

// vitessParameterScope is either a branch's VTGates or, when keyspace is set,
// one keyspace's VTTablets and MySQL servers.
type vitessParameterScope struct {
	branchPath string
	keyspace   string
}

func (s vitessParameterScope) namespaces() []string {
	if s.keyspace == "" {
		return []string{"vtgate"}
	}
	return []string{"vttablet", "mysqld"}
}

func (s vitessParameterScope) configChangesPath() string {
	if s.keyspace == "" {
		return s.branchPath + "/config-changes"
	}
	return s.branchPath + "/keyspaces/" + s.keyspace + "/config-changes"
}

type vitessParameter struct {
	Name     string  `json:"name"`
	Value    *string `json:"value"`
	Override bool    `json:"override"`
}

type vitessParameterChange struct {
	namespace string
	options   map[string]any
}

type vitessParametersResponse struct {
	Parameters     map[string]map[string]string `json:"parameters"`
	ConfigChangeID string                       `json:"config_change_id,omitempty"`
}

func (c *vitessParametersClient) Do(req *http.Request) (*http.Response, error) {
	if req == nil || req.URL == nil {
		return c.client.Do(req)
	}

	if matches := vitessParametersPathPattern.FindStringSubmatch(req.URL.Path); matches != nil {
		scope := vitessParameterScope{branchPath: matches[1], keyspace: matches[2]}
		switch req.Method {
		case http.MethodGet:
			return c.read(req, scope)
		case http.MethodPatch:
			return c.apply(req, scope, false)
		case http.MethodPut:
			return c.apply(req, scope, true)
		}
	}

	if req.Method == http.MethodGet {
		if matches := vitessConfigChangePathPattern.FindStringSubmatch(req.URL.Path); matches != nil {
			return c.configChangeStatus(req, matches[1], strings.Split(matches[2], ","))
		}
	}

	return c.client.Do(req)
}

// read reports every override, including ones set outside Terraform, so a
// configuration that starts managing parameters plans their removal instead
// of resetting them unannounced.
func (c *vitessParametersClient) read(req *http.Request, scope vitessParameterScope) (*http.Response, error) {
	overrides, res, err := c.fetchVitessOverrides(req, scope)
	if err != nil || res != nil {
		return res, err
	}

	return syntheticJSONResponse(req, http.StatusOK, vitessParametersResponse{Parameters: overrides})
}

// apply sends the changes needed to reach the configured parameters. Merging
// (create) only sets configured parameters, because a new branch can inherit
// overrides from the branch it was created from; the next plan then shows
// them. Replacing (update) also resets overrides that the configuration no
// longer lists. The response echoes the
// configured parameters because they are what the API now stores for them.
// Without configured parameters, apply changes nothing and reports the current
// overrides.
func (c *vitessParametersClient) apply(req *http.Request, scope vitessParameterScope, replace bool) (*http.Response, error) {
	var body struct {
		Parameters map[string]map[string]string `json:"parameters"`
	}
	if err := decodeAndClose(req.Body, &body); err != nil {
		return nil, fmt.Errorf("decode Vitess parameters request: %w", err)
	}

	overrides, res, err := c.fetchVitessOverrides(req, scope)
	if err != nil || res != nil {
		return res, err
	}

	if body.Parameters == nil {
		return syntheticJSONResponse(req, http.StatusOK, vitessParametersResponse{
			Parameters:     overrides,
			ConfigChangeID: skippedVitessConfigChangeID,
		})
	}

	changes := vitessParameterChanges(overrides, body.Parameters, scope.namespaces(), replace)
	if len(changes) == 0 {
		return syntheticJSONResponse(req, http.StatusOK, vitessParametersResponse{
			Parameters:     body.Parameters,
			ConfigChangeID: skippedVitessConfigChangeID,
		})
	}

	ids := make([]string, 0, len(changes))
	for _, change := range changes {
		res, err := c.send(req, http.MethodPost, scope.configChangesPath(), map[string]any{
			"change_type": change.namespace,
			"options":     change.options,
		})
		if err != nil {
			c.cancelDrafts(req, scope, ids)
			return nil, err
		}
		if res.StatusCode != http.StatusCreated {
			c.cancelDrafts(req, scope, ids)
			return res, nil
		}

		var draft struct {
			ID string `json:"id"`
		}
		err = decodeAndClose(res.Body, &draft)
		if err == nil && draft.ID == "" {
			err = fmt.Errorf("response has no ID")
		}
		if err != nil {
			c.cancelDrafts(req, scope, ids)
			return nil, fmt.Errorf("decode %s config change draft: %w", change.namespace, err)
		}
		ids = append(ids, draft.ID)
	}

	res, err = c.send(req, http.MethodPost, scope.branchPath+"/config-changes/submit", map[string]any{"ids": ids})
	if err != nil {
		c.cancelDrafts(req, scope, ids)
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		c.cancelDrafts(req, scope, ids)
		return res, nil
	}
	drainAndClose(res.Body)

	return syntheticJSONResponse(req, http.StatusOK, vitessParametersResponse{
		Parameters:     body.Parameters,
		ConfigChangeID: strings.Join(ids, ","),
	})
}

// cancelDrafts removes drafts an apply created before failing. A branch or
// keyspace allows one unfinished change per component, so a leftover draft
// would block every later apply.
func (c *vitessParametersClient) cancelDrafts(req *http.Request, scope vitessParameterScope, ids []string) {
	for _, id := range ids {
		res, err := c.send(req, http.MethodDelete, scope.configChangesPath()+"/"+id, nil)
		if err == nil {
			drainAndClose(res.Body)
		}
	}
}

// configChangeStatus reports the combined state of the config changes an
// apply submitted. Queued changes count as completed: their values are saved
// and roll out in the branch's next maintenance window.
func (c *vitessParametersClient) configChangeStatus(req *http.Request, changesPath string, ids []string) (*http.Response, error) {
	id := strings.Join(ids, ",")
	if id == skippedVitessConfigChangeID {
		return syntheticJSONResponse(req, http.StatusOK, map[string]string{"id": id, "state": "completed"})
	}

	state := "completed"
	for _, changeID := range ids {
		res, err := c.send(req, http.MethodGet, changesPath+"/"+changeID, nil)
		if err != nil {
			return nil, err
		}
		if res.StatusCode != http.StatusOK {
			return res, nil
		}

		var change struct {
			State        string  `json:"state"`
			ChangeType   string  `json:"change_type"`
			ErrorMessage *string `json:"error_message"`
		}
		if err := decodeAndClose(res.Body, &change); err != nil {
			return nil, fmt.Errorf("decode config change %s: %w", changeID, err)
		}

		switch change.State {
		case "completed", "queued":
		case "failed":
			message := "the API did not provide an error"
			if change.ErrorMessage != nil && *change.ErrorMessage != "" {
				message = *change.ErrorMessage
			}
			return nil, fmt.Errorf("%s config change %s failed: %s", change.ChangeType, changeID, message)
		case "canceled":
			return nil, fmt.Errorf("%s config change %s was canceled before it applied", change.ChangeType, changeID)
		default:
			state = change.State
		}
	}

	return syntheticJSONResponse(req, http.StatusOK, map[string]string{"id": id, "state": state})
}

// fetchVitessOverrides reads the branch's parameter list and returns the
// scope's overridden parameters by namespace. A non-nil response is an API
// error the caller returns unchanged. Parameters are read after the branch or
// keyspace itself, so a 404 here is an error rather than a sign that the
// resource is gone, which would make Terraform drop it from state.
func (c *vitessParametersClient) fetchVitessOverrides(req *http.Request, scope vitessParameterScope) (map[string]map[string]string, *http.Response, error) {
	path := scope.branchPath + "/vitess-parameters"
	res, err := c.send(req, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	if res.StatusCode == http.StatusNotFound {
		drainAndClose(res.Body)
		return nil, nil, fmt.Errorf("GET %s returned HTTP 404; refusing to treat the branch as missing", path)
	}
	if res.StatusCode != http.StatusOK {
		return nil, res, nil
	}

	var catalog struct {
		VTGate    []vitessParameter                       `json:"vtgate"`
		Keyspaces map[string]map[string][]vitessParameter `json:"keyspaces"`
	}
	if err := decodeAndClose(res.Body, &catalog); err != nil {
		return nil, nil, fmt.Errorf("decode Vitess parameters response: %w", err)
	}

	if scope.keyspace == "" {
		return vitessParameterOverrides(map[string][]vitessParameter{"vtgate": catalog.VTGate}), nil, nil
	}

	keyspace, ok := catalog.Keyspaces[scope.keyspace]
	if !ok {
		return nil, nil, fmt.Errorf("GET %s did not list keyspace %q", path, scope.keyspace)
	}
	return vitessParameterOverrides(map[string][]vitessParameter{
		"vttablet": keyspace["vttablet"],
		"mysqld":   keyspace["mysqld"],
	}), nil, nil
}

func (c *vitessParametersClient) send(req *http.Request, method string, path string, body any) (*http.Response, error) {
	return sendDerivedRequest(c.client, req, method, path, body)
}

// sendDerivedRequest issues an API request derived from the SDK's request so
// it keeps the SDK's headers, including authorization, and its context.
func sendDerivedRequest(client HTTPClient, req *http.Request, method string, path string, body any) (*http.Response, error) {
	apiReq := req.Clone(req.Context())
	apiReq.Method = method
	apiReq.URL.Path = path
	apiReq.URL.RawPath = ""
	apiReq.URL.RawQuery = ""
	apiReq.Body = http.NoBody
	apiReq.GetBody = nil
	apiReq.ContentLength = 0
	if body != nil {
		if err := replaceRequestBody(apiReq, body); err != nil {
			return nil, err
		}
	}

	res, err := client.Do(apiReq)
	if err == nil && res == nil {
		err = fmt.Errorf("%s %s returned no response", method, path)
	}
	return res, err
}

// vitessParameterOverrides reports the overridden parameters by namespace.
// It is never nil, so a read with no overrides replaces stale state instead of
// leaving it in place.
func vitessParameterOverrides(current map[string][]vitessParameter) map[string]map[string]string {
	parameters := map[string]map[string]string{}
	for namespace, details := range current {
		for _, detail := range details {
			if !detail.Override || detail.Value == nil {
				continue
			}
			if parameters[namespace] == nil {
				parameters[namespace] = map[string]string{}
			}
			parameters[namespace][detail.Name] = *detail.Value
		}
	}

	return parameters
}

// vitessParameterChanges returns, per namespace, the options a config change
// must send: configured values that differ from the current overrides and,
// when replacing, a null for each override the configuration no longer lists.
// A configured value that only matches the current default is still sent, so
// it stays fixed when a resize changes the default.
func vitessParameterChanges(
	overrides map[string]map[string]string,
	desired map[string]map[string]string,
	namespaces []string,
	replace bool,
) []vitessParameterChange {
	var changes []vitessParameterChange
	for _, namespace := range namespaces {
		options := map[string]any{}
		for name, value := range desired[namespace] {
			if override, ok := overrides[namespace][name]; ok && override == value {
				continue
			}
			options[name] = value
		}
		if replace {
			for name := range overrides[namespace] {
				if _, ok := desired[namespace][name]; !ok {
					options[name] = nil
				}
			}
		}

		if len(options) > 0 {
			changes = append(changes, vitessParameterChange{namespace: namespace, options: options})
		}
	}
	return changes
}

func syntheticJSONResponse(req *http.Request, status int, value any) (*http.Response, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")

	return &http.Response{
		Status:        fmt.Sprintf("%d %s", status, http.StatusText(status)),
		StatusCode:    status,
		Header:        headers,
		Body:          io.NopCloser(bytes.NewReader(payload)),
		ContentLength: int64(len(payload)),
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Request:       req,
	}, nil
}
