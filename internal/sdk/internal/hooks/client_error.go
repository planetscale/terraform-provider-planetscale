package hooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// maxClientErrorBodyLength bounds how much of a non-JSON error body is echoed
// back to the user.
const maxClientErrorBodyLength = 512

// ClientErrorHook converts 4xx responses (other than 404) into errors that carry
// the PlanetScale API's `code` and `message`, so Terraform reports why a request
// was rejected instead of a raw request/response dump.
//
// The Terraform generator treats every status code declared in the OpenAPI
// document as a regular response and, for responses without a declared body,
// discards the body before the provider can report it. 404 is deliberately left
// alone so read operations keep detecting deleted resources by status code.
type ClientErrorHook struct{}

var _ afterSuccessHook = (*ClientErrorHook)(nil)

func NewClientErrorHook() *ClientErrorHook {
	return &ClientErrorHook{}
}

// ClientError is the error returned for a 4xx API response.
type ClientError struct {
	StatusCode int
	Code       string
	Message    string
	Body       string
}

func (e *ClientError) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "PlanetScale API returned HTTP %d", e.StatusCode)
	if e.Code != "" {
		fmt.Fprintf(&b, " (%s)", e.Code)
	}
	switch {
	case e.Message != "":
		b.WriteString(": " + e.Message)
	case e.Body != "":
		b.WriteString(": " + e.Body)
	}
	return b.String()
}

func (h *ClientErrorHook) AfterSuccess(_ AfterSuccessContext, res *http.Response) (*http.Response, error) {
	if res == nil || res.StatusCode < 400 || res.StatusCode >= 500 || res.StatusCode == http.StatusNotFound {
		return res, nil
	}

	clientErr := &ClientError{StatusCode: res.StatusCode}

	if res.Body != nil {
		raw, err := io.ReadAll(res.Body)
		_ = res.Body.Close()
		if err != nil {
			return res, nil
		}
		res.Body = io.NopCloser(bytes.NewReader(raw))

		var payload struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		if json.Unmarshal(raw, &payload) == nil && (payload.Code != "" || payload.Message != "") {
			clientErr.Code = payload.Code
			clientErr.Message = payload.Message
		} else {
			clientErr.Body = truncate(strings.TrimSpace(string(raw)), maxClientErrorBodyLength)
		}
	}

	return res, clientErr
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
