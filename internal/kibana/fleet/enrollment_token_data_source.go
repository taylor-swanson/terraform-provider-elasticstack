package fleet

import (
	"context"

	"github.com/elastic/terraform-provider-elasticstack/internal/clients"
	"github.com/elastic/terraform-provider-elasticstack/internal/clients/kibana"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceEnrollmentToken() *schema.Resource {
	enrollmentTokenSchema := map[string]*schema.Schema{
		"id": {
			Description: "Internal identifier of the resource.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"key_id": {
			Description: "The unique identifier of the Enrollment Token.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"api_key": {
			Description: "The API key.",
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
		},
		"api_key_id": {
			Description: "The API key identifier.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"created_at": {
			Description: "The time at which the resource was created.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description: "The name of the resource.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"policy_id": {
			Description: "The identifier of the package policy.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"active": {
			Description: "Indicates if the resource is active.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
	}

	return &schema.Resource{
		Description: "Retrieves Elasticsearch API keys used to enroll Elastic Agents in Fleet. See: https://www.elastic.co/guide/en/fleet/current/fleet-enrollment-tokens.html",

		ReadContext: dataSourceEnrollmentTokenRead,

		Schema: enrollmentTokenSchema,
	}
}

func dataSourceEnrollmentTokenRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client, diags := clients.NewApiClient(d, meta)
	keyID := d.Get("key_id").(string)

	enrollmentToken, diags := kibana.EnrollmentTokenGet(ctx, client, keyID)
	if diags.HasError() {
		return diags
	}

	d.SetId(keyID)

	if err := d.Set("id", enrollmentToken.Id); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("api_key", enrollmentToken.ApiKey); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("api_key_id", enrollmentToken.ApiKeyId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_at", enrollmentToken.CreatedAt); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("active", enrollmentToken.Active); err != nil {
		return diag.FromErr(err)
	}

	if enrollmentToken.Name != nil {
		if err := d.Set("name", *enrollmentToken.Name); err != nil {
			return diag.FromErr(err)
		}
	}
	if enrollmentToken.PolicyId != nil {
		if err := d.Set("policy_id", *enrollmentToken.PolicyId); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}
