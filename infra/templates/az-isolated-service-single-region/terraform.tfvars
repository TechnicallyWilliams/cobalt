# Note to developers: This file shows some examples that you may
# want to use in order to configure this template. It is your
# responsibility to choose the values that make sense for your application.
#
# Note: These values will impact the names of resources. If your deployment
# fails due to a resource name colision, consider using different values for
# the `name` variable.

authn_deployment_targets = [
  {
    app_name                 = "cobalt-frontend-api-1",
    repository               = "",
    dockerfile               = "",
    image_name               = "appsvcsample/echo-server-1",
    image_release_tag_prefix = "release"
  }
]

unauthn_deployment_targets = [
  {
    app_name                 = "cobalt-backend-api-1",
    repository               = "",
    dockerfile               = "",
    image_name               = "appsvcsample/echo-server-1",
    image_release_tag_prefix = "release"
  }
]

# Note: this is configured as such only to test IP Whitelists. This is a well
# known DNS address
resource_ip_whitelist   = ["1.1.1.1"]
ase_name                = "cobalt-static-ase"
name                    = "isolated-service"
ase_resource_group      = "cobalt-static-ase-rg"
ase_vnet_name           = "cobalt-static-ase-vnet"
resource_group_location = "eastus2"
