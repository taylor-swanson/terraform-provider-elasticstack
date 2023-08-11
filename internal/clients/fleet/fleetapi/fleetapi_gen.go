// Package fleetapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.4 DO NOT EDIT.
package fleetapi

import (
	"encoding/json"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
)

const (
	BasicAuthScopes = "basicAuth.Scopes"
)

// Defines values for AgentPolicyMonitoringEnabled.
const (
	AgentPolicyMonitoringEnabledLogs    AgentPolicyMonitoringEnabled = "logs"
	AgentPolicyMonitoringEnabledMetrics AgentPolicyMonitoringEnabled = "metrics"
)

// Defines values for AgentPolicyCreateRequestMonitoringEnabled.
const (
	AgentPolicyCreateRequestMonitoringEnabledLogs    AgentPolicyCreateRequestMonitoringEnabled = "logs"
	AgentPolicyCreateRequestMonitoringEnabledMetrics AgentPolicyCreateRequestMonitoringEnabled = "metrics"
)

// Defines values for AgentPolicyUpdateRequestMonitoringEnabled.
const (
	Logs    AgentPolicyUpdateRequestMonitoringEnabled = "logs"
	Metrics AgentPolicyUpdateRequestMonitoringEnabled = "metrics"
)

// Defines values for ElasticsearchAssetType.
const (
	ComponentTemplate   ElasticsearchAssetType = "component_template"
	DataStreamIlmPolicy ElasticsearchAssetType = "data_stream_ilm_policy"
	IlmPolicy           ElasticsearchAssetType = "ilm_policy"
	IndexTemplate       ElasticsearchAssetType = "index_template"
	IngestPipeline      ElasticsearchAssetType = "ingest_pipeline"
	Transform           ElasticsearchAssetType = "transform"
)

// Defines values for KibanaSavedObjectType.
const (
	CspRuleTemplate KibanaSavedObjectType = "csp_rule_template"
	Dashboard       KibanaSavedObjectType = "dashboard"
	IndexPattern    KibanaSavedObjectType = "index-pattern"
	Lens            KibanaSavedObjectType = "lens"
	Map             KibanaSavedObjectType = "map"
	MlModule        KibanaSavedObjectType = "ml-module"
	Search          KibanaSavedObjectType = "search"
	SecurityRule    KibanaSavedObjectType = "security-rule"
	Visualization   KibanaSavedObjectType = "visualization"
)

// Defines values for OutputType.
const (
	OutputTypeElasticsearch OutputType = "elasticsearch"
	OutputTypeLogstash      OutputType = "logstash"
)

// Defines values for PackageInfoConditionsElasticsearchSubscription.
const (
	Basic      PackageInfoConditionsElasticsearchSubscription = "basic"
	Enterprise PackageInfoConditionsElasticsearchSubscription = "enterprise"
	Gold       PackageInfoConditionsElasticsearchSubscription = "gold"
	Platinum   PackageInfoConditionsElasticsearchSubscription = "platinum"
)

// Defines values for PackageInfoRelease.
const (
	Beta         PackageInfoRelease = "beta"
	Experimental PackageInfoRelease = "experimental"
	Ga           PackageInfoRelease = "ga"
)

// Defines values for PackageInfoSourceLicense.
const (
	Apache20  PackageInfoSourceLicense = "Apache-2.0"
	Elastic20 PackageInfoSourceLicense = "Elastic-2.0"
)

// Defines values for PackageInstallSource.
const (
	Bundled  PackageInstallSource = "bundled"
	Registry PackageInstallSource = "registry"
	Upload   PackageInstallSource = "upload"
)

// Defines values for PackageStatus.
const (
	InstallFailed PackageStatus = "install_failed"
	Installed     PackageStatus = "installed"
	Installing    PackageStatus = "installing"
	NotInstalled  PackageStatus = "not_installed"
)

// Defines values for PostOutputsJSONBodyType.
const (
	PostOutputsJSONBodyTypeElasticsearch PostOutputsJSONBodyType = "elasticsearch"
)

// Defines values for UpdateOutputJSONBodyType.
const (
	Elasticsearch UpdateOutputJSONBodyType = "elasticsearch"
)

// AgentPolicy defines model for agent_policy.
type AgentPolicy struct {
	AgentFeatures *[]struct {
		Enabled bool   `json:"enabled"`
		Name    string `json:"name"`
	} `json:"agent_features,omitempty"`
	Agents            *float32 `json:"agents,omitempty"`
	DataOutputId      *string  `json:"data_output_id"`
	Description       *string  `json:"description,omitempty"`
	DownloadSourceId  *string  `json:"download_source_id"`
	FleetServerHostId *string  `json:"fleet_server_host_id"`
	Id                string   `json:"id"`
	InactivityTimeout *float32 `json:"inactivity_timeout,omitempty"`

	// IsProtected Indicates whether the agent policy has tamper protection enabled. Default false.
	IsProtected        *bool                           `json:"is_protected,omitempty"`
	MonitoringEnabled  *[]AgentPolicyMonitoringEnabled `json:"monitoring_enabled,omitempty"`
	MonitoringOutputId *string                         `json:"monitoring_output_id"`
	Name               string                          `json:"name"`
	Namespace          string                          `json:"namespace"`

	// Overrides Override settings that are defined in the agent policy. Input settings cannot be overridden. The override option should be used only in unusual circumstances and not as a routine procedure.
	Overrides       *map[string]interface{} `json:"overrides"`
	Revision        *float32                `json:"revision,omitempty"`
	UnenrollTimeout *float32                `json:"unenroll_timeout,omitempty"`
	UpdatedBy       *string                 `json:"updated_by,omitempty"`
	UpdatedOn       *time.Time              `json:"updated_on,omitempty"`
}

// AgentPolicyMonitoringEnabled defines model for AgentPolicy.MonitoringEnabled.
type AgentPolicyMonitoringEnabled string

// AgentPolicyCreateRequest defines model for agent_policy_create_request.
type AgentPolicyCreateRequest struct {
	AgentFeatures *[]struct {
		Enabled bool   `json:"enabled"`
		Name    string `json:"name"`
	} `json:"agent_features,omitempty"`
	DataOutputId       *string                                      `json:"data_output_id"`
	Description        *string                                      `json:"description,omitempty"`
	DownloadSourceId   *string                                      `json:"download_source_id"`
	FleetServerHostId  *string                                      `json:"fleet_server_host_id"`
	Id                 *string                                      `json:"id,omitempty"`
	InactivityTimeout  *float32                                     `json:"inactivity_timeout,omitempty"`
	IsProtected        *bool                                        `json:"is_protected,omitempty"`
	MonitoringEnabled  *[]AgentPolicyCreateRequestMonitoringEnabled `json:"monitoring_enabled,omitempty"`
	MonitoringOutputId *string                                      `json:"monitoring_output_id"`
	Name               string                                       `json:"name"`
	Namespace          string                                       `json:"namespace"`
	UnenrollTimeout    *float32                                     `json:"unenroll_timeout,omitempty"`
}

// AgentPolicyCreateRequestMonitoringEnabled defines model for AgentPolicyCreateRequest.MonitoringEnabled.
type AgentPolicyCreateRequestMonitoringEnabled string

// AgentPolicyUpdateRequest defines model for agent_policy_update_request.
type AgentPolicyUpdateRequest struct {
	AgentFeatures *[]struct {
		Enabled bool   `json:"enabled"`
		Name    string `json:"name"`
	} `json:"agent_features,omitempty"`
	DataOutputId       *string                                      `json:"data_output_id"`
	Description        *string                                      `json:"description,omitempty"`
	DownloadSourceId   *string                                      `json:"download_source_id"`
	FleetServerHostId  *string                                      `json:"fleet_server_host_id"`
	InactivityTimeout  *float32                                     `json:"inactivity_timeout,omitempty"`
	IsProtected        *bool                                        `json:"is_protected,omitempty"`
	MonitoringEnabled  *[]AgentPolicyUpdateRequestMonitoringEnabled `json:"monitoring_enabled,omitempty"`
	MonitoringOutputId *string                                      `json:"monitoring_output_id"`
	Name               string                                       `json:"name"`
	Namespace          string                                       `json:"namespace"`
	UnenrollTimeout    *float32                                     `json:"unenroll_timeout,omitempty"`
}

// AgentPolicyUpdateRequestMonitoringEnabled defines model for AgentPolicyUpdateRequest.MonitoringEnabled.
type AgentPolicyUpdateRequestMonitoringEnabled string

// ElasticsearchAssetType defines model for elasticsearch_asset_type.
type ElasticsearchAssetType string

// EnrollmentApiKey defines model for enrollment_api_key.
type EnrollmentApiKey struct {
	Active    bool    `json:"active"`
	ApiKey    string  `json:"api_key"`
	ApiKeyId  string  `json:"api_key_id"`
	CreatedAt string  `json:"created_at"`
	Id        string  `json:"id"`
	Name      *string `json:"name,omitempty"`
	PolicyId  *string `json:"policy_id,omitempty"`
}

// FleetServerHost defines model for fleet_server_host.
type FleetServerHost struct {
	HostUrls        []string `json:"host_urls"`
	Id              string   `json:"id"`
	IsDefault       bool     `json:"is_default"`
	IsPreconfigured bool     `json:"is_preconfigured"`
	Name            *string  `json:"name,omitempty"`
}

// KibanaSavedObjectType defines model for kibana_saved_object_type.
type KibanaSavedObjectType string

// NewPackagePolicy defines model for new_package_policy.
type NewPackagePolicy struct {
	Description *string                       `json:"description,omitempty"`
	Enabled     *bool                         `json:"enabled,omitempty"`
	Inputs      map[string]PackagePolicyInput `json:"inputs"`
	Name        string                        `json:"name"`
	Namespace   *string                       `json:"namespace,omitempty"`
	// Deprecated:
	OutputId *string                   `json:"output_id,omitempty"`
	Package  *PackagePolicyPackageInfo `json:"package,omitempty"`
	PolicyId *string                   `json:"policy_id,omitempty"`
}

// Output defines model for output.
type Output struct {
	CaSha256             *string                 `json:"ca_sha256,omitempty"`
	CaTrustedFingerprint *string                 `json:"ca_trusted_fingerprint,omitempty"`
	Config               *map[string]interface{} `json:"config,omitempty"`
	ConfigYaml           *string                 `json:"config_yaml,omitempty"`
	Hosts                *[]string               `json:"hosts,omitempty"`
	Id                   string                  `json:"id"`
	IsDefault            bool                    `json:"is_default"`
	IsDefaultMonitoring  *bool                   `json:"is_default_monitoring,omitempty"`
	Name                 string                  `json:"name"`
	ProxyId              *string                 `json:"proxy_id,omitempty"`
	Shipper              *struct {
		CompressionLevel            *float32 `json:"compression_level,omitempty"`
		DiskQueueCompressionEnabled *bool    `json:"disk_queue_compression_enabled,omitempty"`
		DiskQueueEnabled            *bool    `json:"disk_queue_enabled,omitempty"`
		DiskQueueEncryptionEnabled  *bool    `json:"disk_queue_encryption_enabled,omitempty"`
		DiskQueueMaxSize            *float32 `json:"disk_queue_max_size,omitempty"`
		DiskQueuePath               *string  `json:"disk_queue_path,omitempty"`
		Loadbalance                 *bool    `json:"loadbalance,omitempty"`
	} `json:"shipper,omitempty"`
	Ssl *struct {
		Certificate            *string   `json:"certificate,omitempty"`
		CertificateAuthorities *[]string `json:"certificate_authorities,omitempty"`
		Key                    *string   `json:"key,omitempty"`
	} `json:"ssl,omitempty"`
	Type OutputType `json:"type"`
}

// OutputType defines model for Output.Type.
type OutputType string

// PackageInfo defines model for package_info.
type PackageInfo struct {
	Assets     []string `json:"assets"`
	Categories []string `json:"categories"`
	Conditions struct {
		Elasticsearch *struct {
			Subscription *PackageInfoConditionsElasticsearchSubscription `json:"subscription,omitempty"`
		} `json:"elasticsearch,omitempty"`
		Kibana *struct {
			Versions *string `json:"versions,omitempty"`
		} `json:"kibana,omitempty"`
	} `json:"conditions"`
	DataStreams *[]struct {
		IngesetPipeline string `json:"ingeset_pipeline"`
		Name            string `json:"name"`
		Package         string `json:"package"`
		Release         string `json:"release"`
		Title           string `json:"title"`
		Type            string `json:"type"`
		Vars            *[]struct {
			Default string `json:"default"`
			Name    string `json:"name"`
		} `json:"vars,omitempty"`
	} `json:"data_streams,omitempty"`
	Description   string `json:"description"`
	Download      string `json:"download"`
	Elasticsearch *struct {
		Privileges *struct {
			Cluster *[]string `json:"cluster,omitempty"`
		} `json:"privileges,omitempty"`
	} `json:"elasticsearch,omitempty"`
	FormatVersion string    `json:"format_version"`
	Icons         *[]string `json:"icons,omitempty"`
	Internal      *bool     `json:"internal,omitempty"`
	Name          string    `json:"name"`
	Path          string    `json:"path"`
	Readme        *string   `json:"readme,omitempty"`

	// Release release label is deprecated, derive from the version instead (packages follow semver)
	// Deprecated:
	Release     *PackageInfoRelease `json:"release,omitempty"`
	Screenshots *[]struct {
		Path  string  `json:"path"`
		Size  *string `json:"size,omitempty"`
		Src   string  `json:"src"`
		Title *string `json:"title,omitempty"`
		Type  *string `json:"type,omitempty"`
	} `json:"screenshots,omitempty"`
	Source *struct {
		License *PackageInfoSourceLicense `json:"license,omitempty"`
	} `json:"source,omitempty"`
	Title   string `json:"title"`
	Type    string `json:"type"`
	Version string `json:"version"`
}

// PackageInfoConditionsElasticsearchSubscription defines model for PackageInfo.Conditions.Elasticsearch.Subscription.
type PackageInfoConditionsElasticsearchSubscription string

// PackageInfoRelease release label is deprecated, derive from the version instead (packages follow semver)
type PackageInfoRelease string

// PackageInfoSourceLicense defines model for PackageInfo.Source.License.
type PackageInfoSourceLicense string

// PackageInstallSource defines model for package_install_source.
type PackageInstallSource string

// PackageItemType defines model for package_item_type.
type PackageItemType struct {
	union json.RawMessage
}

// PackagePolicy defines model for package_policy.
type PackagePolicy struct {
	Description *string                       `json:"description,omitempty"`
	Enabled     *bool                         `json:"enabled,omitempty"`
	Id          string                        `json:"id"`
	Inputs      map[string]PackagePolicyInput `json:"inputs"`
	Name        string                        `json:"name"`
	Namespace   *string                       `json:"namespace,omitempty"`
	// Deprecated:
	OutputId *string                   `json:"output_id,omitempty"`
	Package  *PackagePolicyPackageInfo `json:"package,omitempty"`
	PolicyId *string                   `json:"policy_id,omitempty"`
	Revision float32                   `json:"revision"`
}

// PackagePolicyInput defines model for package_policy_input.
type PackagePolicyInput struct {
	Config     *map[string]interface{} `json:"config,omitempty"`
	Enabled    bool                    `json:"enabled"`
	Processors *[]string               `json:"processors,omitempty"`
	Streams    *map[string]interface{} `json:"streams,omitempty"`
	Type       string                  `json:"type"`
	Vars       *map[string]interface{} `json:"vars,omitempty"`
}

// PackagePolicyPackageInfo defines model for package_policy_package_info.
type PackagePolicyPackageInfo struct {
	Name    string  `json:"name"`
	Title   *string `json:"title,omitempty"`
	Version string  `json:"version"`
}

// PackagePolicyRequest defines model for package_policy_request.
type PackagePolicyRequest struct {
	// Description Package policy description
	Description *string `json:"description,omitempty"`

	// Force Force package policy creation even if package is not verified, or if the agent policy is managed.
	Force *bool `json:"force,omitempty"`

	// Id Package policy unique identifier
	Id *string `json:"id,omitempty"`

	// Inputs Package policy inputs (see integration documentation to know what inputs are available)
	Inputs *map[string]PackagePolicyRequestInput `json:"inputs,omitempty"`

	// Name Package policy name (should be unique)
	Name string `json:"name"`

	// Namespace namespace by default "default"
	Namespace *string `json:"namespace,omitempty"`
	Package   struct {
		// Name Package name
		Name string `json:"name"`

		// Version Package version
		Version string `json:"version"`
	} `json:"package"`

	// PolicyId Agent policy ID where that package policy will be added
	PolicyId string `json:"policy_id"`

	// Vars Package root level variable (see integration documentation for more information)
	Vars *map[string]interface{} `json:"vars,omitempty"`
}

// PackagePolicyRequestInput defines model for package_policy_request_input.
type PackagePolicyRequestInput struct {
	// Enabled enable or disable that input, (default to true)
	Enabled *bool `json:"enabled,omitempty"`

	// Streams Input streams (see integration documentation to know what streams are available)
	Streams *map[string]PackagePolicyRequestInputStream `json:"streams,omitempty"`

	// Vars Input level variable (see integration documentation for more information)
	Vars *map[string]interface{} `json:"vars,omitempty"`
}

// PackagePolicyRequestInputStream defines model for package_policy_request_input_stream.
type PackagePolicyRequestInputStream struct {
	// Enabled enable or disable that stream, (default to true)
	Enabled *bool `json:"enabled,omitempty"`

	// Vars Stream level variable (see integration documentation for more information)
	Vars *map[string]interface{} `json:"vars,omitempty"`
}

// PackageStatus defines model for package_status.
type PackageStatus string

// Error defines model for error.
type Error struct {
	Error      *string  `json:"error,omitempty"`
	Message    *string  `json:"message,omitempty"`
	StatusCode *float32 `json:"statusCode,omitempty"`
}

// DeleteAgentPolicyJSONBody defines parameters for DeleteAgentPolicy.
type DeleteAgentPolicyJSONBody struct {
	AgentPolicyId string `json:"agentPolicyId"`
}

// DeletePackageJSONBody defines parameters for DeletePackage.
type DeletePackageJSONBody struct {
	Force *bool `json:"force,omitempty"`
}

// DeletePackageParams defines parameters for DeletePackage.
type DeletePackageParams struct {
	// IgnoreUnverified Ignore if the package is fails signature verification
	IgnoreUnverified *bool `form:"ignoreUnverified,omitempty" json:"ignoreUnverified,omitempty"`

	// Full Return all fields from the package manifest, not just those supported by the Elastic Package Registry
	Full *bool `form:"full,omitempty" json:"full,omitempty"`

	// Prerelease Whether to return prerelease versions of packages (e.g. beta, rc, preview)
	Prerelease *bool `form:"prerelease,omitempty" json:"prerelease,omitempty"`
}

// GetPackageParams defines parameters for GetPackage.
type GetPackageParams struct {
	// IgnoreUnverified Ignore if the package is fails signature verification
	IgnoreUnverified *bool `form:"ignoreUnverified,omitempty" json:"ignoreUnverified,omitempty"`

	// Full Return all fields from the package manifest, not just those supported by the Elastic Package Registry
	Full *bool `form:"full,omitempty" json:"full,omitempty"`

	// Prerelease Whether to return prerelease versions of packages (e.g. beta, rc, preview)
	Prerelease *bool `form:"prerelease,omitempty" json:"prerelease,omitempty"`
}

// InstallPackageJSONBody defines parameters for InstallPackage.
type InstallPackageJSONBody struct {
	Force             *bool `json:"force,omitempty"`
	IgnoreConstraints *bool `json:"ignore_constraints,omitempty"`
}

// InstallPackageParams defines parameters for InstallPackage.
type InstallPackageParams struct {
	// IgnoreUnverified Ignore if the package is fails signature verification
	IgnoreUnverified *bool `form:"ignoreUnverified,omitempty" json:"ignoreUnverified,omitempty"`

	// Full Return all fields from the package manifest, not just those supported by the Elastic Package Registry
	Full *bool `form:"full,omitempty" json:"full,omitempty"`

	// Prerelease Whether to return prerelease versions of packages (e.g. beta, rc, preview)
	Prerelease *bool `form:"prerelease,omitempty" json:"prerelease,omitempty"`
}

// UpdatePackageJSONBody defines parameters for UpdatePackage.
type UpdatePackageJSONBody struct {
	KeepPoliciesUpToDate *bool `json:"keepPoliciesUpToDate,omitempty"`
}

// UpdatePackageParams defines parameters for UpdatePackage.
type UpdatePackageParams struct {
	// IgnoreUnverified Ignore if the package is fails signature verification
	IgnoreUnverified *bool `form:"ignoreUnverified,omitempty" json:"ignoreUnverified,omitempty"`

	// Full Return all fields from the package manifest, not just those supported by the Elastic Package Registry
	Full *bool `form:"full,omitempty" json:"full,omitempty"`

	// Prerelease Whether to return prerelease versions of packages (e.g. beta, rc, preview)
	Prerelease *bool `form:"prerelease,omitempty" json:"prerelease,omitempty"`
}

// PostFleetServerHostsJSONBody defines parameters for PostFleetServerHosts.
type PostFleetServerHostsJSONBody struct {
	HostUrls  []string `json:"host_urls"`
	Id        *string  `json:"id,omitempty"`
	IsDefault *bool    `json:"is_default,omitempty"`
	Name      string   `json:"name"`
}

// UpdateFleetServerHostsJSONBody defines parameters for UpdateFleetServerHosts.
type UpdateFleetServerHostsJSONBody struct {
	HostUrls  *[]string `json:"host_urls,omitempty"`
	IsDefault *bool     `json:"is_default,omitempty"`
	Name      *string   `json:"name,omitempty"`
}

// PostOutputsJSONBody defines parameters for PostOutputs.
type PostOutputsJSONBody struct {
	CaSha256            *string                 `json:"ca_sha256,omitempty"`
	ConfigYaml          *string                 `json:"config_yaml,omitempty"`
	Hosts               *[]string               `json:"hosts,omitempty"`
	Id                  *string                 `json:"id,omitempty"`
	IsDefault           *bool                   `json:"is_default,omitempty"`
	IsDefaultMonitoring *bool                   `json:"is_default_monitoring,omitempty"`
	Name                string                  `json:"name"`
	Type                PostOutputsJSONBodyType `json:"type"`
}

// PostOutputsJSONBodyType defines parameters for PostOutputs.
type PostOutputsJSONBodyType string

// UpdateOutputJSONBody defines parameters for UpdateOutput.
type UpdateOutputJSONBody struct {
	CaSha256             *string                  `json:"ca_sha256,omitempty"`
	CaTrustedFingerprint *string                  `json:"ca_trusted_fingerprint,omitempty"`
	ConfigYaml           *string                  `json:"config_yaml,omitempty"`
	Hosts                *[]string                `json:"hosts,omitempty"`
	IsDefault            *bool                    `json:"is_default,omitempty"`
	IsDefaultMonitoring  *bool                    `json:"is_default_monitoring,omitempty"`
	Name                 string                   `json:"name"`
	Type                 UpdateOutputJSONBodyType `json:"type"`
}

// UpdateOutputJSONBodyType defines parameters for UpdateOutput.
type UpdateOutputJSONBodyType string

// DeletePackagePolicyParams defines parameters for DeletePackagePolicy.
type DeletePackagePolicyParams struct {
	Force *bool `form:"force,omitempty" json:"force,omitempty"`
}

// CreateAgentPolicyJSONRequestBody defines body for CreateAgentPolicy for application/json ContentType.
type CreateAgentPolicyJSONRequestBody = AgentPolicyCreateRequest

// DeleteAgentPolicyJSONRequestBody defines body for DeleteAgentPolicy for application/json ContentType.
type DeleteAgentPolicyJSONRequestBody DeleteAgentPolicyJSONBody

// UpdateAgentPolicyJSONRequestBody defines body for UpdateAgentPolicy for application/json ContentType.
type UpdateAgentPolicyJSONRequestBody = AgentPolicyUpdateRequest

// DeletePackageJSONRequestBody defines body for DeletePackage for application/json ContentType.
type DeletePackageJSONRequestBody DeletePackageJSONBody

// InstallPackageJSONRequestBody defines body for InstallPackage for application/json ContentType.
type InstallPackageJSONRequestBody InstallPackageJSONBody

// UpdatePackageJSONRequestBody defines body for UpdatePackage for application/json ContentType.
type UpdatePackageJSONRequestBody UpdatePackageJSONBody

// PostFleetServerHostsJSONRequestBody defines body for PostFleetServerHosts for application/json ContentType.
type PostFleetServerHostsJSONRequestBody PostFleetServerHostsJSONBody

// UpdateFleetServerHostsJSONRequestBody defines body for UpdateFleetServerHosts for application/json ContentType.
type UpdateFleetServerHostsJSONRequestBody UpdateFleetServerHostsJSONBody

// PostOutputsJSONRequestBody defines body for PostOutputs for application/json ContentType.
type PostOutputsJSONRequestBody PostOutputsJSONBody

// UpdateOutputJSONRequestBody defines body for UpdateOutput for application/json ContentType.
type UpdateOutputJSONRequestBody UpdateOutputJSONBody

// CreatePackagePolicyJSONRequestBody defines body for CreatePackagePolicy for application/json ContentType.
type CreatePackagePolicyJSONRequestBody = PackagePolicyRequest

// UpdatePackagePolicyJSONRequestBody defines body for UpdatePackagePolicy for application/json ContentType.
type UpdatePackagePolicyJSONRequestBody = PackagePolicyRequest

// AsKibanaSavedObjectType returns the union data inside the PackageItemType as a KibanaSavedObjectType
func (t PackageItemType) AsKibanaSavedObjectType() (KibanaSavedObjectType, error) {
	var body KibanaSavedObjectType
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromKibanaSavedObjectType overwrites any union data inside the PackageItemType as the provided KibanaSavedObjectType
func (t *PackageItemType) FromKibanaSavedObjectType(v KibanaSavedObjectType) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeKibanaSavedObjectType performs a merge with any union data inside the PackageItemType, using the provided KibanaSavedObjectType
func (t *PackageItemType) MergeKibanaSavedObjectType(v KibanaSavedObjectType) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

// AsElasticsearchAssetType returns the union data inside the PackageItemType as a ElasticsearchAssetType
func (t PackageItemType) AsElasticsearchAssetType() (ElasticsearchAssetType, error) {
	var body ElasticsearchAssetType
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromElasticsearchAssetType overwrites any union data inside the PackageItemType as the provided ElasticsearchAssetType
func (t *PackageItemType) FromElasticsearchAssetType(v ElasticsearchAssetType) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeElasticsearchAssetType performs a merge with any union data inside the PackageItemType, using the provided ElasticsearchAssetType
func (t *PackageItemType) MergeElasticsearchAssetType(v ElasticsearchAssetType) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

func (t PackageItemType) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *PackageItemType) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}
