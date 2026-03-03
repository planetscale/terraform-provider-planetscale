package hooks

import (
	"bytes"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

var (
	branchChangePatchPathPattern = regexp.MustCompile(`^/v1/organizations/([^/]+)/databases/([^/]+)/branches/([^/]+)/changes$`)
	branchChangeGetPathPattern   = regexp.MustCompile(`^/v1/organizations/([^/]+)/databases/([^/]+)/branches/([^/]+)/changes/([^/]+)$`)
)

// PostgresBranchNoContentSkipHook wraps the SDK HTTP client and short-circuits
// branch change polling when the previous PATCH /changes call returns 204.
type PostgresBranchNoContentSkipHook struct {
	mu           sync.Mutex
	pendingSkips map[string]struct{}
}

var _ sdkInitHook = (*PostgresBranchNoContentSkipHook)(nil)

func NewPostgresBranchNoContentSkipHook() *PostgresBranchNoContentSkipHook {
	return &PostgresBranchNoContentSkipHook{
		pendingSkips: map[string]struct{}{},
	}
}

func (h *PostgresBranchNoContentSkipHook) SDKInit(baseURL string, client HTTPClient) (string, HTTPClient) {
	if client == nil {
		return baseURL, client
	}

	return baseURL, &postgresBranchNoContentSkipClient{
		client: client,
		state:  h,
	}
}

func (h *PostgresBranchNoContentSkipHook) markSkip(scope string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.pendingSkips[scope] = struct{}{}
}

func (h *PostgresBranchNoContentSkipHook) consumeSkip(scope string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.pendingSkips[scope]; !ok {
		return false
	}

	delete(h.pendingSkips, scope)
	return true
}

type postgresBranchNoContentSkipClient struct {
	client HTTPClient
	state  *PostgresBranchNoContentSkipHook
}

func (c *postgresBranchNoContentSkipClient) Do(req *http.Request) (*http.Response, error) {
	if req == nil || req.URL == nil {
		return c.client.Do(req)
	}

	if req.Method == http.MethodGet {
		if organization, database, branch, _, ok := parseBranchChangeGetPath(req.URL.Path); ok {
			if c.state.consumeSkip(branchScope(organization, database, branch)) {
				return syntheticBranchChangeResponse(req), nil
			}
		}
	}

	res, err := c.client.Do(req)
	if err != nil || res == nil {
		return res, err
	}

	if req.Method != http.MethodPatch {
		return res, nil
	}

	organization, database, branch, ok := parseBranchChangePatchPath(req.URL.Path)
	if !ok || res.StatusCode != http.StatusNoContent {
		return res, nil
	}

	drainAndClose(res.Body)
	c.state.markSkip(branchScope(organization, database, branch))

	return syntheticBranchChangeResponse(req), nil
}

func parseBranchChangePatchPath(path string) (string, string, string, bool) {
	matches := branchChangePatchPathPattern.FindStringSubmatch(path)
	if len(matches) != 4 {
		return "", "", "", false
	}
	return matches[1], matches[2], matches[3], true
}

func parseBranchChangeGetPath(path string) (string, string, string, string, bool) {
	matches := branchChangeGetPathPattern.FindStringSubmatch(path)
	if len(matches) != 5 {
		return "", "", "", "", false
	}
	return matches[1], matches[2], matches[3], matches[4], true
}

func branchScope(organization string, database string, branch string) string {
	return strings.Join([]string{organization, database, branch}, "/")
}

func syntheticBranchChangeResponse(req *http.Request) *http.Response {
	payload := []byte(`{"id":"__speakeasy_skipped_change_request__","state":"completed"}`)

	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")

	return &http.Response{
		Status:        "200 OK",
		StatusCode:    http.StatusOK,
		Header:        headers,
		Body:          io.NopCloser(bytes.NewReader(payload)),
		ContentLength: int64(len(payload)),
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Request:       req,
	}
}

func drainAndClose(body io.ReadCloser) {
	if body == nil {
		return
	}

	_, _ = io.Copy(io.Discard, body)
	_ = body.Close()
}
