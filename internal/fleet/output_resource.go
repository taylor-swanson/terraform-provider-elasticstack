package fleet

import (
	"context"
	"fmt"
	"github.com/elastic/terraform-provider-elasticstack/internal/clients"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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

const (
	outputTypeElasticsearch = "elasticsearch"
	outputTypeLogstash      = "logstash"
)

type outputModel struct {
	ID                   types.String   `tfsdk:"id"`
	Name                 types.String   `tfsdk:"name"`
	Type                 types.String   `tfsdk:"type"`
	Hosts                []types.String `tfsdk:"hosts"`
	CASHA256             types.String   `tfsdk:"ca_sha256"`
	CATrustedFingerprint types.String   `tfsdk:"ca_trusted_fingerprint"`
	DefaultIntegrations  types.Bool     `tfsdk:"default_integrations"`
	DefaultMonitoring    types.Bool     `tfsdk:"default_monitoring"`
	ConfigYAML           types.String   `tfsdk:"config_yaml"`
}

type outputResource struct {
	client *fleet.Client
}

var (
	_ resource.Resource                = &outputResource{}
	_ resource.ResourceWithConfigure   = &outputResource{}
	_ resource.ResourceWithImportState = &outputResource{}
)

func NewOutputResource() resource.Resource {
	return &outputResource{}
}

func (r *outputResource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, res)
}

func (r *outputResource) Configure(_ context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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

func (r *outputResource) Metadata(_ context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_fleet_output"
}

func (r *outputResource) Schema(_ context.Context, _ resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "Unique identifier of the output.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the output.",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The output type.",
				Validators: []validator.String{
					stringvalidator.OneOf(outputTypeElasticsearch, outputTypeLogstash),
				},
			},
			"hosts": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "A list of hosts.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"ca_sha256": schema.StringAttribute{
				Optional:    true,
				Description: "Fingerprint of the Elasticsearch CA certificate.",
			},
			"ca_trusted_fingerprint": schema.StringAttribute{
				Optional:    true,
				Description: "Fingerprint of the trusted CA.",
			},
			"default_integrations": schema.StringAttribute{
				Optional:    true,
				Description: "Make this output the default for agent integrations.",
			},
			"default_monitoring": schema.StringAttribute{
				Optional:    true,
				Description: "Make this output the default for agent monitoring.",
			},
			"config_yaml": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Advanced YAML configuration. YAML settings here will be added to the output section of each agent policy.",
			},
		},
	}
}

func (r *outputResource) createElasticsearch(ctx context.Context, res *resource.CreateResponse, plan *outputModel) {
	// Generate API request from plan.
	apiReqElasticsearch := fleetapi.OutputCreateRequestElasticsearch{
		Id:                   plan.ID.ValueStringPointer(),
		Name:                 plan.Name.ValueString(),
		Type:                 fleetapi.OutputCreateRequestElasticsearchTypeElasticsearch,
		CaSha256:             plan.CASHA256.ValueStringPointer(),
		CaTrustedFingerprint: plan.CATrustedFingerprint.ValueStringPointer(),
		IsDefault:            plan.DefaultIntegrations.ValueBoolPointer(),
		IsDefaultMonitoring:  plan.DefaultMonitoring.ValueBoolPointer(),
		ConfigYaml:           plan.ConfigYAML.ValueStringPointer(),
	}
	var hosts []string
	for _, v := range plan.Hosts {
		hosts = append(hosts, v.ValueString())
	}
	apiReqElasticsearch.Hosts = &hosts

	var apiReq fleetapi.PostOutputsJSONRequestBody
	if err := apiReq.FromOutputCreateRequestElasticsearch(apiReqElasticsearch); err != nil {
		res.Diagnostics.AddError(
			"Error creating Fleet Output",
			"Unable to create Fleet Output: "+err.Error(),
		)
		return
	}

	// Create object via API.
	rawObj, err := fleet.CreateOutput(ctx, r.client, apiReq)
	if err != nil {
		res.Diagnostics.AddError(
			"Error creating Fleet Output",
			"Unable to create Fleet Output: "+err.Error(),
		)
		return
	}
	obj, err := rawObj.AsOutputCreateRequestElasticsearch()
	if err != nil {
		res.Diagnostics.AddError(
			"Error creating Fleet Output",
			"Unable to create Fleet Output: "+err.Error(),
		)
		return
	}
	plan.ID = types.StringValue(*obj.Id)

	// Apply values from new object.
	r.readElasticsearch(&obj, plan)
}

func (r *outputResource) createLogstash(ctx context.Context, res *resource.CreateResponse, plan *outputModel) {
	// Generate API request from plan.
	apiReqLogstash := fleetapi.OutputCreateRequestLogstash{
		Id:                   plan.ID.ValueStringPointer(),
		Name:                 plan.Name.ValueString(),
		Type:                 fleetapi.OutputCreateRequestLogstashTypeLogstash,
		CaSha256:             plan.CASHA256.ValueStringPointer(),
		CaTrustedFingerprint: plan.CATrustedFingerprint.ValueStringPointer(),
		IsDefault:            plan.DefaultIntegrations.ValueBoolPointer(),
		IsDefaultMonitoring:  plan.DefaultMonitoring.ValueBoolPointer(),
		ConfigYaml:           plan.ConfigYAML.ValueStringPointer(),
	}
	for _, v := range plan.Hosts {
		apiReqLogstash.Hosts = append(apiReqLogstash.Hosts, v.ValueString())
	}

	var apiReq fleetapi.PostOutputsJSONRequestBody
	if err := apiReq.FromOutputCreateRequestLogstash(apiReqLogstash); err != nil {
		res.Diagnostics.AddError(
			"Error creating Fleet Output",
			"Unable to create Fleet Output: "+err.Error(),
		)
		return
	}

	// Create object via API.
	rawObj, err := fleet.CreateOutput(ctx, r.client, apiReq)
	if err != nil {
		res.Diagnostics.AddError(
			"Error creating Fleet Output",
			"Unable to create Fleet Output: "+err.Error(),
		)
		return
	}
	obj, err := rawObj.AsOutputCreateRequestLogstash()
	if err != nil {
		res.Diagnostics.AddError(
			"Error creating Fleet Output",
			"Unable to create Fleet Output: "+err.Error(),
		)
		return
	}
	plan.ID = types.StringValue(*obj.Id)

	// Apply values from new object.
	r.readLogstash(&obj, plan)
}

func (r *outputResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	// Retrieve values from plan.
	var plan outputModel
	diags := req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	outputType := plan.Type.ValueString()
	switch outputType {
	case outputTypeElasticsearch:
		r.createElasticsearch(ctx, res, &plan)
	case outputTypeLogstash:
		r.createLogstash(ctx, res, &plan)
	default:
		res.Diagnostics.AddError(
			"Error creating Fleet Output",
			fmt.Sprintf("Unsupported output type: %q", outputType),
		)
		return
	}

	// Set new state.
	diags = res.State.Set(ctx, plan)
	res.Diagnostics.Append(diags...)
}

func (r *outputResource) readElasticsearch(data *fleetapi.OutputCreateRequestElasticsearch, state *outputModel) {
	state.Name = types.StringValue(data.Name)
	if data.Hosts != nil {
		state.Hosts = nil
		for _, v := range *data.Hosts {
			state.Hosts = append(state.Hosts, types.StringValue(v))
		}
	}
	if data.IsDefault != nil {
		state.DefaultIntegrations = types.BoolValue(*data.IsDefault)
	}
	if data.IsDefaultMonitoring != nil {
		state.DefaultMonitoring = types.BoolValue(*data.IsDefaultMonitoring)
	}
	if data.CaSha256 != nil {
		state.CASHA256 = types.StringValue(*data.CaSha256)
	}
	if data.CaTrustedFingerprint != nil {
		state.CATrustedFingerprint = types.StringValue(*data.CaTrustedFingerprint)
	}
	if data.ConfigYaml != nil {
		state.ConfigYAML = types.StringValue(*data.ConfigYaml)
	}
}

func (r *outputResource) readLogstash(data *fleetapi.OutputCreateRequestLogstash, state *outputModel) {
	state.Name = types.StringValue(data.Name)
	state.Hosts = nil
	for _, v := range data.Hosts {
		state.Hosts = append(state.Hosts, types.StringValue(v))
	}
	if data.IsDefault != nil {
		state.DefaultIntegrations = types.BoolValue(*data.IsDefault)
	}
	if data.IsDefaultMonitoring != nil {
		state.DefaultMonitoring = types.BoolValue(*data.IsDefaultMonitoring)
	}
	if data.CaSha256 != nil {
		state.CASHA256 = types.StringValue(*data.CaSha256)
	}
	if data.CaTrustedFingerprint != nil {
		state.CATrustedFingerprint = types.StringValue(*data.CaTrustedFingerprint)
	}
	if data.ConfigYaml != nil {
		state.ConfigYAML = types.StringValue(*data.ConfigYaml)
	}
}

func (r *outputResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	// Retrieve values from state.
	var state outputModel
	diags := req.State.Get(ctx, &state)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	rawOutput, err := fleet.ReadOutput(ctx, r.client, state.ID.ValueString())
	if err != nil {
		res.Diagnostics.AddError(
			"Error reading Fleet Output",
			"Unable to update Fleet Output: "+err.Error(),
		)
	}

	outputType := state.Type.ValueString()
	switch outputType {
	case outputTypeElasticsearch:
		output, err := rawOutput.AsOutputCreateRequestElasticsearch()
		if err != nil {
			res.Diagnostics.AddError(
				"Error reading Fleet Output",
				"Unable to read Fleet Output: "+err.Error(),
			)
			return
		}

		r.readElasticsearch(&output, &state)
	case outputTypeLogstash:
		output, err := rawOutput.AsOutputCreateRequestLogstash()
		if err != nil {
			res.Diagnostics.AddError(
				"Error reading Fleet Output",
				"Unable to read Fleet Output: "+err.Error(),
			)
			return
		}

		r.readLogstash(&output, &state)
	default:
		res.Diagnostics.AddError(
			"Error reading Fleet Output",
			fmt.Sprintf("Unsupported output type: %q", outputType),
		)
		return
	}

	// Set refreshed state.
	diags = res.State.Set(ctx, &state)
	res.Diagnostics.Append(diags...)
}

func (r *outputResource) updateElasticsearch(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse, plan *outputModel) {
	// Generate API request from plan.
	apiReqElasticsearch := fleetapi.OutputUpdateRequestElasticsearch{
		CaSha256:             plan.CASHA256.ValueStringPointer(),
		CaTrustedFingerprint: plan.CATrustedFingerprint.ValueStringPointer(),
		ConfigYaml:           plan.ConfigYAML.ValueStringPointer(),
		IsDefault:            plan.DefaultIntegrations.ValueBoolPointer(),
		IsDefaultMonitoring:  plan.DefaultMonitoring.ValueBoolPointer(),
		Name:                 plan.Name.ValueString(),
		Type:                 fleetapi.OutputUpdateRequestElasticsearchTypeElasticsearch,
	}
	for _, v := range plan.Hosts {
		apiReqElasticsearch.Hosts = append(apiReqElasticsearch.Hosts, v.ValueString())
	}

	var apiReq fleetapi.OutputUpdateRequest
	if err := apiReq.FromOutputUpdateRequestElasticsearch(apiReqElasticsearch); err != nil {
		res.Diagnostics.AddError(
			"Error updating Fleet Output",
			"Unable to update Fleet Output: "+err.Error(),
		)
		return
	}

	// Create object via API.
	rawObj, err := fleet.UpdateOutput(ctx, r.client, plan.ID.ValueString(), apiReq)
	if err != nil {
		res.Diagnostics.AddError(
			"Error updating Fleet Output",
			"Unable to update Fleet Output: "+err.Error(),
		)
		return
	}
	obj, err := rawObj.AsOutputUpdateRequestElasticsearch()
	if err != nil {
		res.Diagnostics.AddError(
			"Error updating Fleet Output",
			"Unable to update Fleet Output: "+err.Error(),
		)
		return
	}

	// Apply values from updated object.
	plan.Name = types.StringValue(obj.Name)
	if obj.Hosts != nil {
		plan.Hosts = nil
		for _, v := range obj.Hosts {
			plan.Hosts = append(plan.Hosts, types.StringValue(v))
		}
	}
	if obj.IsDefault != nil {
		plan.DefaultIntegrations = types.BoolValue(*obj.IsDefault)
	}
	if obj.IsDefaultMonitoring != nil {
		plan.DefaultMonitoring = types.BoolValue(*obj.IsDefaultMonitoring)
	}
	if obj.CaSha256 != nil {
		plan.CASHA256 = types.StringValue(*obj.CaSha256)
	}
	if obj.CaTrustedFingerprint != nil {
		plan.CATrustedFingerprint = types.StringValue(*obj.CaTrustedFingerprint)
	}
	if obj.ConfigYaml != nil {
		plan.ConfigYAML = types.StringValue(*obj.ConfigYaml)
	}
}

func (r *outputResource) updateLogstash(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse, plan *outputModel) {
	// Generate API request from plan.
	apiReqLogstash := fleetapi.OutputUpdateRequestLogstash{
		CaSha256:             plan.CASHA256.ValueStringPointer(),
		CaTrustedFingerprint: plan.CATrustedFingerprint.ValueStringPointer(),
		ConfigYaml:           plan.ConfigYAML.ValueStringPointer(),
		IsDefault:            plan.DefaultIntegrations.ValueBoolPointer(),
		IsDefaultMonitoring:  plan.DefaultMonitoring.ValueBoolPointer(),
		Name:                 plan.Name.ValueString(),
		Type:                 fleetapi.OutputUpdateRequestLogstashTypeLogstash,
	}
	var hosts []string
	for _, v := range plan.Hosts {
		hosts = append(hosts, v.ValueString())
	}
	apiReqLogstash.Hosts = &hosts

	var apiReq fleetapi.OutputUpdateRequest
	if err := apiReq.FromOutputUpdateRequestLogstash(apiReqLogstash); err != nil {
		res.Diagnostics.AddError(
			"Error updating Fleet Output",
			"Unable to update Fleet Output: "+err.Error(),
		)
		return
	}

	// Create object via API.
	rawObj, err := fleet.UpdateOutput(ctx, r.client, plan.ID.ValueString(), apiReq)
	if err != nil {
		res.Diagnostics.AddError(
			"Error updating Fleet Output",
			"Unable to update Fleet Output: "+err.Error(),
		)
		return
	}
	obj, err := rawObj.AsOutputUpdateRequestLogstash()
	if err != nil {
		res.Diagnostics.AddError(
			"Error updating Fleet Output",
			"Unable to update Fleet Output: "+err.Error(),
		)
		return
	}

	// Apply values from updated object.
	plan.Name = types.StringValue(obj.Name)
	if obj.Hosts != nil {
		plan.Hosts = nil
		for _, v := range *obj.Hosts {
			plan.Hosts = append(plan.Hosts, types.StringValue(v))
		}
	}
	if obj.IsDefault != nil {
		plan.DefaultIntegrations = types.BoolValue(*obj.IsDefault)
	}
	if obj.IsDefaultMonitoring != nil {
		plan.DefaultMonitoring = types.BoolValue(*obj.IsDefaultMonitoring)
	}
	if obj.CaSha256 != nil {
		plan.CASHA256 = types.StringValue(*obj.CaSha256)
	}
	if obj.CaTrustedFingerprint != nil {
		plan.CATrustedFingerprint = types.StringValue(*obj.CaTrustedFingerprint)
	}
	if obj.ConfigYaml != nil {
		plan.ConfigYAML = types.StringValue(*obj.ConfigYaml)
	}
}

func (r *outputResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	// Retrieve values from plan.
	var plan outputModel
	diags := req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	outputType := plan.Type.ValueString()
	switch outputType {
	case outputTypeElasticsearch:
		r.updateElasticsearch(ctx, req, res, &plan)
	case outputTypeLogstash:
		r.updateLogstash(ctx, req, res, &plan)
	default:
		res.Diagnostics.AddError(
			"Error updating Fleet Output",
			fmt.Sprintf("Unsupported output type: %q", outputType),
		)
		return
	}

	// Set refreshed state.
	diags = res.State.Set(ctx, &plan)
	res.Diagnostics.Append(diags...)
}

func (r *outputResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	// Retrieve values from state.
	var state outputModel
	diags := req.State.Get(ctx, &state)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	// Delete object via API.
	err := fleet.DeleteOutput(ctx, r.client, state.ID.ValueString())
	if err != nil {
		res.Diagnostics.AddError(
			"Error deleting Fleet Output",
			"Unable to delete Fleet Output: "+err.Error(),
		)
	}
}
