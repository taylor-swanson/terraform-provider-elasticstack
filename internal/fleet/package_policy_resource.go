package fleet

import (
	"context"
	"encoding/json"

	"github.com/elastic/terraform-provider-elasticstack/internal/clients/fleet"
	"github.com/elastic/terraform-provider-elasticstack/internal/clients/fleet/fleetapi"
	"github.com/elastic/terraform-provider-elasticstack/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourcePackagePolicy() *schema.Resource {
	packagePolicySchema := map[string]*schema.Schema{
		"policy_id": {
			Description: "Unique identifier of the package policy.",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			ForceNew:    true,
		},
		"name": {
			Description: "The name of the package policy.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"namespace": {
			Description: "The namespace of the agent policy.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"agent_policy_id": {
			Description: "ID of the agent policy.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"description": {
			Description: "The description of the package policy.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"enabled": {
			Description: "Enable the package policy.",
			Type:        schema.TypeBool,
			Default:     true,
			Optional:    true,
		},
		"force": {
			Description: "Force operations, such as creation and deletion, to occur.",
			Type:        schema.TypeBool,
			Optional:    true,
		},
		"package_name": {
			Description: "The name of the package.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"package_version": {
			Description: "The version of the package.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"input": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"input_id": {
						Description: "The identifier of the input.",
						Type:        schema.TypeString,
						Required:    true,
					},
					"type": {
						Description: "The type of the input.",
						Type:        schema.TypeString,
						Required:    true,
					},
					"enabled": {
						Description: "Enable the input.",
						Type:        schema.TypeBool,
						Default:     true,
						Optional:    true,
					},
					"streams_json": {
						Description:      "Input streams as JSON.",
						Type:             schema.TypeString,
						ValidateFunc:     validation.StringIsJSON,
						DiffSuppressFunc: utils.DiffJsonSuppress,
						Optional:         true,
						Sensitive:        true,
					},
					"vars_json": {
						Description:      "Input variables as JSON.",
						Type:             schema.TypeString,
						ValidateFunc:     validation.StringIsJSON,
						DiffSuppressFunc: utils.DiffJsonSuppress,
						Optional:         true,
						Sensitive:        true,
					},
					"config_json": {
						Description: "Input configuration as JSON.",
						Type:        schema.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"processors_json": {
						Description: "Input processors as JSON.",
						Type:        schema.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
				},
			},
		},
	}

	return &schema.Resource{
		Description: "Creates a new Fleet Package Policy. See https://www.elastic.co/guide/en/fleet/current/agent-policy.html",

		CreateContext: resourcePackagePolicyCreate,
		ReadContext:   resourcePackagePolicyRead,
		UpdateContext: resourcePackagePolicyUpdate,
		DeleteContext: resourcePackagePolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: packagePolicySchema,
	}
}

func resourcePackagePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	fleetClient, diags := getFleetClient(d, meta)
	if diags.HasError() {
		return diags
	}

	req := fleetapi.CreatePackagePolicyJSONRequestBody{
		PolicyId: d.Get("agent_policy_id").(string),
		Name:     d.Get("name").(string),
		Inputs:   nil,
		Vars:     nil,
	}
	req.Package.Name = d.Get("package_name").(string)
	req.Package.Version = d.Get("package_version").(string)

	if value := d.Get("policy_id").(string); value != "" {
		req.Id = &value
	}
	if value := d.Get("namespace").(string); value != "" {
		req.Namespace = &value
	}
	if value := d.Get("description").(string); value != "" {
		req.Description = &value
	}
	if value := d.Get("force").(bool); value {
		req.Force = &value
	}

	if values := d.Get("input").([]interface{}); len(values) > 0 {
		inputMap := map[string]fleetapi.PackagePolicyRequestInput{}

		for _, v := range values {
			var input fleetapi.PackagePolicyRequestInput

			inputData := v.(map[string]interface{})
			inputID := inputData["input_id"].(string)

			enabled, _ := inputData["enabled"].(bool)
			input.Enabled = &enabled

			if streamsRaw, _ := inputData["streams_json"].(string); streamsRaw != "" {
				streams := map[string]fleetapi.PackagePolicyRequestInputStream{}
				if err := json.Unmarshal([]byte(streamsRaw), &streams); err != nil {
					panic(err)
				}
				input.Streams = &streams
			}
			if varsRaw, _ := inputData["vars_json"].(string); varsRaw != "" {
				vars := map[string]interface{}{}
				if err := json.Unmarshal([]byte(varsRaw), &vars); err != nil {
					panic(err)
				}
				input.Vars = &vars
			}

			inputMap[inputID] = input
		}

		req.Inputs = &inputMap
	}

	policy, diags := fleet.CreatePackagePolicy(ctx, fleetClient, req)
	if diags.HasError() {
		return diags
	}

	d.SetId(policy.Id)
	if err := d.Set("policy_id", policy.Id); err != nil {
		return diag.FromErr(err)
	}

	return resourcePackagePolicyRead(ctx, d, meta)
}

func resourcePackagePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	fleetClient, diags := getFleetClient(d, meta)
	if diags.HasError() {
		return diags
	}

	id := d.Get("policy_id").(string)
	d.SetId(id)

	req := fleetapi.UpdatePackagePolicyJSONRequestBody{
		PolicyId: d.Get("agent_policy_id").(string),
		Name:     d.Get("name").(string),
		Inputs:   nil,
		Vars:     nil,
	}
	req.Package.Name = d.Get("package_name").(string)
	req.Package.Version = d.Get("package_version").(string)

	if value := d.Get("policy_id").(string); value != "" {
		req.Id = &value
	}
	if value := d.Get("namespace").(string); value != "" {
		req.Namespace = &value
	}
	if value := d.Get("description").(string); value != "" {
		req.Description = &value
	}
	if value := d.Get("force").(bool); value {
		req.Force = &value
	}

	if values := d.Get("input").([]interface{}); len(values) > 0 {
		inputMap := map[string]fleetapi.PackagePolicyRequestInput{}

		for _, v := range values {
			var input fleetapi.PackagePolicyRequestInput

			inputData := v.(map[string]interface{})
			inputID := inputData["input_id"].(string)

			enabled, _ := inputData["enabled"].(bool)
			input.Enabled = &enabled

			if streamsRaw, _ := inputData["streams_json"].(string); streamsRaw != "" {
				streams := map[string]fleetapi.PackagePolicyRequestInputStream{}
				if err := json.Unmarshal([]byte(streamsRaw), &streams); err != nil {
					panic(err)
				}
				input.Streams = &streams
			}
			if varsRaw, _ := inputData["vars_json"].(string); varsRaw != "" {
				vars := map[string]interface{}{}
				if err := json.Unmarshal([]byte(varsRaw), &vars); err != nil {
					panic(err)
				}
				input.Vars = &vars
			}

			inputMap[inputID] = input
		}

		req.Inputs = &inputMap
	}

	_, diags = fleet.UpdatePackagePolicy(ctx, fleetClient, id, req)
	if diags.HasError() {
		return diags
	}

	return resourcePackagePolicyRead(ctx, d, meta)
}

func resourcePackagePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	fleetClient, diags := getFleetClient(d, meta)
	if diags.HasError() {
		return diags
	}

	id := d.Get("policy_id").(string)
	d.SetId(id)

	pkgPolicy, diags := fleet.ReadPackagePolicy(ctx, fleetClient, id)
	if diags.HasError() {
		return diags
	}

	// Not found.
	if pkgPolicy == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("name", pkgPolicy.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("namespace", pkgPolicy.Namespace); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("package_name", pkgPolicy.Package.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("package_version", pkgPolicy.Package.Version); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("agent_policy_id", pkgPolicy.PolicyId); err != nil {
		return diag.FromErr(err)
	}
	if pkgPolicy.Description != nil {
		if err := d.Set("description", *pkgPolicy.Description); err != nil {
			return diag.FromErr(err)
		}
	}

	// TODO: Inputs and Vars.

	return nil
}

func resourcePackagePolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	fleetClient, diags := getFleetClient(d, meta)
	if diags.HasError() {
		return diags
	}

	id := d.Get("policy_id").(string)
	d.SetId(id)

	force := d.Get("force").(bool)

	if diags = fleet.DeletePackagePolicy(ctx, fleetClient, id, force); diags.HasError() {
		return diags
	}
	d.SetId("")

	return diags
}
