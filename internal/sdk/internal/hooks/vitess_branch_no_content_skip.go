package hooks

import (
	"bytes"
	"io"
	"net/http"
	"regexp"
)

const skippedResizeRequestID = "__speakeasy_skipped_resize_request__"

var (
	branchResizePutPathPattern = regexp.MustCompile(`^/v1/organizations/[^/]+/databases/[^/]+/branches/[^/]+/resizes$`)
	branchResizeGetPathPattern = regexp.MustCompile(`^/v1/organizations/[^/]+/databases/[^/]+/branches/[^/]+/resizes/([^/]+)$`)
)

// VitessBranchNoContentSkipHook short-circuits resize polling when the
// preceding PUT /resizes call reports that no changes are needed.
type VitessBranchNoContentSkipHook struct{}

var _ sdkInitHook = (*VitessBranchNoContentSkipHook)(nil)

func NewVitessBranchNoContentSkipHook() *VitessBranchNoContentSkipHook {
	return &VitessBranchNoContentSkipHook{}
}

func (h *VitessBranchNoContentSkipHook) SDKInit(baseURL string, client HTTPClient) (string, HTTPClient) {
	if client == nil {
		return baseURL, client
	}

	return baseURL, &vitessBranchNoContentSkipClient{client: client}
}

type vitessBranchNoContentSkipClient struct {
	client HTTPClient
}

func (c *vitessBranchNoContentSkipClient) Do(req *http.Request) (*http.Response, error) {
	if req == nil || req.URL == nil {
		return c.client.Do(req)
	}

	if req.Method == http.MethodGet && isSkippedResizeRequestPath(req.URL.Path) {
		return syntheticResizeResponse(req), nil
	}

	res, err := c.client.Do(req)
	if err != nil || res == nil {
		return res, err
	}

	if req.Method != http.MethodPut || !branchResizePutPathPattern.MatchString(req.URL.Path) || res.StatusCode != http.StatusNoContent {
		return res, nil
	}

	drainAndClose(res.Body)
	return syntheticResizeResponse(req), nil
}

func isSkippedResizeRequestPath(path string) bool {
	matches := branchResizeGetPathPattern.FindStringSubmatch(path)
	return len(matches) == 2 && matches[1] == skippedResizeRequestID
}

func syntheticResizeResponse(req *http.Request) *http.Response {
	payload := []byte(`{"id":"` + skippedResizeRequestID + `","state":"completed"}`)

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
