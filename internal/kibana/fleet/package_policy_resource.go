package fleet

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/terraform-provider-elasticstack/internal/clients"
	"github.com/elastic/terraform-provider-elasticstack/internal/clients/kibana"
	"github.com/elastic/terraform-provider-elasticstack/internal/clients/kibana/fleetapi"
	"github.com/elastic/terraform-provider-elasticstack/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourcePackagePolicy() *schema.Resource {
	packagePolicySchema := map[string]*schema.Schema{
		"id": {
			Description: "Internal identifier of the resource.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"package_policy_id": {
			Description: "Unique identifier of the Package Policy.",
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
		"agent_policy_id": {
			Description: "The identifier of the associated agent policy.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"enabled": {
			Description: "Flag indicating if the package policy is enabled.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"package": {
			Description: "Integration package to configure.",
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			MaxItems:    1,
			ForceNew:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Description: "Name of the package.",
						Type:        schema.TypeString,
						Required:    true,
					},
					"version": {
						Description: "Version of the package.",
						Type:        schema.TypeString,
						Required:    true,
					},
				},
			},
		},
		"vars_json": {
			Description:      "JSON-encoded string containing root level variables.",
			Type:             schema.TypeString,
			Optional:         true,
			ForceNew:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: utils.DiffJsonSuppress,
		},
		"input": {
			Description: "List of inputs to configure.",
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			ForceNew:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"policy_template": {
						Description: "Name of the policy template containing the input (see the integration's manifest.yml).",
						Type:        schema.TypeString,
						Required:    true,
					},
					"type": {
						Description: "Input type.",
						Type:        schema.TypeString,
						Required:    true,
					},
					"enabled": {
						Description: "Enable the input.",
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     true,
					},
					"vars_json": {
						Description:      "JSON-encoded string containing input level variables.",
						Type:             schema.TypeString,
						Optional:         true,
						Computed:         true,
						ValidateFunc:     validation.StringIsJSON,
						DiffSuppressFunc: utils.DiffJsonSuppress,
					},
					"stream": {
						Description: "Input level variables.",
						Type:        schema.TypeList,
						Required:    true,
						MinItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"data_stream": {
									Description: `Name of the data_stream within the integration (e.g. "log").`,
									Type:        schema.TypeString,
									Required:    true,
								},
								"enabled": {
									Description: "Enabled enable or disable that stream.",
									Type:        schema.TypeBool,
									Optional:    true,
									Computed:    true,
								},
								"vars_json": {
									Description:      "JSON-encoded string containing stream level variables.",
									Type:             schema.TypeString,
									Optional:         true,
									Computed:         true,
									ValidateFunc:     validation.StringIsJSON,
									DiffSuppressFunc: utils.DiffJsonSuppress,
								},
								"compiled_stream": {
									Description: "JSON-encoded string containing final configuration for the stream.",
									Type:        schema.TypeString,
									Computed:    true,
								},
								"vars": {
									Description: "Contains the stream level variables used by Fleet. This is merged merged with the defaults from the package manifests.",
									Type:        schema.TypeString,
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	}

	return &schema.Resource{
		Description: "Creates a Fleet Package Policy.",

		CreateContext: resourcePackagePolicyUpsert,
		ReadContext:   resourcePackagePolicyRead,
		UpdateContext: resourcePackagePolicyUpsert,
		DeleteContext: resourcePackagePolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: packagePolicySchema,
	}
}

func resourcePackagePolicyUpsert(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, diags := clients.NewApiClient(d, meta)
	if diags.HasError() {
		return diags
	}

	req := fleetapi.PackagePolicyRequest{
		Name:     d.Get("name").(string),
		PolicyId: d.Get("agent_policy_id").(string),
	}

	if value := d.Get("package_policy_id").(string); value != "" {
		req.Id = &value
	}
	if value := d.Get("description").(string); value != "" {
		req.Description = &value
	}
	if value := d.Get("namespace").(string); value != "" {
		req.Namespace = &value
	}
	if value := d.Get("package.0.name").(string); value != "" {
		req.Package.Name = value
	}
	if value := d.Get("package.0.version").(string); value != "" {
		req.Package.Version = value
	}
	if value := d.Get("agent_policy_id").(string); value != "" {
		req.PolicyId = value
	}
	if value, _ := d.Get("vars_json").(string); value != "" {
		if err := json.Unmarshal([]byte(value), &req.Vars); err != nil {
			return diag.FromErr(err)
		}
	}
	if value, _ := d.Get("input").([]interface{}); value != nil {
		inputsMap := make(map[string]fleetapi.PackagePolicyRequestInput, len(value))
		for i, input := range value {
			inputMap := input.(map[string]any)

			var policyTemplate, inputType string
			var reqInput fleetapi.PackagePolicyRequestInput
			if v, ok := inputMap["policy_template"].(string); ok {
				policyTemplate = v
			}
			if v, ok := inputMap["type"].(string); ok {
				inputType = v
			}
			if v, ok := inputMap["enabled"].(bool); ok {
				reqInput.Enabled = &v
			}
			if v, ok := inputMap["vars_json"].(string); ok && v != "" {
				if err := json.Unmarshal([]byte(v), &reqInput.Vars); err != nil {
					return diag.FromErr(fmt.Errorf("failed unmarshaling input.%d.vars_json: %w", i, err))
				}
			}

			streamList := inputMap["stream"].([]any)
			streams := make(map[string]fleetapi.PackagePolicyRequestInputStream, len(streamList))
			for j, stream := range streamList {
				streamMap := stream.(map[string]any)

				var stream fleetapi.PackagePolicyRequestInputStream
				var dataStream string
				if v, ok := streamMap["data_stream"].(string); ok {
					dataStream = v
				}
				if v, ok := streamMap["enabled"].(bool); ok {
					stream.Enabled = &v
				}
				if v, ok := streamMap["vars_json"].(string); ok && v != "" {
					if err := json.Unmarshal([]byte(v), &stream.Vars); err != nil {
						return diag.FromErr(fmt.Errorf("failed unmarshaling input.%d.stream.%d.vars_json: %w", i, j, err))
					}
				}

				streams[req.Package.Name+"."+dataStream] = stream
			}
			reqInput.Streams = &streams

			inputsMap[policyTemplate+"-"+inputType] = reqInput
		}

		req.Inputs = &inputsMap
	}

	var policyID string
	if d.IsNewResource() {
		policyID, diags = kibana.PackagePolicyCreate(ctx, client, &req)
		if err := d.Set("package_policy_id", policyID); err != nil {
			return diag.FromErr(err)
		}
	} else {
		policyID = d.Get("package_policy_id").(string)
		diags = kibana.PackagePolicyUpdate(ctx, client, policyID, &req)
	}
	if diags.HasError() {
		return diags
	}
	d.SetId(policyID)

	return resourcePackagePolicyRead(ctx, d, meta)
}

func resourcePackagePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, diags := clients.NewApiClient(d, meta)
	if diags.HasError() {
		return diags
	}

	id := d.Get("package_policy_id").(string)
	d.SetId(id)

	packagePolicy, diags := kibana.PackagePolicyRead(ctx, client, id)
	if diags.HasError() {
		return diags
	}

	if err := d.Set("name", packagePolicy.Name); err != nil {
		return diag.FromErr(err)
	}
	if packagePolicy.Namespace != nil {
		if err := d.Set("namespace", *packagePolicy.Namespace); err != nil {
			return diag.FromErr(err)
		}
	}
	if packagePolicy.Description != nil {
		if err := d.Set("description", *packagePolicy.Description); err != nil {
			return diag.FromErr(err)
		}
	}
	if packagePolicy.PolicyId != nil {
		if err := d.Set("agent_policy_id", *packagePolicy.PolicyId); err != nil {
			return diag.FromErr(err)
		}
	}
	if packagePolicy.Package != nil {
		if err := d.Set("package", []map[string]any{
			{
				"name":    packagePolicy.Package.Name,
				"version": packagePolicy.Package.Version,
			},
		}); err != nil {
			return diag.FromErr(err)
		}
	}

	var inputs []map[string]any
	for inputIdx, input := range packagePolicy.Inputs {
		if !input.Enabled {
			continue
		}

		var streams []map[string]any
		if input.Streams != nil {
			for streamIdx, stream := range *input.Streams {
				streamMap := map[string]any{}

				if stream.Enabled != nil {
					streamMap["enabled"] = *stream.Enabled
					if !*stream.Enabled {
						continue
					}
				}
				if stream.DataStream != nil && stream.DataStream.Dataset != nil {
					if parts := strings.SplitN(*stream.DataStream.Dataset, ".", 2); len(parts) == 2 {
						streamMap["data_stream"] = parts[1]
					}
				}
				if stream.Vars != nil {
					vars, err := json.Marshal(*stream.Vars)
					if err != nil {
						return diag.FromErr(err)
					}

					streamMap["vars"] = string(vars)

					flatVars := map[string]any{}
					for k, v := range *stream.Vars {
						obj := v.(map[string]any)
						if value, ok := obj["value"]; ok {
							flatVars[k] = value
						}
					}

					flatVarsJSON, err := json.Marshal(flatVars)
					if err != nil {
						return diag.FromErr(err)
					}

					if configured, ok := d.GetOk(fmt.Sprintf("input.%d.stream.%d.vars_json", inputIdx, streamIdx)); ok {
						flatVarsJSON, err = filterUnspecifiedVarsJSONKeys(configured.(string), string(flatVarsJSON))
						if err != nil {
							return diag.FromErr(err)
						}
					}
					streamMap["vars_json"] = string(flatVarsJSON)
				}
			}
		}

		inputs = append(inputs, map[string]any{
			"type":            input.Type,
			"enabled":         input.Enabled,
			"policy_template": input.PolicyTemplate,
			"stream":          streams,
		})
	}
	if err := d.Set("input", inputs); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourcePackagePolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, diags := clients.NewApiClient(d, meta)
	if diags.HasError() {
		return diags
	}

	id := d.Get("package_policy_id").(string)
	d.SetId(id)

	if diags = kibana.PackagePolicyDelete(ctx, client, id); diags.HasError() {
		return diags
	}

	d.SetId("")

	return diags
}

func filterUnspecifiedVarsJSONKeys(old, new string) ([]byte, error) {
	var oldMap, newMap map[string]any
	if err := json.Unmarshal([]byte(old), &oldMap); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(new), &newMap); err != nil {
		return nil, err
	}

	for k := range newMap {
		if _, found := oldMap[k]; !found {
			delete(newMap, k)
		}
	}

	jsonNewMap, err := json.Marshal(newMap)
	if err != nil {
		return nil, err
	}

	return jsonNewMap, nil
}
