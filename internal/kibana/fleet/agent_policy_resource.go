package fleet

import (
	"context"

	"github.com/elastic/terraform-provider-elasticstack/internal/clients"
	"github.com/elastic/terraform-provider-elasticstack/internal/clients/kibana"
	"github.com/elastic/terraform-provider-elasticstack/internal/clients/kibana/fleetapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAgentPolicy() *schema.Resource {
	agentPolicySchema := map[string]*schema.Schema{
		"id": {
			Description: "Internal identifier of the resource.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"agent_policy_id": {
			Description: "Unique identifier of the Agent Policy.",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		"name": {
			Description: "The name of the resource.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"namespace": {
			Description: "The namespace of the resource.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"description": {
			Description: "The description of the resource.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"data_output_id": {
			Description: "The identifier for the output.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"monitoring_output_id": {
			Description: "The identifier for monitoring.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"fleet_server_host_id": {
			Description: "The identifier for the Fleet Server.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"download_source_id": {
			Description: "The identifier for the Download Server.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"sys_monitoring": {
			Description: "Enable system log gathering.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"monitoring_enabled": {
			Description: "Monitoring of Elastic Agent (logs, metrics).",
			Type:        schema.TypeSet,
			Optional:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}

	return &schema.Resource{
		Description: "Creates a new Fleet Agent Policy. See https://www.elastic.co/guide/en/fleet/current/agent-policy.html",

		CreateContext: resourceAgentPolicyCreate,
		ReadContext:   resourceAgentPolicyRead,
		UpdateContext: resourceAgentPolicyUpdate,
		DeleteContext: resourceAgentPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: agentPolicySchema,
	}
}

func resourceAgentPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, diags := clients.NewApiClient(d, meta)
	if diags.HasError() {
		return diags
	}

	req := fleetapi.AgentPolicyCreateRequest{
		Name:      d.Get("name").(string),
		Namespace: d.Get("namespace").(string),
	}

	if value := d.Get("agent_policy_id").(string); value != "" {
		req.Id = &value
	}
	if value := d.Get("description").(string); value != "" {
		req.Description = &value
	}
	if value := d.Get("data_output_id").(string); value != "" {
		req.DataOutputId = &value
	}
	if value := d.Get("download_source_id").(string); value != "" {
		req.DownloadSourceId = &value
	}
	if value := d.Get("fleet_server_host_id").(string); value != "" {
		req.FleetServerHostId = &value
	}
	if value := d.Get("monitoring_output_id").(string); value != "" {
		req.MonitoringOutputId = &value
	}
	if value := d.Get("monitoring_enabled").(*schema.Set); value != nil {
		var values []fleetapi.AgentPolicyCreateRequestMonitoringEnabled
		for _, v := range value.List() {
			values = append(values, fleetapi.AgentPolicyCreateRequestMonitoringEnabled(v.(string)))
		}
		req.MonitoringEnabled = &values
	}

	policyID, diags := kibana.AgentPolicyCreate(ctx, client, &req)
	if diags.HasError() {
		return diags
	}

	d.SetId(policyID)
	if err := d.Set("agent_policy_id", policyID); err != nil {
		return diag.FromErr(err)
	}

	return resourceAgentPolicyRead(ctx, d, meta)
}

func resourceAgentPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, diags := clients.NewApiClient(d, meta)
	if diags.HasError() {
		return diags
	}

	id := d.Get("agent_policy_id").(string)
	d.SetId(id)

	req := fleetapi.AgentPolicyUpdateRequest{
		Name:      d.Get("name").(string),
		Namespace: d.Get("namespace").(string),
	}

	if value := d.Get("description").(string); value != "" {
		req.Description = &value
	}
	if value := d.Get("data_output_id").(string); value != "" {
		req.DataOutputId = &value
	}
	if value := d.Get("download_source_id").(string); value != "" {
		req.DownloadSourceId = &value
	}
	if value := d.Get("fleet_server_host_id").(string); value != "" {
		req.FleetServerHostId = &value
	}
	if value := d.Get("monitoring_output_id").(string); value != "" {
		req.MonitoringOutputId = &value
	}
	if value := d.Get("monitoring_enabled").(*schema.Set); value != nil {
		var values []fleetapi.AgentPolicyUpdateRequestMonitoringEnabled
		for _, v := range value.List() {
			values = append(values, fleetapi.AgentPolicyUpdateRequestMonitoringEnabled(v.(string)))
		}
		req.MonitoringEnabled = &values
	}

	diags = kibana.AgentPolicyUpdate(ctx, client, id, &req)
	if diags.HasError() {
		return diags
	}

	return resourceAgentPolicyRead(ctx, d, meta)
}

func resourceAgentPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, diags := clients.NewApiClient(d, meta)
	if diags.HasError() {
		return diags
	}

	id := d.Get("agent_policy_id").(string)
	d.SetId(id)

	agentPolicy, diags := kibana.AgentPolicyRead(ctx, client, id)
	if diags.HasError() {
		return diags
	}

	if err := d.Set("name", agentPolicy.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("namespace", agentPolicy.Namespace); err != nil {
		return diag.FromErr(err)
	}
	if agentPolicy.Description != nil {
		if err := d.Set("description", *agentPolicy.Description); err != nil {
			return diag.FromErr(err)
		}
	}
	if agentPolicy.DataOutputId != nil {
		if err := d.Set("data_output_id", *agentPolicy.DataOutputId); err != nil {
			return diag.FromErr(err)
		}
	}
	if agentPolicy.DownloadSourceId != nil {
		if err := d.Set("download_source_id", *agentPolicy.DownloadSourceId); err != nil {
			return diag.FromErr(err)
		}
	}
	if agentPolicy.FleetServerHostId != nil {
		if err := d.Set("fleet_server_host_id", *agentPolicy.FleetServerHostId); err != nil {
			return diag.FromErr(err)
		}
	}
	if agentPolicy.MonitoringOutputId != nil {
		if err := d.Set("monitoring_output_id", *agentPolicy.MonitoringOutputId); err != nil {
			return diag.FromErr(err)
		}
	}
	if agentPolicy.MonitoringEnabled != nil {
		var values []string
		for _, v := range *agentPolicy.MonitoringEnabled {
			values = append(values, string(v))
		}
		if err := d.Set("monitoring_enabled", values); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceAgentPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, diags := clients.NewApiClient(d, meta)
	if diags.HasError() {
		return diags
	}

	id := d.Get("agent_policy_id").(string)
	d.SetId(id)

	if diags = kibana.AgentPolicyDelete(ctx, client, id); diags.HasError() {
		return diags
	}
	d.SetId("")

	return diags
}
