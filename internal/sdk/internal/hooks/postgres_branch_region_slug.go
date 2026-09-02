package hooks

import (
	"bytes"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

var postgresBranchGetPathPattern = regexp.MustCompile(`^/v1/organizations/[^/]+/databases/[^/]+/branches/[^/]+$`)

// PostgresBranchRegionSlugHook copies the nested region slug to a synthetic
// top-level field on PostgreSQL branch GET responses. The generated Terraform
// resource maps that field to the region attribute used by create and import.
type PostgresBranchRegionSlugHook struct{}

var _ sdkInitHook = (*PostgresBranchRegionSlugHook)(nil)

func NewPostgresBranchRegionSlugHook() *PostgresBranchRegionSlugHook {
	return &PostgresBranchRegionSlugHook{}
}

func (h *PostgresBranchRegionSlugHook) SDKInit(baseURL string, client HTTPClient) (string, HTTPClient) {
	if client == nil {
		return baseURL, client
	}

	return baseURL, &postgresBranchRegionSlugClient{client: client}
}

type postgresBranchRegionSlugClient struct {
	client HTTPClient
}

func (c *postgresBranchRegionSlugClient) Do(req *http.Request) (*http.Response, error) {
	res, err := c.client.Do(req)
	if err != nil || res == nil || req == nil || req.URL == nil {
		return res, err
	}

	if req.Method != http.MethodGet || res.StatusCode != http.StatusOK || !postgresBranchGetPathPattern.MatchString(req.URL.Path) {
		return res, nil
	}

	body, readErr := io.ReadAll(res.Body)
	drainAndClose(res.Body)
	if readErr != nil {
		return res, readErr
	}

	if rewritten, ok := hoistRegionSlug(body); ok {
		body = rewritten
	}

	res.Body = io.NopCloser(bytes.NewReader(body))
	res.ContentLength = int64(len(body))
	if res.Header == nil {
		res.Header = make(http.Header)
	}
	res.Header.Set("Content-Length", strconv.Itoa(len(body)))

	return res, nil
}
