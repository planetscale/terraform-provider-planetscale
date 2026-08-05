package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/planetscale/terraform-provider-planetscale/internal/sdk/internal/hooks"
	"github.com/planetscale/terraform-provider-planetscale/internal/sdk/internal/utils"
	sdkerrors "github.com/planetscale/terraform-provider-planetscale/internal/sdk/models/errors"
)

// VitessBranchSettings is the stable subset of branch state used by the
// custom Terraform settings resources.
type VitessBranchSettings struct {
	ID             string `json:"id"`
	SafeMigrations bool   `json:"safe_migrations"`
	VTGateSize     string `json:"vtgate_size"`
	VTGateCount    int64  `json:"vtgate_count"`
}

// VitessBranchResizeRequest represents the VTGate configuration captured by a
// branch resize request.
type VitessBranchResizeRequest struct {
	ID                                 string `json:"id"`
	State                              string `json:"state"`
	VTGateSize                         string `json:"vtgate_size"`
	VTGateName                         string `json:"vtgate_name"`
	PreviousVTGateSize                 string `json:"previous_vtgate_size"`
	PreviousVTGateName                 string `json:"previous_vtgate_name"`
	VTGateCount                        int64  `json:"vtgate_count"`
	PreviousVTGateCount                int64  `json:"previous_vtgate_count"`
	VTGateMaxCount                     *int64 `json:"vtgate_max_count"`
	PreviousVTGateMaxCount             *int64 `json:"previous_vtgate_max_count"`
	VTGateAutoscaling                  bool   `json:"vtgate_autoscaling"`
	PreviousVTGateAutoscaling          bool   `json:"previous_vtgate_autoscaling"`
	VTGateTargetCPUUtilization         *int64 `json:"vtgate_target_cpu_utilization"`
	PreviousVTGateTargetCPUUtilization *int64 `json:"previous_vtgate_target_cpu_utilization"`
}

// UpdateVitessBranchVTGateConfigurationRequest contains the optional VTGate
// values accepted by the branch resize API.
type UpdateVitessBranchVTGateConfigurationRequest struct {
	VTGateSize                 *string `json:"vtgate_size,omitempty"`
	VTGateCount                *int64  `json:"vtgate_count,omitempty"`
	VTGateMaxCount             *int64  `json:"vtgate_max_count,omitempty"`
	VTGateAutoscaling          *bool   `json:"vtgate_autoscaling,omitempty"`
	VTGateTargetCPUUtilization *int64  `json:"vtgate_target_cpu_utilization,omitempty"`
}

type vitessBranchResizeList struct {
	Data []VitessBranchResizeRequest `json:"data"`
}

// GetVitessBranchSettings reads the current branch settings.
func (s *DatabaseBranches) GetVitessBranchSettings(ctx context.Context, organization, database, branch string) (*VitessBranchSettings, error) {
	var settings VitessBranchSettings
	if err := s.doCustomBranchRequest(ctx, http.MethodGet, organization, database, branch, "", nil, &settings, "get_vitess_branch_settings"); err != nil {
		return nil, err
	}
	return &settings, nil
}

// SetVitessBranchSafeMigrations enables or disables safe migrations.
func (s *DatabaseBranches) SetVitessBranchSafeMigrations(ctx context.Context, organization, database, branch string, enabled bool) (*VitessBranchSettings, error) {
	method := http.MethodDelete
	if enabled {
		method = http.MethodPost
	}

	var settings VitessBranchSettings
	if err := s.doCustomBranchRequest(ctx, method, organization, database, branch, "safe-migrations", nil, &settings, "set_vitess_branch_safe_migrations"); err != nil {
		return nil, err
	}
	return &settings, nil
}

// ListVitessBranchResizeRequests returns resize requests newest first.
func (s *DatabaseBranches) ListVitessBranchResizeRequests(ctx context.Context, organization, database, branch string) ([]VitessBranchResizeRequest, error) {
	var response vitessBranchResizeList
	if err := s.doCustomBranchRequest(ctx, http.MethodGet, organization, database, branch, "resizes", nil, &response, "list_vitess_branch_resize_requests"); err != nil {
		return nil, err
	}
	return response.Data, nil
}

// UpdateVitessBranchVTGateConfiguration queues a VTGate configuration change.
func (s *DatabaseBranches) UpdateVitessBranchVTGateConfiguration(ctx context.Context, organization, database, branch string, update UpdateVitessBranchVTGateConfigurationRequest) (*VitessBranchResizeRequest, error) {
	var resize VitessBranchResizeRequest
	if err := s.doCustomBranchRequest(ctx, http.MethodPut, organization, database, branch, "resizes", update, &resize, "update_vitess_branch_vtgate_configuration"); err != nil {
		return nil, err
	}
	return &resize, nil
}

func (s *DatabaseBranches) doCustomBranchRequest(
	ctx context.Context,
	method string,
	organization string,
	database string,
	branch string,
	suffix string,
	body any,
	result any,
	operationID string,
) error {
	baseURL := utils.ReplaceParameters(s.sdkConfiguration.GetServerDetails())
	path := fmt.Sprintf(
		"/organizations/%s/databases/%s/branches/%s",
		url.PathEscape(organization),
		url.PathEscape(database),
		url.PathEscape(branch),
	)
	if suffix != "" {
		path += "/" + suffix
	}

	var bodyReader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal %s request: %w", operationID, err)
		}
		bodyReader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, strings.TrimSuffix(baseURL, "/")+path, bodyReader)
	if err != nil {
		return fmt.Errorf("create %s request: %w", operationID, err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", s.sdkConfiguration.UserAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if err := utils.PopulateSecurity(ctx, req, s.sdkConfiguration.Security); err != nil {
		return err
	}

	hookContext := hooks.HookContext{
		SDK:              s.rootSDK,
		SDKConfiguration: s.sdkConfiguration,
		BaseURL:          baseURL,
		Context:          ctx,
		OperationID:      operationID,
		SecuritySource:   s.sdkConfiguration.Security,
	}
	req, err = s.hooks.BeforeRequest(hooks.BeforeRequestContext{HookContext: hookContext}, req)
	if err != nil {
		return err
	}

	response, err := s.sdkConfiguration.Client.Do(req)
	if err != nil || response == nil {
		if err == nil {
			err = fmt.Errorf("no response")
		}
		_, hookErr := s.hooks.AfterError(hooks.AfterErrorContext{HookContext: hookContext}, response, err)
		if hookErr != nil {
			return hookErr
		}
		return fmt.Errorf("send %s request: %w", operationID, err)
	}

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		response, err = s.hooks.AfterError(hooks.AfterErrorContext{HookContext: hookContext}, response, nil)
		if err != nil {
			return err
		}
		rawBody, readErr := utils.ConsumeRawBody(response)
		if readErr != nil {
			return readErr
		}
		return sdkerrors.NewAPIError(http.StatusText(response.StatusCode), response.StatusCode, string(rawBody), response)
	}

	response, err = s.hooks.AfterSuccess(hooks.AfterSuccessContext{HookContext: hookContext}, response)
	if err != nil {
		return err
	}
	if result == nil || response.StatusCode == http.StatusNoContent {
		utils.DrainBody(response)
		return nil
	}

	rawBody, err := utils.ConsumeRawBody(response)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(rawBody, result); err != nil {
		return fmt.Errorf("decode %s response: %w", operationID, err)
	}
	return nil
}
