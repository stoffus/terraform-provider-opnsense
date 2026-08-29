package interfaces

import (
	"regexp"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/interfaces"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type assignmentResourceModel struct {
	Device      types.String `tfsdk:"device"`
	Description types.String `tfsdk:"description"`
	Lock        types.Bool   `tfsdk:"lock"`

	Id types.String `tfsdk:"id"`
}

func assignmentResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Assigns a physical or virtual network device to an OPNsense interface. This resource manages the assignment, description, and lock; enabling the interface and configuring its addresses are outside its scope.",
		Version:             1,

		Attributes: map[string]schema.Attribute{
			"device": schema.StringAttribute{
				MarkdownDescription: "Network device to assign, such as `vlan0100` or `vtnet1`.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Short description used to identify the interface.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9_]+$`), "must contain only letters, numbers, and underscores"),
				},
			},
			"lock": schema.BoolAttribute{
				MarkdownDescription: "When enabled, prevent the interface from being removed accidentally. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "OPNsense interface identifier, such as `opt1`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func assignmentDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Retrieves an assigned OPNsense interface. This data source exposes the assignment, description, and lock only.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "OPNsense interface identifier, such as `opt1`.",
				Required:            true,
			},
			"device": dschema.StringAttribute{
				MarkdownDescription: "Assigned physical or virtual network device.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "Short description used to identify the interface.",
				Computed:            true,
			},
			"lock": dschema.BoolAttribute{
				MarkdownDescription: "Whether the interface is protected from removal.",
				Computed:            true,
			},
		},
	}
}

func convertAssignmentSchemaToStruct(d *assignmentResourceModel) (*interfaces.Assignment, error) {
	return &interfaces.Assignment{
		Device:      api.SelectedMap(d.Device.ValueString()),
		Description: d.Description.ValueString(),
		Lock:        tools.BoolToString(d.Lock.ValueBool()),
	}, nil
}

func convertAssignmentStructToSchema(d *interfaces.Assignment) (*assignmentResourceModel, error) {
	return &assignmentResourceModel{
		Device:      types.StringValue(d.Device.String()),
		Description: types.StringValue(d.Description),
		Lock:        types.BoolValue(tools.StringToBool(d.Lock)),
	}, nil
}
