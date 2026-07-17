package hooks

import (
	"bytes"
	"io"
	"net/http"
	"regexp"
)

var bouncerTerraformChangesPutPathPattern = regexp.MustCompile(`^/v1/organizations/[^/]+/databases/[^/]+/branches/[^/]+/bouncers/[^/]+/terraform-changes$`)

// PostgresBouncerNoContentSkipHook wraps the SDK HTTP client and turns the
// 204 returned by PUT .../bouncers/{bouncer}/terraform-changes when nothing
// changed into a synthetic completed resize request, since the generated
// create/update code always expects a response body.
type PostgresBouncerNoContentSkipHook struct{}

var _ sdkInitHook = (*PostgresBouncerNoContentSkipHook)(nil)

func NewPostgresBouncerNoContentSkipHook() *PostgresBouncerNoContentSkipHook {
	return &PostgresBouncerNoContentSkipHook{}
}

func (h *PostgresBouncerNoContentSkipHook) SDKInit(baseURL string, client HTTPClient) (string, HTTPClient) {
	if client == nil {
		return baseURL, client
	}

	return baseURL, &postgresBouncerNoContentSkipClient{client: client}
}

type postgresBouncerNoContentSkipClient struct {
	client HTTPClient
}

func (c *postgresBouncerNoContentSkipClient) Do(req *http.Request) (*http.Response, error) {
	res, err := c.client.Do(req)
	if err != nil || res == nil || req == nil || req.URL == nil {
		return res, err
	}

	if req.Method != http.MethodPut ||
		res.StatusCode != http.StatusNoContent ||
		!bouncerTerraformChangesPutPathPattern.MatchString(req.URL.Path) {
		return res, nil
	}

	drainAndClose(res.Body)

	return syntheticBouncerResizeResponse(req), nil
}

func syntheticBouncerResizeResponse(req *http.Request) *http.Response {
	payload := []byte(`{"id":"__speakeasy_skipped_bouncer_resize__","state":"completed"}`)

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
