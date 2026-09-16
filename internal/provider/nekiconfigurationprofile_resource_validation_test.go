package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestNekiExtensionsConflictWithSharedPreloadLibraries(t *testing.T) {
	diagnostics := validateNekiExtensions(t, []types.String{types.StringValue("hll")}, map[string]types.String{
		"shared_preload_libraries": types.StringValue("pg_cron"),
	})

	require.Len(t, diagnostics.Errors(), 1)
	require.Equal(t, `Attribute "parameters[\"pgconf\"][\"shared_preload_libraries\"]" cannot be specified when "extensions" is specified`, diagnostics.Errors()[0].Detail())
}

func TestEmptyNekiExtensionsConflictWithSessionPreloadLibraries(t *testing.T) {
	diagnostics := validateNekiExtensions(t, []types.String{}, map[string]types.String{
		"session_preload_libraries": types.StringValue("hll"),
	})

	require.Len(t, diagnostics.Errors(), 1)
	require.Equal(t, `Attribute "parameters[\"pgconf\"][\"session_preload_libraries\"]" cannot be specified when "extensions" is specified`, diagnostics.Errors()[0].Detail())
}

func TestNekiExtensionsAllowExtensionSettings(t *testing.T) {
	diagnostics := validateNekiExtensions(t, []types.String{types.StringValue("hll")}, map[string]types.String{
		"hll.force_groupagg": types.StringValue("on"),
	})

	require.Empty(t, diagnostics.Errors())
}

func TestNekiPreloadLibrariesAllowOmittedExtensions(t *testing.T) {
	diagnostics := validateNekiExtensions(t, nil, map[string]types.String{
		"shared_preload_libraries":  types.StringValue("pg_cron"),
		"session_preload_libraries": types.StringValue("hll"),
	})

	require.Empty(t, diagnostics.Errors())
}

func validateNekiExtensions(t *testing.T, extensions []types.String, pgconf map[string]types.String) diag.Diagnostics {
	t.Helper()
	ctx := context.Background()
	var response resource.SchemaResponse
	(&NekiConfigurationProfileResource{}).Schema(ctx, resource.SchemaRequest{}, &response)
	config, diagnostics := types.ObjectValueFrom(ctx, response.Schema.Type().(types.ObjectType).AttrTypes, &NekiConfigurationProfileResourceModel{
		Extensions: extensions,
		Parameters: map[string]map[string]types.String{"pgconf": pgconf},
	})
	require.False(t, diagnostics.HasError(), diagnostics.Errors())
	raw, err := config.ToTerraformValue(ctx)
	require.NoError(t, err)
	value, diagnostics := types.ListValueFrom(ctx, types.StringType, extensions)
	require.False(t, diagnostics.HasError(), diagnostics.Errors())
	request := validator.ListRequest{
		Path:           path.Root("extensions"),
		PathExpression: path.MatchRoot("extensions"),
		Config:         tfsdk.Config{Raw: raw, Schema: response.Schema},
		ConfigValue:    value,
	}

	var result validator.ListResponse
	for _, check := range response.Schema.Attributes["extensions"].(schema.ListAttribute).Validators {
		check.ValidateList(ctx, request, &result)
	}
	return result.Diagnostics
}
