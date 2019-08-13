package integration

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/microsoft/cobalt/test-harness/infratests"
	"github.com/microsoft/cobalt/test-harness/terratest-extensions/modules/azure"
	"os"
	"testing"
)

var region = "eastus2"
var workspace = ""
var aseName = "cobalt-static-ase"
var aseVnetName = "cobalt-static-ase-vnet"
var aseResourceGroup = "cobalt-static-ase-rg"
var adminSubscription = os.Getenv("ARM_SUBSCRIPTION_ID")

var unauthn_deploymentTargets = []map[string]string{
	map[string]string{
		"app_name":                 "cobalt-backend-api-1",
		"repository":               "https://github.com/erikschlegel/echo-server.git",
		"dockerfile":               "Dockerfile",
		"image_name":               "appsvcsample/echo-server-2",
		"image_release_tag_prefix": "release",
	},
}

var authn_deploymentTargets = []map[string]string{
	map[string]string{
		"app_name":                 "cobalt-frontend-api-1",
		"repository":               "https://github.com/erikschlegel/echo-server.git",
		"dockerfile":               "Dockerfile",
		"image_name":               "appsvcsample/echo-server-1",
		"image_release_tag_prefix": "release",
	},
}

var tfOptions = &terraform.Options{
	TerraformDir: "../../",
	Upgrade:      true,
	Vars: map[string]interface{}{
		"resource_group_location":    region,
		"ase_subscription_id":        adminSubscription,
		"ase_name":                   aseName,
		"ase_resource_group":         aseResourceGroup,
		"unauthn_deployment_targets": unauthn_deploymentTargets,
		"authn_deployment_targets":   authn_deploymentTargets,
	},
	BackendConfig: map[string]interface{}{
		"storage_account_name": os.Getenv("TF_VAR_remote_state_account"),
		"container_name":       os.Getenv("TF_VAR_remote_state_container"),
	},
}

func TestAzureSimple(t *testing.T) {
	workspace = terraform.RunTerraformCommand(t, tfOptions, "workspace", "show")

	// Note: creating an App Service Plan configured with an Isolated SKU can take > 1.5
	// hours. In order to prevent a very long test cycle this test uses a static environment
	// that lives beyond the lifetime of this test. This is achieved using the
	// `SkipCleanupAfterTest` option.
	//
	// Be aware that this breaks testing isolation such that it is
	// possible that the environment is messed up by a test, which causes future tests
	// to fail for unrelated reasons.
	testFixture := infratests.IntegrationTestFixture{
		GoTest:                t,
		TfOptions:             tfOptions,
		ExpectedTfOutputCount: 9,
		ExpectedTfOutput: infratests.TerraformOutput{
			"fqdns": []string{
				"http://cobalt-frontend-api-1-" + workspace + "." + aseName + ".p.azurewebsites.net",
				"http://cobalt-backend-api-2-" + workspace + "." + aseName + ".p.azurewebsites.net",
			},
		},
		TfOutputAssertions: []infratests.TerraformOutputValidation{
			verifyVnetIntegrationForKeyVault,
			verifyVnetIntegrationForACR,
			verifyCDHooksConfiguredProperly,
			verifyCorrectWebhookEndpointForApps,
			verifyCorrectDeploymentTargetForApps,
		},
	}

	azure.CliServicePrincipalLogin(t)
	infratests.RunIntegrationTests(&testFixture)
}
