package hooks

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/planetscale/terraform-provider-planetscale/internal/redact"
)

// RedactedRolePasswordHook strips generated passwords from responses for the
// postgres redacted branch role. That resource exists so the password never
// enters Terraform; the API still returns it, and generated error reporting
// includes the raw response body.
type RedactedRolePasswordHook struct{}

var _ afterSuccessHook = (*RedactedRolePasswordHook)(nil)

func NewRedactedRolePasswordHook() *RedactedRolePasswordHook {
	return &RedactedRolePasswordHook{}
}

func (h *RedactedRolePasswordHook) AfterSuccess(ctx AfterSuccessContext, res *http.Response) (*http.Response, error) {
	if res == nil || res.Body == nil || !redactedRoleOperation(ctx.OperationID) {
		return res, nil
	}

	raw, err := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if err != nil {
		res.Body = io.NopCloser(bytes.NewReader(nil))
		return res, nil
	}

	redacted := []byte(redact.RedactCredentialFields(string(raw)))
	res.Body = io.NopCloser(bytes.NewReader(redacted))
	res.ContentLength = int64(len(redacted))
	if res.Header != nil && res.Header.Get("Content-Length") != "" {
		res.Header.Set("Content-Length", strconv.Itoa(len(redacted)))
	}
	return res, nil
}

func redactedRoleOperation(operationID string) bool {
	switch operationID {
	case "create_redacted_role", "get_redacted_role", "update_redacted_role":
		return true
	default:
		return false
	}
}
