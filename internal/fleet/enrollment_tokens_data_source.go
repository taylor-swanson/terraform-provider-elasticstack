package fleet

import (
	"context"
	"github.com/elastic/terraform-provider-elasticstack/internal/clients"
	"github.com/hashicorp/terraform-plugin-framework/attr"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/elastic/terraform-provider-elasticstack/internal/clients/fleet"
	"github.com/elastic/terraform-provider-elasticstack/internal/utils"
)

type enrollmentTokenModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	PolicyID  types.String `tfsdk:"policy_id"`
	APIKey    types.String `tfsdk:"api_key"`
	APIKeyID  types.String `tfsdk:"api_key_id"`
	CreatedAt types.String `tfsdk:"created_at"`
	Active    types.Bool   `tfsdk:"active"`
}

type enrollmentTokensModel struct {
	Tokens   []enrollmentTokenModel `tfsdk:"tokens"`
	PolicyID types.String           `tfsdk:"policy_id"`
	ID       types.String           `tfsdk:"id"`
}

type enrollmentTokenDataSource struct {
	client *fleet.Client
}

var (
	_ datasource.DataSource              = &enrollmentTokenDataSource{}
	_ datasource.DataSourceWithConfigure = &enrollmentTokenDataSource{}
)

func NewEnrollmentTokenDataSource() datasource.DataSource {
	return &enrollmentTokenDataSource{}
}

func (d *enrollmentTokenDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, res *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	var err error
	apiClient := req.ProviderData.(*clients.ApiClient)
	d.client, err = apiClient.GetFleetClient()
	if err != nil {
		res.Diagnostics.AddError(
			"Provider setup error",
			"Unable to get Fleet client from provider: "+err.Error(),
		)
	}
}

func (d *enrollmentTokenDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, res *datasource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_fleet_enrollment_tokens"
}

func (d *enrollmentTokenDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, res *datasource.SchemaResponse) {
	res.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Unique identifier of this object",
			},
			"policy_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The identifier of the target agent policy. When provided, only the enrollment tokens associated with this agent policy will be selected. Omit this value to select all enrollment tokens.",
			},
			"tokens": schema.ListAttribute{
				Computed:    true,
				Description: "A list of enrollment tokens.",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"id":         types.StringType,
						"name":       types.StringType,
						"policy_id":  types.StringType,
						"api_key":    types.StringType,
						"api_key_id": types.StringType,
						"created_at": types.StringType,
						"active":     types.BoolType,
					},
				},
			},
		},
	}
}

func (d *enrollmentTokenDataSource) Read(ctx context.Context, req datasource.ReadRequest, res *datasource.ReadResponse) {
	// Retrieve values from state.
	var state enrollmentTokensModel
	diags := req.Config.Get(ctx, &state)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	// Get all tokens from API.
	allTokens, err := fleet.AllEnrollmentTokens(ctx, d.client)
	if err != nil {
		res.Diagnostics.AddError(
			"Error reading Fleet Enrollment Tokens",
			"Unable to read Fleet Enrollment Tokens: "+err.Error(),
		)
		return
	}

	// Set tokens in state.
	state.Tokens = nil
	policyID := state.PolicyID.ValueString()
	for _, v := range allTokens {
		if policyID != "" && v.PolicyId != nil && *v.PolicyId != policyID {
			continue
		}

		tokenState := enrollmentTokenModel{
			Active:    types.BoolValue(v.Active),
			APIKey:    types.StringValue(v.ApiKey),
			APIKeyID:  types.StringValue(v.ApiKeyId),
			CreatedAt: types.StringValue(v.CreatedAt),
			ID:        types.StringValue(v.Id),
			Name:      types.StringPointerValue(v.Name),
			PolicyID:  types.StringPointerValue(v.PolicyId),
		}

		state.Tokens = append(state.Tokens, tokenState)
	}

	// Compute ID for state. Uses the top-level policy ID, if provided,
	// otherwise computes and uses the SHA1 hash of the Fleet Client URL.
	if policyID != "" {
		state.ID = types.StringValue(policyID)
	} else {
		hash, err := utils.StringToHash(d.client.URL)
		if err != nil {
			res.Diagnostics.AddError(
				"Error reading Fleet Enrollment Tokens",
				"Unable to compute ID hash: "+err.Error(),
			)
		}
		state.ID = types.StringValue(*hash)
	}

	// Set state
	diags = res.State.Set(ctx, &state)
	res.Diagnostics.Append(diags...)
}
