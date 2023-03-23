package kibana

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/elastic/terraform-provider-elasticstack/internal/clients"
	"github.com/elastic/terraform-provider-elasticstack/internal/clients/kibana/fleetapi"
	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func EnrollmentTokenGet(_ context.Context, apiClient *clients.ApiClient, id string) (*fleetapi.EnrollmentApiKey, diag.Diagnostics) {
	var diags diag.Diagnostics
	var enrollmentToken fleetapi.EnrollmentApiKey

	kibanaClient, err := apiClient.GetKibanaClient()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	res, err := kibanaClient.Client.NewRequest().Get(fmt.Sprintf("/api/fleet/enrollment_api_keys/%s", id))
	if err != nil {
		return nil, diag.FromErr(err)
	}

	if res.StatusCode() != http.StatusOK {
		if res.StatusCode() == http.StatusNotFound {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to find Fleet Enrollment Token",
				Detail:   fmt.Sprintf("Unable to find Fleet Enrollment Token with ID %q", id),
			})
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to find Fleet Enrollment Token",
				Detail:   fmt.Sprintf("Encountered HTTP status %d when retrieving Fleet Enrollment Token with ID %q", res.StatusCode(), id),
			})
		}

		return nil, diags
	}

	if err = json.Unmarshal(res.Body(), &enrollmentToken); err != nil {
		return nil, diag.FromErr(err)
	}

	return &enrollmentToken, diags
}

func AgentPolicyRead(_ context.Context, apiClient *clients.ApiClient, id string) (*fleetapi.AgentPolicy, diag.Diagnostics) {
	type responseData struct {
		Item fleetapi.AgentPolicy `json:"item"`
	}

	kibanaClient, err := apiClient.GetKibanaClient()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	res, err := kibanaClient.Client.NewRequest().Get("/api/fleet/agent_policies/" + id)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	var diags diag.Diagnostics
	if res.StatusCode() != http.StatusOK {
		if errMsg := getErrorMessage(res); errMsg != "" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to find Fleet Agent Policy",
				Detail:   fmt.Sprintf("Unable to find Fleet Agent Policy with ID %q: %s", id, errMsg),
			})
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to find Fleet Agent Policy",
				Detail:   fmt.Sprintf("Encountered HTTP status %d when retrieving Fleet Agent Policy with ID %q", res.StatusCode(), id),
			})
		}

		return nil, diags
	}

	var resInfo responseData
	if err = json.Unmarshal(res.Body(), &resInfo); err != nil {
		return nil, diag.FromErr(err)
	}

	return &resInfo.Item, diags
}

func AgentPolicyCreate(_ context.Context, apiClient *clients.ApiClient, req *fleetapi.AgentPolicyCreateRequest) (string, diag.Diagnostics) {
	type responseData struct {
		Item struct {
			Id string `json:"id"`
		} `json:"item"`
	}

	kibanaClient, err := apiClient.GetKibanaClient()
	if err != nil {
		return "", diag.FromErr(err)
	}

	reqData, err := json.Marshal(req)
	if err != nil {
		return "", diag.FromErr(err)
	}

	res, err := kibanaClient.Client.NewRequest().SetBody(reqData).Post("/api/fleet/agent_policies")
	if err != nil {
		return "", diag.FromErr(err)
	}

	var diags diag.Diagnostics
	if res.StatusCode() != http.StatusOK {
		if errMsg := getErrorMessage(res); errMsg != "" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create Fleet Agent Policy",
				Detail:   fmt.Sprintf("Unable to create Fleet Agent Policy: %s", errMsg),
			})
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to delete Fleet Agent Policy",
				Detail:   fmt.Sprintf("Encountered HTTP status %d when creating Fleet Agent Policy", res.StatusCode()),
			})
		}

		return "", diags
	}

	var resInfo responseData
	if err = json.Unmarshal(res.Body(), &resInfo); err != nil {
		return "", diag.FromErr(err)
	}

	return resInfo.Item.Id, diags
}

func AgentPolicyUpdate(_ context.Context, apiClient *clients.ApiClient, id string, req *fleetapi.AgentPolicyUpdateRequest) diag.Diagnostics {
	kibanaClient, err := apiClient.GetKibanaClient()
	if err != nil {
		return diag.FromErr(err)
	}

	reqData, err := json.Marshal(req)
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := kibanaClient.Client.NewRequest().SetBody(reqData).Put("/api/fleet/agent_policies/" + id)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics
	if res.StatusCode() != http.StatusOK {
		if errMsg := getErrorMessage(res); errMsg != "" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to update Fleet Agent Policy",
				Detail:   fmt.Sprintf("Unable to update Fleet Agent Policy with ID %q: %s", id, errMsg),
			})
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to update Fleet Agent Policy",
				Detail:   fmt.Sprintf("Encountered HTTP status %d when updating Fleet Agent Policy with ID %q", res.StatusCode(), id),
			})
		}

		return diags
	}

	return diags
}

func AgentPolicyDelete(_ context.Context, apiClient *clients.ApiClient, id string) diag.Diagnostics {
	kibanaClient, err := apiClient.GetKibanaClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := fleetapi.DeleteAgentPolicyJSONBody{
		AgentPolicyId: id,
	}

	reqData, err := json.Marshal(req)
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := kibanaClient.Client.NewRequest().SetBody(reqData).Post("/api/fleet/agent_policies/delete")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics
	if res.StatusCode() != http.StatusOK {
		if errMsg := getErrorMessage(res); errMsg != "" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to delete Fleet Agent Policy",
				Detail:   fmt.Sprintf("Unable to delete Fleet Agent Policy with ID %q: %s", id, errMsg),
			})
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to delete Fleet Agent Policy",
				Detail:   fmt.Sprintf("Encountered HTTP status %d when deleting Fleet Agent Policy with ID %q", res.StatusCode(), id),
			})
		}

		return diags
	}

	return diags
}

func PackagePolicyRead(_ context.Context, apiClient *clients.ApiClient, id string) (*fleetapi.PackagePolicy, diag.Diagnostics) {
	type responseData struct {
		Item fleetapi.PackagePolicy `json:"item"`
	}

	kibanaClient, err := apiClient.GetKibanaClient()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	res, err := kibanaClient.Client.NewRequest().Get("/api/fleet/package_policies/" + id)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	var diags diag.Diagnostics
	if res.StatusCode() != http.StatusOK {
		if errMsg := getErrorMessage(res); errMsg != "" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to get Fleet Package Policy",
				Detail:   fmt.Sprintf("Unable to get Fleet Package Policy with ID %q: %s", id, errMsg),
			})
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to get Fleet Package Policy",
				Detail:   fmt.Sprintf("Encountered HTTP status %d when getting Fleet Package Policy with ID %q", res.StatusCode(), id),
			})
		}

		return nil, diags
	}

	var resInfo responseData
	if err = json.Unmarshal(res.Body(), &resInfo); err != nil {
		return nil, diag.FromErr(err)
	}

	return &resInfo.Item, diags
}

func PackagePolicyCreate(_ context.Context, apiClient *clients.ApiClient, req *fleetapi.PackagePolicyRequest) (string, diag.Diagnostics) {
	type responseData struct {
		Item struct {
			Id string `json:"id"`
		} `json:"item"`
	}

	kibanaClient, err := apiClient.GetKibanaClient()
	if err != nil {
		return "", diag.FromErr(err)
	}

	reqData, err := json.Marshal(req)
	if err != nil {
		return "", diag.FromErr(err)
	}

	res, err := kibanaClient.Client.NewRequest().SetBody(reqData).Post("/api/fleet/package_policies")
	if err != nil {
		return "", diag.FromErr(err)
	}

	var diags diag.Diagnostics
	if res.StatusCode() != http.StatusOK {
		if errMsg := getErrorMessage(res); errMsg != "" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create Fleet Package Policy",
				Detail:   fmt.Sprintf("Unable to create Fleet Package Policy: %s", errMsg),
			})
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to delete Fleet Package Policy",
				Detail:   fmt.Sprintf("Encountered HTTP status %d when creating Fleet Package Policy", res.StatusCode()),
			})
		}

		return "", diags
	}

	var resInfo responseData
	if err = json.Unmarshal(res.Body(), &resInfo); err != nil {
		return "", diag.FromErr(err)
	}

	return resInfo.Item.Id, diags
}

func PackagePolicyUpdate(_ context.Context, apiClient *clients.ApiClient, id string, req *fleetapi.PackagePolicyRequest) diag.Diagnostics {
	kibanaClient, err := apiClient.GetKibanaClient()
	if err != nil {
		return diag.FromErr(err)
	}

	reqData, err := json.Marshal(req)
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := kibanaClient.Client.NewRequest().SetBody(reqData).Put("/api/fleet/package_policies/" + id)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics
	if res.StatusCode() != http.StatusOK {
		if errMsg := getErrorMessage(res); errMsg != "" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to update Fleet Package Policy",
				Detail:   fmt.Sprintf("Unable to update Fleet Package Policy with ID %q: %s", id, errMsg),
			})
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to delete Fleet Package Policy",
				Detail:   fmt.Sprintf("Encountered HTTP status %d when updating Fleet Package Policy with ID %q", res.StatusCode(), id),
			})
		}

		return diags
	}

	return diags
}

func PackagePolicyDelete(_ context.Context, apiClient *clients.ApiClient, id string) diag.Diagnostics {
	kibanaClient, err := apiClient.GetKibanaClient()
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := kibanaClient.Client.NewRequest().Delete("/api/fleet/package_policies/" + id)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics
	if res.StatusCode() != http.StatusOK {
		if errMsg := getErrorMessage(res); errMsg != "" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to delete Fleet Package Policy",
				Detail:   fmt.Sprintf("Unable to delete Fleet Package Policy with ID %q: %s", id, errMsg),
			})
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to delete Fleet Package Policy",
				Detail:   fmt.Sprintf("Encountered HTTP status %d when deleting Fleet Package Policy with ID %q", res.StatusCode(), id),
			})
		}

		return diags
	}

	return diags
}

func getErrorMessage(res *resty.Response) string {
	var errData fleetapi.Error
	if err := json.Unmarshal(res.Body(), &errData); err != nil {
		return ""
	}
	if errData.Error == nil && errData.StatusCode == nil && errData.Message == nil {
		return ""
	}

	var errStr string
	var errCode int
	var errMsg string

	if errData.Error != nil {
		errStr = *errData.Error
	}
	if errData.StatusCode != nil {
		errCode = int(*errData.StatusCode)
	}
	if errData.Message != nil {
		errMsg = *errData.Message
	}

	return fmt.Sprintf("%s (%d): %s", errStr, errCode, errMsg)
}
