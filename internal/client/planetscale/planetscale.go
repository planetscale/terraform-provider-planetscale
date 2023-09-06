package planetscale

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Client struct {
	httpCl  *http.Client
	baseURL *url.URL
}

func NewClient(httpCl *http.Client, baseURL *url.URL) *Client {
	if baseURL == nil {
		baseURL = &url.URL{Scheme: "https", Host: "api.planetscale.com", Path: "/v1"}
	}
	if !strings.HasSuffix(baseURL.Path, "/") {
		baseURL.Path = baseURL.Path + "/"
	}
	return &Client{httpCl: httpCl, baseURL: baseURL}
}

type ListOrganizationsRes500 struct{}
type ListOrganizationsRes200_DataItem_Features struct {
	Insights      *bool `json:"insights,omitempty" tfsdk:"insights"`
	SingleTenancy *bool `json:"single_tenancy,omitempty" tfsdk:"single_tenancy"`
	Sso           *bool `json:"sso,omitempty" tfsdk:"sso"`
}
type ListOrganizationsRes200_DataItem_Flags struct {
	ExampleFlag *string `json:"example_flag,omitempty" tfsdk:"example_flag"`
}
type ListOrganizationsRes200_DataItem struct {
	AdminOnlyProductionAccess bool                                       `json:"admin_only_production_access" tfsdk:"admin_only_production_access"`
	BillingEmail              *string                                    `json:"billing_email,omitempty" tfsdk:"billing_email"`
	CanCreateDatabases        bool                                       `json:"can_create_databases" tfsdk:"can_create_databases"`
	CreatedAt                 string                                     `json:"created_at" tfsdk:"created_at"`
	DatabaseCount             float64                                    `json:"database_count" tfsdk:"database_count"`
	Features                  *ListOrganizationsRes200_DataItem_Features `json:"features,omitempty" tfsdk:"features"`
	Flags                     *ListOrganizationsRes200_DataItem_Flags    `json:"flags,omitempty" tfsdk:"flags"`
	FreeDatabasesRemaining    float64                                    `json:"free_databases_remaining" tfsdk:"free_databases_remaining"`
	HasPastDueInvoices        bool                                       `json:"has_past_due_invoices" tfsdk:"has_past_due_invoices"`
	Id                        string                                     `json:"id" tfsdk:"id"`
	Name                      string                                     `json:"name" tfsdk:"name"`
	Plan                      string                                     `json:"plan" tfsdk:"plan"`
	SingleTenancy             bool                                       `json:"single_tenancy" tfsdk:"single_tenancy"`
	SleepingDatabaseCount     float64                                    `json:"sleeping_database_count" tfsdk:"sleeping_database_count"`
	Sso                       bool                                       `json:"sso" tfsdk:"sso"`
	SsoDirectory              bool                                       `json:"sso_directory" tfsdk:"sso_directory"`
	SsoPortalUrl              *string                                    `json:"sso_portal_url,omitempty" tfsdk:"sso_portal_url"`
	UpdatedAt                 string                                     `json:"updated_at" tfsdk:"updated_at"`
	ValidBillingInfo          bool                                       `json:"valid_billing_info" tfsdk:"valid_billing_info"`
}
type ListOrganizationsRes200 struct {
	Data []ListOrganizationsRes200_DataItem `json:"data" tfsdk:"data"`
}
type ListOrganizationsRes404 struct{}
type ListOrganizationsRes403 struct{}

func (cl *Client) ListOrganizations(ctx context.Context, page *int, perPage *int) (res200 *ListOrganizationsRes200, res403 *ListOrganizationsRes403, res404 *ListOrganizationsRes404, res500 *ListOrganizationsRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations"})
	q := u.Query()
	if page != nil {
		q.Set("page", strconv.Itoa(*page))
	}
	if perPage != nil {
		q.Set("per_page", strconv.Itoa(*perPage))
	}
	u.RawQuery = q.Encode()
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListOrganizationsRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(ListOrganizationsRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(ListOrganizationsRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(ListOrganizationsRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type GetOrganizationRes404 struct{}
type GetOrganizationRes403 struct{}
type GetOrganizationRes500 struct{}
type GetOrganizationRes200_Features struct {
	Insights      *bool `json:"insights,omitempty" tfsdk:"insights"`
	SingleTenancy *bool `json:"single_tenancy,omitempty" tfsdk:"single_tenancy"`
	Sso           *bool `json:"sso,omitempty" tfsdk:"sso"`
}
type GetOrganizationRes200_Flags struct {
	ExampleFlag *string `json:"example_flag,omitempty" tfsdk:"example_flag"`
}
type GetOrganizationRes200 struct {
	AdminOnlyProductionAccess bool                            `json:"admin_only_production_access" tfsdk:"admin_only_production_access"`
	BillingEmail              *string                         `json:"billing_email,omitempty" tfsdk:"billing_email"`
	CanCreateDatabases        bool                            `json:"can_create_databases" tfsdk:"can_create_databases"`
	CreatedAt                 string                          `json:"created_at" tfsdk:"created_at"`
	DatabaseCount             float64                         `json:"database_count" tfsdk:"database_count"`
	Features                  *GetOrganizationRes200_Features `json:"features,omitempty" tfsdk:"features"`
	Flags                     *GetOrganizationRes200_Flags    `json:"flags,omitempty" tfsdk:"flags"`
	FreeDatabasesRemaining    float64                         `json:"free_databases_remaining" tfsdk:"free_databases_remaining"`
	HasPastDueInvoices        bool                            `json:"has_past_due_invoices" tfsdk:"has_past_due_invoices"`
	Id                        string                          `json:"id" tfsdk:"id"`
	Name                      string                          `json:"name" tfsdk:"name"`
	Plan                      string                          `json:"plan" tfsdk:"plan"`
	SingleTenancy             bool                            `json:"single_tenancy" tfsdk:"single_tenancy"`
	SleepingDatabaseCount     float64                         `json:"sleeping_database_count" tfsdk:"sleeping_database_count"`
	Sso                       bool                            `json:"sso" tfsdk:"sso"`
	SsoDirectory              bool                            `json:"sso_directory" tfsdk:"sso_directory"`
	SsoPortalUrl              *string                         `json:"sso_portal_url,omitempty" tfsdk:"sso_portal_url"`
	UpdatedAt                 string                          `json:"updated_at" tfsdk:"updated_at"`
	ValidBillingInfo          bool                            `json:"valid_billing_info" tfsdk:"valid_billing_info"`
}

func (cl *Client) GetOrganization(ctx context.Context, name string) (res200 *GetOrganizationRes200, res403 *GetOrganizationRes403, res404 *GetOrganizationRes404, res500 *GetOrganizationRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + name})
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(GetOrganizationRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(GetOrganizationRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(GetOrganizationRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(GetOrganizationRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type UpdateOrganizationReq struct {
	BillingEmail                    *string `json:"billing_email,omitempty" tfsdk:"billing_email"`
	IdpManagedRoles                 *bool   `json:"idp_managed_roles,omitempty" tfsdk:"idp_managed_roles"`
	RequireAdminForProductionAccess *bool   `json:"require_admin_for_production_access,omitempty" tfsdk:"require_admin_for_production_access"`
}
type UpdateOrganizationRes200_Features struct {
	Insights      *bool `json:"insights,omitempty" tfsdk:"insights"`
	SingleTenancy *bool `json:"single_tenancy,omitempty" tfsdk:"single_tenancy"`
	Sso           *bool `json:"sso,omitempty" tfsdk:"sso"`
}
type UpdateOrganizationRes200_Flags struct {
	ExampleFlag *string `json:"example_flag,omitempty" tfsdk:"example_flag"`
}
type UpdateOrganizationRes200 struct {
	AdminOnlyProductionAccess bool                               `json:"admin_only_production_access" tfsdk:"admin_only_production_access"`
	BillingEmail              *string                            `json:"billing_email,omitempty" tfsdk:"billing_email"`
	CanCreateDatabases        bool                               `json:"can_create_databases" tfsdk:"can_create_databases"`
	CreatedAt                 string                             `json:"created_at" tfsdk:"created_at"`
	DatabaseCount             float64                            `json:"database_count" tfsdk:"database_count"`
	Features                  *UpdateOrganizationRes200_Features `json:"features,omitempty" tfsdk:"features"`
	Flags                     *UpdateOrganizationRes200_Flags    `json:"flags,omitempty" tfsdk:"flags"`
	FreeDatabasesRemaining    float64                            `json:"free_databases_remaining" tfsdk:"free_databases_remaining"`
	HasPastDueInvoices        bool                               `json:"has_past_due_invoices" tfsdk:"has_past_due_invoices"`
	Id                        string                             `json:"id" tfsdk:"id"`
	Name                      string                             `json:"name" tfsdk:"name"`
	Plan                      string                             `json:"plan" tfsdk:"plan"`
	SingleTenancy             bool                               `json:"single_tenancy" tfsdk:"single_tenancy"`
	SleepingDatabaseCount     float64                            `json:"sleeping_database_count" tfsdk:"sleeping_database_count"`
	Sso                       bool                               `json:"sso" tfsdk:"sso"`
	SsoDirectory              bool                               `json:"sso_directory" tfsdk:"sso_directory"`
	SsoPortalUrl              *string                            `json:"sso_portal_url,omitempty" tfsdk:"sso_portal_url"`
	UpdatedAt                 string                             `json:"updated_at" tfsdk:"updated_at"`
	ValidBillingInfo          bool                               `json:"valid_billing_info" tfsdk:"valid_billing_info"`
}
type UpdateOrganizationRes404 struct{}
type UpdateOrganizationRes403 struct{}
type UpdateOrganizationRes500 struct{}

func (cl *Client) UpdateOrganization(ctx context.Context, name string, req UpdateOrganizationReq) (res200 *UpdateOrganizationRes200, res403 *UpdateOrganizationRes403, res404 *UpdateOrganizationRes404, res500 *UpdateOrganizationRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + name})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res200, res403, res404, res500, err
	}
	r, err := http.NewRequestWithContext(ctx, "PATCH", u.String(), body)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(UpdateOrganizationRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(UpdateOrganizationRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(UpdateOrganizationRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(UpdateOrganizationRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type ListRegionsForOrganizationRes404 struct{}
type ListRegionsForOrganizationRes403 struct{}
type ListRegionsForOrganizationRes500 struct{}
type ListRegionsForOrganizationRes200_DataItem struct {
	DisplayName       string   `json:"display_name" tfsdk:"display_name"`
	Enabled           bool     `json:"enabled" tfsdk:"enabled"`
	Id                string   `json:"id" tfsdk:"id"`
	Location          string   `json:"location" tfsdk:"location"`
	Provider          string   `json:"provider" tfsdk:"provider"`
	PublicIpAddresses []string `json:"public_ip_addresses" tfsdk:"public_ip_addresses"`
	Slug              string   `json:"slug" tfsdk:"slug"`
}
type ListRegionsForOrganizationRes200 struct {
	Data []ListRegionsForOrganizationRes200_DataItem `json:"data" tfsdk:"data"`
}

func (cl *Client) ListRegionsForOrganization(ctx context.Context, name string, page *int, perPage *int) (res200 *ListRegionsForOrganizationRes200, res403 *ListRegionsForOrganizationRes403, res404 *ListRegionsForOrganizationRes404, res500 *ListRegionsForOrganizationRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + name + "/region"})
	q := u.Query()
	if page != nil {
		q.Set("page", strconv.Itoa(*page))
	}
	if perPage != nil {
		q.Set("per_page", strconv.Itoa(*perPage))
	}
	u.RawQuery = q.Encode()
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListRegionsForOrganizationRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(ListRegionsForOrganizationRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(ListRegionsForOrganizationRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(ListRegionsForOrganizationRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type ListDatabasesRes404 struct{}
type ListDatabasesRes403 struct{}
type ListDatabasesRes500 struct{}
type ListDatabasesRes200_DataItem_DataImport_DataSource struct {
	Database string `json:"database" tfsdk:"database"`
	Hostname string `json:"hostname" tfsdk:"hostname"`
	Port     string `json:"port" tfsdk:"port"`
}
type ListDatabasesRes200_DataItem_DataImport struct {
	DataSource        ListDatabasesRes200_DataItem_DataImport_DataSource `json:"data_source" tfsdk:"data_source"`
	FinishedAt        string                                             `json:"finished_at" tfsdk:"finished_at"`
	ImportCheckErrors string                                             `json:"import_check_errors" tfsdk:"import_check_errors"`
	StartedAt         string                                             `json:"started_at" tfsdk:"started_at"`
	State             string                                             `json:"state" tfsdk:"state"`
}
type ListDatabasesRes200_DataItem_Region struct {
	DisplayName       string   `json:"display_name" tfsdk:"display_name"`
	Enabled           bool     `json:"enabled" tfsdk:"enabled"`
	Id                string   `json:"id" tfsdk:"id"`
	Location          string   `json:"location" tfsdk:"location"`
	Provider          string   `json:"provider" tfsdk:"provider"`
	PublicIpAddresses []string `json:"public_ip_addresses" tfsdk:"public_ip_addresses"`
	Slug              string   `json:"slug" tfsdk:"slug"`
}
type ListDatabasesRes200_DataItem struct {
	AllowDataBranching                bool                                     `json:"allow_data_branching" tfsdk:"allow_data_branching"`
	AtBackupRestoreBranchesLimit      bool                                     `json:"at_backup_restore_branches_limit" tfsdk:"at_backup_restore_branches_limit"`
	AtDevelopmentBranchLimit          bool                                     `json:"at_development_branch_limit" tfsdk:"at_development_branch_limit"`
	AutomaticMigrations               bool                                     `json:"automatic_migrations" tfsdk:"automatic_migrations"`
	BranchesCount                     float64                                  `json:"branches_count" tfsdk:"branches_count"`
	BranchesUrl                       string                                   `json:"branches_url" tfsdk:"branches_url"`
	CreatedAt                         string                                   `json:"created_at" tfsdk:"created_at"`
	DataImport                        *ListDatabasesRes200_DataItem_DataImport `json:"data_import,omitempty" tfsdk:"data_import"`
	DefaultBranch                     string                                   `json:"default_branch" tfsdk:"default_branch"`
	DefaultBranchReadOnlyRegionsCount float64                                  `json:"default_branch_read_only_regions_count" tfsdk:"default_branch_read_only_regions_count"`
	DefaultBranchShardCount           float64                                  `json:"default_branch_shard_count" tfsdk:"default_branch_shard_count"`
	DefaultBranchTableCount           float64                                  `json:"default_branch_table_count" tfsdk:"default_branch_table_count"`
	DevelopmentBranchesCount          float64                                  `json:"development_branches_count" tfsdk:"development_branches_count"`
	HtmlUrl                           string                                   `json:"html_url" tfsdk:"html_url"`
	Id                                string                                   `json:"id" tfsdk:"id"`
	InsightsRawQueries                bool                                     `json:"insights_raw_queries" tfsdk:"insights_raw_queries"`
	IssuesCount                       float64                                  `json:"issues_count" tfsdk:"issues_count"`
	MigrationFramework                *string                                  `json:"migration_framework,omitempty" tfsdk:"migration_framework"`
	MigrationTableName                *string                                  `json:"migration_table_name,omitempty" tfsdk:"migration_table_name"`
	MultipleAdminsRequiredForDeletion bool                                     `json:"multiple_admins_required_for_deletion" tfsdk:"multiple_admins_required_for_deletion"`
	Name                              string                                   `json:"name" tfsdk:"name"`
	Notes                             *string                                  `json:"notes,omitempty" tfsdk:"notes"`
	Plan                              string                                   `json:"plan" tfsdk:"plan"`
	ProductionBranchWebConsole        bool                                     `json:"production_branch_web_console" tfsdk:"production_branch_web_console"`
	ProductionBranchesCount           float64                                  `json:"production_branches_count" tfsdk:"production_branches_count"`
	Ready                             bool                                     `json:"ready" tfsdk:"ready"`
	Region                            ListDatabasesRes200_DataItem_Region      `json:"region" tfsdk:"region"`
	RequireApprovalForDeploy          bool                                     `json:"require_approval_for_deploy" tfsdk:"require_approval_for_deploy"`
	RestrictBranchRegion              bool                                     `json:"restrict_branch_region" tfsdk:"restrict_branch_region"`
	SchemaLastUpdatedAt               *string                                  `json:"schema_last_updated_at,omitempty" tfsdk:"schema_last_updated_at"`
	Sharded                           bool                                     `json:"sharded" tfsdk:"sharded"`
	State                             string                                   `json:"state" tfsdk:"state"`
	Type                              string                                   `json:"type" tfsdk:"type"`
	UpdatedAt                         string                                   `json:"updated_at" tfsdk:"updated_at"`
	Url                               string                                   `json:"url" tfsdk:"url"`
}
type ListDatabasesRes200 struct {
	Data []ListDatabasesRes200_DataItem `json:"data" tfsdk:"data"`
}

func (cl *Client) ListDatabases(ctx context.Context, organization string, page *int, perPage *int) (res200 *ListDatabasesRes200, res403 *ListDatabasesRes403, res404 *ListDatabasesRes404, res500 *ListDatabasesRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases"})
	q := u.Query()
	if page != nil {
		q.Set("page", strconv.Itoa(*page))
	}
	if perPage != nil {
		q.Set("per_page", strconv.Itoa(*perPage))
	}
	u.RawQuery = q.Encode()
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListDatabasesRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(ListDatabasesRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(ListDatabasesRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(ListDatabasesRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type CreateDatabaseReq struct {
	ClusterSize *string `json:"cluster_size,omitempty" tfsdk:"cluster_size"`
	Name        string  `json:"name" tfsdk:"name"`
	Notes       *string `json:"notes,omitempty" tfsdk:"notes"`
	Plan        *string `json:"plan,omitempty" tfsdk:"plan"`
	Region      *string `json:"region,omitempty" tfsdk:"region"`
}
type CreateDatabaseRes404 struct{}
type CreateDatabaseRes403 struct{}
type CreateDatabaseRes500 struct{}
type CreateDatabaseRes201_DataImport_DataSource struct {
	Database string `json:"database" tfsdk:"database"`
	Hostname string `json:"hostname" tfsdk:"hostname"`
	Port     string `json:"port" tfsdk:"port"`
}
type CreateDatabaseRes201_DataImport struct {
	DataSource        CreateDatabaseRes201_DataImport_DataSource `json:"data_source" tfsdk:"data_source"`
	FinishedAt        string                                     `json:"finished_at" tfsdk:"finished_at"`
	ImportCheckErrors string                                     `json:"import_check_errors" tfsdk:"import_check_errors"`
	StartedAt         string                                     `json:"started_at" tfsdk:"started_at"`
	State             string                                     `json:"state" tfsdk:"state"`
}
type CreateDatabaseRes201_Region struct {
	DisplayName       string   `json:"display_name" tfsdk:"display_name"`
	Enabled           bool     `json:"enabled" tfsdk:"enabled"`
	Id                string   `json:"id" tfsdk:"id"`
	Location          string   `json:"location" tfsdk:"location"`
	Provider          string   `json:"provider" tfsdk:"provider"`
	PublicIpAddresses []string `json:"public_ip_addresses" tfsdk:"public_ip_addresses"`
	Slug              string   `json:"slug" tfsdk:"slug"`
}
type CreateDatabaseRes201 struct {
	AllowDataBranching                bool                             `json:"allow_data_branching" tfsdk:"allow_data_branching"`
	AtBackupRestoreBranchesLimit      bool                             `json:"at_backup_restore_branches_limit" tfsdk:"at_backup_restore_branches_limit"`
	AtDevelopmentBranchLimit          bool                             `json:"at_development_branch_limit" tfsdk:"at_development_branch_limit"`
	AutomaticMigrations               bool                             `json:"automatic_migrations" tfsdk:"automatic_migrations"`
	BranchesCount                     float64                          `json:"branches_count" tfsdk:"branches_count"`
	BranchesUrl                       string                           `json:"branches_url" tfsdk:"branches_url"`
	CreatedAt                         string                           `json:"created_at" tfsdk:"created_at"`
	DataImport                        *CreateDatabaseRes201_DataImport `json:"data_import,omitempty" tfsdk:"data_import"`
	DefaultBranch                     string                           `json:"default_branch" tfsdk:"default_branch"`
	DefaultBranchReadOnlyRegionsCount float64                          `json:"default_branch_read_only_regions_count" tfsdk:"default_branch_read_only_regions_count"`
	DefaultBranchShardCount           float64                          `json:"default_branch_shard_count" tfsdk:"default_branch_shard_count"`
	DefaultBranchTableCount           float64                          `json:"default_branch_table_count" tfsdk:"default_branch_table_count"`
	DevelopmentBranchesCount          float64                          `json:"development_branches_count" tfsdk:"development_branches_count"`
	HtmlUrl                           string                           `json:"html_url" tfsdk:"html_url"`
	Id                                string                           `json:"id" tfsdk:"id"`
	InsightsRawQueries                bool                             `json:"insights_raw_queries" tfsdk:"insights_raw_queries"`
	IssuesCount                       float64                          `json:"issues_count" tfsdk:"issues_count"`
	MigrationFramework                *string                          `json:"migration_framework,omitempty" tfsdk:"migration_framework"`
	MigrationTableName                *string                          `json:"migration_table_name,omitempty" tfsdk:"migration_table_name"`
	MultipleAdminsRequiredForDeletion bool                             `json:"multiple_admins_required_for_deletion" tfsdk:"multiple_admins_required_for_deletion"`
	Name                              string                           `json:"name" tfsdk:"name"`
	Notes                             *string                          `json:"notes,omitempty" tfsdk:"notes"`
	Plan                              string                           `json:"plan" tfsdk:"plan"`
	ProductionBranchWebConsole        bool                             `json:"production_branch_web_console" tfsdk:"production_branch_web_console"`
	ProductionBranchesCount           float64                          `json:"production_branches_count" tfsdk:"production_branches_count"`
	Ready                             bool                             `json:"ready" tfsdk:"ready"`
	Region                            CreateDatabaseRes201_Region      `json:"region" tfsdk:"region"`
	RequireApprovalForDeploy          bool                             `json:"require_approval_for_deploy" tfsdk:"require_approval_for_deploy"`
	RestrictBranchRegion              bool                             `json:"restrict_branch_region" tfsdk:"restrict_branch_region"`
	SchemaLastUpdatedAt               *string                          `json:"schema_last_updated_at,omitempty" tfsdk:"schema_last_updated_at"`
	Sharded                           bool                             `json:"sharded" tfsdk:"sharded"`
	State                             string                           `json:"state" tfsdk:"state"`
	Type                              string                           `json:"type" tfsdk:"type"`
	UpdatedAt                         string                           `json:"updated_at" tfsdk:"updated_at"`
	Url                               string                           `json:"url" tfsdk:"url"`
}

func (cl *Client) CreateDatabase(ctx context.Context, organization string, req CreateDatabaseReq) (res201 *CreateDatabaseRes201, res403 *CreateDatabaseRes403, res404 *CreateDatabaseRes404, res500 *CreateDatabaseRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases"})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res201, res403, res404, res500, err
	}
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return res201, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res201, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 201:
		res201 = new(CreateDatabaseRes201)
		err = json.NewDecoder(res.Body).Decode(&res201)
	case 403:
		res403 = new(CreateDatabaseRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(CreateDatabaseRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(CreateDatabaseRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res201, res403, res404, res500, err
}

type ListBranchesRes500 struct{}
type ListBranchesRes200_DataItem_ApiActor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type ListBranchesRes200_DataItem_PlanetscaleRegion struct {
	DisplayName       string   `json:"display_name" tfsdk:"display_name"`
	Enabled           bool     `json:"enabled" tfsdk:"enabled"`
	Id                string   `json:"id" tfsdk:"id"`
	Location          string   `json:"location" tfsdk:"location"`
	Provider          string   `json:"provider" tfsdk:"provider"`
	PublicIpAddresses []string `json:"public_ip_addresses" tfsdk:"public_ip_addresses"`
	Slug              string   `json:"slug" tfsdk:"slug"`
}
type ListBranchesRes200_DataItem_RestoredFromBranch struct {
	CreatedAt string `json:"created_at" tfsdk:"created_at"`
	DeletedAt string `json:"deleted_at" tfsdk:"deleted_at"`
	Id        string `json:"id" tfsdk:"id"`
	Name      string `json:"name" tfsdk:"name"`
	UpdatedAt string `json:"updated_at" tfsdk:"updated_at"`
}
type ListBranchesRes200_DataItem struct {
	AccessHostUrl               *string                                         `json:"access_host_url,omitempty" tfsdk:"access_host_url"`
	ApiActor                    *ListBranchesRes200_DataItem_ApiActor           `json:"api_actor,omitempty" tfsdk:"api_actor"`
	CreatedAt                   string                                          `json:"created_at" tfsdk:"created_at"`
	HtmlUrl                     string                                          `json:"html_url" tfsdk:"html_url"`
	Id                          string                                          `json:"id" tfsdk:"id"`
	InitialRestoreId            *string                                         `json:"initial_restore_id,omitempty" tfsdk:"initial_restore_id"`
	MysqlAddress                string                                          `json:"mysql_address" tfsdk:"mysql_address"`
	MysqlEdgeAddress            string                                          `json:"mysql_edge_address" tfsdk:"mysql_edge_address"`
	Name                        string                                          `json:"name" tfsdk:"name"`
	ParentBranch                string                                          `json:"parent_branch" tfsdk:"parent_branch"`
	PlanetscaleRegion           *ListBranchesRes200_DataItem_PlanetscaleRegion  `json:"planetscale_region,omitempty" tfsdk:"planetscale_region"`
	Production                  bool                                            `json:"production" tfsdk:"production"`
	Ready                       bool                                            `json:"ready" tfsdk:"ready"`
	RestoreChecklistCompletedAt *string                                         `json:"restore_checklist_completed_at,omitempty" tfsdk:"restore_checklist_completed_at"`
	RestoredFromBranch          *ListBranchesRes200_DataItem_RestoredFromBranch `json:"restored_from_branch,omitempty" tfsdk:"restored_from_branch"`
	SchemaLastUpdatedAt         string                                          `json:"schema_last_updated_at" tfsdk:"schema_last_updated_at"`
	ShardCount                  *float64                                        `json:"shard_count,omitempty" tfsdk:"shard_count"`
	Sharded                     bool                                            `json:"sharded" tfsdk:"sharded"`
	UpdatedAt                   string                                          `json:"updated_at" tfsdk:"updated_at"`
}
type ListBranchesRes200 struct {
	Data []ListBranchesRes200_DataItem `json:"data" tfsdk:"data"`
}
type ListBranchesRes404 struct{}
type ListBranchesRes403 struct{}

func (cl *Client) ListBranches(ctx context.Context, organization string, database string, page *int, perPage *int) (res200 *ListBranchesRes200, res403 *ListBranchesRes403, res404 *ListBranchesRes404, res500 *ListBranchesRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches"})
	q := u.Query()
	if page != nil {
		q.Set("page", strconv.Itoa(*page))
	}
	if perPage != nil {
		q.Set("per_page", strconv.Itoa(*perPage))
	}
	u.RawQuery = q.Encode()
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListBranchesRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(ListBranchesRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(ListBranchesRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(ListBranchesRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type CreateBranchReq struct {
	BackupId     *string `json:"backup_id,omitempty" tfsdk:"backup_id"`
	Name         string  `json:"name" tfsdk:"name"`
	ParentBranch string  `json:"parent_branch" tfsdk:"parent_branch"`
}
type CreateBranchRes500 struct{}
type CreateBranchRes201_ApiActor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type CreateBranchRes201_PlanetscaleRegion struct {
	DisplayName       string   `json:"display_name" tfsdk:"display_name"`
	Enabled           bool     `json:"enabled" tfsdk:"enabled"`
	Id                string   `json:"id" tfsdk:"id"`
	Location          string   `json:"location" tfsdk:"location"`
	Provider          string   `json:"provider" tfsdk:"provider"`
	PublicIpAddresses []string `json:"public_ip_addresses" tfsdk:"public_ip_addresses"`
	Slug              string   `json:"slug" tfsdk:"slug"`
}
type CreateBranchRes201_RestoredFromBranch struct {
	CreatedAt string `json:"created_at" tfsdk:"created_at"`
	DeletedAt string `json:"deleted_at" tfsdk:"deleted_at"`
	Id        string `json:"id" tfsdk:"id"`
	Name      string `json:"name" tfsdk:"name"`
	UpdatedAt string `json:"updated_at" tfsdk:"updated_at"`
}
type CreateBranchRes201 struct {
	AccessHostUrl               *string                                `json:"access_host_url,omitempty" tfsdk:"access_host_url"`
	ApiActor                    *CreateBranchRes201_ApiActor           `json:"api_actor,omitempty" tfsdk:"api_actor"`
	CreatedAt                   string                                 `json:"created_at" tfsdk:"created_at"`
	HtmlUrl                     string                                 `json:"html_url" tfsdk:"html_url"`
	Id                          string                                 `json:"id" tfsdk:"id"`
	InitialRestoreId            *string                                `json:"initial_restore_id,omitempty" tfsdk:"initial_restore_id"`
	MysqlAddress                string                                 `json:"mysql_address" tfsdk:"mysql_address"`
	MysqlEdgeAddress            string                                 `json:"mysql_edge_address" tfsdk:"mysql_edge_address"`
	Name                        string                                 `json:"name" tfsdk:"name"`
	ParentBranch                string                                 `json:"parent_branch" tfsdk:"parent_branch"`
	PlanetscaleRegion           *CreateBranchRes201_PlanetscaleRegion  `json:"planetscale_region,omitempty" tfsdk:"planetscale_region"`
	Production                  bool                                   `json:"production" tfsdk:"production"`
	Ready                       bool                                   `json:"ready" tfsdk:"ready"`
	RestoreChecklistCompletedAt *string                                `json:"restore_checklist_completed_at,omitempty" tfsdk:"restore_checklist_completed_at"`
	RestoredFromBranch          *CreateBranchRes201_RestoredFromBranch `json:"restored_from_branch,omitempty" tfsdk:"restored_from_branch"`
	SchemaLastUpdatedAt         string                                 `json:"schema_last_updated_at" tfsdk:"schema_last_updated_at"`
	ShardCount                  *float64                               `json:"shard_count,omitempty" tfsdk:"shard_count"`
	Sharded                     bool                                   `json:"sharded" tfsdk:"sharded"`
	UpdatedAt                   string                                 `json:"updated_at" tfsdk:"updated_at"`
}
type CreateBranchRes404 struct{}
type CreateBranchRes403 struct{}

func (cl *Client) CreateBranch(ctx context.Context, organization string, database string, req CreateBranchReq) (res201 *CreateBranchRes201, res403 *CreateBranchRes403, res404 *CreateBranchRes404, res500 *CreateBranchRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches"})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res201, res403, res404, res500, err
	}
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return res201, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res201, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 201:
		res201 = new(CreateBranchRes201)
		err = json.NewDecoder(res.Body).Decode(&res201)
	case 403:
		res403 = new(CreateBranchRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(CreateBranchRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(CreateBranchRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res201, res403, res404, res500, err
}

type ListBackupsRes403 struct{}
type ListBackupsRes500 struct{}
type ListBackupsRes200_DataItem_Actor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type ListBackupsRes200_DataItem_BackupPolicy struct {
	CreatedAt      string  `json:"created_at" tfsdk:"created_at"`
	FrequencyUnit  string  `json:"frequency_unit" tfsdk:"frequency_unit"`
	FrequencyValue float64 `json:"frequency_value" tfsdk:"frequency_value"`
	Id             string  `json:"id" tfsdk:"id"`
	LastRanAt      string  `json:"last_ran_at" tfsdk:"last_ran_at"`
	Name           string  `json:"name" tfsdk:"name"`
	NextRunAt      string  `json:"next_run_at" tfsdk:"next_run_at"`
	RetentionUnit  string  `json:"retention_unit" tfsdk:"retention_unit"`
	RetentionValue float64 `json:"retention_value" tfsdk:"retention_value"`
	ScheduleDay    string  `json:"schedule_day" tfsdk:"schedule_day"`
	ScheduleWeek   string  `json:"schedule_week" tfsdk:"schedule_week"`
	Target         string  `json:"target" tfsdk:"target"`
	UpdatedAt      string  `json:"updated_at" tfsdk:"updated_at"`
}
type ListBackupsRes200_DataItem_SchemaSnapshot struct {
	CreatedAt string `json:"created_at" tfsdk:"created_at"`
	Id        string `json:"id" tfsdk:"id"`
	Name      string `json:"name" tfsdk:"name"`
	UpdatedAt string `json:"updated_at" tfsdk:"updated_at"`
	Url       string `json:"url" tfsdk:"url"`
}
type ListBackupsRes200_DataItem struct {
	Actor                ListBackupsRes200_DataItem_Actor          `json:"actor" tfsdk:"actor"`
	BackupPolicy         ListBackupsRes200_DataItem_BackupPolicy   `json:"backup_policy" tfsdk:"backup_policy"`
	CreatedAt            string                                    `json:"created_at" tfsdk:"created_at"`
	EstimatedStorageCost float64                                   `json:"estimated_storage_cost" tfsdk:"estimated_storage_cost"`
	Id                   string                                    `json:"id" tfsdk:"id"`
	Name                 string                                    `json:"name" tfsdk:"name"`
	Required             bool                                      `json:"required" tfsdk:"required"`
	RestoredBranches     *[]string                                 `json:"restored_branches,omitempty" tfsdk:"restored_branches"`
	SchemaSnapshot       ListBackupsRes200_DataItem_SchemaSnapshot `json:"schema_snapshot" tfsdk:"schema_snapshot"`
	Size                 float64                                   `json:"size" tfsdk:"size"`
	State                string                                    `json:"state" tfsdk:"state"`
	UpdatedAt            string                                    `json:"updated_at" tfsdk:"updated_at"`
}
type ListBackupsRes200 struct {
	Data []ListBackupsRes200_DataItem `json:"data" tfsdk:"data"`
}
type ListBackupsRes404 struct{}

func (cl *Client) ListBackups(ctx context.Context, organization string, database string, branch string, page *int, perPage *int) (res200 *ListBackupsRes200, res403 *ListBackupsRes403, res404 *ListBackupsRes404, res500 *ListBackupsRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + branch + "/backups"})
	q := u.Query()
	if page != nil {
		q.Set("page", strconv.Itoa(*page))
	}
	if perPage != nil {
		q.Set("per_page", strconv.Itoa(*perPage))
	}
	u.RawQuery = q.Encode()
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListBackupsRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(ListBackupsRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(ListBackupsRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(ListBackupsRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type CreateBackupReq struct {
	Name           *string  `json:"name,omitempty" tfsdk:"name"`
	RetentionUnit  *string  `json:"retention_unit,omitempty" tfsdk:"retention_unit"`
	RetentionValue *float64 `json:"retention_value,omitempty" tfsdk:"retention_value"`
}
type CreateBackupRes404 struct{}
type CreateBackupRes403 struct{}
type CreateBackupRes500 struct{}
type CreateBackupRes200_DataItem_Actor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type CreateBackupRes200_DataItem_BackupPolicy struct {
	CreatedAt      string  `json:"created_at" tfsdk:"created_at"`
	FrequencyUnit  string  `json:"frequency_unit" tfsdk:"frequency_unit"`
	FrequencyValue float64 `json:"frequency_value" tfsdk:"frequency_value"`
	Id             string  `json:"id" tfsdk:"id"`
	LastRanAt      string  `json:"last_ran_at" tfsdk:"last_ran_at"`
	Name           string  `json:"name" tfsdk:"name"`
	NextRunAt      string  `json:"next_run_at" tfsdk:"next_run_at"`
	RetentionUnit  string  `json:"retention_unit" tfsdk:"retention_unit"`
	RetentionValue float64 `json:"retention_value" tfsdk:"retention_value"`
	ScheduleDay    string  `json:"schedule_day" tfsdk:"schedule_day"`
	ScheduleWeek   string  `json:"schedule_week" tfsdk:"schedule_week"`
	Target         string  `json:"target" tfsdk:"target"`
	UpdatedAt      string  `json:"updated_at" tfsdk:"updated_at"`
}
type CreateBackupRes200_DataItem_SchemaSnapshot struct {
	CreatedAt string `json:"created_at" tfsdk:"created_at"`
	Id        string `json:"id" tfsdk:"id"`
	Name      string `json:"name" tfsdk:"name"`
	UpdatedAt string `json:"updated_at" tfsdk:"updated_at"`
	Url       string `json:"url" tfsdk:"url"`
}
type CreateBackupRes200_DataItem struct {
	Actor                CreateBackupRes200_DataItem_Actor          `json:"actor" tfsdk:"actor"`
	BackupPolicy         CreateBackupRes200_DataItem_BackupPolicy   `json:"backup_policy" tfsdk:"backup_policy"`
	CreatedAt            string                                     `json:"created_at" tfsdk:"created_at"`
	EstimatedStorageCost float64                                    `json:"estimated_storage_cost" tfsdk:"estimated_storage_cost"`
	Id                   string                                     `json:"id" tfsdk:"id"`
	Name                 string                                     `json:"name" tfsdk:"name"`
	Required             bool                                       `json:"required" tfsdk:"required"`
	RestoredBranches     *[]string                                  `json:"restored_branches,omitempty" tfsdk:"restored_branches"`
	SchemaSnapshot       CreateBackupRes200_DataItem_SchemaSnapshot `json:"schema_snapshot" tfsdk:"schema_snapshot"`
	Size                 float64                                    `json:"size" tfsdk:"size"`
	State                string                                     `json:"state" tfsdk:"state"`
	UpdatedAt            string                                     `json:"updated_at" tfsdk:"updated_at"`
}
type CreateBackupRes200 struct {
	Data []CreateBackupRes200_DataItem `json:"data" tfsdk:"data"`
}

func (cl *Client) CreateBackup(ctx context.Context, organization string, database string, branch string, req CreateBackupReq) (res200 *CreateBackupRes200, res403 *CreateBackupRes403, res404 *CreateBackupRes404, res500 *CreateBackupRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + branch + "/backups"})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res200, res403, res404, res500, err
	}
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(CreateBackupRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(CreateBackupRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(CreateBackupRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(CreateBackupRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type GetBackupRes500 struct{}
type GetBackupRes200_Actor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type GetBackupRes200_BackupPolicy struct {
	CreatedAt      string  `json:"created_at" tfsdk:"created_at"`
	FrequencyUnit  string  `json:"frequency_unit" tfsdk:"frequency_unit"`
	FrequencyValue float64 `json:"frequency_value" tfsdk:"frequency_value"`
	Id             string  `json:"id" tfsdk:"id"`
	LastRanAt      string  `json:"last_ran_at" tfsdk:"last_ran_at"`
	Name           string  `json:"name" tfsdk:"name"`
	NextRunAt      string  `json:"next_run_at" tfsdk:"next_run_at"`
	RetentionUnit  string  `json:"retention_unit" tfsdk:"retention_unit"`
	RetentionValue float64 `json:"retention_value" tfsdk:"retention_value"`
	ScheduleDay    string  `json:"schedule_day" tfsdk:"schedule_day"`
	ScheduleWeek   string  `json:"schedule_week" tfsdk:"schedule_week"`
	Target         string  `json:"target" tfsdk:"target"`
	UpdatedAt      string  `json:"updated_at" tfsdk:"updated_at"`
}
type GetBackupRes200_SchemaSnapshot struct {
	CreatedAt string `json:"created_at" tfsdk:"created_at"`
	Id        string `json:"id" tfsdk:"id"`
	Name      string `json:"name" tfsdk:"name"`
	UpdatedAt string `json:"updated_at" tfsdk:"updated_at"`
	Url       string `json:"url" tfsdk:"url"`
}
type GetBackupRes200 struct {
	Actor                GetBackupRes200_Actor          `json:"actor" tfsdk:"actor"`
	BackupPolicy         GetBackupRes200_BackupPolicy   `json:"backup_policy" tfsdk:"backup_policy"`
	CreatedAt            string                         `json:"created_at" tfsdk:"created_at"`
	EstimatedStorageCost float64                        `json:"estimated_storage_cost" tfsdk:"estimated_storage_cost"`
	Id                   string                         `json:"id" tfsdk:"id"`
	Name                 string                         `json:"name" tfsdk:"name"`
	Required             bool                           `json:"required" tfsdk:"required"`
	RestoredBranches     *[]string                      `json:"restored_branches,omitempty" tfsdk:"restored_branches"`
	SchemaSnapshot       GetBackupRes200_SchemaSnapshot `json:"schema_snapshot" tfsdk:"schema_snapshot"`
	Size                 float64                        `json:"size" tfsdk:"size"`
	State                string                         `json:"state" tfsdk:"state"`
	UpdatedAt            string                         `json:"updated_at" tfsdk:"updated_at"`
}
type GetBackupRes404 struct{}
type GetBackupRes403 struct{}

func (cl *Client) GetBackup(ctx context.Context, organization string, database string, branch string, id string) (res200 *GetBackupRes200, res403 *GetBackupRes403, res404 *GetBackupRes404, res500 *GetBackupRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + branch + "/backups/" + id})
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(GetBackupRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(GetBackupRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(GetBackupRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(GetBackupRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type DeleteBackupRes204 struct{}
type DeleteBackupRes404 struct{}
type DeleteBackupRes403 struct{}
type DeleteBackupRes500 struct{}

func (cl *Client) DeleteBackup(ctx context.Context, organization string, database string, branch string, id string) (res204 *DeleteBackupRes204, res403 *DeleteBackupRes403, res404 *DeleteBackupRes404, res500 *DeleteBackupRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + branch + "/backups/" + id})
	r, err := http.NewRequestWithContext(ctx, "DELETE", u.String(), nil)
	if err != nil {
		return res204, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res204, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 204:
		res204 = new(DeleteBackupRes204)
		err = json.NewDecoder(res.Body).Decode(&res204)
	case 403:
		res403 = new(DeleteBackupRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(DeleteBackupRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(DeleteBackupRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res204, res403, res404, res500, err
}

type ListPasswordsRes200_DataItem_Actor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type ListPasswordsRes200_DataItem_DatabaseBranch struct {
	AccessHostUrl    string `json:"access_host_url" tfsdk:"access_host_url"`
	Id               string `json:"id" tfsdk:"id"`
	MysqlEdgeAddress string `json:"mysql_edge_address" tfsdk:"mysql_edge_address"`
	Name             string `json:"name" tfsdk:"name"`
	Production       bool   `json:"production" tfsdk:"production"`
}
type ListPasswordsRes200_DataItem_Region struct {
	DisplayName       string   `json:"display_name" tfsdk:"display_name"`
	Enabled           bool     `json:"enabled" tfsdk:"enabled"`
	Id                string   `json:"id" tfsdk:"id"`
	Location          string   `json:"location" tfsdk:"location"`
	Provider          string   `json:"provider" tfsdk:"provider"`
	PublicIpAddresses []string `json:"public_ip_addresses" tfsdk:"public_ip_addresses"`
	Slug              string   `json:"slug" tfsdk:"slug"`
}
type ListPasswordsRes200_DataItem struct {
	AccessHostUrl  string                                      `json:"access_host_url" tfsdk:"access_host_url"`
	Actor          *ListPasswordsRes200_DataItem_Actor         `json:"actor,omitempty" tfsdk:"actor"`
	CreatedAt      string                                      `json:"created_at" tfsdk:"created_at"`
	DatabaseBranch ListPasswordsRes200_DataItem_DatabaseBranch `json:"database_branch" tfsdk:"database_branch"`
	DeletedAt      *string                                     `json:"deleted_at,omitempty" tfsdk:"deleted_at"`
	ExpiresAt      *string                                     `json:"expires_at,omitempty" tfsdk:"expires_at"`
	Id             string                                      `json:"id" tfsdk:"id"`
	Integrations   []string                                    `json:"integrations" tfsdk:"integrations"`
	Name           string                                      `json:"name" tfsdk:"name"`
	Region         *ListPasswordsRes200_DataItem_Region        `json:"region,omitempty" tfsdk:"region"`
	Renewable      bool                                        `json:"renewable" tfsdk:"renewable"`
	Role           string                                      `json:"role" tfsdk:"role"`
	TtlSeconds     float64                                     `json:"ttl_seconds" tfsdk:"ttl_seconds"`
	Username       *string                                     `json:"username,omitempty" tfsdk:"username"`
}
type ListPasswordsRes200 struct {
	Data []ListPasswordsRes200_DataItem `json:"data" tfsdk:"data"`
}
type ListPasswordsRes404 struct{}
type ListPasswordsRes403 struct{}
type ListPasswordsRes500 struct{}

func (cl *Client) ListPasswords(ctx context.Context, organization string, database string, branch string, readOnlyRegionId *string, page *int, perPage *int) (res200 *ListPasswordsRes200, res403 *ListPasswordsRes403, res404 *ListPasswordsRes404, res500 *ListPasswordsRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + branch + "/passwords"})
	q := u.Query()
	if readOnlyRegionId != nil {
		q.Set("read_only_region_id", *readOnlyRegionId)
	}
	if page != nil {
		q.Set("page", strconv.Itoa(*page))
	}
	if perPage != nil {
		q.Set("per_page", strconv.Itoa(*perPage))
	}
	u.RawQuery = q.Encode()
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListPasswordsRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(ListPasswordsRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(ListPasswordsRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(ListPasswordsRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type CreatePasswordReq struct {
	Name *string  `json:"name,omitempty" tfsdk:"name"`
	Role *string  `json:"role,omitempty" tfsdk:"role"`
	Ttl  *float64 `json:"ttl,omitempty" tfsdk:"ttl"`
}
type CreatePasswordRes404 struct{}
type CreatePasswordRes403 struct{}
type CreatePasswordRes500 struct{}
type CreatePasswordRes201_Actor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type CreatePasswordRes201_DatabaseBranch struct {
	AccessHostUrl    string `json:"access_host_url" tfsdk:"access_host_url"`
	Id               string `json:"id" tfsdk:"id"`
	MysqlEdgeAddress string `json:"mysql_edge_address" tfsdk:"mysql_edge_address"`
	Name             string `json:"name" tfsdk:"name"`
	Production       bool   `json:"production" tfsdk:"production"`
}
type CreatePasswordRes201_Region struct {
	DisplayName       string   `json:"display_name" tfsdk:"display_name"`
	Enabled           bool     `json:"enabled" tfsdk:"enabled"`
	Id                string   `json:"id" tfsdk:"id"`
	Location          string   `json:"location" tfsdk:"location"`
	Provider          string   `json:"provider" tfsdk:"provider"`
	PublicIpAddresses []string `json:"public_ip_addresses" tfsdk:"public_ip_addresses"`
	Slug              string   `json:"slug" tfsdk:"slug"`
}
type CreatePasswordRes201 struct {
	AccessHostUrl  string                              `json:"access_host_url" tfsdk:"access_host_url"`
	Actor          *CreatePasswordRes201_Actor         `json:"actor,omitempty" tfsdk:"actor"`
	CreatedAt      string                              `json:"created_at" tfsdk:"created_at"`
	DatabaseBranch CreatePasswordRes201_DatabaseBranch `json:"database_branch" tfsdk:"database_branch"`
	DeletedAt      *string                             `json:"deleted_at,omitempty" tfsdk:"deleted_at"`
	ExpiresAt      *string                             `json:"expires_at,omitempty" tfsdk:"expires_at"`
	Id             string                              `json:"id" tfsdk:"id"`
	Integrations   []string                            `json:"integrations" tfsdk:"integrations"`
	Name           string                              `json:"name" tfsdk:"name"`
	PlainText      string                              `json:"plain_text" tfsdk:"plain_text"`
	Region         *CreatePasswordRes201_Region        `json:"region,omitempty" tfsdk:"region"`
	Renewable      bool                                `json:"renewable" tfsdk:"renewable"`
	Role           string                              `json:"role" tfsdk:"role"`
	TtlSeconds     float64                             `json:"ttl_seconds" tfsdk:"ttl_seconds"`
	Username       *string                             `json:"username,omitempty" tfsdk:"username"`
}

func (cl *Client) CreatePassword(ctx context.Context, organization string, database string, branch string, req CreatePasswordReq) (res201 *CreatePasswordRes201, res403 *CreatePasswordRes403, res404 *CreatePasswordRes404, res500 *CreatePasswordRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + branch + "/passwords"})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res201, res403, res404, res500, err
	}
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return res201, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res201, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 201:
		res201 = new(CreatePasswordRes201)
		err = json.NewDecoder(res.Body).Decode(&res201)
	case 403:
		res403 = new(CreatePasswordRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(CreatePasswordRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(CreatePasswordRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res201, res403, res404, res500, err
}

type GetPasswordRes500 struct{}
type GetPasswordRes200_Actor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type GetPasswordRes200_DatabaseBranch struct {
	AccessHostUrl    string `json:"access_host_url" tfsdk:"access_host_url"`
	Id               string `json:"id" tfsdk:"id"`
	MysqlEdgeAddress string `json:"mysql_edge_address" tfsdk:"mysql_edge_address"`
	Name             string `json:"name" tfsdk:"name"`
	Production       bool   `json:"production" tfsdk:"production"`
}
type GetPasswordRes200_Region struct {
	DisplayName       string   `json:"display_name" tfsdk:"display_name"`
	Enabled           bool     `json:"enabled" tfsdk:"enabled"`
	Id                string   `json:"id" tfsdk:"id"`
	Location          string   `json:"location" tfsdk:"location"`
	Provider          string   `json:"provider" tfsdk:"provider"`
	PublicIpAddresses []string `json:"public_ip_addresses" tfsdk:"public_ip_addresses"`
	Slug              string   `json:"slug" tfsdk:"slug"`
}
type GetPasswordRes200 struct {
	AccessHostUrl  string                           `json:"access_host_url" tfsdk:"access_host_url"`
	Actor          *GetPasswordRes200_Actor         `json:"actor,omitempty" tfsdk:"actor"`
	CreatedAt      string                           `json:"created_at" tfsdk:"created_at"`
	DatabaseBranch GetPasswordRes200_DatabaseBranch `json:"database_branch" tfsdk:"database_branch"`
	DeletedAt      *string                          `json:"deleted_at,omitempty" tfsdk:"deleted_at"`
	ExpiresAt      *string                          `json:"expires_at,omitempty" tfsdk:"expires_at"`
	Id             string                           `json:"id" tfsdk:"id"`
	Integrations   []string                         `json:"integrations" tfsdk:"integrations"`
	Name           string                           `json:"name" tfsdk:"name"`
	Region         *GetPasswordRes200_Region        `json:"region,omitempty" tfsdk:"region"`
	Renewable      bool                             `json:"renewable" tfsdk:"renewable"`
	Role           string                           `json:"role" tfsdk:"role"`
	TtlSeconds     float64                          `json:"ttl_seconds" tfsdk:"ttl_seconds"`
	Username       *string                          `json:"username,omitempty" tfsdk:"username"`
}
type GetPasswordRes404 struct{}
type GetPasswordRes403 struct{}

func (cl *Client) GetPassword(ctx context.Context, organization string, database string, branch string, id string, readOnlyRegionId *string) (res200 *GetPasswordRes200, res403 *GetPasswordRes403, res404 *GetPasswordRes404, res500 *GetPasswordRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + branch + "/passwords/" + id})
	q := u.Query()
	if readOnlyRegionId != nil {
		q.Set("read_only_region_id", *readOnlyRegionId)
	}
	u.RawQuery = q.Encode()
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(GetPasswordRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(GetPasswordRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(GetPasswordRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(GetPasswordRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type DeletePasswordRes404 struct{}
type DeletePasswordRes403 struct{}
type DeletePasswordRes500 struct{}
type DeletePasswordRes204 struct{}

func (cl *Client) DeletePassword(ctx context.Context, organization string, database string, branch string, id string) (res204 *DeletePasswordRes204, res403 *DeletePasswordRes403, res404 *DeletePasswordRes404, res500 *DeletePasswordRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + branch + "/passwords/" + id})
	r, err := http.NewRequestWithContext(ctx, "DELETE", u.String(), nil)
	if err != nil {
		return res204, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res204, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 204:
		res204 = new(DeletePasswordRes204)
		err = json.NewDecoder(res.Body).Decode(&res204)
	case 403:
		res403 = new(DeletePasswordRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(DeletePasswordRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(DeletePasswordRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res204, res403, res404, res500, err
}

type UpdatePasswordReq struct {
	Name string `json:"name" tfsdk:"name"`
}
type UpdatePasswordRes404 struct{}
type UpdatePasswordRes403 struct{}
type UpdatePasswordRes500 struct{}
type UpdatePasswordRes200_Actor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type UpdatePasswordRes200_DatabaseBranch struct {
	AccessHostUrl    string `json:"access_host_url" tfsdk:"access_host_url"`
	Id               string `json:"id" tfsdk:"id"`
	MysqlEdgeAddress string `json:"mysql_edge_address" tfsdk:"mysql_edge_address"`
	Name             string `json:"name" tfsdk:"name"`
	Production       bool   `json:"production" tfsdk:"production"`
}
type UpdatePasswordRes200_Region struct {
	DisplayName       string   `json:"display_name" tfsdk:"display_name"`
	Enabled           bool     `json:"enabled" tfsdk:"enabled"`
	Id                string   `json:"id" tfsdk:"id"`
	Location          string   `json:"location" tfsdk:"location"`
	Provider          string   `json:"provider" tfsdk:"provider"`
	PublicIpAddresses []string `json:"public_ip_addresses" tfsdk:"public_ip_addresses"`
	Slug              string   `json:"slug" tfsdk:"slug"`
}
type UpdatePasswordRes200 struct {
	AccessHostUrl  string                              `json:"access_host_url" tfsdk:"access_host_url"`
	Actor          *UpdatePasswordRes200_Actor         `json:"actor,omitempty" tfsdk:"actor"`
	CreatedAt      string                              `json:"created_at" tfsdk:"created_at"`
	DatabaseBranch UpdatePasswordRes200_DatabaseBranch `json:"database_branch" tfsdk:"database_branch"`
	DeletedAt      *string                             `json:"deleted_at,omitempty" tfsdk:"deleted_at"`
	ExpiresAt      *string                             `json:"expires_at,omitempty" tfsdk:"expires_at"`
	Id             string                              `json:"id" tfsdk:"id"`
	Integrations   []string                            `json:"integrations" tfsdk:"integrations"`
	Name           string                              `json:"name" tfsdk:"name"`
	Region         *UpdatePasswordRes200_Region        `json:"region,omitempty" tfsdk:"region"`
	Renewable      bool                                `json:"renewable" tfsdk:"renewable"`
	Role           string                              `json:"role" tfsdk:"role"`
	TtlSeconds     float64                             `json:"ttl_seconds" tfsdk:"ttl_seconds"`
	Username       *string                             `json:"username,omitempty" tfsdk:"username"`
}

func (cl *Client) UpdatePassword(ctx context.Context, organization string, database string, branch string, id string, req UpdatePasswordReq) (res200 *UpdatePasswordRes200, res403 *UpdatePasswordRes403, res404 *UpdatePasswordRes404, res500 *UpdatePasswordRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + branch + "/passwords/" + id})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res200, res403, res404, res500, err
	}
	r, err := http.NewRequestWithContext(ctx, "PATCH", u.String(), body)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(UpdatePasswordRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(UpdatePasswordRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(UpdatePasswordRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(UpdatePasswordRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type RenewPasswordReq struct {
	ReadOnlyRegionId *string `json:"read_only_region_id,omitempty" tfsdk:"read_only_region_id"`
}
type RenewPasswordRes404 struct{}
type RenewPasswordRes403 struct{}
type RenewPasswordRes500 struct{}
type RenewPasswordRes200_Actor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type RenewPasswordRes200_DatabaseBranch struct {
	AccessHostUrl    string `json:"access_host_url" tfsdk:"access_host_url"`
	Id               string `json:"id" tfsdk:"id"`
	MysqlEdgeAddress string `json:"mysql_edge_address" tfsdk:"mysql_edge_address"`
	Name             string `json:"name" tfsdk:"name"`
	Production       bool   `json:"production" tfsdk:"production"`
}
type RenewPasswordRes200_Region struct {
	DisplayName       string   `json:"display_name" tfsdk:"display_name"`
	Enabled           bool     `json:"enabled" tfsdk:"enabled"`
	Id                string   `json:"id" tfsdk:"id"`
	Location          string   `json:"location" tfsdk:"location"`
	Provider          string   `json:"provider" tfsdk:"provider"`
	PublicIpAddresses []string `json:"public_ip_addresses" tfsdk:"public_ip_addresses"`
	Slug              string   `json:"slug" tfsdk:"slug"`
}
type RenewPasswordRes200 struct {
	AccessHostUrl  string                             `json:"access_host_url" tfsdk:"access_host_url"`
	Actor          *RenewPasswordRes200_Actor         `json:"actor,omitempty" tfsdk:"actor"`
	CreatedAt      string                             `json:"created_at" tfsdk:"created_at"`
	DatabaseBranch RenewPasswordRes200_DatabaseBranch `json:"database_branch" tfsdk:"database_branch"`
	DeletedAt      *string                            `json:"deleted_at,omitempty" tfsdk:"deleted_at"`
	ExpiresAt      *string                            `json:"expires_at,omitempty" tfsdk:"expires_at"`
	Id             string                             `json:"id" tfsdk:"id"`
	Integrations   []string                           `json:"integrations" tfsdk:"integrations"`
	Name           string                             `json:"name" tfsdk:"name"`
	PlainText      string                             `json:"plain_text" tfsdk:"plain_text"`
	Region         *RenewPasswordRes200_Region        `json:"region,omitempty" tfsdk:"region"`
	Renewable      bool                               `json:"renewable" tfsdk:"renewable"`
	Role           string                             `json:"role" tfsdk:"role"`
	TtlSeconds     float64                            `json:"ttl_seconds" tfsdk:"ttl_seconds"`
	Username       *string                            `json:"username,omitempty" tfsdk:"username"`
}

func (cl *Client) RenewPassword(ctx context.Context, organization string, database string, branch string, id string, req RenewPasswordReq) (res200 *RenewPasswordRes200, res403 *RenewPasswordRes403, res404 *RenewPasswordRes404, res500 *RenewPasswordRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + branch + "/passwords/" + id + "/renew"})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res200, res403, res404, res500, err
	}
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(RenewPasswordRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(RenewPasswordRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(RenewPasswordRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(RenewPasswordRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type GetBranchRes404 struct{}
type GetBranchRes403 struct{}
type GetBranchRes500 struct{}
type GetBranchRes200_ApiActor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type GetBranchRes200_PlanetscaleRegion struct {
	DisplayName       string   `json:"display_name" tfsdk:"display_name"`
	Enabled           bool     `json:"enabled" tfsdk:"enabled"`
	Id                string   `json:"id" tfsdk:"id"`
	Location          string   `json:"location" tfsdk:"location"`
	Provider          string   `json:"provider" tfsdk:"provider"`
	PublicIpAddresses []string `json:"public_ip_addresses" tfsdk:"public_ip_addresses"`
	Slug              string   `json:"slug" tfsdk:"slug"`
}
type GetBranchRes200_RestoredFromBranch struct {
	CreatedAt string `json:"created_at" tfsdk:"created_at"`
	DeletedAt string `json:"deleted_at" tfsdk:"deleted_at"`
	Id        string `json:"id" tfsdk:"id"`
	Name      string `json:"name" tfsdk:"name"`
	UpdatedAt string `json:"updated_at" tfsdk:"updated_at"`
}
type GetBranchRes200 struct {
	AccessHostUrl               *string                             `json:"access_host_url,omitempty" tfsdk:"access_host_url"`
	ApiActor                    *GetBranchRes200_ApiActor           `json:"api_actor,omitempty" tfsdk:"api_actor"`
	CreatedAt                   string                              `json:"created_at" tfsdk:"created_at"`
	HtmlUrl                     string                              `json:"html_url" tfsdk:"html_url"`
	Id                          string                              `json:"id" tfsdk:"id"`
	InitialRestoreId            *string                             `json:"initial_restore_id,omitempty" tfsdk:"initial_restore_id"`
	MysqlAddress                string                              `json:"mysql_address" tfsdk:"mysql_address"`
	MysqlEdgeAddress            string                              `json:"mysql_edge_address" tfsdk:"mysql_edge_address"`
	Name                        string                              `json:"name" tfsdk:"name"`
	ParentBranch                string                              `json:"parent_branch" tfsdk:"parent_branch"`
	PlanetscaleRegion           *GetBranchRes200_PlanetscaleRegion  `json:"planetscale_region,omitempty" tfsdk:"planetscale_region"`
	Production                  bool                                `json:"production" tfsdk:"production"`
	Ready                       bool                                `json:"ready" tfsdk:"ready"`
	RestoreChecklistCompletedAt *string                             `json:"restore_checklist_completed_at,omitempty" tfsdk:"restore_checklist_completed_at"`
	RestoredFromBranch          *GetBranchRes200_RestoredFromBranch `json:"restored_from_branch,omitempty" tfsdk:"restored_from_branch"`
	SchemaLastUpdatedAt         string                              `json:"schema_last_updated_at" tfsdk:"schema_last_updated_at"`
	ShardCount                  *float64                            `json:"shard_count,omitempty" tfsdk:"shard_count"`
	Sharded                     bool                                `json:"sharded" tfsdk:"sharded"`
	UpdatedAt                   string                              `json:"updated_at" tfsdk:"updated_at"`
}

func (cl *Client) GetBranch(ctx context.Context, organization string, database string, name string) (res200 *GetBranchRes200, res403 *GetBranchRes403, res404 *GetBranchRes404, res500 *GetBranchRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + name})
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(GetBranchRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(GetBranchRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(GetBranchRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(GetBranchRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type DeleteBranchRes500 struct{}
type DeleteBranchRes204 struct{}
type DeleteBranchRes404 struct{}
type DeleteBranchRes403 struct{}

func (cl *Client) DeleteBranch(ctx context.Context, organization string, database string, name string) (res204 *DeleteBranchRes204, res403 *DeleteBranchRes403, res404 *DeleteBranchRes404, res500 *DeleteBranchRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + name})
	r, err := http.NewRequestWithContext(ctx, "DELETE", u.String(), nil)
	if err != nil {
		return res204, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res204, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 204:
		res204 = new(DeleteBranchRes204)
		err = json.NewDecoder(res.Body).Decode(&res204)
	case 403:
		res403 = new(DeleteBranchRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(DeleteBranchRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(DeleteBranchRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res204, res403, res404, res500, err
}

type DemoteBranchRes403 struct{}
type DemoteBranchRes500 struct{}
type DemoteBranchRes200_ApiActor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type DemoteBranchRes200_PlanetscaleRegion struct {
	DisplayName       string   `json:"display_name" tfsdk:"display_name"`
	Enabled           bool     `json:"enabled" tfsdk:"enabled"`
	Id                string   `json:"id" tfsdk:"id"`
	Location          string   `json:"location" tfsdk:"location"`
	Provider          string   `json:"provider" tfsdk:"provider"`
	PublicIpAddresses []string `json:"public_ip_addresses" tfsdk:"public_ip_addresses"`
	Slug              string   `json:"slug" tfsdk:"slug"`
}
type DemoteBranchRes200_RestoredFromBranch struct {
	CreatedAt string `json:"created_at" tfsdk:"created_at"`
	DeletedAt string `json:"deleted_at" tfsdk:"deleted_at"`
	Id        string `json:"id" tfsdk:"id"`
	Name      string `json:"name" tfsdk:"name"`
	UpdatedAt string `json:"updated_at" tfsdk:"updated_at"`
}
type DemoteBranchRes200 struct {
	AccessHostUrl               *string                                `json:"access_host_url,omitempty" tfsdk:"access_host_url"`
	ApiActor                    *DemoteBranchRes200_ApiActor           `json:"api_actor,omitempty" tfsdk:"api_actor"`
	CreatedAt                   string                                 `json:"created_at" tfsdk:"created_at"`
	HtmlUrl                     string                                 `json:"html_url" tfsdk:"html_url"`
	Id                          string                                 `json:"id" tfsdk:"id"`
	InitialRestoreId            *string                                `json:"initial_restore_id,omitempty" tfsdk:"initial_restore_id"`
	MysqlAddress                string                                 `json:"mysql_address" tfsdk:"mysql_address"`
	MysqlEdgeAddress            string                                 `json:"mysql_edge_address" tfsdk:"mysql_edge_address"`
	Name                        string                                 `json:"name" tfsdk:"name"`
	ParentBranch                string                                 `json:"parent_branch" tfsdk:"parent_branch"`
	PlanetscaleRegion           *DemoteBranchRes200_PlanetscaleRegion  `json:"planetscale_region,omitempty" tfsdk:"planetscale_region"`
	Production                  bool                                   `json:"production" tfsdk:"production"`
	Ready                       bool                                   `json:"ready" tfsdk:"ready"`
	RestoreChecklistCompletedAt *string                                `json:"restore_checklist_completed_at,omitempty" tfsdk:"restore_checklist_completed_at"`
	RestoredFromBranch          *DemoteBranchRes200_RestoredFromBranch `json:"restored_from_branch,omitempty" tfsdk:"restored_from_branch"`
	SchemaLastUpdatedAt         string                                 `json:"schema_last_updated_at" tfsdk:"schema_last_updated_at"`
	ShardCount                  *float64                               `json:"shard_count,omitempty" tfsdk:"shard_count"`
	Sharded                     bool                                   `json:"sharded" tfsdk:"sharded"`
	UpdatedAt                   string                                 `json:"updated_at" tfsdk:"updated_at"`
}
type DemoteBranchRes404 struct{}

func (cl *Client) DemoteBranch(ctx context.Context, organization string, database string, name string) (res200 *DemoteBranchRes200, res403 *DemoteBranchRes403, res404 *DemoteBranchRes404, res500 *DemoteBranchRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + name + "/demote"})
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), nil)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(DemoteBranchRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(DemoteBranchRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(DemoteBranchRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(DemoteBranchRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type PromoteBranchRes500 struct{}
type PromoteBranchRes200_ApiActor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type PromoteBranchRes200_PlanetscaleRegion struct {
	DisplayName       string   `json:"display_name" tfsdk:"display_name"`
	Enabled           bool     `json:"enabled" tfsdk:"enabled"`
	Id                string   `json:"id" tfsdk:"id"`
	Location          string   `json:"location" tfsdk:"location"`
	Provider          string   `json:"provider" tfsdk:"provider"`
	PublicIpAddresses []string `json:"public_ip_addresses" tfsdk:"public_ip_addresses"`
	Slug              string   `json:"slug" tfsdk:"slug"`
}
type PromoteBranchRes200_RestoredFromBranch struct {
	CreatedAt string `json:"created_at" tfsdk:"created_at"`
	DeletedAt string `json:"deleted_at" tfsdk:"deleted_at"`
	Id        string `json:"id" tfsdk:"id"`
	Name      string `json:"name" tfsdk:"name"`
	UpdatedAt string `json:"updated_at" tfsdk:"updated_at"`
}
type PromoteBranchRes200 struct {
	AccessHostUrl               *string                                 `json:"access_host_url,omitempty" tfsdk:"access_host_url"`
	ApiActor                    *PromoteBranchRes200_ApiActor           `json:"api_actor,omitempty" tfsdk:"api_actor"`
	CreatedAt                   string                                  `json:"created_at" tfsdk:"created_at"`
	HtmlUrl                     string                                  `json:"html_url" tfsdk:"html_url"`
	Id                          string                                  `json:"id" tfsdk:"id"`
	InitialRestoreId            *string                                 `json:"initial_restore_id,omitempty" tfsdk:"initial_restore_id"`
	MysqlAddress                string                                  `json:"mysql_address" tfsdk:"mysql_address"`
	MysqlEdgeAddress            string                                  `json:"mysql_edge_address" tfsdk:"mysql_edge_address"`
	Name                        string                                  `json:"name" tfsdk:"name"`
	ParentBranch                string                                  `json:"parent_branch" tfsdk:"parent_branch"`
	PlanetscaleRegion           *PromoteBranchRes200_PlanetscaleRegion  `json:"planetscale_region,omitempty" tfsdk:"planetscale_region"`
	Production                  bool                                    `json:"production" tfsdk:"production"`
	Ready                       bool                                    `json:"ready" tfsdk:"ready"`
	RestoreChecklistCompletedAt *string                                 `json:"restore_checklist_completed_at,omitempty" tfsdk:"restore_checklist_completed_at"`
	RestoredFromBranch          *PromoteBranchRes200_RestoredFromBranch `json:"restored_from_branch,omitempty" tfsdk:"restored_from_branch"`
	SchemaLastUpdatedAt         string                                  `json:"schema_last_updated_at" tfsdk:"schema_last_updated_at"`
	ShardCount                  *float64                                `json:"shard_count,omitempty" tfsdk:"shard_count"`
	Sharded                     bool                                    `json:"sharded" tfsdk:"sharded"`
	UpdatedAt                   string                                  `json:"updated_at" tfsdk:"updated_at"`
}
type PromoteBranchRes404 struct{}
type PromoteBranchRes403 struct{}

func (cl *Client) PromoteBranch(ctx context.Context, organization string, database string, name string) (res200 *PromoteBranchRes200, res403 *PromoteBranchRes403, res404 *PromoteBranchRes404, res500 *PromoteBranchRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + name + "/promote"})
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), nil)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(PromoteBranchRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(PromoteBranchRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(PromoteBranchRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(PromoteBranchRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type EnableSafeMigrationsForBranchRes404 struct{}
type EnableSafeMigrationsForBranchRes403 struct{}
type EnableSafeMigrationsForBranchRes500 struct{}
type EnableSafeMigrationsForBranchRes200_ApiActor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type EnableSafeMigrationsForBranchRes200_PlanetscaleRegion struct {
	DisplayName       string   `json:"display_name" tfsdk:"display_name"`
	Enabled           bool     `json:"enabled" tfsdk:"enabled"`
	Id                string   `json:"id" tfsdk:"id"`
	Location          string   `json:"location" tfsdk:"location"`
	Provider          string   `json:"provider" tfsdk:"provider"`
	PublicIpAddresses []string `json:"public_ip_addresses" tfsdk:"public_ip_addresses"`
	Slug              string   `json:"slug" tfsdk:"slug"`
}
type EnableSafeMigrationsForBranchRes200_RestoredFromBranch struct {
	CreatedAt string `json:"created_at" tfsdk:"created_at"`
	DeletedAt string `json:"deleted_at" tfsdk:"deleted_at"`
	Id        string `json:"id" tfsdk:"id"`
	Name      string `json:"name" tfsdk:"name"`
	UpdatedAt string `json:"updated_at" tfsdk:"updated_at"`
}
type EnableSafeMigrationsForBranchRes200 struct {
	AccessHostUrl               *string                                                 `json:"access_host_url,omitempty" tfsdk:"access_host_url"`
	ApiActor                    *EnableSafeMigrationsForBranchRes200_ApiActor           `json:"api_actor,omitempty" tfsdk:"api_actor"`
	CreatedAt                   string                                                  `json:"created_at" tfsdk:"created_at"`
	HtmlUrl                     string                                                  `json:"html_url" tfsdk:"html_url"`
	Id                          string                                                  `json:"id" tfsdk:"id"`
	InitialRestoreId            *string                                                 `json:"initial_restore_id,omitempty" tfsdk:"initial_restore_id"`
	MysqlAddress                string                                                  `json:"mysql_address" tfsdk:"mysql_address"`
	MysqlEdgeAddress            string                                                  `json:"mysql_edge_address" tfsdk:"mysql_edge_address"`
	Name                        string                                                  `json:"name" tfsdk:"name"`
	ParentBranch                string                                                  `json:"parent_branch" tfsdk:"parent_branch"`
	PlanetscaleRegion           *EnableSafeMigrationsForBranchRes200_PlanetscaleRegion  `json:"planetscale_region,omitempty" tfsdk:"planetscale_region"`
	Production                  bool                                                    `json:"production" tfsdk:"production"`
	Ready                       bool                                                    `json:"ready" tfsdk:"ready"`
	RestoreChecklistCompletedAt *string                                                 `json:"restore_checklist_completed_at,omitempty" tfsdk:"restore_checklist_completed_at"`
	RestoredFromBranch          *EnableSafeMigrationsForBranchRes200_RestoredFromBranch `json:"restored_from_branch,omitempty" tfsdk:"restored_from_branch"`
	SchemaLastUpdatedAt         string                                                  `json:"schema_last_updated_at" tfsdk:"schema_last_updated_at"`
	ShardCount                  *float64                                                `json:"shard_count,omitempty" tfsdk:"shard_count"`
	Sharded                     bool                                                    `json:"sharded" tfsdk:"sharded"`
	UpdatedAt                   string                                                  `json:"updated_at" tfsdk:"updated_at"`
}

func (cl *Client) EnableSafeMigrationsForBranch(ctx context.Context, organization string, database string, name string) (res200 *EnableSafeMigrationsForBranchRes200, res403 *EnableSafeMigrationsForBranchRes403, res404 *EnableSafeMigrationsForBranchRes404, res500 *EnableSafeMigrationsForBranchRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + name + "/safe-migrations"})
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), nil)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(EnableSafeMigrationsForBranchRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(EnableSafeMigrationsForBranchRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(EnableSafeMigrationsForBranchRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(EnableSafeMigrationsForBranchRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type DisableSafeMigrationsForBranchRes403 struct{}
type DisableSafeMigrationsForBranchRes500 struct{}
type DisableSafeMigrationsForBranchRes200_ApiActor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type DisableSafeMigrationsForBranchRes200_PlanetscaleRegion struct {
	DisplayName       string   `json:"display_name" tfsdk:"display_name"`
	Enabled           bool     `json:"enabled" tfsdk:"enabled"`
	Id                string   `json:"id" tfsdk:"id"`
	Location          string   `json:"location" tfsdk:"location"`
	Provider          string   `json:"provider" tfsdk:"provider"`
	PublicIpAddresses []string `json:"public_ip_addresses" tfsdk:"public_ip_addresses"`
	Slug              string   `json:"slug" tfsdk:"slug"`
}
type DisableSafeMigrationsForBranchRes200_RestoredFromBranch struct {
	CreatedAt string `json:"created_at" tfsdk:"created_at"`
	DeletedAt string `json:"deleted_at" tfsdk:"deleted_at"`
	Id        string `json:"id" tfsdk:"id"`
	Name      string `json:"name" tfsdk:"name"`
	UpdatedAt string `json:"updated_at" tfsdk:"updated_at"`
}
type DisableSafeMigrationsForBranchRes200 struct {
	AccessHostUrl               *string                                                  `json:"access_host_url,omitempty" tfsdk:"access_host_url"`
	ApiActor                    *DisableSafeMigrationsForBranchRes200_ApiActor           `json:"api_actor,omitempty" tfsdk:"api_actor"`
	CreatedAt                   string                                                   `json:"created_at" tfsdk:"created_at"`
	HtmlUrl                     string                                                   `json:"html_url" tfsdk:"html_url"`
	Id                          string                                                   `json:"id" tfsdk:"id"`
	InitialRestoreId            *string                                                  `json:"initial_restore_id,omitempty" tfsdk:"initial_restore_id"`
	MysqlAddress                string                                                   `json:"mysql_address" tfsdk:"mysql_address"`
	MysqlEdgeAddress            string                                                   `json:"mysql_edge_address" tfsdk:"mysql_edge_address"`
	Name                        string                                                   `json:"name" tfsdk:"name"`
	ParentBranch                string                                                   `json:"parent_branch" tfsdk:"parent_branch"`
	PlanetscaleRegion           *DisableSafeMigrationsForBranchRes200_PlanetscaleRegion  `json:"planetscale_region,omitempty" tfsdk:"planetscale_region"`
	Production                  bool                                                     `json:"production" tfsdk:"production"`
	Ready                       bool                                                     `json:"ready" tfsdk:"ready"`
	RestoreChecklistCompletedAt *string                                                  `json:"restore_checklist_completed_at,omitempty" tfsdk:"restore_checklist_completed_at"`
	RestoredFromBranch          *DisableSafeMigrationsForBranchRes200_RestoredFromBranch `json:"restored_from_branch,omitempty" tfsdk:"restored_from_branch"`
	SchemaLastUpdatedAt         string                                                   `json:"schema_last_updated_at" tfsdk:"schema_last_updated_at"`
	ShardCount                  *float64                                                 `json:"shard_count,omitempty" tfsdk:"shard_count"`
	Sharded                     bool                                                     `json:"sharded" tfsdk:"sharded"`
	UpdatedAt                   string                                                   `json:"updated_at" tfsdk:"updated_at"`
}
type DisableSafeMigrationsForBranchRes404 struct{}

func (cl *Client) DisableSafeMigrationsForBranch(ctx context.Context, organization string, database string, name string) (res200 *DisableSafeMigrationsForBranchRes200, res403 *DisableSafeMigrationsForBranchRes403, res404 *DisableSafeMigrationsForBranchRes404, res500 *DisableSafeMigrationsForBranchRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + name + "/safe-migrations"})
	r, err := http.NewRequestWithContext(ctx, "DELETE", u.String(), nil)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(DisableSafeMigrationsForBranchRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(DisableSafeMigrationsForBranchRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(DisableSafeMigrationsForBranchRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(DisableSafeMigrationsForBranchRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type GetBranchSchemaRes403 struct{}
type GetBranchSchemaRes500 struct{}
type GetBranchSchemaRes200_DataItem struct {
	Html string `json:"html" tfsdk:"html"`
	Name string `json:"name" tfsdk:"name"`
	Raw  string `json:"raw" tfsdk:"raw"`
}
type GetBranchSchemaRes200 struct {
	Data []GetBranchSchemaRes200_DataItem `json:"data" tfsdk:"data"`
}
type GetBranchSchemaRes404 struct{}

func (cl *Client) GetBranchSchema(ctx context.Context, organization string, database string, name string, keyspace *string) (res200 *GetBranchSchemaRes200, res403 *GetBranchSchemaRes403, res404 *GetBranchSchemaRes404, res500 *GetBranchSchemaRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + name + "/schema"})
	q := u.Query()
	if keyspace != nil {
		q.Set("keyspace", *keyspace)
	}
	u.RawQuery = q.Encode()
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(GetBranchSchemaRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(GetBranchSchemaRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(GetBranchSchemaRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(GetBranchSchemaRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type LintBranchSchemaRes500 struct{}
type LintBranchSchemaRes200_DataItem struct {
	AutoIncrementColumnNames []string `json:"auto_increment_column_names" tfsdk:"auto_increment_column_names"`
	CharsetName              string   `json:"charset_name" tfsdk:"charset_name"`
	CheckConstraintName      string   `json:"check_constraint_name" tfsdk:"check_constraint_name"`
	ColumnName               string   `json:"column_name" tfsdk:"column_name"`
	DocsUrl                  string   `json:"docs_url" tfsdk:"docs_url"`
	EngineName               string   `json:"engine_name" tfsdk:"engine_name"`
	EnumValue                string   `json:"enum_value" tfsdk:"enum_value"`
	ErrorDescription         string   `json:"error_description" tfsdk:"error_description"`
	ForeignKeyColumnNames    []string `json:"foreign_key_column_names" tfsdk:"foreign_key_column_names"`
	JsonPath                 string   `json:"json_path" tfsdk:"json_path"`
	KeyspaceName             string   `json:"keyspace_name" tfsdk:"keyspace_name"`
	LintError                string   `json:"lint_error" tfsdk:"lint_error"`
	PartitionName            string   `json:"partition_name" tfsdk:"partition_name"`
	PartitioningType         string   `json:"partitioning_type" tfsdk:"partitioning_type"`
	SubjectType              string   `json:"subject_type" tfsdk:"subject_type"`
	TableName                string   `json:"table_name" tfsdk:"table_name"`
	VindexName               string   `json:"vindex_name" tfsdk:"vindex_name"`
}
type LintBranchSchemaRes200 struct {
	Data []LintBranchSchemaRes200_DataItem `json:"data" tfsdk:"data"`
}
type LintBranchSchemaRes404 struct{}
type LintBranchSchemaRes403 struct{}

func (cl *Client) LintBranchSchema(ctx context.Context, organization string, database string, name string, page *int, perPage *int) (res200 *LintBranchSchemaRes200, res403 *LintBranchSchemaRes403, res404 *LintBranchSchemaRes404, res500 *LintBranchSchemaRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + name + "/schema/lint"})
	q := u.Query()
	if page != nil {
		q.Set("page", strconv.Itoa(*page))
	}
	if perPage != nil {
		q.Set("per_page", strconv.Itoa(*perPage))
	}
	u.RawQuery = q.Encode()
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(LintBranchSchemaRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(LintBranchSchemaRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(LintBranchSchemaRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(LintBranchSchemaRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type GetTheDeployQueueRes200_DataItem struct {
	AutoCutover       bool    `json:"auto_cutover" tfsdk:"auto_cutover"`
	CreatedAt         string  `json:"created_at" tfsdk:"created_at"`
	CutoverAt         *string `json:"cutover_at,omitempty" tfsdk:"cutover_at"`
	CutoverExpiring   bool    `json:"cutover_expiring" tfsdk:"cutover_expiring"`
	DeployCheckErrors *string `json:"deploy_check_errors,omitempty" tfsdk:"deploy_check_errors"`
	FinishedAt        *string `json:"finished_at,omitempty" tfsdk:"finished_at"`
	Id                string  `json:"id" tfsdk:"id"`
	QueuedAt          *string `json:"queued_at,omitempty" tfsdk:"queued_at"`
	ReadyToCutoverAt  *string `json:"ready_to_cutover_at,omitempty" tfsdk:"ready_to_cutover_at"`
	StartedAt         *string `json:"started_at,omitempty" tfsdk:"started_at"`
	State             string  `json:"state" tfsdk:"state"`
	SubmittedAt       string  `json:"submitted_at" tfsdk:"submitted_at"`
	UpdatedAt         string  `json:"updated_at" tfsdk:"updated_at"`
}
type GetTheDeployQueueRes200 struct {
	Data []GetTheDeployQueueRes200_DataItem `json:"data" tfsdk:"data"`
}

func (cl *Client) GetTheDeployQueue(ctx context.Context, organization string, database string) (res200 *GetTheDeployQueueRes200, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-queue"})
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(GetTheDeployQueueRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type ListDeployRequestsRes200_DataItem_Actor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type ListDeployRequestsRes200_DataItem_BranchDeletedBy struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type ListDeployRequestsRes200_DataItem_ClosedBy struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type ListDeployRequestsRes200_DataItem struct {
	Actor                ListDeployRequestsRes200_DataItem_Actor           `json:"actor" tfsdk:"actor"`
	Approved             bool                                              `json:"approved" tfsdk:"approved"`
	Branch               string                                            `json:"branch" tfsdk:"branch"`
	BranchDeleted        bool                                              `json:"branch_deleted" tfsdk:"branch_deleted"`
	BranchDeletedAt      string                                            `json:"branch_deleted_at" tfsdk:"branch_deleted_at"`
	BranchDeletedBy      ListDeployRequestsRes200_DataItem_BranchDeletedBy `json:"branch_deleted_by" tfsdk:"branch_deleted_by"`
	ClosedAt             string                                            `json:"closed_at" tfsdk:"closed_at"`
	ClosedBy             ListDeployRequestsRes200_DataItem_ClosedBy        `json:"closed_by" tfsdk:"closed_by"`
	CreatedAt            string                                            `json:"created_at" tfsdk:"created_at"`
	DeployedAt           string                                            `json:"deployed_at" tfsdk:"deployed_at"`
	DeploymentState      string                                            `json:"deployment_state" tfsdk:"deployment_state"`
	HtmlBody             string                                            `json:"html_body" tfsdk:"html_body"`
	HtmlUrl              string                                            `json:"html_url" tfsdk:"html_url"`
	Id                   string                                            `json:"id" tfsdk:"id"`
	IntoBranch           string                                            `json:"into_branch" tfsdk:"into_branch"`
	IntoBranchShardCount float64                                           `json:"into_branch_shard_count" tfsdk:"into_branch_shard_count"`
	IntoBranchSharded    bool                                              `json:"into_branch_sharded" tfsdk:"into_branch_sharded"`
	Notes                string                                            `json:"notes" tfsdk:"notes"`
	Number               float64                                           `json:"number" tfsdk:"number"`
	State                string                                            `json:"state" tfsdk:"state"`
	UpdatedAt            string                                            `json:"updated_at" tfsdk:"updated_at"`
}
type ListDeployRequestsRes200 struct {
	Data []ListDeployRequestsRes200_DataItem `json:"data" tfsdk:"data"`
}

func (cl *Client) ListDeployRequests(ctx context.Context, organization string, database string, page *int, perPage *int, state *string, branch *string, intoBranch *string) (res200 *ListDeployRequestsRes200, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests"})
	q := u.Query()
	if page != nil {
		q.Set("page", strconv.Itoa(*page))
	}
	if perPage != nil {
		q.Set("per_page", strconv.Itoa(*perPage))
	}
	if state != nil {
		q.Set("state", *state)
	}
	if branch != nil {
		q.Set("branch", *branch)
	}
	if intoBranch != nil {
		q.Set("into_branch", *intoBranch)
	}
	u.RawQuery = q.Encode()
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListDeployRequestsRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type CreateDeployRequestReq struct {
	Branch     *string `json:"branch,omitempty" tfsdk:"branch"`
	IntoBranch *string `json:"into_branch,omitempty" tfsdk:"into_branch"`
	Notes      *string `json:"notes,omitempty" tfsdk:"notes"`
}
type CreateDeployRequestRes201_Actor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type CreateDeployRequestRes201_BranchDeletedBy struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type CreateDeployRequestRes201_ClosedBy struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type CreateDeployRequestRes201_Deployment struct {
	AutoCutover       bool    `json:"auto_cutover" tfsdk:"auto_cutover"`
	CreatedAt         string  `json:"created_at" tfsdk:"created_at"`
	CutoverAt         *string `json:"cutover_at,omitempty" tfsdk:"cutover_at"`
	CutoverExpiring   bool    `json:"cutover_expiring" tfsdk:"cutover_expiring"`
	DeployCheckErrors *string `json:"deploy_check_errors,omitempty" tfsdk:"deploy_check_errors"`
	FinishedAt        *string `json:"finished_at,omitempty" tfsdk:"finished_at"`
	Id                string  `json:"id" tfsdk:"id"`
	QueuedAt          *string `json:"queued_at,omitempty" tfsdk:"queued_at"`
	ReadyToCutoverAt  *string `json:"ready_to_cutover_at,omitempty" tfsdk:"ready_to_cutover_at"`
	StartedAt         *string `json:"started_at,omitempty" tfsdk:"started_at"`
	State             string  `json:"state" tfsdk:"state"`
	SubmittedAt       string  `json:"submitted_at" tfsdk:"submitted_at"`
	UpdatedAt         string  `json:"updated_at" tfsdk:"updated_at"`
}
type CreateDeployRequestRes201 struct {
	Actor                CreateDeployRequestRes201_Actor           `json:"actor" tfsdk:"actor"`
	Approved             bool                                      `json:"approved" tfsdk:"approved"`
	Branch               string                                    `json:"branch" tfsdk:"branch"`
	BranchDeleted        bool                                      `json:"branch_deleted" tfsdk:"branch_deleted"`
	BranchDeletedAt      string                                    `json:"branch_deleted_at" tfsdk:"branch_deleted_at"`
	BranchDeletedBy      CreateDeployRequestRes201_BranchDeletedBy `json:"branch_deleted_by" tfsdk:"branch_deleted_by"`
	ClosedAt             string                                    `json:"closed_at" tfsdk:"closed_at"`
	ClosedBy             CreateDeployRequestRes201_ClosedBy        `json:"closed_by" tfsdk:"closed_by"`
	CreatedAt            string                                    `json:"created_at" tfsdk:"created_at"`
	DeployedAt           string                                    `json:"deployed_at" tfsdk:"deployed_at"`
	Deployment           CreateDeployRequestRes201_Deployment      `json:"deployment" tfsdk:"deployment"`
	DeploymentState      string                                    `json:"deployment_state" tfsdk:"deployment_state"`
	HtmlBody             string                                    `json:"html_body" tfsdk:"html_body"`
	HtmlUrl              string                                    `json:"html_url" tfsdk:"html_url"`
	Id                   string                                    `json:"id" tfsdk:"id"`
	IntoBranch           string                                    `json:"into_branch" tfsdk:"into_branch"`
	IntoBranchShardCount float64                                   `json:"into_branch_shard_count" tfsdk:"into_branch_shard_count"`
	IntoBranchSharded    bool                                      `json:"into_branch_sharded" tfsdk:"into_branch_sharded"`
	Notes                string                                    `json:"notes" tfsdk:"notes"`
	Number               float64                                   `json:"number" tfsdk:"number"`
	State                string                                    `json:"state" tfsdk:"state"`
	UpdatedAt            string                                    `json:"updated_at" tfsdk:"updated_at"`
}

func (cl *Client) CreateDeployRequest(ctx context.Context, organization string, database string, req CreateDeployRequestReq) (res201 *CreateDeployRequestRes201, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests"})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res201, err
	}
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return res201, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res201, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 201:
		res201 = new(CreateDeployRequestRes201)
		err = json.NewDecoder(res.Body).Decode(&res201)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res201, err
}

type GetDeployRequestRes200_Actor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type GetDeployRequestRes200_BranchDeletedBy struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type GetDeployRequestRes200_ClosedBy struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type GetDeployRequestRes200_Deployment struct {
	AutoCutover       bool    `json:"auto_cutover" tfsdk:"auto_cutover"`
	CreatedAt         string  `json:"created_at" tfsdk:"created_at"`
	CutoverAt         *string `json:"cutover_at,omitempty" tfsdk:"cutover_at"`
	CutoverExpiring   bool    `json:"cutover_expiring" tfsdk:"cutover_expiring"`
	DeployCheckErrors *string `json:"deploy_check_errors,omitempty" tfsdk:"deploy_check_errors"`
	FinishedAt        *string `json:"finished_at,omitempty" tfsdk:"finished_at"`
	Id                string  `json:"id" tfsdk:"id"`
	QueuedAt          *string `json:"queued_at,omitempty" tfsdk:"queued_at"`
	ReadyToCutoverAt  *string `json:"ready_to_cutover_at,omitempty" tfsdk:"ready_to_cutover_at"`
	StartedAt         *string `json:"started_at,omitempty" tfsdk:"started_at"`
	State             string  `json:"state" tfsdk:"state"`
	SubmittedAt       string  `json:"submitted_at" tfsdk:"submitted_at"`
	UpdatedAt         string  `json:"updated_at" tfsdk:"updated_at"`
}
type GetDeployRequestRes200 struct {
	Actor                GetDeployRequestRes200_Actor           `json:"actor" tfsdk:"actor"`
	Approved             bool                                   `json:"approved" tfsdk:"approved"`
	Branch               string                                 `json:"branch" tfsdk:"branch"`
	BranchDeleted        bool                                   `json:"branch_deleted" tfsdk:"branch_deleted"`
	BranchDeletedAt      string                                 `json:"branch_deleted_at" tfsdk:"branch_deleted_at"`
	BranchDeletedBy      GetDeployRequestRes200_BranchDeletedBy `json:"branch_deleted_by" tfsdk:"branch_deleted_by"`
	ClosedAt             string                                 `json:"closed_at" tfsdk:"closed_at"`
	ClosedBy             GetDeployRequestRes200_ClosedBy        `json:"closed_by" tfsdk:"closed_by"`
	CreatedAt            string                                 `json:"created_at" tfsdk:"created_at"`
	DeployedAt           string                                 `json:"deployed_at" tfsdk:"deployed_at"`
	Deployment           GetDeployRequestRes200_Deployment      `json:"deployment" tfsdk:"deployment"`
	DeploymentState      string                                 `json:"deployment_state" tfsdk:"deployment_state"`
	HtmlBody             string                                 `json:"html_body" tfsdk:"html_body"`
	HtmlUrl              string                                 `json:"html_url" tfsdk:"html_url"`
	Id                   string                                 `json:"id" tfsdk:"id"`
	IntoBranch           string                                 `json:"into_branch" tfsdk:"into_branch"`
	IntoBranchShardCount float64                                `json:"into_branch_shard_count" tfsdk:"into_branch_shard_count"`
	IntoBranchSharded    bool                                   `json:"into_branch_sharded" tfsdk:"into_branch_sharded"`
	Notes                string                                 `json:"notes" tfsdk:"notes"`
	Number               float64                                `json:"number" tfsdk:"number"`
	State                string                                 `json:"state" tfsdk:"state"`
	UpdatedAt            string                                 `json:"updated_at" tfsdk:"updated_at"`
}

func (cl *Client) GetDeployRequest(ctx context.Context, organization string, database string, number string) (res200 *GetDeployRequestRes200, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests/" + number})
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(GetDeployRequestRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type CloseDeployRequestReq struct {
	State *string `json:"state,omitempty" tfsdk:"state"`
}
type CloseDeployRequestRes200_Actor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type CloseDeployRequestRes200_BranchDeletedBy struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type CloseDeployRequestRes200_ClosedBy struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type CloseDeployRequestRes200_Deployment struct {
	AutoCutover       bool    `json:"auto_cutover" tfsdk:"auto_cutover"`
	CreatedAt         string  `json:"created_at" tfsdk:"created_at"`
	CutoverAt         *string `json:"cutover_at,omitempty" tfsdk:"cutover_at"`
	CutoverExpiring   bool    `json:"cutover_expiring" tfsdk:"cutover_expiring"`
	DeployCheckErrors *string `json:"deploy_check_errors,omitempty" tfsdk:"deploy_check_errors"`
	FinishedAt        *string `json:"finished_at,omitempty" tfsdk:"finished_at"`
	Id                string  `json:"id" tfsdk:"id"`
	QueuedAt          *string `json:"queued_at,omitempty" tfsdk:"queued_at"`
	ReadyToCutoverAt  *string `json:"ready_to_cutover_at,omitempty" tfsdk:"ready_to_cutover_at"`
	StartedAt         *string `json:"started_at,omitempty" tfsdk:"started_at"`
	State             string  `json:"state" tfsdk:"state"`
	SubmittedAt       string  `json:"submitted_at" tfsdk:"submitted_at"`
	UpdatedAt         string  `json:"updated_at" tfsdk:"updated_at"`
}
type CloseDeployRequestRes200 struct {
	Actor                CloseDeployRequestRes200_Actor           `json:"actor" tfsdk:"actor"`
	Approved             bool                                     `json:"approved" tfsdk:"approved"`
	Branch               string                                   `json:"branch" tfsdk:"branch"`
	BranchDeleted        bool                                     `json:"branch_deleted" tfsdk:"branch_deleted"`
	BranchDeletedAt      string                                   `json:"branch_deleted_at" tfsdk:"branch_deleted_at"`
	BranchDeletedBy      CloseDeployRequestRes200_BranchDeletedBy `json:"branch_deleted_by" tfsdk:"branch_deleted_by"`
	ClosedAt             string                                   `json:"closed_at" tfsdk:"closed_at"`
	ClosedBy             CloseDeployRequestRes200_ClosedBy        `json:"closed_by" tfsdk:"closed_by"`
	CreatedAt            string                                   `json:"created_at" tfsdk:"created_at"`
	DeployedAt           string                                   `json:"deployed_at" tfsdk:"deployed_at"`
	Deployment           CloseDeployRequestRes200_Deployment      `json:"deployment" tfsdk:"deployment"`
	DeploymentState      string                                   `json:"deployment_state" tfsdk:"deployment_state"`
	HtmlBody             string                                   `json:"html_body" tfsdk:"html_body"`
	HtmlUrl              string                                   `json:"html_url" tfsdk:"html_url"`
	Id                   string                                   `json:"id" tfsdk:"id"`
	IntoBranch           string                                   `json:"into_branch" tfsdk:"into_branch"`
	IntoBranchShardCount float64                                  `json:"into_branch_shard_count" tfsdk:"into_branch_shard_count"`
	IntoBranchSharded    bool                                     `json:"into_branch_sharded" tfsdk:"into_branch_sharded"`
	Notes                string                                   `json:"notes" tfsdk:"notes"`
	Number               float64                                  `json:"number" tfsdk:"number"`
	State                string                                   `json:"state" tfsdk:"state"`
	UpdatedAt            string                                   `json:"updated_at" tfsdk:"updated_at"`
}

func (cl *Client) CloseDeployRequest(ctx context.Context, organization string, database string, number string, req CloseDeployRequestReq) (res200 *CloseDeployRequestRes200, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests/" + number})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res200, err
	}
	r, err := http.NewRequestWithContext(ctx, "PATCH", u.String(), body)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(CloseDeployRequestRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type CompleteGatedDeployRequestRes200_Actor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type CompleteGatedDeployRequestRes200_BranchDeletedBy struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type CompleteGatedDeployRequestRes200_ClosedBy struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type CompleteGatedDeployRequestRes200 struct {
	Actor                CompleteGatedDeployRequestRes200_Actor           `json:"actor" tfsdk:"actor"`
	Approved             bool                                             `json:"approved" tfsdk:"approved"`
	Branch               string                                           `json:"branch" tfsdk:"branch"`
	BranchDeleted        bool                                             `json:"branch_deleted" tfsdk:"branch_deleted"`
	BranchDeletedAt      string                                           `json:"branch_deleted_at" tfsdk:"branch_deleted_at"`
	BranchDeletedBy      CompleteGatedDeployRequestRes200_BranchDeletedBy `json:"branch_deleted_by" tfsdk:"branch_deleted_by"`
	ClosedAt             string                                           `json:"closed_at" tfsdk:"closed_at"`
	ClosedBy             CompleteGatedDeployRequestRes200_ClosedBy        `json:"closed_by" tfsdk:"closed_by"`
	CreatedAt            string                                           `json:"created_at" tfsdk:"created_at"`
	DeployedAt           string                                           `json:"deployed_at" tfsdk:"deployed_at"`
	DeploymentState      string                                           `json:"deployment_state" tfsdk:"deployment_state"`
	HtmlBody             string                                           `json:"html_body" tfsdk:"html_body"`
	HtmlUrl              string                                           `json:"html_url" tfsdk:"html_url"`
	Id                   string                                           `json:"id" tfsdk:"id"`
	IntoBranch           string                                           `json:"into_branch" tfsdk:"into_branch"`
	IntoBranchShardCount float64                                          `json:"into_branch_shard_count" tfsdk:"into_branch_shard_count"`
	IntoBranchSharded    bool                                             `json:"into_branch_sharded" tfsdk:"into_branch_sharded"`
	Notes                string                                           `json:"notes" tfsdk:"notes"`
	Number               float64                                          `json:"number" tfsdk:"number"`
	State                string                                           `json:"state" tfsdk:"state"`
	UpdatedAt            string                                           `json:"updated_at" tfsdk:"updated_at"`
}

func (cl *Client) CompleteGatedDeployRequest(ctx context.Context, organization string, database string, number string) (res200 *CompleteGatedDeployRequestRes200, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests/" + number + "/apply-deploy"})
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(CompleteGatedDeployRequestRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type UpdateAutoApplyForDeployRequestReq struct {
	Enable *bool `json:"enable,omitempty" tfsdk:"enable"`
}
type UpdateAutoApplyForDeployRequestRes200_Actor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type UpdateAutoApplyForDeployRequestRes200_BranchDeletedBy struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type UpdateAutoApplyForDeployRequestRes200_ClosedBy struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type UpdateAutoApplyForDeployRequestRes200 struct {
	Actor                UpdateAutoApplyForDeployRequestRes200_Actor           `json:"actor" tfsdk:"actor"`
	Approved             bool                                                  `json:"approved" tfsdk:"approved"`
	Branch               string                                                `json:"branch" tfsdk:"branch"`
	BranchDeleted        bool                                                  `json:"branch_deleted" tfsdk:"branch_deleted"`
	BranchDeletedAt      string                                                `json:"branch_deleted_at" tfsdk:"branch_deleted_at"`
	BranchDeletedBy      UpdateAutoApplyForDeployRequestRes200_BranchDeletedBy `json:"branch_deleted_by" tfsdk:"branch_deleted_by"`
	ClosedAt             string                                                `json:"closed_at" tfsdk:"closed_at"`
	ClosedBy             UpdateAutoApplyForDeployRequestRes200_ClosedBy        `json:"closed_by" tfsdk:"closed_by"`
	CreatedAt            string                                                `json:"created_at" tfsdk:"created_at"`
	DeployedAt           string                                                `json:"deployed_at" tfsdk:"deployed_at"`
	DeploymentState      string                                                `json:"deployment_state" tfsdk:"deployment_state"`
	HtmlBody             string                                                `json:"html_body" tfsdk:"html_body"`
	HtmlUrl              string                                                `json:"html_url" tfsdk:"html_url"`
	Id                   string                                                `json:"id" tfsdk:"id"`
	IntoBranch           string                                                `json:"into_branch" tfsdk:"into_branch"`
	IntoBranchShardCount float64                                               `json:"into_branch_shard_count" tfsdk:"into_branch_shard_count"`
	IntoBranchSharded    bool                                                  `json:"into_branch_sharded" tfsdk:"into_branch_sharded"`
	Notes                string                                                `json:"notes" tfsdk:"notes"`
	Number               float64                                               `json:"number" tfsdk:"number"`
	State                string                                                `json:"state" tfsdk:"state"`
	UpdatedAt            string                                                `json:"updated_at" tfsdk:"updated_at"`
}

func (cl *Client) UpdateAutoApplyForDeployRequest(ctx context.Context, organization string, database string, number string, req UpdateAutoApplyForDeployRequestReq) (res200 *UpdateAutoApplyForDeployRequestRes200, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests/" + number + "/auto-apply"})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res200, err
	}
	r, err := http.NewRequestWithContext(ctx, "PUT", u.String(), body)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(UpdateAutoApplyForDeployRequestRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type CancelQueuedDeployRequestRes200_Actor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type CancelQueuedDeployRequestRes200_BranchDeletedBy struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type CancelQueuedDeployRequestRes200_ClosedBy struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type CancelQueuedDeployRequestRes200 struct {
	Actor                CancelQueuedDeployRequestRes200_Actor           `json:"actor" tfsdk:"actor"`
	Approved             bool                                            `json:"approved" tfsdk:"approved"`
	Branch               string                                          `json:"branch" tfsdk:"branch"`
	BranchDeleted        bool                                            `json:"branch_deleted" tfsdk:"branch_deleted"`
	BranchDeletedAt      string                                          `json:"branch_deleted_at" tfsdk:"branch_deleted_at"`
	BranchDeletedBy      CancelQueuedDeployRequestRes200_BranchDeletedBy `json:"branch_deleted_by" tfsdk:"branch_deleted_by"`
	ClosedAt             string                                          `json:"closed_at" tfsdk:"closed_at"`
	ClosedBy             CancelQueuedDeployRequestRes200_ClosedBy        `json:"closed_by" tfsdk:"closed_by"`
	CreatedAt            string                                          `json:"created_at" tfsdk:"created_at"`
	DeployedAt           string                                          `json:"deployed_at" tfsdk:"deployed_at"`
	DeploymentState      string                                          `json:"deployment_state" tfsdk:"deployment_state"`
	HtmlBody             string                                          `json:"html_body" tfsdk:"html_body"`
	HtmlUrl              string                                          `json:"html_url" tfsdk:"html_url"`
	Id                   string                                          `json:"id" tfsdk:"id"`
	IntoBranch           string                                          `json:"into_branch" tfsdk:"into_branch"`
	IntoBranchShardCount float64                                         `json:"into_branch_shard_count" tfsdk:"into_branch_shard_count"`
	IntoBranchSharded    bool                                            `json:"into_branch_sharded" tfsdk:"into_branch_sharded"`
	Notes                string                                          `json:"notes" tfsdk:"notes"`
	Number               float64                                         `json:"number" tfsdk:"number"`
	State                string                                          `json:"state" tfsdk:"state"`
	UpdatedAt            string                                          `json:"updated_at" tfsdk:"updated_at"`
}

func (cl *Client) CancelQueuedDeployRequest(ctx context.Context, organization string, database string, number string) (res200 *CancelQueuedDeployRequestRes200, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests/" + number + "/cancel"})
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(CancelQueuedDeployRequestRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type CompleteErroredDeployRes200_Actor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type CompleteErroredDeployRes200_BranchDeletedBy struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type CompleteErroredDeployRes200_ClosedBy struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type CompleteErroredDeployRes200 struct {
	Actor                CompleteErroredDeployRes200_Actor           `json:"actor" tfsdk:"actor"`
	Approved             bool                                        `json:"approved" tfsdk:"approved"`
	Branch               string                                      `json:"branch" tfsdk:"branch"`
	BranchDeleted        bool                                        `json:"branch_deleted" tfsdk:"branch_deleted"`
	BranchDeletedAt      string                                      `json:"branch_deleted_at" tfsdk:"branch_deleted_at"`
	BranchDeletedBy      CompleteErroredDeployRes200_BranchDeletedBy `json:"branch_deleted_by" tfsdk:"branch_deleted_by"`
	ClosedAt             string                                      `json:"closed_at" tfsdk:"closed_at"`
	ClosedBy             CompleteErroredDeployRes200_ClosedBy        `json:"closed_by" tfsdk:"closed_by"`
	CreatedAt            string                                      `json:"created_at" tfsdk:"created_at"`
	DeployedAt           string                                      `json:"deployed_at" tfsdk:"deployed_at"`
	DeploymentState      string                                      `json:"deployment_state" tfsdk:"deployment_state"`
	HtmlBody             string                                      `json:"html_body" tfsdk:"html_body"`
	HtmlUrl              string                                      `json:"html_url" tfsdk:"html_url"`
	Id                   string                                      `json:"id" tfsdk:"id"`
	IntoBranch           string                                      `json:"into_branch" tfsdk:"into_branch"`
	IntoBranchShardCount float64                                     `json:"into_branch_shard_count" tfsdk:"into_branch_shard_count"`
	IntoBranchSharded    bool                                        `json:"into_branch_sharded" tfsdk:"into_branch_sharded"`
	Notes                string                                      `json:"notes" tfsdk:"notes"`
	Number               float64                                     `json:"number" tfsdk:"number"`
	State                string                                      `json:"state" tfsdk:"state"`
	UpdatedAt            string                                      `json:"updated_at" tfsdk:"updated_at"`
}

func (cl *Client) CompleteErroredDeploy(ctx context.Context, organization string, database string, number string) (res200 *CompleteErroredDeployRes200, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests/" + number + "/complete-deploy"})
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(CompleteErroredDeployRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type QueueDeployRequestRes200_Actor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type QueueDeployRequestRes200_BranchDeletedBy struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type QueueDeployRequestRes200_ClosedBy struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type QueueDeployRequestRes200 struct {
	Actor                QueueDeployRequestRes200_Actor           `json:"actor" tfsdk:"actor"`
	Approved             bool                                     `json:"approved" tfsdk:"approved"`
	Branch               string                                   `json:"branch" tfsdk:"branch"`
	BranchDeleted        bool                                     `json:"branch_deleted" tfsdk:"branch_deleted"`
	BranchDeletedAt      string                                   `json:"branch_deleted_at" tfsdk:"branch_deleted_at"`
	BranchDeletedBy      QueueDeployRequestRes200_BranchDeletedBy `json:"branch_deleted_by" tfsdk:"branch_deleted_by"`
	ClosedAt             string                                   `json:"closed_at" tfsdk:"closed_at"`
	ClosedBy             QueueDeployRequestRes200_ClosedBy        `json:"closed_by" tfsdk:"closed_by"`
	CreatedAt            string                                   `json:"created_at" tfsdk:"created_at"`
	DeployedAt           string                                   `json:"deployed_at" tfsdk:"deployed_at"`
	DeploymentState      string                                   `json:"deployment_state" tfsdk:"deployment_state"`
	HtmlBody             string                                   `json:"html_body" tfsdk:"html_body"`
	HtmlUrl              string                                   `json:"html_url" tfsdk:"html_url"`
	Id                   string                                   `json:"id" tfsdk:"id"`
	IntoBranch           string                                   `json:"into_branch" tfsdk:"into_branch"`
	IntoBranchShardCount float64                                  `json:"into_branch_shard_count" tfsdk:"into_branch_shard_count"`
	IntoBranchSharded    bool                                     `json:"into_branch_sharded" tfsdk:"into_branch_sharded"`
	Notes                string                                   `json:"notes" tfsdk:"notes"`
	Number               float64                                  `json:"number" tfsdk:"number"`
	State                string                                   `json:"state" tfsdk:"state"`
	UpdatedAt            string                                   `json:"updated_at" tfsdk:"updated_at"`
}

func (cl *Client) QueueDeployRequest(ctx context.Context, organization string, database string, number string) (res200 *QueueDeployRequestRes200, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests/" + number + "/deploy"})
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(QueueDeployRequestRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type GetDeploymentRes200 struct {
	AutoCutover       bool    `json:"auto_cutover" tfsdk:"auto_cutover"`
	CreatedAt         string  `json:"created_at" tfsdk:"created_at"`
	CutoverAt         *string `json:"cutover_at,omitempty" tfsdk:"cutover_at"`
	CutoverExpiring   bool    `json:"cutover_expiring" tfsdk:"cutover_expiring"`
	DeployCheckErrors *string `json:"deploy_check_errors,omitempty" tfsdk:"deploy_check_errors"`
	FinishedAt        *string `json:"finished_at,omitempty" tfsdk:"finished_at"`
	Id                string  `json:"id" tfsdk:"id"`
	QueuedAt          *string `json:"queued_at,omitempty" tfsdk:"queued_at"`
	ReadyToCutoverAt  *string `json:"ready_to_cutover_at,omitempty" tfsdk:"ready_to_cutover_at"`
	StartedAt         *string `json:"started_at,omitempty" tfsdk:"started_at"`
	State             string  `json:"state" tfsdk:"state"`
	SubmittedAt       string  `json:"submitted_at" tfsdk:"submitted_at"`
	UpdatedAt         string  `json:"updated_at" tfsdk:"updated_at"`
}

func (cl *Client) GetDeployment(ctx context.Context, organization string, database string, number string) (res200 *GetDeploymentRes200, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests/" + number + "/deployment"})
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(GetDeploymentRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type ListDeployOperationsRes200_DataItem struct {
	CanDropData          bool     `json:"can_drop_data" tfsdk:"can_drop_data"`
	CreatedAt            string   `json:"created_at" tfsdk:"created_at"`
	DdlStatement         string   `json:"ddl_statement" tfsdk:"ddl_statement"`
	DeployErrorDocsUrl   string   `json:"deploy_error_docs_url" tfsdk:"deploy_error_docs_url"`
	DeployErrors         []string `json:"deploy_errors" tfsdk:"deploy_errors"`
	EtaSeconds           float64  `json:"eta_seconds" tfsdk:"eta_seconds"`
	Id                   string   `json:"id" tfsdk:"id"`
	KeyspaceName         string   `json:"keyspace_name" tfsdk:"keyspace_name"`
	OperationName        string   `json:"operation_name" tfsdk:"operation_name"`
	ProgressPercentage   float64  `json:"progress_percentage" tfsdk:"progress_percentage"`
	State                string   `json:"state" tfsdk:"state"`
	SyntaxHighlightedDdl string   `json:"syntax_highlighted_ddl" tfsdk:"syntax_highlighted_ddl"`
	TableName            string   `json:"table_name" tfsdk:"table_name"`
	TableRecentlyUsed    bool     `json:"table_recently_used" tfsdk:"table_recently_used"`
	TableRecentlyUsedAt  string   `json:"table_recently_used_at" tfsdk:"table_recently_used_at"`
	UpdatedAt            string   `json:"updated_at" tfsdk:"updated_at"`
}
type ListDeployOperationsRes200 struct {
	Data []ListDeployOperationsRes200_DataItem `json:"data" tfsdk:"data"`
}

func (cl *Client) ListDeployOperations(ctx context.Context, organization string, database string, number string, page *int, perPage *int) (res200 *ListDeployOperationsRes200, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests/" + number + "/operations"})
	q := u.Query()
	if page != nil {
		q.Set("page", strconv.Itoa(*page))
	}
	if perPage != nil {
		q.Set("per_page", strconv.Itoa(*perPage))
	}
	u.RawQuery = q.Encode()
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListDeployOperationsRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type CompleteRevertRes200_Actor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type CompleteRevertRes200_BranchDeletedBy struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type CompleteRevertRes200_ClosedBy struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type CompleteRevertRes200 struct {
	Actor                CompleteRevertRes200_Actor           `json:"actor" tfsdk:"actor"`
	Approved             bool                                 `json:"approved" tfsdk:"approved"`
	Branch               string                               `json:"branch" tfsdk:"branch"`
	BranchDeleted        bool                                 `json:"branch_deleted" tfsdk:"branch_deleted"`
	BranchDeletedAt      string                               `json:"branch_deleted_at" tfsdk:"branch_deleted_at"`
	BranchDeletedBy      CompleteRevertRes200_BranchDeletedBy `json:"branch_deleted_by" tfsdk:"branch_deleted_by"`
	ClosedAt             string                               `json:"closed_at" tfsdk:"closed_at"`
	ClosedBy             CompleteRevertRes200_ClosedBy        `json:"closed_by" tfsdk:"closed_by"`
	CreatedAt            string                               `json:"created_at" tfsdk:"created_at"`
	DeployedAt           string                               `json:"deployed_at" tfsdk:"deployed_at"`
	DeploymentState      string                               `json:"deployment_state" tfsdk:"deployment_state"`
	HtmlBody             string                               `json:"html_body" tfsdk:"html_body"`
	HtmlUrl              string                               `json:"html_url" tfsdk:"html_url"`
	Id                   string                               `json:"id" tfsdk:"id"`
	IntoBranch           string                               `json:"into_branch" tfsdk:"into_branch"`
	IntoBranchShardCount float64                              `json:"into_branch_shard_count" tfsdk:"into_branch_shard_count"`
	IntoBranchSharded    bool                                 `json:"into_branch_sharded" tfsdk:"into_branch_sharded"`
	Notes                string                               `json:"notes" tfsdk:"notes"`
	Number               float64                              `json:"number" tfsdk:"number"`
	State                string                               `json:"state" tfsdk:"state"`
	UpdatedAt            string                               `json:"updated_at" tfsdk:"updated_at"`
}

func (cl *Client) CompleteRevert(ctx context.Context, organization string, database string, number string) (res200 *CompleteRevertRes200, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests/" + number + "/revert"})
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(CompleteRevertRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type ListDeployRequestReviewsRes200_DataItem_Actor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type ListDeployRequestReviewsRes200_DataItem struct {
	Actor     ListDeployRequestReviewsRes200_DataItem_Actor `json:"actor" tfsdk:"actor"`
	Body      string                                        `json:"body" tfsdk:"body"`
	CreatedAt string                                        `json:"created_at" tfsdk:"created_at"`
	HtmlBody  string                                        `json:"html_body" tfsdk:"html_body"`
	Id        string                                        `json:"id" tfsdk:"id"`
	State     string                                        `json:"state" tfsdk:"state"`
	UpdatedAt string                                        `json:"updated_at" tfsdk:"updated_at"`
}
type ListDeployRequestReviewsRes200 struct {
	Data []ListDeployRequestReviewsRes200_DataItem `json:"data" tfsdk:"data"`
}

func (cl *Client) ListDeployRequestReviews(ctx context.Context, organization string, database string, number string) (res200 *ListDeployRequestReviewsRes200, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests/" + number + "/reviews"})
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListDeployRequestReviewsRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type ReviewDeployRequestReq struct {
	Body  *string `json:"body,omitempty" tfsdk:"body"`
	State *string `json:"state,omitempty" tfsdk:"state"`
}
type ReviewDeployRequestRes201_Actor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type ReviewDeployRequestRes201 struct {
	Actor     ReviewDeployRequestRes201_Actor `json:"actor" tfsdk:"actor"`
	Body      string                          `json:"body" tfsdk:"body"`
	CreatedAt string                          `json:"created_at" tfsdk:"created_at"`
	HtmlBody  string                          `json:"html_body" tfsdk:"html_body"`
	Id        string                          `json:"id" tfsdk:"id"`
	State     string                          `json:"state" tfsdk:"state"`
	UpdatedAt string                          `json:"updated_at" tfsdk:"updated_at"`
}

func (cl *Client) ReviewDeployRequest(ctx context.Context, organization string, database string, number string, req ReviewDeployRequestReq) (res201 *ReviewDeployRequestRes201, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests/" + number + "/reviews"})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res201, err
	}
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return res201, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res201, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 201:
		res201 = new(ReviewDeployRequestRes201)
		err = json.NewDecoder(res.Body).Decode(&res201)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res201, err
}

type SkipRevertPeriodRes200_Actor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type SkipRevertPeriodRes200_BranchDeletedBy struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type SkipRevertPeriodRes200_ClosedBy struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type SkipRevertPeriodRes200 struct {
	Actor                SkipRevertPeriodRes200_Actor           `json:"actor" tfsdk:"actor"`
	Approved             bool                                   `json:"approved" tfsdk:"approved"`
	Branch               string                                 `json:"branch" tfsdk:"branch"`
	BranchDeleted        bool                                   `json:"branch_deleted" tfsdk:"branch_deleted"`
	BranchDeletedAt      string                                 `json:"branch_deleted_at" tfsdk:"branch_deleted_at"`
	BranchDeletedBy      SkipRevertPeriodRes200_BranchDeletedBy `json:"branch_deleted_by" tfsdk:"branch_deleted_by"`
	ClosedAt             string                                 `json:"closed_at" tfsdk:"closed_at"`
	ClosedBy             SkipRevertPeriodRes200_ClosedBy        `json:"closed_by" tfsdk:"closed_by"`
	CreatedAt            string                                 `json:"created_at" tfsdk:"created_at"`
	DeployedAt           string                                 `json:"deployed_at" tfsdk:"deployed_at"`
	DeploymentState      string                                 `json:"deployment_state" tfsdk:"deployment_state"`
	HtmlBody             string                                 `json:"html_body" tfsdk:"html_body"`
	HtmlUrl              string                                 `json:"html_url" tfsdk:"html_url"`
	Id                   string                                 `json:"id" tfsdk:"id"`
	IntoBranch           string                                 `json:"into_branch" tfsdk:"into_branch"`
	IntoBranchShardCount float64                                `json:"into_branch_shard_count" tfsdk:"into_branch_shard_count"`
	IntoBranchSharded    bool                                   `json:"into_branch_sharded" tfsdk:"into_branch_sharded"`
	Notes                string                                 `json:"notes" tfsdk:"notes"`
	Number               float64                                `json:"number" tfsdk:"number"`
	State                string                                 `json:"state" tfsdk:"state"`
	UpdatedAt            string                                 `json:"updated_at" tfsdk:"updated_at"`
}

func (cl *Client) SkipRevertPeriod(ctx context.Context, organization string, database string, number string) (res200 *SkipRevertPeriodRes200, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests/" + number + "/skip-revert"})
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(SkipRevertPeriodRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type GetDatabaseRes404 struct{}
type GetDatabaseRes403 struct{}
type GetDatabaseRes500 struct{}
type GetDatabaseRes200_DataImport_DataSource struct {
	Database string `json:"database" tfsdk:"database"`
	Hostname string `json:"hostname" tfsdk:"hostname"`
	Port     string `json:"port" tfsdk:"port"`
}
type GetDatabaseRes200_DataImport struct {
	DataSource        GetDatabaseRes200_DataImport_DataSource `json:"data_source" tfsdk:"data_source"`
	FinishedAt        string                                  `json:"finished_at" tfsdk:"finished_at"`
	ImportCheckErrors string                                  `json:"import_check_errors" tfsdk:"import_check_errors"`
	StartedAt         string                                  `json:"started_at" tfsdk:"started_at"`
	State             string                                  `json:"state" tfsdk:"state"`
}
type GetDatabaseRes200_Region struct {
	DisplayName       string   `json:"display_name" tfsdk:"display_name"`
	Enabled           bool     `json:"enabled" tfsdk:"enabled"`
	Id                string   `json:"id" tfsdk:"id"`
	Location          string   `json:"location" tfsdk:"location"`
	Provider          string   `json:"provider" tfsdk:"provider"`
	PublicIpAddresses []string `json:"public_ip_addresses" tfsdk:"public_ip_addresses"`
	Slug              string   `json:"slug" tfsdk:"slug"`
}
type GetDatabaseRes200 struct {
	AllowDataBranching                bool                          `json:"allow_data_branching" tfsdk:"allow_data_branching"`
	AtBackupRestoreBranchesLimit      bool                          `json:"at_backup_restore_branches_limit" tfsdk:"at_backup_restore_branches_limit"`
	AtDevelopmentBranchLimit          bool                          `json:"at_development_branch_limit" tfsdk:"at_development_branch_limit"`
	AutomaticMigrations               bool                          `json:"automatic_migrations" tfsdk:"automatic_migrations"`
	BranchesCount                     float64                       `json:"branches_count" tfsdk:"branches_count"`
	BranchesUrl                       string                        `json:"branches_url" tfsdk:"branches_url"`
	CreatedAt                         string                        `json:"created_at" tfsdk:"created_at"`
	DataImport                        *GetDatabaseRes200_DataImport `json:"data_import,omitempty" tfsdk:"data_import"`
	DefaultBranch                     string                        `json:"default_branch" tfsdk:"default_branch"`
	DefaultBranchReadOnlyRegionsCount float64                       `json:"default_branch_read_only_regions_count" tfsdk:"default_branch_read_only_regions_count"`
	DefaultBranchShardCount           float64                       `json:"default_branch_shard_count" tfsdk:"default_branch_shard_count"`
	DefaultBranchTableCount           float64                       `json:"default_branch_table_count" tfsdk:"default_branch_table_count"`
	DevelopmentBranchesCount          float64                       `json:"development_branches_count" tfsdk:"development_branches_count"`
	HtmlUrl                           string                        `json:"html_url" tfsdk:"html_url"`
	Id                                string                        `json:"id" tfsdk:"id"`
	InsightsRawQueries                bool                          `json:"insights_raw_queries" tfsdk:"insights_raw_queries"`
	IssuesCount                       float64                       `json:"issues_count" tfsdk:"issues_count"`
	MigrationFramework                *string                       `json:"migration_framework,omitempty" tfsdk:"migration_framework"`
	MigrationTableName                *string                       `json:"migration_table_name,omitempty" tfsdk:"migration_table_name"`
	MultipleAdminsRequiredForDeletion bool                          `json:"multiple_admins_required_for_deletion" tfsdk:"multiple_admins_required_for_deletion"`
	Name                              string                        `json:"name" tfsdk:"name"`
	Notes                             *string                       `json:"notes,omitempty" tfsdk:"notes"`
	Plan                              string                        `json:"plan" tfsdk:"plan"`
	ProductionBranchWebConsole        bool                          `json:"production_branch_web_console" tfsdk:"production_branch_web_console"`
	ProductionBranchesCount           float64                       `json:"production_branches_count" tfsdk:"production_branches_count"`
	Ready                             bool                          `json:"ready" tfsdk:"ready"`
	Region                            GetDatabaseRes200_Region      `json:"region" tfsdk:"region"`
	RequireApprovalForDeploy          bool                          `json:"require_approval_for_deploy" tfsdk:"require_approval_for_deploy"`
	RestrictBranchRegion              bool                          `json:"restrict_branch_region" tfsdk:"restrict_branch_region"`
	SchemaLastUpdatedAt               *string                       `json:"schema_last_updated_at,omitempty" tfsdk:"schema_last_updated_at"`
	Sharded                           bool                          `json:"sharded" tfsdk:"sharded"`
	State                             string                        `json:"state" tfsdk:"state"`
	Type                              string                        `json:"type" tfsdk:"type"`
	UpdatedAt                         string                        `json:"updated_at" tfsdk:"updated_at"`
	Url                               string                        `json:"url" tfsdk:"url"`
}

func (cl *Client) GetDatabase(ctx context.Context, organization string, name string) (res200 *GetDatabaseRes200, res403 *GetDatabaseRes403, res404 *GetDatabaseRes404, res500 *GetDatabaseRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + name})
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(GetDatabaseRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(GetDatabaseRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(GetDatabaseRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(GetDatabaseRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type DeleteDatabaseRes404 struct{}
type DeleteDatabaseRes403 struct{}
type DeleteDatabaseRes500 struct{}
type DeleteDatabaseRes204 struct{}

func (cl *Client) DeleteDatabase(ctx context.Context, organization string, name string) (res204 *DeleteDatabaseRes204, res403 *DeleteDatabaseRes403, res404 *DeleteDatabaseRes404, res500 *DeleteDatabaseRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + name})
	r, err := http.NewRequestWithContext(ctx, "DELETE", u.String(), nil)
	if err != nil {
		return res204, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res204, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 204:
		res204 = new(DeleteDatabaseRes204)
		err = json.NewDecoder(res.Body).Decode(&res204)
	case 403:
		res403 = new(DeleteDatabaseRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(DeleteDatabaseRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(DeleteDatabaseRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res204, res403, res404, res500, err
}

type UpdateDatabaseSettingsReq struct {
	AllowDataBranching         *bool   `json:"allow_data_branching,omitempty" tfsdk:"allow_data_branching"`
	AutomaticMigrations        *bool   `json:"automatic_migrations,omitempty" tfsdk:"automatic_migrations"`
	DefaultBranch              *string `json:"default_branch,omitempty" tfsdk:"default_branch"`
	InsightsRawQueries         *bool   `json:"insights_raw_queries,omitempty" tfsdk:"insights_raw_queries"`
	MigrationFramework         *string `json:"migration_framework,omitempty" tfsdk:"migration_framework"`
	MigrationTableName         *string `json:"migration_table_name,omitempty" tfsdk:"migration_table_name"`
	Notes                      *string `json:"notes,omitempty" tfsdk:"notes"`
	ProductionBranchWebConsole *bool   `json:"production_branch_web_console,omitempty" tfsdk:"production_branch_web_console"`
	RequireApprovalForDeploy   *bool   `json:"require_approval_for_deploy,omitempty" tfsdk:"require_approval_for_deploy"`
	RestrictBranchRegion       *bool   `json:"restrict_branch_region,omitempty" tfsdk:"restrict_branch_region"`
}
type UpdateDatabaseSettingsRes500 struct{}
type UpdateDatabaseSettingsRes200_DataImport_DataSource struct {
	Database string `json:"database" tfsdk:"database"`
	Hostname string `json:"hostname" tfsdk:"hostname"`
	Port     string `json:"port" tfsdk:"port"`
}
type UpdateDatabaseSettingsRes200_DataImport struct {
	DataSource        UpdateDatabaseSettingsRes200_DataImport_DataSource `json:"data_source" tfsdk:"data_source"`
	FinishedAt        string                                             `json:"finished_at" tfsdk:"finished_at"`
	ImportCheckErrors string                                             `json:"import_check_errors" tfsdk:"import_check_errors"`
	StartedAt         string                                             `json:"started_at" tfsdk:"started_at"`
	State             string                                             `json:"state" tfsdk:"state"`
}
type UpdateDatabaseSettingsRes200_Region struct {
	DisplayName       string   `json:"display_name" tfsdk:"display_name"`
	Enabled           bool     `json:"enabled" tfsdk:"enabled"`
	Id                string   `json:"id" tfsdk:"id"`
	Location          string   `json:"location" tfsdk:"location"`
	Provider          string   `json:"provider" tfsdk:"provider"`
	PublicIpAddresses []string `json:"public_ip_addresses" tfsdk:"public_ip_addresses"`
	Slug              string   `json:"slug" tfsdk:"slug"`
}
type UpdateDatabaseSettingsRes200 struct {
	AllowDataBranching                bool                                     `json:"allow_data_branching" tfsdk:"allow_data_branching"`
	AtBackupRestoreBranchesLimit      bool                                     `json:"at_backup_restore_branches_limit" tfsdk:"at_backup_restore_branches_limit"`
	AtDevelopmentBranchLimit          bool                                     `json:"at_development_branch_limit" tfsdk:"at_development_branch_limit"`
	AutomaticMigrations               bool                                     `json:"automatic_migrations" tfsdk:"automatic_migrations"`
	BranchesCount                     float64                                  `json:"branches_count" tfsdk:"branches_count"`
	BranchesUrl                       string                                   `json:"branches_url" tfsdk:"branches_url"`
	CreatedAt                         string                                   `json:"created_at" tfsdk:"created_at"`
	DataImport                        *UpdateDatabaseSettingsRes200_DataImport `json:"data_import,omitempty" tfsdk:"data_import"`
	DefaultBranch                     string                                   `json:"default_branch" tfsdk:"default_branch"`
	DefaultBranchReadOnlyRegionsCount float64                                  `json:"default_branch_read_only_regions_count" tfsdk:"default_branch_read_only_regions_count"`
	DefaultBranchShardCount           float64                                  `json:"default_branch_shard_count" tfsdk:"default_branch_shard_count"`
	DefaultBranchTableCount           float64                                  `json:"default_branch_table_count" tfsdk:"default_branch_table_count"`
	DevelopmentBranchesCount          float64                                  `json:"development_branches_count" tfsdk:"development_branches_count"`
	HtmlUrl                           string                                   `json:"html_url" tfsdk:"html_url"`
	Id                                string                                   `json:"id" tfsdk:"id"`
	InsightsRawQueries                bool                                     `json:"insights_raw_queries" tfsdk:"insights_raw_queries"`
	IssuesCount                       float64                                  `json:"issues_count" tfsdk:"issues_count"`
	MigrationFramework                *string                                  `json:"migration_framework,omitempty" tfsdk:"migration_framework"`
	MigrationTableName                *string                                  `json:"migration_table_name,omitempty" tfsdk:"migration_table_name"`
	MultipleAdminsRequiredForDeletion bool                                     `json:"multiple_admins_required_for_deletion" tfsdk:"multiple_admins_required_for_deletion"`
	Name                              string                                   `json:"name" tfsdk:"name"`
	Notes                             *string                                  `json:"notes,omitempty" tfsdk:"notes"`
	Plan                              string                                   `json:"plan" tfsdk:"plan"`
	ProductionBranchWebConsole        bool                                     `json:"production_branch_web_console" tfsdk:"production_branch_web_console"`
	ProductionBranchesCount           float64                                  `json:"production_branches_count" tfsdk:"production_branches_count"`
	Ready                             bool                                     `json:"ready" tfsdk:"ready"`
	Region                            UpdateDatabaseSettingsRes200_Region      `json:"region" tfsdk:"region"`
	RequireApprovalForDeploy          bool                                     `json:"require_approval_for_deploy" tfsdk:"require_approval_for_deploy"`
	RestrictBranchRegion              bool                                     `json:"restrict_branch_region" tfsdk:"restrict_branch_region"`
	SchemaLastUpdatedAt               *string                                  `json:"schema_last_updated_at,omitempty" tfsdk:"schema_last_updated_at"`
	Sharded                           bool                                     `json:"sharded" tfsdk:"sharded"`
	State                             string                                   `json:"state" tfsdk:"state"`
	Type                              string                                   `json:"type" tfsdk:"type"`
	UpdatedAt                         string                                   `json:"updated_at" tfsdk:"updated_at"`
	Url                               string                                   `json:"url" tfsdk:"url"`
}
type UpdateDatabaseSettingsRes404 struct{}
type UpdateDatabaseSettingsRes403 struct{}

func (cl *Client) UpdateDatabaseSettings(ctx context.Context, organization string, name string, req UpdateDatabaseSettingsReq) (res200 *UpdateDatabaseSettingsRes200, res403 *UpdateDatabaseSettingsRes403, res404 *UpdateDatabaseSettingsRes404, res500 *UpdateDatabaseSettingsRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + name})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res200, res403, res404, res500, err
	}
	r, err := http.NewRequestWithContext(ctx, "PATCH", u.String(), body)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(UpdateDatabaseSettingsRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(UpdateDatabaseSettingsRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(UpdateDatabaseSettingsRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(UpdateDatabaseSettingsRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type ListReadOnlyRegionsRes500 struct{}
type ListReadOnlyRegionsRes200_Actor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type ListReadOnlyRegionsRes200_Region struct {
	DisplayName       string   `json:"display_name" tfsdk:"display_name"`
	Enabled           bool     `json:"enabled" tfsdk:"enabled"`
	Id                string   `json:"id" tfsdk:"id"`
	Location          string   `json:"location" tfsdk:"location"`
	Provider          string   `json:"provider" tfsdk:"provider"`
	PublicIpAddresses []string `json:"public_ip_addresses" tfsdk:"public_ip_addresses"`
	Slug              string   `json:"slug" tfsdk:"slug"`
}
type ListReadOnlyRegionsRes200 struct {
	Actor       ListReadOnlyRegionsRes200_Actor  `json:"actor" tfsdk:"actor"`
	CreatedAt   string                           `json:"created_at" tfsdk:"created_at"`
	DisplayName string                           `json:"display_name" tfsdk:"display_name"`
	Id          string                           `json:"id" tfsdk:"id"`
	Ready       bool                             `json:"ready" tfsdk:"ready"`
	ReadyAt     string                           `json:"ready_at" tfsdk:"ready_at"`
	Region      ListReadOnlyRegionsRes200_Region `json:"region" tfsdk:"region"`
	UpdatedAt   string                           `json:"updated_at" tfsdk:"updated_at"`
}
type ListReadOnlyRegionsRes404 struct{}
type ListReadOnlyRegionsRes403 struct{}

func (cl *Client) ListReadOnlyRegions(ctx context.Context, organization string, name string, page *int, perPage *int) (res200 *ListReadOnlyRegionsRes200, res403 *ListReadOnlyRegionsRes403, res404 *ListReadOnlyRegionsRes404, res500 *ListReadOnlyRegionsRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + name + "/read-only-regions"})
	q := u.Query()
	if page != nil {
		q.Set("page", strconv.Itoa(*page))
	}
	if perPage != nil {
		q.Set("per_page", strconv.Itoa(*perPage))
	}
	u.RawQuery = q.Encode()
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListReadOnlyRegionsRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(ListReadOnlyRegionsRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(ListReadOnlyRegionsRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(ListReadOnlyRegionsRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type ListDatabaseRegionsRes404 struct{}
type ListDatabaseRegionsRes403 struct{}
type ListDatabaseRegionsRes500 struct{}
type ListDatabaseRegionsRes200_DataItem struct {
	DisplayName       string   `json:"display_name" tfsdk:"display_name"`
	Enabled           bool     `json:"enabled" tfsdk:"enabled"`
	Id                string   `json:"id" tfsdk:"id"`
	Location          string   `json:"location" tfsdk:"location"`
	Provider          string   `json:"provider" tfsdk:"provider"`
	PublicIpAddresses []string `json:"public_ip_addresses" tfsdk:"public_ip_addresses"`
	Slug              string   `json:"slug" tfsdk:"slug"`
}
type ListDatabaseRegionsRes200 struct {
	Data []ListDatabaseRegionsRes200_DataItem `json:"data" tfsdk:"data"`
}

func (cl *Client) ListDatabaseRegions(ctx context.Context, organization string, name string, page *int, perPage *int) (res200 *ListDatabaseRegionsRes200, res403 *ListDatabaseRegionsRes403, res404 *ListDatabaseRegionsRes404, res500 *ListDatabaseRegionsRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + name + "/regions"})
	q := u.Query()
	if page != nil {
		q.Set("page", strconv.Itoa(*page))
	}
	if perPage != nil {
		q.Set("per_page", strconv.Itoa(*perPage))
	}
	u.RawQuery = q.Encode()
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListDatabaseRegionsRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(ListDatabaseRegionsRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(ListDatabaseRegionsRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(ListDatabaseRegionsRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type ListOauthApplicationsRes404 struct{}
type ListOauthApplicationsRes403 struct{}
type ListOauthApplicationsRes500 struct{}
type ListOauthApplicationsRes200_DataItem struct {
	Avatar      *string  `json:"avatar,omitempty" tfsdk:"avatar"`
	ClientId    string   `json:"client_id" tfsdk:"client_id"`
	CreatedAt   string   `json:"created_at" tfsdk:"created_at"`
	Domain      string   `json:"domain" tfsdk:"domain"`
	Id          string   `json:"id" tfsdk:"id"`
	Name        string   `json:"name" tfsdk:"name"`
	RedirectUri string   `json:"redirect_uri" tfsdk:"redirect_uri"`
	Scopes      []string `json:"scopes" tfsdk:"scopes"`
	Tokens      float64  `json:"tokens" tfsdk:"tokens"`
	UpdatedAt   string   `json:"updated_at" tfsdk:"updated_at"`
}
type ListOauthApplicationsRes200 struct {
	Data []ListOauthApplicationsRes200_DataItem `json:"data" tfsdk:"data"`
}

func (cl *Client) ListOauthApplications(ctx context.Context, organization string, page *int, perPage *int) (res200 *ListOauthApplicationsRes200, res403 *ListOauthApplicationsRes403, res404 *ListOauthApplicationsRes404, res500 *ListOauthApplicationsRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/oauth-applications"})
	q := u.Query()
	if page != nil {
		q.Set("page", strconv.Itoa(*page))
	}
	if perPage != nil {
		q.Set("per_page", strconv.Itoa(*perPage))
	}
	u.RawQuery = q.Encode()
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListOauthApplicationsRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(ListOauthApplicationsRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(ListOauthApplicationsRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(ListOauthApplicationsRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type GetOauthApplicationRes404 struct{}
type GetOauthApplicationRes403 struct{}
type GetOauthApplicationRes500 struct{}
type GetOauthApplicationRes200 struct {
	Avatar      *string  `json:"avatar,omitempty" tfsdk:"avatar"`
	ClientId    string   `json:"client_id" tfsdk:"client_id"`
	CreatedAt   string   `json:"created_at" tfsdk:"created_at"`
	Domain      string   `json:"domain" tfsdk:"domain"`
	Id          string   `json:"id" tfsdk:"id"`
	Name        string   `json:"name" tfsdk:"name"`
	RedirectUri string   `json:"redirect_uri" tfsdk:"redirect_uri"`
	Scopes      []string `json:"scopes" tfsdk:"scopes"`
	Tokens      float64  `json:"tokens" tfsdk:"tokens"`
	UpdatedAt   string   `json:"updated_at" tfsdk:"updated_at"`
}

func (cl *Client) GetOauthApplication(ctx context.Context, organization string, applicationId string) (res200 *GetOauthApplicationRes200, res403 *GetOauthApplicationRes403, res404 *GetOauthApplicationRes404, res500 *GetOauthApplicationRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/oauth-applications/" + applicationId})
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(GetOauthApplicationRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(GetOauthApplicationRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(GetOauthApplicationRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(GetOauthApplicationRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type ListOauthTokensRes200_DataItem struct {
	ActorDisplayName string `json:"actor_display_name" tfsdk:"actor_display_name"`
	ActorId          string `json:"actor_id" tfsdk:"actor_id"`
	ActorType        string `json:"actor_type" tfsdk:"actor_type"`
	AvatarUrl        string `json:"avatar_url" tfsdk:"avatar_url"`
	CreatedAt        string `json:"created_at" tfsdk:"created_at"`
	DisplayName      string `json:"display_name" tfsdk:"display_name"`
	ExpiresAt        string `json:"expires_at" tfsdk:"expires_at"`
	Id               string `json:"id" tfsdk:"id"`
	LastUsedAt       string `json:"last_used_at" tfsdk:"last_used_at"`
	Name             string `json:"name" tfsdk:"name"`
	UpdatedAt        string `json:"updated_at" tfsdk:"updated_at"`
}
type ListOauthTokensRes200 struct {
	Data []ListOauthTokensRes200_DataItem `json:"data" tfsdk:"data"`
}
type ListOauthTokensRes404 struct{}
type ListOauthTokensRes403 struct{}
type ListOauthTokensRes500 struct{}

func (cl *Client) ListOauthTokens(ctx context.Context, organization string, applicationId string, page *int, perPage *int) (res200 *ListOauthTokensRes200, res403 *ListOauthTokensRes403, res404 *ListOauthTokensRes404, res500 *ListOauthTokensRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/oauth-applications/" + applicationId + "/tokens"})
	q := u.Query()
	if page != nil {
		q.Set("page", strconv.Itoa(*page))
	}
	if perPage != nil {
		q.Set("per_page", strconv.Itoa(*perPage))
	}
	u.RawQuery = q.Encode()
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListOauthTokensRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(ListOauthTokensRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(ListOauthTokensRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(ListOauthTokensRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type GetOauthTokenRes404 struct{}
type GetOauthTokenRes403 struct{}
type GetOauthTokenRes500 struct{}
type GetOauthTokenRes200_OauthAccessesByResource_Branch struct {
	Accesses []string `json:"accesses" tfsdk:"accesses"`
	Branches []string `json:"branches" tfsdk:"branches"`
}
type GetOauthTokenRes200_OauthAccessesByResource_Database struct {
	Accesses  []string `json:"accesses" tfsdk:"accesses"`
	Databases []string `json:"databases" tfsdk:"databases"`
}
type GetOauthTokenRes200_OauthAccessesByResource_Organization struct {
	Accesses      []string `json:"accesses" tfsdk:"accesses"`
	Organizations []string `json:"organizations" tfsdk:"organizations"`
}
type GetOauthTokenRes200_OauthAccessesByResource_User struct {
	Accesses []string `json:"accesses" tfsdk:"accesses"`
	Users    []string `json:"users" tfsdk:"users"`
}
type GetOauthTokenRes200_OauthAccessesByResource struct {
	Branch       GetOauthTokenRes200_OauthAccessesByResource_Branch       `json:"branch" tfsdk:"branch"`
	Database     GetOauthTokenRes200_OauthAccessesByResource_Database     `json:"database" tfsdk:"database"`
	Organization GetOauthTokenRes200_OauthAccessesByResource_Organization `json:"organization" tfsdk:"organization"`
	User         GetOauthTokenRes200_OauthAccessesByResource_User         `json:"user" tfsdk:"user"`
}
type GetOauthTokenRes200 struct {
	ActorDisplayName        string                                      `json:"actor_display_name" tfsdk:"actor_display_name"`
	ActorId                 string                                      `json:"actor_id" tfsdk:"actor_id"`
	ActorType               string                                      `json:"actor_type" tfsdk:"actor_type"`
	AvatarUrl               string                                      `json:"avatar_url" tfsdk:"avatar_url"`
	CreatedAt               string                                      `json:"created_at" tfsdk:"created_at"`
	DisplayName             string                                      `json:"display_name" tfsdk:"display_name"`
	ExpiresAt               string                                      `json:"expires_at" tfsdk:"expires_at"`
	Id                      string                                      `json:"id" tfsdk:"id"`
	LastUsedAt              string                                      `json:"last_used_at" tfsdk:"last_used_at"`
	Name                    string                                      `json:"name" tfsdk:"name"`
	OauthAccessesByResource GetOauthTokenRes200_OauthAccessesByResource `json:"oauth_accesses_by_resource" tfsdk:"oauth_accesses_by_resource"`
	UpdatedAt               string                                      `json:"updated_at" tfsdk:"updated_at"`
}

func (cl *Client) GetOauthToken(ctx context.Context, organization string, applicationId string, tokenId string) (res200 *GetOauthTokenRes200, res403 *GetOauthTokenRes403, res404 *GetOauthTokenRes404, res500 *GetOauthTokenRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/oauth-applications/" + applicationId + "/tokens/" + tokenId})
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(GetOauthTokenRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(GetOauthTokenRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(GetOauthTokenRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(GetOauthTokenRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}

type DeleteOauthTokenRes500 struct{}
type DeleteOauthTokenRes204 struct{}
type DeleteOauthTokenRes404 struct{}
type DeleteOauthTokenRes403 struct{}

func (cl *Client) DeleteOauthToken(ctx context.Context, organization string, applicationId string, tokenId string) (res204 *DeleteOauthTokenRes204, res403 *DeleteOauthTokenRes403, res404 *DeleteOauthTokenRes404, res500 *DeleteOauthTokenRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/oauth-applications/" + applicationId + "/tokens/" + tokenId})
	r, err := http.NewRequestWithContext(ctx, "DELETE", u.String(), nil)
	if err != nil {
		return res204, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res204, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 204:
		res204 = new(DeleteOauthTokenRes204)
		err = json.NewDecoder(res.Body).Decode(&res204)
	case 403:
		res403 = new(DeleteOauthTokenRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(DeleteOauthTokenRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(DeleteOauthTokenRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res204, res403, res404, res500, err
}

type CreateOrRenewOauthTokenReq struct {
	ClientId     string  `json:"client_id" tfsdk:"client_id"`
	ClientSecret string  `json:"client_secret" tfsdk:"client_secret"`
	Code         *string `json:"code,omitempty" tfsdk:"code"`
	GrantType    string  `json:"grant_type" tfsdk:"grant_type"`
	RedirectUri  *string `json:"redirect_uri,omitempty" tfsdk:"redirect_uri"`
	RefreshToken *string `json:"refresh_token,omitempty" tfsdk:"refresh_token"`
}
type CreateOrRenewOauthTokenRes422 struct{}
type CreateOrRenewOauthTokenRes500 struct{}
type CreateOrRenewOauthTokenRes200 struct {
	ActorDisplayName      *string   `json:"actor_display_name,omitempty" tfsdk:"actor_display_name"`
	ActorId               *string   `json:"actor_id,omitempty" tfsdk:"actor_id"`
	DisplayName           *string   `json:"display_name,omitempty" tfsdk:"display_name"`
	Name                  *string   `json:"name,omitempty" tfsdk:"name"`
	PlainTextRefreshToken *string   `json:"plain_text_refresh_token,omitempty" tfsdk:"plain_text_refresh_token"`
	ServiceTokenAccesses  *[]string `json:"service_token_accesses,omitempty" tfsdk:"service_token_accesses"`
	Token                 *string   `json:"token,omitempty" tfsdk:"token"`
}
type CreateOrRenewOauthTokenRes404 struct{}
type CreateOrRenewOauthTokenRes403 struct{}

func (cl *Client) CreateOrRenewOauthToken(ctx context.Context, organization string, id string, req CreateOrRenewOauthTokenReq) (res200 *CreateOrRenewOauthTokenRes200, res403 *CreateOrRenewOauthTokenRes403, res404 *CreateOrRenewOauthTokenRes404, res422 *CreateOrRenewOauthTokenRes422, res500 *CreateOrRenewOauthTokenRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/oauth-applications/" + id + "/token"})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res200, res403, res404, res422, res500, err
	}
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return res200, res403, res404, res422, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res422, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(CreateOrRenewOauthTokenRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(CreateOrRenewOauthTokenRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(CreateOrRenewOauthTokenRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 422:
		res422 = new(CreateOrRenewOauthTokenRes422)
		err = json.NewDecoder(res.Body).Decode(&res422)
	case 500:
		res500 = new(CreateOrRenewOauthTokenRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res422, res500, err
}

type GetCurrentUserRes500 struct{}
type GetCurrentUserRes200 struct {
	AvatarUrl               *string `json:"avatar_url,omitempty" tfsdk:"avatar_url"`
	CreatedAt               *string `json:"created_at,omitempty" tfsdk:"created_at"`
	DefaultOrganizationId   *string `json:"default_organization_id,omitempty" tfsdk:"default_organization_id"`
	DirectoryManaged        *bool   `json:"directory_managed,omitempty" tfsdk:"directory_managed"`
	DisplayName             *string `json:"display_name,omitempty" tfsdk:"display_name"`
	Email                   *string `json:"email,omitempty" tfsdk:"email"`
	EmailVerified           *bool   `json:"email_verified,omitempty" tfsdk:"email_verified"`
	Id                      *string `json:"id,omitempty" tfsdk:"id"`
	Managed                 *bool   `json:"managed,omitempty" tfsdk:"managed"`
	Name                    *string `json:"name,omitempty" tfsdk:"name"`
	Sso                     *bool   `json:"sso,omitempty" tfsdk:"sso"`
	TwoFactorAuthConfigured *bool   `json:"two_factor_auth_configured,omitempty" tfsdk:"two_factor_auth_configured"`
	UpdatedAt               *string `json:"updated_at,omitempty" tfsdk:"updated_at"`
}
type GetCurrentUserRes404 struct{}
type GetCurrentUserRes403 struct{}

func (cl *Client) GetCurrentUser(ctx context.Context) (res200 *GetCurrentUserRes200, res403 *GetCurrentUserRes403, res404 *GetCurrentUserRes404, res500 *GetCurrentUserRes500, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "user"})
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, res403, res404, res500, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(GetCurrentUserRes200)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 = new(GetCurrentUserRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
	case 404:
		res404 = new(GetCurrentUserRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
	case 500:
		res500 = new(GetCurrentUserRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
	default:
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, res403, res404, res500, err
}
