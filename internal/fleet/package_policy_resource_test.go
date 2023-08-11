package fleet_test

import (
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/elastic/terraform-provider-elasticstack/internal/acctest"
)

var minVersionPackagePolicy = version.Must(version.NewVersion("8.10.0"))

const packagePolicyConfig = `
provider "elasticstack" {
  elasticsearch {}
  kibana {}
}

resource "elasticstack_fleet_package" "test_policy" {
  name    = "windows"
  version = "1.31.0"
  force   = true
}

resource "elasticstack_fleet_agent_policy" "test_policy" {
  name            = "PackagePolicyTest Agent Policy"
  namespace       = "default"
  description     = "PackagePolicyTest Agent Policy"
  monitor_logs    = true
  monitor_metrics = true
  skip_destroy    = false
}

data "elasticstack_fleet_enrollment_tokens" "test_policy" {
  policy_id = elasticstack_fleet_agent_policy.test_policy.policy_id
}

resource "elasticstack_fleet_package_policy" "test_policy" {
  name            = "PackagePolicyTest Policy"
  namespace       = "default"
  description     = "PackagePolicyTest Policy"
  agent_policy_id = elasticstack_fleet_agent_policy.test_policy.policy_id
  package_name    = elasticstack_fleet_package.test_policy.name
  package_version = elasticstack_fleet_package.test_policy.version

  input {
    input_id = "windows-winlog"
    type = "winlog"
	streams_json = jsonencode({
		"windows.applocker_exe_and_dll": {
		  "enabled": false,
		  "vars": {
			"preserve_original_event": false,
			"event_id": null,
			"ignore_older": "72h",
			"language": 0,
			"tags": []
		  }
		},
		"windows.applocker_msi_and_script": {
		  "enabled": false,
		  "vars": {
			"preserve_original_event": false,
			"event_id": null,
			"ignore_older": "72h",
			"language": 0,
			"tags": []
		  }
		}
	})
  }

  input {
    input_id = "windows-httpjson"
    type = "httpjson"
    vars_json = jsonencode({
        "url": "https://server.example.com:8089",
        "ssl": "#certificate_authorities:\n#  - |\n#    -----BEGIN CERTIFICATE-----\n#    MIIDCjCCAfKgAwIBAgITJ706Mu2wJlKckpIvkWxEHvEyijANBgkqhkiG9w0BAQsF\n#    ADAUMRIwEAYDVQQDDAlsb2NhbGhvc3QwIBcNMTkwNzIyMTkyOTA0WhgPMjExOTA2\n#    MjgxOTI5MDRaMBQxEjAQBgNVBAMMCWxvY2FsaG9zdDCCASIwDQYJKoZIhvcNAQEB\n#    BQADggEPADCCAQoCggEBANce58Y/JykI58iyOXpxGfw0/gMvF0hUQAcUrSMxEO6n\n#    fZRA49b4OV4SwWmA3395uL2eB2NB8y8qdQ9muXUdPBWE4l9rMZ6gmfu90N5B5uEl\n#    94NcfBfYOKi1fJQ9i7WKhTjlRkMCgBkWPkUokvBZFRt8RtF7zI77BSEorHGQCk9t\n#    /D7BS0GJyfVEhftbWcFEAG3VRcoMhF7kUzYwp+qESoriFRYLeDWv68ZOvG7eoWnP\n#    PsvZStEVEimjvK5NSESEQa9xWyJOmlOKXhkdymtcUd/nXnx6UTCFgnkgzSdTWV41\n#    CI6B6aJ9svCTI2QuoIq2HxX/ix7OvW1huVmcyHVxyUECAwEAAaNTMFEwHQYDVR0O\n#    BBYEFPwN1OceFGm9v6ux8G+DZ3TUDYxqMB8GA1UdIwQYMBaAFPwN1OceFGm9v6ux\n#    8G+DZ3TUDYxqMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEBAG5D\n#    874A4YI7YUwOVsVAdbWtgp1d0zKcPRR+r2OdSbTAV5/gcS3jgBJ3i1BN34JuDVFw\n#    3DeJSYT3nxy2Y56lLnxDeF8CUTUtVQx3CuGkRg1ouGAHpO/6OqOhwLLorEmxi7tA\n#    H2O8mtT0poX5AnOAhzVy7QW0D/k4WaoLyckM5hUa6RtvgvLxOwA0U+VGurCDoctu\n#    8F4QOgTAWyh8EZIwaKCliFRSynDpv3JTUwtfZkxo6K6nce1RhCWFAsMvDZL8Dgc0\n#    yvgJ38BRsFOtkRuAGSf6ZUwTO8JJRRIFnpUzXflAnGivK9M13D5GEQMmIl6U9Pvk\n#    sxSmbIUfc2SGJGCJD4I=\n#    -----END CERTIFICATE-----\n"
    })
    streams_json = jsonencode({
		"windows.applocker_exe_and_dll": {
          "enabled": true,
          "vars": {
            "interval": "10s",
            "search": "search sourcetype=\"XmlWinEventLog:Microsoft-Windows-AppLocker/EXE and DLL\"",
            "tags": [
              "forwarded"
            ],
            "preserve_original_event": false
          }
        },
        "windows.applocker_msi_and_script": {
          "enabled": true,
          "vars": {
            "interval": "10s",
            "search": "search sourcetype=\"XmlWinEventLog:Microsoft-Windows-AppLocker/MSI and Script\"",
            "tags": [
              "forwarded"
            ],
            "preserve_original_event": false
          }
        }
	})
  }
}
`

func TestAccResourcePackagePolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		CheckDestroy:             checkResourcePackagePolicyDestroy,
		ProtoV5ProviderFactories: acctest.Providers,
		Steps: []resource.TestStep{
			{
				//SkipFunc: versionutils.CheckIfVersionIsUnsupported(minVersionPackagePolicy),
				Config: packagePolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_fleet_package_policy.test_policy", "name", "PackagePolicyTest Policy"),
					resource.TestCheckResourceAttr("elasticstack_fleet_package_policy.test_policy", "package_name", "windows"),
					resource.TestCheckResourceAttr("elasticstack_fleet_package_policy.test_policy", "package_version", "1.31.0"),
				),
			},
		},
	})
}

func checkResourcePackagePolicyDestroy(s *terraform.State) error {
	return nil
}
