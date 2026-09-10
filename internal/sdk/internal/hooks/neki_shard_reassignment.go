package hooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const reassignNekiShardOperationID = "reassign_neki_shard"

// NekiShardReassignmentHook adapts Terraform's singular shard update to the
// API's bulk shard-assignment request and response.
type NekiShardReassignmentHook struct{}

var _ beforeRequestHook = (*NekiShardReassignmentHook)(nil)
var _ afterSuccessHook = (*NekiShardReassignmentHook)(nil)

func NewNekiShardReassignmentHook() *NekiShardReassignmentHook {
	return &NekiShardReassignmentHook{}
}

func (h *NekiShardReassignmentHook) BeforeRequest(
	hookCtx BeforeRequestContext,
	req *http.Request,
) (*http.Request, error) {
	if hookCtx.OperationID != reassignNekiShardOperationID || req == nil {
		return req, nil
	}

	var singular struct {
		ID string `json:"id"`
	}
	if err := decodeAndClose(req.Body, &singular); err != nil {
		return req, fmt.Errorf("decode Neki shard reassignment request: %w", err)
	}
	if singular.ID == "" {
		return req, fmt.Errorf("neki shard reassignment request is missing the shard ID")
	}

	return req, replaceRequestBody(req, struct {
		ShardIDs []string `json:"shard_ids"`
	}{
		ShardIDs: []string{singular.ID},
	})
}

func (h *NekiShardReassignmentHook) AfterSuccess(
	hookCtx AfterSuccessContext,
	res *http.Response,
) (*http.Response, error) {
	if hookCtx.OperationID != reassignNekiShardOperationID ||
		res == nil ||
		res.StatusCode != http.StatusOK {
		return res, nil
	}

	var assignments []struct {
		ID     string          `json:"id"`
		Status string          `json:"status"`
		Error  json.RawMessage `json:"error"`
	}
	if err := decodeAndClose(res.Body, &assignments); err != nil {
		return res, fmt.Errorf("decode Neki shard reassignment response: %w", err)
	}
	if len(assignments) != 1 {
		return res, fmt.Errorf(
			"expected one Neki shard reassignment result, received %d",
			len(assignments),
		)
	}

	assignment := assignments[0]
	if assignment.Status == "failed" {
		return res, fmt.Errorf(
			"neki shard %q reassignment failed: %s",
			assignment.ID,
			assignmentError(assignment.Error),
		)
	}
	if assignment.Status != "assigned" && assignment.Status != "unchanged" {
		return res, fmt.Errorf(
			"neki shard %q reassignment returned unexpected status %q",
			assignment.ID,
			assignment.Status,
		)
	}

	return res, replaceResponseBody(res, struct {
		ID string `json:"id"`
	}{
		ID: assignment.ID,
	})
}

func decodeAndClose(body io.ReadCloser, value any) error {
	if body == nil {
		return fmt.Errorf("body is empty")
	}
	defer func() {
		_ = body.Close()
	}()

	return json.NewDecoder(body).Decode(value)
}

func replaceRequestBody(req *http.Request, value any) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}

	req.Body = io.NopCloser(bytes.NewReader(payload))
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(payload)), nil
	}
	req.ContentLength = int64(len(payload))
	req.Header.Set("Content-Type", "application/json")
	return nil
}

func replaceResponseBody(res *http.Response, value any) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}

	res.Body = io.NopCloser(bytes.NewReader(payload))
	res.ContentLength = int64(len(payload))
	res.Header.Set("Content-Length", fmt.Sprintf("%d", len(payload)))
	return nil
}

func assignmentError(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return "the API did not provide an error"
	}

	var detail struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if json.Unmarshal(raw, &detail) == nil {
		switch {
		case detail.Code != "" && detail.Message != "":
			return detail.Code + ": " + detail.Message
		case detail.Message != "":
			return detail.Message
		case detail.Code != "":
			return detail.Code
		}
	}

	return string(raw)
}
