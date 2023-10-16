package fleet

import (
	"context"
	fleetapi "github.com/elastic/terraform-provider-elasticstack/generated/fleet"
	"github.com/elastic/terraform-provider-elasticstack/internal/clients"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/elastic/terraform-provider-elasticstack/internal/clients/fleet"
)

const (
	monitorLogs    = "logs"
	monitorMetrics = "metrics"
)

type agentPolicyModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Namespace          types.String `tfsdk:"namespace"`
	Description        types.String `tfsdk:"description"`
	DataOutputID       types.String `tfsdk:"data_output_id"`
	MonitoringOutputID types.String `tfsdk:"monitoring_output_id"`
	FleetServerHostID  types.String `tfsdk:"fleet_server_host_id"`
	DownloadSourceID   types.String `tfsdk:"download_source_id"`
	SysMonitoring      types.Bool   `tfsdk:"sys_monitoring"`
	MonitorLogs        types.Bool   `tfsdk:"monitor_logs"`
	MonitorMetrics     types.Bool   `tfsdk:"monitor_metrics"`
	SkipDestroy        types.Bool   `tfsdk:"skip_destroy"`
}

type agentPolicyResource struct {
	client *fleet.Client
}

var (
	_ resource.Resource                = &agentPolicyResource{}
	_ resource.ResourceWithConfigure   = &agentPolicyResource{}
	_ resource.ResourceWithImportState = &agentPolicyResource{}
)

// NewAgentPolicyResource is a helper function to simplify the provider implementation.
func NewAgentPolicyResource() resource.Resource {
	return &agentPolicyResource{}
}

func (r *agentPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, res)

}

func (r *agentPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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

func (r *agentPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_fleet_agent_policy"
}

func (r *agentPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "Unique identifier of the agent policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the agent policy.",
			},
			"namespace": schema.StringAttribute{
				Required:    true,
				Description: "The namespace of the agent policy.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the agent policy.",
			},
			"data_output_id": schema.StringAttribute{
				Optional:    true,
				Description: "The identifier for the data output.",
			},
			"monitoring_output_id": schema.StringAttribute{
				Optional:    true,
				Description: "The identifier for monitoring output.",
			},
			"fleet_server_host_id": schema.StringAttribute{
				Optional:    true,
				Description: "The identifier for the Fleet server host.",
			},
			"download_source_id": schema.StringAttribute{
				Optional:    true,
				Description: "The identifier for the Elastic Agent binary download server.",
			},
			"sys_monitoring": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable collection of system logs and metrics.",
			},
			"monitor_logs": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable collection of agent logs.",
			},
			"monitor_metrics": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable collection of agent metrics.",
			},
			"skip_destroy": schema.BoolAttribute{
				Optional:    true,
				Description: "Set to true if you do not wish the agent policy to be deleted at destroy time, and instead just remove the agent policy from the Terraform state.",
			},
		},
	}
}

func (r *agentPolicyResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	// Retrieve values from plan.
	var plan agentPolicyModel
	diags := req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	// Generate API request from plan.
	apiReq := fleetapi.AgentPolicyCreateRequest{
		DataOutputId:       plan.DataOutputID.ValueStringPointer(),
		Description:        plan.Description.ValueStringPointer(),
		DownloadSourceId:   plan.DownloadSourceID.ValueStringPointer(),
		FleetServerHostId:  plan.FleetServerHostID.ValueStringPointer(),
		Id:                 plan.ID.ValueStringPointer(),
		MonitoringOutputId: plan.MonitoringOutputID.ValueStringPointer(),
		Name:               plan.Name.ValueString(),
		Namespace:          plan.Namespace.ValueString(),
	}
	monitoringValues := make([]fleetapi.AgentPolicyCreateRequestMonitoringEnabled, 0, 2)
	if plan.MonitorLogs.ValueBool() {
		monitoringValues = append(monitoringValues, monitorLogs)
	}
	if plan.MonitorMetrics.ValueBool() {
		monitoringValues = append(monitoringValues, monitorMetrics)
	}
	apiReq.MonitoringEnabled = &monitoringValues

	// Create object via API.
	obj, err := fleet.CreateAgentPolicy(ctx, r.client, apiReq)
	if err != nil {
		res.Diagnostics.AddError(
			"Error creating Fleet Agent Policy",
			"Unable to create Fleet Agent Policy: "+err.Error(),
		)
		return
	}
	plan.ID = types.StringValue(obj.Id)

	// Apply values from new object.
	plan.Name = types.StringValue(obj.Name)
	plan.Namespace = types.StringValue(obj.Namespace)
	plan.Description = types.StringPointerValue(obj.Description)
	plan.DataOutputID = types.StringPointerValue(obj.DataOutputId)
	plan.DownloadSourceID = types.StringPointerValue(obj.DownloadSourceId)
	plan.FleetServerHostID = types.StringPointerValue(obj.FleetServerHostId)
	plan.MonitoringOutputID = types.StringPointerValue(obj.MonitoringOutputId)

	plan.MonitorLogs = types.BoolValue(false)
	plan.MonitorMetrics = types.BoolValue(false)
	if obj.MonitoringEnabled != nil {
		for _, v := range *obj.MonitoringEnabled {
			switch v {
			case monitorLogs:
				plan.MonitorLogs = types.BoolValue(true)
			case monitorMetrics:
				plan.MonitorMetrics = types.BoolValue(true)
			}
		}
	}

	// Set new state.
	diags = res.State.Set(ctx, plan)
	res.Diagnostics.Append(diags...)
}

func (r *agentPolicyResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	// Retrieve values from state.
	var state agentPolicyModel
	diags := req.State.Get(ctx, &state)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	// Get object via API.
	obj, err := fleet.ReadAgentPolicy(ctx, r.client, state.ID.ValueString())
	if err != nil {
		res.Diagnostics.AddError(
			"Error reading Fleet Agent Policy",
			"Unable to read Fleet Agent Policy: "+err.Error(),
		)
		return
	}

	// Set retrieved values.
	state.Name = types.StringValue(obj.Name)
	state.Namespace = types.StringValue(obj.Namespace)
	state.Description = types.StringPointerValue(obj.Description)
	state.DataOutputID = types.StringPointerValue(obj.DataOutputId)
	state.DownloadSourceID = types.StringPointerValue(obj.DownloadSourceId)
	state.FleetServerHostID = types.StringPointerValue(obj.FleetServerHostId)
	state.MonitoringOutputID = types.StringPointerValue(obj.MonitoringOutputId)

	state.MonitorLogs = types.BoolValue(false)
	state.MonitorMetrics = types.BoolValue(false)
	if obj.MonitoringEnabled != nil {
		for _, v := range *obj.MonitoringEnabled {
			switch v {
			case monitorLogs:
				state.MonitorLogs = types.BoolValue(true)
			case monitorMetrics:
				state.MonitorMetrics = types.BoolValue(true)
			}
		}
	}

	// Set refreshed state.
	diags = res.State.Set(ctx, &state)
	res.Diagnostics.Append(diags...)
}

func (r *agentPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	// Retrieve values from plan.
	var plan agentPolicyModel
	diags := req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	// Generate API request from plan.
	apiReq := fleetapi.AgentPolicyUpdateRequest{
		DataOutputId:       plan.DataOutputID.ValueStringPointer(),
		Description:        plan.Description.ValueStringPointer(),
		DownloadSourceId:   plan.DownloadSourceID.ValueStringPointer(),
		FleetServerHostId:  plan.FleetServerHostID.ValueStringPointer(),
		MonitoringOutputId: plan.MonitoringOutputID.ValueStringPointer(),
		Name:               plan.Name.ValueString(),
		Namespace:          plan.Namespace.ValueString(),
	}
	monitoringValues := make([]fleetapi.AgentPolicyUpdateRequestMonitoringEnabled, 0, 2)
	if plan.MonitorLogs.ValueBool() {
		monitoringValues = append(monitoringValues, monitorLogs)
	}
	if plan.MonitorMetrics.ValueBool() {
		monitoringValues = append(monitoringValues, monitorMetrics)
	}
	apiReq.MonitoringEnabled = &monitoringValues

	// Create object via API.
	obj, err := fleet.UpdateAgentPolicy(ctx, r.client, plan.ID.ValueString(), apiReq)
	if err != nil {
		res.Diagnostics.AddError(
			"Error creating Fleet Agent Policy",
			"Unable to create Fleet Agent Policy: "+err.Error(),
		)
		return
	}
	plan.ID = types.StringValue(obj.Id)

	// Apply values from new object.
	plan.Name = types.StringValue(obj.Name)
	plan.Namespace = types.StringValue(obj.Namespace)
	plan.Description = types.StringPointerValue(obj.Description)
	plan.DataOutputID = types.StringPointerValue(obj.DataOutputId)
	plan.DownloadSourceID = types.StringPointerValue(obj.DownloadSourceId)
	plan.FleetServerHostID = types.StringPointerValue(obj.FleetServerHostId)
	plan.MonitoringOutputID = types.StringPointerValue(obj.MonitoringOutputId)

	plan.MonitorLogs = types.BoolValue(false)
	plan.MonitorMetrics = types.BoolValue(false)
	if obj.MonitoringEnabled != nil {
		for _, v := range *obj.MonitoringEnabled {
			switch v {
			case monitorLogs:
				plan.MonitorLogs = types.BoolValue(true)
			case monitorMetrics:
				plan.MonitorMetrics = types.BoolValue(true)
			}
		}
	}
}

func (r *agentPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	// Retrieve values from state.
	var state agentPolicyModel
	diags := req.State.Get(ctx, &state)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	// Delete object via API.
	err := fleet.DeleteAgentPolicy(ctx, r.client, state.ID.ValueString())
	if err != nil {
		res.Diagnostics.AddError(
			"Error deleting Fleet Agent Policy",
			"Unable to delete Fleet Agent Policy: "+err.Error(),
		)
	}
}
