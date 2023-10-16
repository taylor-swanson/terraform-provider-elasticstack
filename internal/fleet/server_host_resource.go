package fleet

import (
	"context"
	"github.com/elastic/terraform-provider-elasticstack/internal/clients"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	fleetapi "github.com/elastic/terraform-provider-elasticstack/generated/fleet"
	"github.com/elastic/terraform-provider-elasticstack/internal/clients/fleet"
)

type serverHostModel struct {
	ID      types.String   `tfsdk:"id"`
	Name    types.String   `tfsdk:"name"`
	Hosts   []types.String `tfsdk:"hosts"`
	Default types.Bool     `tfsdk:"default"`
}

type serverHostResource struct {
	client *fleet.Client
}

var (
	_ resource.Resource                = &serverHostResource{}
	_ resource.ResourceWithConfigure   = &serverHostResource{}
	_ resource.ResourceWithImportState = &serverHostResource{}
)

func NewServerHostResource() resource.Resource {
	return &serverHostResource{}
}

func (r *serverHostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, res)
}

func (r *serverHostResource) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	var err error
	apiClient := req.ProviderData.(*clients.ApiClient)
	r.client, err = apiClient.GetFleetClient()
	if err != nil {
		res.Diagnostics.AddError(
			"Provider setup error",
			"Unable to get Fleet client from provider: "+err.Error(),
		)
	}
}

func (r *serverHostResource) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_fleet_server_host"
}

func (r *serverHostResource) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "Unique identifier of the Fleet server host.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Fleet server host.",
			},
			"default": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Set as default.",
			},
			"hosts": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "A list of hosts.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
		},
	}
}

func (r *serverHostResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	// Retrieve values from plan.
	var plan serverHostModel
	diags := req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	// Generate API request from plan.
	apiReq := fleetapi.PostFleetServerHostsJSONRequestBody{
		Id:        plan.ID.ValueStringPointer(),
		IsDefault: plan.Default.ValueBoolPointer(),
		Name:      plan.Name.ValueString(),
	}
	for _, v := range plan.Hosts {
		apiReq.HostUrls = append(apiReq.HostUrls, v.ValueString())
	}

	// Create object via API.
	obj, err := fleet.CreateFleetServerHost(ctx, r.client, apiReq)
	if err != nil {
		res.Diagnostics.AddError(
			"Error creating Fleet Server Host",
			"Unable to create Fleet Server Host: "+err.Error(),
		)
		return
	}
	plan.ID = types.StringValue(obj.Id)

	// Apply values from new object.
	if obj.Name != nil {
		plan.Name = types.StringValue(*obj.Name)
	}
	plan.Default = types.BoolValue(obj.IsDefault)
	for _, v := range obj.HostUrls {
		plan.Hosts = append(plan.Hosts, types.StringValue(v))
	}

	// Set new state.
	diags = res.State.Set(ctx, plan)
	res.Diagnostics.Append(diags...)
}

func (r *serverHostResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	// Retrieve values from state.
	var state serverHostModel
	diags := req.State.Get(ctx, &state)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	// Get object via API.
	obj, err := fleet.ReadFleetServerHost(ctx, r.client, state.ID.ValueString())
	if err != nil {
		res.Diagnostics.AddError(
			"Error reading Fleet Server Host",
			"Unable to read Fleet Server Host: "+err.Error(),
		)
		return
	}

	// Set retrieved values.
	if obj.Name != nil {
		state.Name = types.StringValue(*obj.Name)
	}
	state.Default = types.BoolValue(obj.IsDefault)
	for _, v := range obj.HostUrls {
		state.Hosts = append(state.Hosts, types.StringValue(v))
	}

	// Set refreshed state.
	diags = res.State.Set(ctx, &state)
	res.Diagnostics.Append(diags...)
}

func (r *serverHostResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	// Retrieve values from plan.
	var plan serverHostModel
	diags := req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	// Generate API request from plan.
	apiReq := fleetapi.UpdateFleetServerHostsJSONRequestBody{
		Name:      plan.Name.ValueStringPointer(),
		IsDefault: plan.Default.ValueBoolPointer(),
	}
	var hosts []string
	for _, v := range plan.Hosts {
		hosts = append(hosts, v.ValueString())
	}
	if hosts != nil {
		apiReq.HostUrls = &hosts
	}

	// Update object via API.
	obj, err := fleet.UpdateFleetServerHost(ctx, r.client, plan.ID.ValueString(), apiReq)
	if err != nil {
		res.Diagnostics.AddError(
			"Error updating Fleet Server Host",
			"Unable to update Fleet Server Host: "+err.Error(),
		)
		return
	}

	// Set retrieved values.
	if obj.Name != nil {
		plan.Name = types.StringValue(*obj.Name)
	}
	plan.Default = types.BoolValue(obj.IsDefault)
	for _, v := range obj.HostUrls {
		plan.Hosts = append(plan.Hosts, types.StringValue(v))
	}

	// Set refreshed state.
	diags = res.State.Set(ctx, &plan)
	res.Diagnostics.Append(diags...)
}

func (r *serverHostResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	// Retrieve values from state.
	var state serverHostModel
	diags := req.State.Get(ctx, &state)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	// Delete object via API.
	err := fleet.DeleteFleetServerHost(ctx, r.client, state.ID.ValueString())
	if err != nil {
		res.Diagnostics.AddError(
			"Error deleting Fleet Server Host",
			"Unable to delete Fleet Server Host: "+err.Error(),
		)
	}
}
