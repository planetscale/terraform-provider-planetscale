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

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (err *ErrorResponse) Error() string {
	return fmt.Sprintf("error %s: %s", err.Code, err.Message)
}

type Region struct {
	DisplayName       string   `json:"display_name" tfsdk:"display_name"`
	Enabled           bool     `json:"enabled" tfsdk:"enabled"`
	Id                string   `json:"id" tfsdk:"id"`
	Location          string   `json:"location" tfsdk:"location"`
	Provider          string   `json:"provider" tfsdk:"provider"`
	PublicIpAddresses []string `json:"public_ip_addresses" tfsdk:"public_ip_addresses"`
	Slug              string   `json:"slug" tfsdk:"slug"`
}
type Actor struct {
	AvatarUrl   string `json:"avatar_url" tfsdk:"avatar_url"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
}
type Features struct {
	Insights      *bool `json:"insights,omitempty" tfsdk:"insights"`
	SingleTenancy *bool `json:"single_tenancy,omitempty" tfsdk:"single_tenancy"`
	Sso           *bool `json:"sso,omitempty" tfsdk:"sso"`
}
type OauthApplication struct {
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
type OauthDatabaseAccesses struct {
	Accesses  []string `json:"accesses" tfsdk:"accesses"`
	Databases []string `json:"databases" tfsdk:"databases"`
}
type OauthUserAccesses struct {
	Accesses []string `json:"accesses" tfsdk:"accesses"`
	Users    []string `json:"users" tfsdk:"users"`
}
type SchemaSnapshot struct {
	CreatedAt string `json:"created_at" tfsdk:"created_at"`
	Id        string `json:"id" tfsdk:"id"`
	Name      string `json:"name" tfsdk:"name"`
	UpdatedAt string `json:"updated_at" tfsdk:"updated_at"`
	Url       string `json:"url" tfsdk:"url"`
}
type Backup struct {
	Actor                Actor          `json:"actor" tfsdk:"actor"`
	BackupPolicy         BackupPolicy   `json:"backup_policy" tfsdk:"backup_policy"`
	CreatedAt            string         `json:"created_at" tfsdk:"created_at"`
	EstimatedStorageCost float64        `json:"estimated_storage_cost" tfsdk:"estimated_storage_cost"`
	Id                   string         `json:"id" tfsdk:"id"`
	Name                 string         `json:"name" tfsdk:"name"`
	Required             bool           `json:"required" tfsdk:"required"`
	RestoredBranches     *[]string      `json:"restored_branches,omitempty" tfsdk:"restored_branches"`
	SchemaSnapshot       SchemaSnapshot `json:"schema_snapshot" tfsdk:"schema_snapshot"`
	Size                 float64        `json:"size" tfsdk:"size"`
	State                string         `json:"state" tfsdk:"state"`
	UpdatedAt            string         `json:"updated_at" tfsdk:"updated_at"`
}
type DataImport struct {
	DataSource        DataSource `json:"data_source" tfsdk:"data_source"`
	FinishedAt        string     `json:"finished_at" tfsdk:"finished_at"`
	ImportCheckErrors string     `json:"import_check_errors" tfsdk:"import_check_errors"`
	StartedAt         string     `json:"started_at" tfsdk:"started_at"`
	State             string     `json:"state" tfsdk:"state"`
}
type DeployOperation struct {
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
type OauthAccessesByResource struct {
	Branch       OauthBranchAccesses       `json:"branch" tfsdk:"branch"`
	Database     OauthDatabaseAccesses     `json:"database" tfsdk:"database"`
	Organization OauthOrganizationAccesses `json:"organization" tfsdk:"organization"`
	User         OauthUserAccesses         `json:"user" tfsdk:"user"`
}
type BranchForPassword struct {
	AccessHostUrl    string `json:"access_host_url" tfsdk:"access_host_url"`
	Id               string `json:"id" tfsdk:"id"`
	MysqlEdgeAddress string `json:"mysql_edge_address" tfsdk:"mysql_edge_address"`
	Name             string `json:"name" tfsdk:"name"`
	Production       bool   `json:"production" tfsdk:"production"`
}
type DataSource struct {
	Database string `json:"database" tfsdk:"database"`
	Hostname string `json:"hostname" tfsdk:"hostname"`
	Port     string `json:"port" tfsdk:"port"`
}
type Flags struct {
	ExampleFlag *string `json:"example_flag,omitempty" tfsdk:"example_flag"`
}
type OauthTokenWithDetails struct {
	ActorDisplayName        string                  `json:"actor_display_name" tfsdk:"actor_display_name"`
	ActorId                 string                  `json:"actor_id" tfsdk:"actor_id"`
	ActorType               string                  `json:"actor_type" tfsdk:"actor_type"`
	AvatarUrl               string                  `json:"avatar_url" tfsdk:"avatar_url"`
	CreatedAt               string                  `json:"created_at" tfsdk:"created_at"`
	DisplayName             string                  `json:"display_name" tfsdk:"display_name"`
	ExpiresAt               string                  `json:"expires_at" tfsdk:"expires_at"`
	Id                      string                  `json:"id" tfsdk:"id"`
	LastUsedAt              string                  `json:"last_used_at" tfsdk:"last_used_at"`
	Name                    string                  `json:"name" tfsdk:"name"`
	OauthAccessesByResource OauthAccessesByResource `json:"oauth_accesses_by_resource" tfsdk:"oauth_accesses_by_resource"`
	UpdatedAt               string                  `json:"updated_at" tfsdk:"updated_at"`
}
type ReadOnlyRegion struct {
	Actor       Actor  `json:"actor" tfsdk:"actor"`
	CreatedAt   string `json:"created_at" tfsdk:"created_at"`
	DisplayName string `json:"display_name" tfsdk:"display_name"`
	Id          string `json:"id" tfsdk:"id"`
	Ready       bool   `json:"ready" tfsdk:"ready"`
	ReadyAt     string `json:"ready_at" tfsdk:"ready_at"`
	Region      Region `json:"region" tfsdk:"region"`
	UpdatedAt   string `json:"updated_at" tfsdk:"updated_at"`
}
type BackupPolicy struct {
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
type OauthToken struct {
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
type Organization struct {
	AdminOnlyProductionAccess bool      `json:"admin_only_production_access" tfsdk:"admin_only_production_access"`
	BillingEmail              *string   `json:"billing_email,omitempty" tfsdk:"billing_email"`
	CanCreateDatabases        bool      `json:"can_create_databases" tfsdk:"can_create_databases"`
	CreatedAt                 string    `json:"created_at" tfsdk:"created_at"`
	DatabaseCount             float64   `json:"database_count" tfsdk:"database_count"`
	Features                  *Features `json:"features,omitempty" tfsdk:"features"`
	Flags                     *Flags    `json:"flags,omitempty" tfsdk:"flags"`
	FreeDatabasesRemaining    float64   `json:"free_databases_remaining" tfsdk:"free_databases_remaining"`
	HasPastDueInvoices        bool      `json:"has_past_due_invoices" tfsdk:"has_past_due_invoices"`
	Id                        string    `json:"id" tfsdk:"id"`
	IdpManagedRoles           bool      `json:"idp_managed_roles" tfsdk:"idp_managed_roles"`
	Name                      string    `json:"name" tfsdk:"name"`
	Plan                      string    `json:"plan" tfsdk:"plan"`
	SingleTenancy             bool      `json:"single_tenancy" tfsdk:"single_tenancy"`
	SleepingDatabaseCount     float64   `json:"sleeping_database_count" tfsdk:"sleeping_database_count"`
	Sso                       bool      `json:"sso" tfsdk:"sso"`
	SsoDirectory              bool      `json:"sso_directory" tfsdk:"sso_directory"`
	SsoPortalUrl              *string   `json:"sso_portal_url,omitempty" tfsdk:"sso_portal_url"`
	UpdatedAt                 string    `json:"updated_at" tfsdk:"updated_at"`
	ValidBillingInfo          bool      `json:"valid_billing_info" tfsdk:"valid_billing_info"`
}
type QueuedDeployRequest struct {
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
type DeployRequest struct {
	Actor                Actor   `json:"actor" tfsdk:"actor"`
	Approved             bool    `json:"approved" tfsdk:"approved"`
	Branch               string  `json:"branch" tfsdk:"branch"`
	BranchDeleted        bool    `json:"branch_deleted" tfsdk:"branch_deleted"`
	BranchDeletedAt      string  `json:"branch_deleted_at" tfsdk:"branch_deleted_at"`
	BranchDeletedBy      Actor   `json:"branch_deleted_by" tfsdk:"branch_deleted_by"`
	ClosedAt             string  `json:"closed_at" tfsdk:"closed_at"`
	ClosedBy             Actor   `json:"closed_by" tfsdk:"closed_by"`
	CreatedAt            string  `json:"created_at" tfsdk:"created_at"`
	DeployedAt           string  `json:"deployed_at" tfsdk:"deployed_at"`
	DeploymentState      string  `json:"deployment_state" tfsdk:"deployment_state"`
	HtmlBody             string  `json:"html_body" tfsdk:"html_body"`
	HtmlUrl              string  `json:"html_url" tfsdk:"html_url"`
	Id                   string  `json:"id" tfsdk:"id"`
	IntoBranch           string  `json:"into_branch" tfsdk:"into_branch"`
	IntoBranchShardCount float64 `json:"into_branch_shard_count" tfsdk:"into_branch_shard_count"`
	IntoBranchSharded    bool    `json:"into_branch_sharded" tfsdk:"into_branch_sharded"`
	Notes                string  `json:"notes" tfsdk:"notes"`
	Number               float64 `json:"number" tfsdk:"number"`
	State                string  `json:"state" tfsdk:"state"`
	UpdatedAt            string  `json:"updated_at" tfsdk:"updated_at"`
}
type DeployReview struct {
	Actor     Actor  `json:"actor" tfsdk:"actor"`
	Body      string `json:"body" tfsdk:"body"`
	CreatedAt string `json:"created_at" tfsdk:"created_at"`
	HtmlBody  string `json:"html_body" tfsdk:"html_body"`
	Id        string `json:"id" tfsdk:"id"`
	State     string `json:"state" tfsdk:"state"`
	UpdatedAt string `json:"updated_at" tfsdk:"updated_at"`
}
type OauthOrganizationAccesses struct {
	Accesses      []string `json:"accesses" tfsdk:"accesses"`
	Organizations []string `json:"organizations" tfsdk:"organizations"`
}
type TableSchema struct {
	Html string `json:"html" tfsdk:"html"`
	Name string `json:"name" tfsdk:"name"`
	Raw  string `json:"raw" tfsdk:"raw"`
}
type Database struct {
	AllowDataBranching                bool        `json:"allow_data_branching" tfsdk:"allow_data_branching"`
	AtBackupRestoreBranchesLimit      bool        `json:"at_backup_restore_branches_limit" tfsdk:"at_backup_restore_branches_limit"`
	AtDevelopmentBranchLimit          bool        `json:"at_development_branch_limit" tfsdk:"at_development_branch_limit"`
	AutomaticMigrations               *bool       `json:"automatic_migrations,omitempty" tfsdk:"automatic_migrations"`
	BranchesCount                     float64     `json:"branches_count" tfsdk:"branches_count"`
	BranchesUrl                       string      `json:"branches_url" tfsdk:"branches_url"`
	CreatedAt                         string      `json:"created_at" tfsdk:"created_at"`
	DataImport                        *DataImport `json:"data_import,omitempty" tfsdk:"data_import"`
	DefaultBranch                     string      `json:"default_branch" tfsdk:"default_branch"`
	DefaultBranchReadOnlyRegionsCount float64     `json:"default_branch_read_only_regions_count" tfsdk:"default_branch_read_only_regions_count"`
	DefaultBranchShardCount           float64     `json:"default_branch_shard_count" tfsdk:"default_branch_shard_count"`
	DefaultBranchTableCount           float64     `json:"default_branch_table_count" tfsdk:"default_branch_table_count"`
	DevelopmentBranchesCount          float64     `json:"development_branches_count" tfsdk:"development_branches_count"`
	HtmlUrl                           string      `json:"html_url" tfsdk:"html_url"`
	Id                                string      `json:"id" tfsdk:"id"`
	InsightsRawQueries                bool        `json:"insights_raw_queries" tfsdk:"insights_raw_queries"`
	IssuesCount                       float64     `json:"issues_count" tfsdk:"issues_count"`
	MigrationFramework                *string     `json:"migration_framework,omitempty" tfsdk:"migration_framework"`
	MigrationTableName                *string     `json:"migration_table_name,omitempty" tfsdk:"migration_table_name"`
	MultipleAdminsRequiredForDeletion bool        `json:"multiple_admins_required_for_deletion" tfsdk:"multiple_admins_required_for_deletion"`
	Name                              string      `json:"name" tfsdk:"name"`
	Plan                              string      `json:"plan" tfsdk:"plan"`
	ProductionBranchWebConsole        bool        `json:"production_branch_web_console" tfsdk:"production_branch_web_console"`
	ProductionBranchesCount           float64     `json:"production_branches_count" tfsdk:"production_branches_count"`
	Ready                             bool        `json:"ready" tfsdk:"ready"`
	Region                            Region      `json:"region" tfsdk:"region"`
	RequireApprovalForDeploy          bool        `json:"require_approval_for_deploy" tfsdk:"require_approval_for_deploy"`
	RestrictBranchRegion              bool        `json:"restrict_branch_region" tfsdk:"restrict_branch_region"`
	SchemaLastUpdatedAt               *string     `json:"schema_last_updated_at,omitempty" tfsdk:"schema_last_updated_at"`
	Sharded                           bool        `json:"sharded" tfsdk:"sharded"`
	State                             string      `json:"state" tfsdk:"state"`
	Type                              string      `json:"type" tfsdk:"type"`
	UpdatedAt                         string      `json:"updated_at" tfsdk:"updated_at"`
	Url                               string      `json:"url" tfsdk:"url"`
}
type DeployRequestWithDeployment struct {
	Actor                Actor      `json:"actor" tfsdk:"actor"`
	Approved             bool       `json:"approved" tfsdk:"approved"`
	Branch               string     `json:"branch" tfsdk:"branch"`
	BranchDeleted        bool       `json:"branch_deleted" tfsdk:"branch_deleted"`
	BranchDeletedAt      string     `json:"branch_deleted_at" tfsdk:"branch_deleted_at"`
	BranchDeletedBy      Actor      `json:"branch_deleted_by" tfsdk:"branch_deleted_by"`
	ClosedAt             string     `json:"closed_at" tfsdk:"closed_at"`
	ClosedBy             Actor      `json:"closed_by" tfsdk:"closed_by"`
	CreatedAt            string     `json:"created_at" tfsdk:"created_at"`
	DeployedAt           string     `json:"deployed_at" tfsdk:"deployed_at"`
	Deployment           Deployment `json:"deployment" tfsdk:"deployment"`
	DeploymentState      string     `json:"deployment_state" tfsdk:"deployment_state"`
	HtmlBody             string     `json:"html_body" tfsdk:"html_body"`
	HtmlUrl              string     `json:"html_url" tfsdk:"html_url"`
	Id                   string     `json:"id" tfsdk:"id"`
	IntoBranch           string     `json:"into_branch" tfsdk:"into_branch"`
	IntoBranchShardCount float64    `json:"into_branch_shard_count" tfsdk:"into_branch_shard_count"`
	IntoBranchSharded    bool       `json:"into_branch_sharded" tfsdk:"into_branch_sharded"`
	Notes                string     `json:"notes" tfsdk:"notes"`
	Number               float64    `json:"number" tfsdk:"number"`
	State                string     `json:"state" tfsdk:"state"`
	UpdatedAt            string     `json:"updated_at" tfsdk:"updated_at"`
}
type Deployment struct {
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
type Branch struct {
	AccessHostUrl               *string             `json:"access_host_url,omitempty" tfsdk:"access_host_url"`
	Actor                       *Actor              `json:"actor,omitempty" tfsdk:"actor"`
	ClusterRateName             string              `json:"cluster_rate_name" tfsdk:"cluster_rate_name"`
	CreatedAt                   string              `json:"created_at" tfsdk:"created_at"`
	HtmlUrl                     string              `json:"html_url" tfsdk:"html_url"`
	Id                          string              `json:"id" tfsdk:"id"`
	InitialRestoreId            *string             `json:"initial_restore_id,omitempty" tfsdk:"initial_restore_id"`
	MysqlAddress                string              `json:"mysql_address" tfsdk:"mysql_address"`
	MysqlEdgeAddress            string              `json:"mysql_edge_address" tfsdk:"mysql_edge_address"`
	Name                        string              `json:"name" tfsdk:"name"`
	ParentBranch                *string             `json:"parent_branch,omitempty" tfsdk:"parent_branch"`
	Production                  bool                `json:"production" tfsdk:"production"`
	Ready                       bool                `json:"ready" tfsdk:"ready"`
	Region                      *Region             `json:"region,omitempty" tfsdk:"region"`
	RestoreChecklistCompletedAt *string             `json:"restore_checklist_completed_at,omitempty" tfsdk:"restore_checklist_completed_at"`
	RestoredFromBranch          *RestoredFromBranch `json:"restored_from_branch,omitempty" tfsdk:"restored_from_branch"`
	SchemaLastUpdatedAt         string              `json:"schema_last_updated_at" tfsdk:"schema_last_updated_at"`
	ShardCount                  *float64            `json:"shard_count,omitempty" tfsdk:"shard_count"`
	Sharded                     bool                `json:"sharded" tfsdk:"sharded"`
	UpdatedAt                   string              `json:"updated_at" tfsdk:"updated_at"`
}
type OauthBranchAccesses struct {
	Accesses []string `json:"accesses" tfsdk:"accesses"`
	Branches []string `json:"branches" tfsdk:"branches"`
}
type Password struct {
	AccessHostUrl  string            `json:"access_host_url" tfsdk:"access_host_url"`
	Actor          *Actor            `json:"actor,omitempty" tfsdk:"actor"`
	CreatedAt      string            `json:"created_at" tfsdk:"created_at"`
	DatabaseBranch BranchForPassword `json:"database_branch" tfsdk:"database_branch"`
	DeletedAt      *string           `json:"deleted_at,omitempty" tfsdk:"deleted_at"`
	ExpiresAt      *string           `json:"expires_at,omitempty" tfsdk:"expires_at"`
	Id             string            `json:"id" tfsdk:"id"`
	Integrations   []string          `json:"integrations" tfsdk:"integrations"`
	Name           string            `json:"name" tfsdk:"name"`
	Region         *Region           `json:"region,omitempty" tfsdk:"region"`
	Renewable      bool              `json:"renewable" tfsdk:"renewable"`
	Role           string            `json:"role" tfsdk:"role"`
	TtlSeconds     float64           `json:"ttl_seconds" tfsdk:"ttl_seconds"`
	Username       *string           `json:"username,omitempty" tfsdk:"username"`
}
type RestoredFromBranch struct {
	CreatedAt string `json:"created_at" tfsdk:"created_at"`
	DeletedAt string `json:"deleted_at" tfsdk:"deleted_at"`
	Id        string `json:"id" tfsdk:"id"`
	Name      string `json:"name" tfsdk:"name"`
	UpdatedAt string `json:"updated_at" tfsdk:"updated_at"`
}
type CreatedOauthToken struct {
	ActorDisplayName      *string   `json:"actor_display_name,omitempty" tfsdk:"actor_display_name"`
	ActorId               *string   `json:"actor_id,omitempty" tfsdk:"actor_id"`
	DisplayName           *string   `json:"display_name,omitempty" tfsdk:"display_name"`
	Name                  *string   `json:"name,omitempty" tfsdk:"name"`
	PlainTextRefreshToken *string   `json:"plain_text_refresh_token,omitempty" tfsdk:"plain_text_refresh_token"`
	ServiceTokenAccesses  *[]string `json:"service_token_accesses,omitempty" tfsdk:"service_token_accesses"`
	Token                 *string   `json:"token,omitempty" tfsdk:"token"`
}
type LintError struct {
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
type PasswordWithPlaintext struct {
	AccessHostUrl  string            `json:"access_host_url" tfsdk:"access_host_url"`
	Actor          *Actor            `json:"actor,omitempty" tfsdk:"actor"`
	CreatedAt      string            `json:"created_at" tfsdk:"created_at"`
	DatabaseBranch BranchForPassword `json:"database_branch" tfsdk:"database_branch"`
	DeletedAt      *string           `json:"deleted_at,omitempty" tfsdk:"deleted_at"`
	ExpiresAt      *string           `json:"expires_at,omitempty" tfsdk:"expires_at"`
	Id             string            `json:"id" tfsdk:"id"`
	Integrations   []string          `json:"integrations" tfsdk:"integrations"`
	Name           string            `json:"name" tfsdk:"name"`
	PlainText      string            `json:"plain_text" tfsdk:"plain_text"`
	Region         *Region           `json:"region,omitempty" tfsdk:"region"`
	Renewable      bool              `json:"renewable" tfsdk:"renewable"`
	Role           string            `json:"role" tfsdk:"role"`
	TtlSeconds     float64           `json:"ttl_seconds" tfsdk:"ttl_seconds"`
	Username       *string           `json:"username,omitempty" tfsdk:"username"`
}
type User struct {
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
type ListOrganizationsRes500 struct {
	*ErrorResponse
}
type ListOrganizationsRes struct {
	Data []Organization `json:"data" tfsdk:"data"`
}
type ListOrganizationsRes401 struct {
	*ErrorResponse
}
type ListOrganizationsRes403 struct {
	*ErrorResponse
}
type ListOrganizationsRes404 struct {
	*ErrorResponse
}

func (cl *Client) ListOrganizations(ctx context.Context, page *int, perPage *int) (res200 *ListOrganizationsRes, err error) {
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
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListOrganizationsRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(ListOrganizationsRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(ListOrganizationsRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(ListOrganizationsRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(ListOrganizationsRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type GetOrganizationRes struct {
	Organization
}
type GetOrganizationRes401 struct {
	*ErrorResponse
}
type GetOrganizationRes403 struct {
	*ErrorResponse
}
type GetOrganizationRes404 struct {
	*ErrorResponse
}
type GetOrganizationRes500 struct {
	*ErrorResponse
}

func (cl *Client) GetOrganization(ctx context.Context, name string) (res200 *GetOrganizationRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + name})
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(GetOrganizationRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(GetOrganizationRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(GetOrganizationRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(GetOrganizationRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(GetOrganizationRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type UpdateOrganizationReq struct {
	BillingEmail                    *string `json:"billing_email,omitempty" tfsdk:"billing_email"`
	IdpManagedRoles                 *bool   `json:"idp_managed_roles,omitempty" tfsdk:"idp_managed_roles"`
	RequireAdminForProductionAccess *bool   `json:"require_admin_for_production_access,omitempty" tfsdk:"require_admin_for_production_access"`
}
type UpdateOrganizationRes500 struct {
	*ErrorResponse
}
type UpdateOrganizationRes struct {
	Organization
}
type UpdateOrganizationRes401 struct {
	*ErrorResponse
}
type UpdateOrganizationRes403 struct {
	*ErrorResponse
}
type UpdateOrganizationRes404 struct {
	*ErrorResponse
}

func (cl *Client) UpdateOrganization(ctx context.Context, name string, req UpdateOrganizationReq) (res200 *UpdateOrganizationRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + name})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res200, err
	}
	r, err := http.NewRequestWithContext(ctx, "PATCH", u.String(), body)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(UpdateOrganizationRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(UpdateOrganizationRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(UpdateOrganizationRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(UpdateOrganizationRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(UpdateOrganizationRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type ListRegionsForOrganizationRes403 struct {
	*ErrorResponse
}
type ListRegionsForOrganizationRes404 struct {
	*ErrorResponse
}
type ListRegionsForOrganizationRes500 struct {
	*ErrorResponse
}
type ListRegionsForOrganizationRes struct {
	Data []Region `json:"data" tfsdk:"data"`
}
type ListRegionsForOrganizationRes401 struct {
	*ErrorResponse
}

func (cl *Client) ListRegionsForOrganization(ctx context.Context, name string, page *int, perPage *int) (res200 *ListRegionsForOrganizationRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + name + "/regions"})
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
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListRegionsForOrganizationRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(ListRegionsForOrganizationRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(ListRegionsForOrganizationRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(ListRegionsForOrganizationRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(ListRegionsForOrganizationRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type ListDatabasesRes struct {
	Data []Database `json:"data" tfsdk:"data"`
}
type ListDatabasesRes401 struct {
	*ErrorResponse
}
type ListDatabasesRes403 struct {
	*ErrorResponse
}
type ListDatabasesRes404 struct {
	*ErrorResponse
}
type ListDatabasesRes500 struct {
	*ErrorResponse
}

func (cl *Client) ListDatabases(ctx context.Context, organization string, page *int, perPage *int) (res200 *ListDatabasesRes, err error) {
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
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListDatabasesRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(ListDatabasesRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(ListDatabasesRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(ListDatabasesRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(ListDatabasesRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type CreateDatabaseReq struct {
	ClusterSize *string `json:"cluster_size,omitempty" tfsdk:"cluster_size"`
	Name        string  `json:"name" tfsdk:"name"`
	Notes       *string `json:"notes,omitempty" tfsdk:"notes"`
	Plan        *string `json:"plan,omitempty" tfsdk:"plan"`
	Region      *string `json:"region,omitempty" tfsdk:"region"`
}
type CreateDatabaseRes404 struct {
	*ErrorResponse
}
type CreateDatabaseRes500 struct {
	*ErrorResponse
}
type CreateDatabaseRes struct {
	Database
}
type CreateDatabaseRes401 struct {
	*ErrorResponse
}
type CreateDatabaseRes403 struct {
	*ErrorResponse
}

func (cl *Client) CreateDatabase(ctx context.Context, organization string, req CreateDatabaseReq) (res201 *CreateDatabaseRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases"})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res201, err
	}
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return res201, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res201, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 201:
		res201 = new(CreateDatabaseRes)
		err = json.NewDecoder(res.Body).Decode(&res201)
	case 401:
		res401 := new(CreateDatabaseRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(CreateDatabaseRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(CreateDatabaseRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(CreateDatabaseRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res201, err
}

type ListBranchesRes struct {
	Data []Branch `json:"data" tfsdk:"data"`
}
type ListBranchesRes401 struct {
	*ErrorResponse
}
type ListBranchesRes403 struct {
	*ErrorResponse
}
type ListBranchesRes404 struct {
	*ErrorResponse
}
type ListBranchesRes500 struct {
	*ErrorResponse
}

func (cl *Client) ListBranches(ctx context.Context, organization string, database string, page *int, perPage *int) (res200 *ListBranchesRes, err error) {
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
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListBranchesRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(ListBranchesRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(ListBranchesRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(ListBranchesRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(ListBranchesRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type CreateBranchReq struct {
	BackupId     *string `json:"backup_id,omitempty" tfsdk:"backup_id"`
	Name         string  `json:"name" tfsdk:"name"`
	ParentBranch string  `json:"parent_branch" tfsdk:"parent_branch"`
}
type CreateBranchRes401 struct {
	*ErrorResponse
}
type CreateBranchRes403 struct {
	*ErrorResponse
}
type CreateBranchRes404 struct {
	*ErrorResponse
}
type CreateBranchRes500 struct {
	*ErrorResponse
}
type CreateBranchRes struct {
	Branch
}

func (cl *Client) CreateBranch(ctx context.Context, organization string, database string, req CreateBranchReq) (res201 *CreateBranchRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches"})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res201, err
	}
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return res201, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res201, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 201:
		res201 = new(CreateBranchRes)
		err = json.NewDecoder(res.Body).Decode(&res201)
	case 401:
		res401 := new(CreateBranchRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(CreateBranchRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(CreateBranchRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(CreateBranchRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res201, err
}

type ListBackupsRes500 struct {
	*ErrorResponse
}
type ListBackupsRes struct {
	Data []Backup `json:"data" tfsdk:"data"`
}
type ListBackupsRes401 struct {
	*ErrorResponse
}
type ListBackupsRes403 struct {
	*ErrorResponse
}
type ListBackupsRes404 struct {
	*ErrorResponse
}

func (cl *Client) ListBackups(ctx context.Context, organization string, database string, branch string, page *int, perPage *int) (res200 *ListBackupsRes, err error) {
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
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListBackupsRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(ListBackupsRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(ListBackupsRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(ListBackupsRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(ListBackupsRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type CreateBackupReq struct {
	Name           *string  `json:"name,omitempty" tfsdk:"name"`
	RetentionUnit  *string  `json:"retention_unit,omitempty" tfsdk:"retention_unit"`
	RetentionValue *float64 `json:"retention_value,omitempty" tfsdk:"retention_value"`
}
type CreateBackupRes struct {
	Backup
}
type CreateBackupRes401 struct {
	*ErrorResponse
}
type CreateBackupRes403 struct {
	*ErrorResponse
}
type CreateBackupRes404 struct {
	*ErrorResponse
}
type CreateBackupRes500 struct {
	*ErrorResponse
}

func (cl *Client) CreateBackup(ctx context.Context, organization string, database string, branch string, req CreateBackupReq) (res201 *CreateBackupRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + branch + "/backups"})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res201, err
	}
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return res201, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res201, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 201:
		res201 = new(CreateBackupRes)
		err = json.NewDecoder(res.Body).Decode(&res201)
	case 401:
		res401 := new(CreateBackupRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(CreateBackupRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(CreateBackupRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(CreateBackupRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res201, err
}

type GetBackupRes500 struct {
	*ErrorResponse
}
type GetBackupRes struct {
	Backup
}
type GetBackupRes401 struct {
	*ErrorResponse
}
type GetBackupRes403 struct {
	*ErrorResponse
}
type GetBackupRes404 struct {
	*ErrorResponse
}

func (cl *Client) GetBackup(ctx context.Context, organization string, database string, branch string, id string) (res200 *GetBackupRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + branch + "/backups/" + id})
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(GetBackupRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(GetBackupRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(GetBackupRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(GetBackupRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(GetBackupRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type DeleteBackupRes404 struct {
	*ErrorResponse
}
type DeleteBackupRes500 struct {
	*ErrorResponse
}
type DeleteBackupRes struct{}
type DeleteBackupRes401 struct {
	*ErrorResponse
}
type DeleteBackupRes403 struct {
	*ErrorResponse
}

func (cl *Client) DeleteBackup(ctx context.Context, organization string, database string, branch string, id string) (res204 *DeleteBackupRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + branch + "/backups/" + id})
	r, err := http.NewRequestWithContext(ctx, "DELETE", u.String(), nil)
	if err != nil {
		return res204, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res204, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 204:
		res204 = new(DeleteBackupRes)
		err = json.NewDecoder(res.Body).Decode(&res204)
	case 401:
		res401 := new(DeleteBackupRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(DeleteBackupRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(DeleteBackupRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(DeleteBackupRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res204, err
}

type ListPasswordsRes403 struct {
	*ErrorResponse
}
type ListPasswordsRes404 struct {
	*ErrorResponse
}
type ListPasswordsRes500 struct {
	*ErrorResponse
}
type ListPasswordsRes struct {
	Data []Password `json:"data" tfsdk:"data"`
}
type ListPasswordsRes401 struct {
	*ErrorResponse
}

func (cl *Client) ListPasswords(ctx context.Context, organization string, database string, branch string, readOnlyRegionId *string, page *int, perPage *int) (res200 *ListPasswordsRes, err error) {
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
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListPasswordsRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(ListPasswordsRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(ListPasswordsRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(ListPasswordsRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(ListPasswordsRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type CreatePasswordReq struct {
	Name *string  `json:"name,omitempty" tfsdk:"name"`
	Role *string  `json:"role,omitempty" tfsdk:"role"`
	Ttl  *float64 `json:"ttl,omitempty" tfsdk:"ttl"`
}
type CreatePasswordRes500 struct {
	*ErrorResponse
}
type CreatePasswordRes struct {
	PasswordWithPlaintext
}
type CreatePasswordRes401 struct {
	*ErrorResponse
}
type CreatePasswordRes403 struct {
	*ErrorResponse
}
type CreatePasswordRes404 struct {
	*ErrorResponse
}
type CreatePasswordRes422 struct {
	*ErrorResponse
}

func (cl *Client) CreatePassword(ctx context.Context, organization string, database string, branch string, req CreatePasswordReq) (res201 *CreatePasswordRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + branch + "/passwords"})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res201, err
	}
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return res201, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res201, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 201:
		res201 = new(CreatePasswordRes)
		err = json.NewDecoder(res.Body).Decode(&res201)
	case 401:
		res401 := new(CreatePasswordRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(CreatePasswordRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(CreatePasswordRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 422:
		res422 := new(CreatePasswordRes422)
		err = json.NewDecoder(res.Body).Decode(&res422)
		if err == nil {
			err = res422
		}
	case 500:
		res500 := new(CreatePasswordRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res201, err
}

type GetPasswordRes struct {
	Password
}
type GetPasswordRes401 struct {
	*ErrorResponse
}
type GetPasswordRes403 struct {
	*ErrorResponse
}
type GetPasswordRes404 struct {
	*ErrorResponse
}
type GetPasswordRes500 struct {
	*ErrorResponse
}

func (cl *Client) GetPassword(ctx context.Context, organization string, database string, branch string, id string, readOnlyRegionId *string) (res200 *GetPasswordRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + branch + "/passwords/" + id})
	q := u.Query()
	if readOnlyRegionId != nil {
		q.Set("read_only_region_id", *readOnlyRegionId)
	}
	u.RawQuery = q.Encode()
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(GetPasswordRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(GetPasswordRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(GetPasswordRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(GetPasswordRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(GetPasswordRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type DeletePasswordRes403 struct {
	*ErrorResponse
}
type DeletePasswordRes404 struct {
	*ErrorResponse
}
type DeletePasswordRes500 struct {
	*ErrorResponse
}
type DeletePasswordRes struct{}
type DeletePasswordRes401 struct {
	*ErrorResponse
}

func (cl *Client) DeletePassword(ctx context.Context, organization string, database string, branch string, id string) (res204 *DeletePasswordRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + branch + "/passwords/" + id})
	r, err := http.NewRequestWithContext(ctx, "DELETE", u.String(), nil)
	if err != nil {
		return res204, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res204, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 204:
		res204 = new(DeletePasswordRes)
		err = json.NewDecoder(res.Body).Decode(&res204)
	case 401:
		res401 := new(DeletePasswordRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(DeletePasswordRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(DeletePasswordRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(DeletePasswordRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res204, err
}

type UpdatePasswordReq struct {
	Name string `json:"name" tfsdk:"name"`
}
type UpdatePasswordRes401 struct {
	*ErrorResponse
}
type UpdatePasswordRes403 struct {
	*ErrorResponse
}
type UpdatePasswordRes404 struct {
	*ErrorResponse
}
type UpdatePasswordRes500 struct {
	*ErrorResponse
}
type UpdatePasswordRes struct {
	Password
}

func (cl *Client) UpdatePassword(ctx context.Context, organization string, database string, branch string, id string, req UpdatePasswordReq) (res200 *UpdatePasswordRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + branch + "/passwords/" + id})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res200, err
	}
	r, err := http.NewRequestWithContext(ctx, "PATCH", u.String(), body)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(UpdatePasswordRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(UpdatePasswordRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(UpdatePasswordRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(UpdatePasswordRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(UpdatePasswordRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type RenewPasswordReq struct {
	ReadOnlyRegionId *string `json:"read_only_region_id,omitempty" tfsdk:"read_only_region_id"`
}
type RenewPasswordRes500 struct {
	*ErrorResponse
}
type RenewPasswordRes struct {
	PasswordWithPlaintext
}
type RenewPasswordRes401 struct {
	*ErrorResponse
}
type RenewPasswordRes403 struct {
	*ErrorResponse
}
type RenewPasswordRes404 struct {
	*ErrorResponse
}

func (cl *Client) RenewPassword(ctx context.Context, organization string, database string, branch string, id string, req RenewPasswordReq) (res200 *RenewPasswordRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + branch + "/passwords/" + id + "/renew"})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res200, err
	}
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(RenewPasswordRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(RenewPasswordRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(RenewPasswordRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(RenewPasswordRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(RenewPasswordRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type GetBranchRes401 struct {
	*ErrorResponse
}
type GetBranchRes403 struct {
	*ErrorResponse
}
type GetBranchRes404 struct {
	*ErrorResponse
}
type GetBranchRes500 struct {
	*ErrorResponse
}
type GetBranchRes struct {
	Branch
}

func (cl *Client) GetBranch(ctx context.Context, organization string, database string, name string) (res200 *GetBranchRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + name})
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(GetBranchRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(GetBranchRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(GetBranchRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(GetBranchRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(GetBranchRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type DeleteBranchRes401 struct {
	*ErrorResponse
}
type DeleteBranchRes403 struct {
	*ErrorResponse
}
type DeleteBranchRes404 struct {
	*ErrorResponse
}
type DeleteBranchRes500 struct {
	*ErrorResponse
}
type DeleteBranchRes struct{}

func (cl *Client) DeleteBranch(ctx context.Context, organization string, database string, name string) (res204 *DeleteBranchRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + name})
	r, err := http.NewRequestWithContext(ctx, "DELETE", u.String(), nil)
	if err != nil {
		return res204, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res204, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 204:
		res204 = new(DeleteBranchRes)
		err = json.NewDecoder(res.Body).Decode(&res204)
	case 401:
		res401 := new(DeleteBranchRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(DeleteBranchRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(DeleteBranchRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(DeleteBranchRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res204, err
}

type DemoteBranchRes struct {
	Branch
}
type DemoteBranchRes401 struct {
	*ErrorResponse
}
type DemoteBranchRes403 struct {
	*ErrorResponse
}
type DemoteBranchRes404 struct {
	*ErrorResponse
}
type DemoteBranchRes500 struct {
	*ErrorResponse
}

func (cl *Client) DemoteBranch(ctx context.Context, organization string, database string, name string) (res200 *DemoteBranchRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + name + "/demote"})
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(DemoteBranchRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(DemoteBranchRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(DemoteBranchRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(DemoteBranchRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(DemoteBranchRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type PromoteBranchRes401 struct {
	*ErrorResponse
}
type PromoteBranchRes403 struct {
	*ErrorResponse
}
type PromoteBranchRes404 struct {
	*ErrorResponse
}
type PromoteBranchRes500 struct {
	*ErrorResponse
}
type PromoteBranchRes struct {
	Branch
}

func (cl *Client) PromoteBranch(ctx context.Context, organization string, database string, name string) (res200 *PromoteBranchRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + name + "/promote"})
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(PromoteBranchRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(PromoteBranchRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(PromoteBranchRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(PromoteBranchRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(PromoteBranchRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type EnableSafeMigrationsForBranchRes struct {
	Branch
}
type EnableSafeMigrationsForBranchRes401 struct {
	*ErrorResponse
}
type EnableSafeMigrationsForBranchRes403 struct {
	*ErrorResponse
}
type EnableSafeMigrationsForBranchRes404 struct {
	*ErrorResponse
}
type EnableSafeMigrationsForBranchRes500 struct {
	*ErrorResponse
}

func (cl *Client) EnableSafeMigrationsForBranch(ctx context.Context, organization string, database string, name string) (res200 *EnableSafeMigrationsForBranchRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + name + "/safe-migrations"})
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(EnableSafeMigrationsForBranchRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(EnableSafeMigrationsForBranchRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(EnableSafeMigrationsForBranchRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(EnableSafeMigrationsForBranchRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(EnableSafeMigrationsForBranchRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type DisableSafeMigrationsForBranchRes struct {
	Branch
}
type DisableSafeMigrationsForBranchRes401 struct {
	*ErrorResponse
}
type DisableSafeMigrationsForBranchRes403 struct {
	*ErrorResponse
}
type DisableSafeMigrationsForBranchRes404 struct {
	*ErrorResponse
}
type DisableSafeMigrationsForBranchRes500 struct {
	*ErrorResponse
}

func (cl *Client) DisableSafeMigrationsForBranch(ctx context.Context, organization string, database string, name string) (res200 *DisableSafeMigrationsForBranchRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + name + "/safe-migrations"})
	r, err := http.NewRequestWithContext(ctx, "DELETE", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(DisableSafeMigrationsForBranchRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(DisableSafeMigrationsForBranchRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(DisableSafeMigrationsForBranchRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(DisableSafeMigrationsForBranchRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(DisableSafeMigrationsForBranchRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type GetBranchSchemaRes403 struct {
	*ErrorResponse
}
type GetBranchSchemaRes404 struct {
	*ErrorResponse
}
type GetBranchSchemaRes500 struct {
	*ErrorResponse
}
type GetBranchSchemaRes struct {
	Data []TableSchema `json:"data" tfsdk:"data"`
}
type GetBranchSchemaRes401 struct {
	*ErrorResponse
}

func (cl *Client) GetBranchSchema(ctx context.Context, organization string, database string, name string, keyspace *string) (res200 *GetBranchSchemaRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/branches/" + name + "/schema"})
	q := u.Query()
	if keyspace != nil {
		q.Set("keyspace", *keyspace)
	}
	u.RawQuery = q.Encode()
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(GetBranchSchemaRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(GetBranchSchemaRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(GetBranchSchemaRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(GetBranchSchemaRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(GetBranchSchemaRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type LintBranchSchemaRes struct {
	Data []LintError `json:"data" tfsdk:"data"`
}
type LintBranchSchemaRes401 struct {
	*ErrorResponse
}
type LintBranchSchemaRes403 struct {
	*ErrorResponse
}
type LintBranchSchemaRes404 struct {
	*ErrorResponse
}
type LintBranchSchemaRes500 struct {
	*ErrorResponse
}

func (cl *Client) LintBranchSchema(ctx context.Context, organization string, database string, name string, page *int, perPage *int) (res200 *LintBranchSchemaRes, err error) {
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
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(LintBranchSchemaRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(LintBranchSchemaRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(LintBranchSchemaRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(LintBranchSchemaRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(LintBranchSchemaRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type GetTheDeployQueueRes struct {
	Data []QueuedDeployRequest `json:"data" tfsdk:"data"`
}

func (cl *Client) GetTheDeployQueue(ctx context.Context, organization string, database string) (res200 *GetTheDeployQueueRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-queue"})
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(GetTheDeployQueueRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type ListDeployRequestsRes struct {
	Data []DeployRequest `json:"data" tfsdk:"data"`
}

func (cl *Client) ListDeployRequests(ctx context.Context, organization string, database string, page *int, perPage *int, state *string, branch *string, intoBranch *string) (res200 *ListDeployRequestsRes, err error) {
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
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListDeployRequestsRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
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
type CreateDeployRequestRes struct {
	DeployRequestWithDeployment
}

func (cl *Client) CreateDeployRequest(ctx context.Context, organization string, database string, req CreateDeployRequestReq) (res201 *CreateDeployRequestRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests"})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res201, err
	}
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return res201, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res201, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 201:
		res201 = new(CreateDeployRequestRes)
		err = json.NewDecoder(res.Body).Decode(&res201)
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res201, err
}

type GetDeployRequestRes struct {
	DeployRequestWithDeployment
}

func (cl *Client) GetDeployRequest(ctx context.Context, organization string, database string, number string) (res200 *GetDeployRequestRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests/" + number})
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(GetDeployRequestRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type CloseDeployRequestReq struct {
	State *string `json:"state,omitempty" tfsdk:"state"`
}
type CloseDeployRequestRes struct {
	DeployRequestWithDeployment
}

func (cl *Client) CloseDeployRequest(ctx context.Context, organization string, database string, number string, req CloseDeployRequestReq) (res200 *CloseDeployRequestRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests/" + number})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res200, err
	}
	r, err := http.NewRequestWithContext(ctx, "PATCH", u.String(), body)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(CloseDeployRequestRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type CompleteGatedDeployRequestRes struct {
	DeployRequest
}

func (cl *Client) CompleteGatedDeployRequest(ctx context.Context, organization string, database string, number string) (res200 *CompleteGatedDeployRequestRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests/" + number + "/apply-deploy"})
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(CompleteGatedDeployRequestRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type UpdateAutoApplyForDeployRequestReq struct {
	Enable *bool `json:"enable,omitempty" tfsdk:"enable"`
}
type UpdateAutoApplyForDeployRequestRes struct {
	DeployRequest
}

func (cl *Client) UpdateAutoApplyForDeployRequest(ctx context.Context, organization string, database string, number string, req UpdateAutoApplyForDeployRequestReq) (res200 *UpdateAutoApplyForDeployRequestRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests/" + number + "/auto-apply"})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res200, err
	}
	r, err := http.NewRequestWithContext(ctx, "PUT", u.String(), body)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(UpdateAutoApplyForDeployRequestRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type CancelQueuedDeployRequestRes struct {
	DeployRequest
}

func (cl *Client) CancelQueuedDeployRequest(ctx context.Context, organization string, database string, number string) (res200 *CancelQueuedDeployRequestRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests/" + number + "/cancel"})
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(CancelQueuedDeployRequestRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type CompleteErroredDeployRes struct {
	DeployRequest
}

func (cl *Client) CompleteErroredDeploy(ctx context.Context, organization string, database string, number string) (res200 *CompleteErroredDeployRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests/" + number + "/complete-deploy"})
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(CompleteErroredDeployRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type QueueDeployRequestRes struct {
	DeployRequest
}

func (cl *Client) QueueDeployRequest(ctx context.Context, organization string, database string, number string) (res200 *QueueDeployRequestRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests/" + number + "/deploy"})
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(QueueDeployRequestRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type GetDeploymentRes struct {
	Deployment
}

func (cl *Client) GetDeployment(ctx context.Context, organization string, database string, number string) (res200 *GetDeploymentRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests/" + number + "/deployment"})
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(GetDeploymentRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type ListDeployOperationsRes struct {
	Data []DeployOperation `json:"data" tfsdk:"data"`
}

func (cl *Client) ListDeployOperations(ctx context.Context, organization string, database string, number string, page *int, perPage *int) (res200 *ListDeployOperationsRes, err error) {
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
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListDeployOperationsRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type CompleteRevertRes struct {
	DeployRequest
}

func (cl *Client) CompleteRevert(ctx context.Context, organization string, database string, number string) (res200 *CompleteRevertRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests/" + number + "/revert"})
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(CompleteRevertRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type ListDeployRequestReviewsRes struct {
	Data []DeployReview `json:"data" tfsdk:"data"`
}

func (cl *Client) ListDeployRequestReviews(ctx context.Context, organization string, database string, number string) (res200 *ListDeployRequestReviewsRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests/" + number + "/reviews"})
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListDeployRequestReviewsRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
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
type ReviewDeployRequestRes struct {
	DeployReview
}

func (cl *Client) ReviewDeployRequest(ctx context.Context, organization string, database string, number string, req ReviewDeployRequestReq) (res201 *ReviewDeployRequestRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests/" + number + "/reviews"})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res201, err
	}
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return res201, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res201, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 201:
		res201 = new(ReviewDeployRequestRes)
		err = json.NewDecoder(res.Body).Decode(&res201)
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res201, err
}

type SkipRevertPeriodRes struct {
	DeployRequest
}

func (cl *Client) SkipRevertPeriod(ctx context.Context, organization string, database string, number string) (res200 *SkipRevertPeriodRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + database + "/deploy-requests/" + number + "/skip-revert"})
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(SkipRevertPeriodRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type GetDatabaseRes404 struct {
	*ErrorResponse
}
type GetDatabaseRes500 struct {
	*ErrorResponse
}
type GetDatabaseRes struct {
	Database
}
type GetDatabaseRes401 struct {
	*ErrorResponse
}
type GetDatabaseRes403 struct {
	*ErrorResponse
}

func (cl *Client) GetDatabase(ctx context.Context, organization string, name string) (res200 *GetDatabaseRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + name})
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(GetDatabaseRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(GetDatabaseRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(GetDatabaseRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(GetDatabaseRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(GetDatabaseRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type DeleteDatabaseRes500 struct {
	*ErrorResponse
}
type DeleteDatabaseRes struct{}
type DeleteDatabaseRes401 struct {
	*ErrorResponse
}
type DeleteDatabaseRes403 struct {
	*ErrorResponse
}
type DeleteDatabaseRes404 struct {
	*ErrorResponse
}

func (cl *Client) DeleteDatabase(ctx context.Context, organization string, name string) (res204 *DeleteDatabaseRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + name})
	r, err := http.NewRequestWithContext(ctx, "DELETE", u.String(), nil)
	if err != nil {
		return res204, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res204, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 204:
		res204 = new(DeleteDatabaseRes)
		err = json.NewDecoder(res.Body).Decode(&res204)
	case 401:
		res401 := new(DeleteDatabaseRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(DeleteDatabaseRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(DeleteDatabaseRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(DeleteDatabaseRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res204, err
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
type UpdateDatabaseSettingsRes403 struct {
	*ErrorResponse
}
type UpdateDatabaseSettingsRes404 struct {
	*ErrorResponse
}
type UpdateDatabaseSettingsRes500 struct {
	*ErrorResponse
}
type UpdateDatabaseSettingsRes struct {
	Database
}
type UpdateDatabaseSettingsRes401 struct {
	*ErrorResponse
}

func (cl *Client) UpdateDatabaseSettings(ctx context.Context, organization string, name string, req UpdateDatabaseSettingsReq) (res200 *UpdateDatabaseSettingsRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/databases/" + name})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res200, err
	}
	r, err := http.NewRequestWithContext(ctx, "PATCH", u.String(), body)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(UpdateDatabaseSettingsRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(UpdateDatabaseSettingsRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(UpdateDatabaseSettingsRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(UpdateDatabaseSettingsRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(UpdateDatabaseSettingsRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type ListReadOnlyRegionsRes404 struct {
	*ErrorResponse
}
type ListReadOnlyRegionsRes500 struct {
	*ErrorResponse
}
type ListReadOnlyRegionsRes struct {
	Data []ReadOnlyRegion `json:"data" tfsdk:"data"`
}
type ListReadOnlyRegionsRes401 struct {
	*ErrorResponse
}
type ListReadOnlyRegionsRes403 struct {
	*ErrorResponse
}

func (cl *Client) ListReadOnlyRegions(ctx context.Context, organization string, name string, page *int, perPage *int) (res200 *ListReadOnlyRegionsRes, err error) {
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
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListReadOnlyRegionsRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(ListReadOnlyRegionsRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(ListReadOnlyRegionsRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(ListReadOnlyRegionsRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(ListReadOnlyRegionsRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type ListDatabaseRegionsRes401 struct {
	*ErrorResponse
}
type ListDatabaseRegionsRes403 struct {
	*ErrorResponse
}
type ListDatabaseRegionsRes404 struct {
	*ErrorResponse
}
type ListDatabaseRegionsRes500 struct {
	*ErrorResponse
}
type ListDatabaseRegionsRes struct {
	Data []Region `json:"data" tfsdk:"data"`
}

func (cl *Client) ListDatabaseRegions(ctx context.Context, organization string, name string, page *int, perPage *int) (res200 *ListDatabaseRegionsRes, err error) {
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
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListDatabaseRegionsRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(ListDatabaseRegionsRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(ListDatabaseRegionsRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(ListDatabaseRegionsRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(ListDatabaseRegionsRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type ListOauthApplicationsRes403 struct {
	*ErrorResponse
}
type ListOauthApplicationsRes404 struct {
	*ErrorResponse
}
type ListOauthApplicationsRes500 struct {
	*ErrorResponse
}
type ListOauthApplicationsRes struct {
	Data []OauthApplication `json:"data" tfsdk:"data"`
}
type ListOauthApplicationsRes401 struct {
	*ErrorResponse
}

func (cl *Client) ListOauthApplications(ctx context.Context, organization string, page *int, perPage *int) (res200 *ListOauthApplicationsRes, err error) {
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
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListOauthApplicationsRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(ListOauthApplicationsRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(ListOauthApplicationsRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(ListOauthApplicationsRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(ListOauthApplicationsRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type GetOauthApplicationRes401 struct {
	*ErrorResponse
}
type GetOauthApplicationRes403 struct {
	*ErrorResponse
}
type GetOauthApplicationRes404 struct {
	*ErrorResponse
}
type GetOauthApplicationRes500 struct {
	*ErrorResponse
}
type GetOauthApplicationRes struct {
	OauthApplication
}

func (cl *Client) GetOauthApplication(ctx context.Context, organization string, applicationId string) (res200 *GetOauthApplicationRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/oauth-applications/" + applicationId})
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(GetOauthApplicationRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(GetOauthApplicationRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(GetOauthApplicationRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(GetOauthApplicationRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(GetOauthApplicationRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type ListOauthTokensRes401 struct {
	*ErrorResponse
}
type ListOauthTokensRes403 struct {
	*ErrorResponse
}
type ListOauthTokensRes404 struct {
	*ErrorResponse
}
type ListOauthTokensRes500 struct {
	*ErrorResponse
}
type ListOauthTokensRes struct {
	Data []OauthToken `json:"data" tfsdk:"data"`
}

func (cl *Client) ListOauthTokens(ctx context.Context, organization string, applicationId string, page *int, perPage *int) (res200 *ListOauthTokensRes, err error) {
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
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(ListOauthTokensRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(ListOauthTokensRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(ListOauthTokensRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(ListOauthTokensRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(ListOauthTokensRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type GetOauthTokenRes401 struct {
	*ErrorResponse
}
type GetOauthTokenRes403 struct {
	*ErrorResponse
}
type GetOauthTokenRes404 struct {
	*ErrorResponse
}
type GetOauthTokenRes500 struct {
	*ErrorResponse
}
type GetOauthTokenRes struct {
	OauthTokenWithDetails
}

func (cl *Client) GetOauthToken(ctx context.Context, organization string, applicationId string, tokenId string) (res200 *GetOauthTokenRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/oauth-applications/" + applicationId + "/tokens/" + tokenId})
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(GetOauthTokenRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(GetOauthTokenRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(GetOauthTokenRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(GetOauthTokenRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(GetOauthTokenRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type DeleteOauthTokenRes404 struct {
	*ErrorResponse
}
type DeleteOauthTokenRes500 struct {
	*ErrorResponse
}
type DeleteOauthTokenRes struct{}
type DeleteOauthTokenRes401 struct {
	*ErrorResponse
}
type DeleteOauthTokenRes403 struct {
	*ErrorResponse
}

func (cl *Client) DeleteOauthToken(ctx context.Context, organization string, applicationId string, tokenId string) (res204 *DeleteOauthTokenRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/oauth-applications/" + applicationId + "/tokens/" + tokenId})
	r, err := http.NewRequestWithContext(ctx, "DELETE", u.String(), nil)
	if err != nil {
		return res204, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res204, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 204:
		res204 = new(DeleteOauthTokenRes)
		err = json.NewDecoder(res.Body).Decode(&res204)
	case 401:
		res401 := new(DeleteOauthTokenRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(DeleteOauthTokenRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(DeleteOauthTokenRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(DeleteOauthTokenRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res204, err
}

type CreateOrRenewOauthTokenReq struct {
	ClientId     string  `json:"client_id" tfsdk:"client_id"`
	ClientSecret string  `json:"client_secret" tfsdk:"client_secret"`
	Code         *string `json:"code,omitempty" tfsdk:"code"`
	GrantType    string  `json:"grant_type" tfsdk:"grant_type"`
	RedirectUri  *string `json:"redirect_uri,omitempty" tfsdk:"redirect_uri"`
	RefreshToken *string `json:"refresh_token,omitempty" tfsdk:"refresh_token"`
}
type CreateOrRenewOauthTokenRes404 struct {
	*ErrorResponse
}
type CreateOrRenewOauthTokenRes422 struct {
	*ErrorResponse
}
type CreateOrRenewOauthTokenRes500 struct {
	*ErrorResponse
}
type CreateOrRenewOauthTokenRes struct {
	CreatedOauthToken
}
type CreateOrRenewOauthTokenRes403 struct {
	*ErrorResponse
}

func (cl *Client) CreateOrRenewOauthToken(ctx context.Context, organization string, id string, req CreateOrRenewOauthTokenReq) (res200 *CreateOrRenewOauthTokenRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "organizations/" + organization + "/oauth-applications/" + id + "/token"})
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(req); err != nil {
		return res200, err
	}
	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(CreateOrRenewOauthTokenRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 403:
		res403 := new(CreateOrRenewOauthTokenRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(CreateOrRenewOauthTokenRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 422:
		res422 := new(CreateOrRenewOauthTokenRes422)
		err = json.NewDecoder(res.Body).Decode(&res422)
		if err == nil {
			err = res422
		}
	case 500:
		res500 := new(CreateOrRenewOauthTokenRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}

type GetCurrentUserRes500 struct {
	*ErrorResponse
}
type GetCurrentUserRes struct {
	User
}
type GetCurrentUserRes401 struct {
	*ErrorResponse
}
type GetCurrentUserRes403 struct {
	*ErrorResponse
}
type GetCurrentUserRes404 struct {
	*ErrorResponse
}

func (cl *Client) GetCurrentUser(ctx context.Context) (res200 *GetCurrentUserRes, err error) {
	u := cl.baseURL.ResolveReference(&url.URL{Path: "user"})
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return res200, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	res, err := cl.httpCl.Do(r)
	if err != nil {
		return res200, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		res200 = new(GetCurrentUserRes)
		err = json.NewDecoder(res.Body).Decode(&res200)
	case 401:
		res401 := new(GetCurrentUserRes401)
		err = json.NewDecoder(res.Body).Decode(&res401)
		if err == nil {
			err = res401
		}
	case 403:
		res403 := new(GetCurrentUserRes403)
		err = json.NewDecoder(res.Body).Decode(&res403)
		if err == nil {
			err = res403
		}
	case 404:
		res404 := new(GetCurrentUserRes404)
		err = json.NewDecoder(res.Body).Decode(&res404)
		if err == nil {
			err = res404
		}
	case 500:
		res500 := new(GetCurrentUserRes500)
		err = json.NewDecoder(res.Body).Decode(&res500)
		if err == nil {
			err = res500
		}
	default:
		var errBody *ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody != nil {
			err = errBody
		} else {
			err = fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return res200, err
}
