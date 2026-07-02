package firewall

import (
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
	opnfirewall "github.com/browningluke/opnsense-go/pkg/firewall"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestConvertAliasSchemaToStructIgnoresIPProtocolForHostAlias(t *testing.T) {
	data := &aliasResourceModel{
		Enabled:    types.BoolValue(true),
		Name:       types.StringValue("ip_cameras"),
		Type:       types.StringValue("host"),
		IPProtocol: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("IPv4")}),
		Interface:  types.StringValue(""),
		Content: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("10.0.40.3"),
			types.StringValue("10.0.40.4"),
		}),
		Categories:  types.SetValueMust(types.StringType, []attr.Value{}),
		UpdateFreq:  types.Float64Value(-1),
		Statistics:  types.BoolValue(false),
		Description: types.StringNull(),
	}

	result, err := convertAliasSchemaToStruct(data)

	require.NoError(t, err)
	require.Equal(t, "host", result.Type.String())
	require.Empty(t, result.IPProtocol)
}

func TestConvertAliasSchemaToStructDefaultsIPProtocolForApplicableAlias(t *testing.T) {
	data := &aliasResourceModel{
		Enabled:     types.BoolValue(true),
		Name:        types.StringValue("external_alias"),
		Type:        types.StringValue("external"),
		IPProtocol:  types.SetNull(types.StringType),
		Interface:   types.StringValue(""),
		Content:     types.SetValueMust(types.StringType, []attr.Value{}),
		Categories:  types.SetValueMust(types.StringType, []attr.Value{}),
		UpdateFreq:  types.Float64Value(-1),
		Statistics:  types.BoolValue(false),
		Description: types.StringNull(),
	}

	result, err := convertAliasSchemaToStruct(data)

	require.NoError(t, err)
	require.Equal(t, "external", result.Type.String())
	require.Equal(t, api.SelectedMapList{"IPv4"}, result.IPProtocol)
}

func TestConvertAliasStructToSchemaIgnoresIPProtocolForHostAlias(t *testing.T) {
	result, err := convertAliasStructToSchema(&opnfirewall.Alias{
		Enabled:    "1",
		Name:       "ip_cameras",
		Type:       api.SelectedMap("host"),
		IPProtocol: api.SelectedMapList{"IPv4"},
		Interface:  api.SelectedMap(""),
		Content: api.SelectedMapListNL{
			"10.0.40.3",
			"10.0.40.4",
		},
		Categories: []string{},
		UpdateFreq: "-1",
		Statistics: "0",
	})

	require.NoError(t, err)
	require.Equal(t, "host", result.Type.ValueString())
	require.Empty(t, result.IPProtocol.Elements())
}

func TestConvertAliasStructToSchemaDefaultsIPProtocolForApplicableAlias(t *testing.T) {
	result, err := convertAliasStructToSchema(&opnfirewall.Alias{
		Enabled:    "1",
		Name:       "external_alias",
		Type:       api.SelectedMap("external"),
		IPProtocol: api.SelectedMapList{},
		Interface:  api.SelectedMap(""),
		Content:    api.SelectedMapListNL{},
		Categories: []string{},
		UpdateFreq: "-1",
		Statistics: "0",
	})

	require.NoError(t, err)
	require.Equal(t, "external", result.Type.ValueString())
	require.Len(t, result.IPProtocol.Elements(), 1)
	require.Contains(t, result.IPProtocol.Elements(), types.StringValue("IPv4"))
}
