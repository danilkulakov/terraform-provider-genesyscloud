package auth_role

import (
	"fmt"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/provider"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/util"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAuthRole(t *testing.T) {
	var (
		roleResourceLabel   = "auth-role"
		roleDataSourceLabel = "auth-role-data"
		roleName            = "Terraform Role-" + uuid.NewString()
		roleDesc            = "Terraform test role"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { util.TestAccPreCheck(t) },
		ProviderFactories: provider.GetProviderFactories(providerResources, providerDataSources),
		Steps: []resource.TestStep{
			{
				Config: GenerateAuthRoleResource(
					roleResourceLabel,
					roleName,
					roleDesc,
				) + GenerateAuthRoleDataSource(
					roleDataSourceLabel,
					"genesyscloud_auth_role."+roleResourceLabel+".name",
					"genesyscloud_auth_role."+roleResourceLabel,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.genesyscloud_auth_role."+roleDataSourceLabel, "id", "genesyscloud_auth_role."+roleResourceLabel, "id"),
				),
			},
		},
	})
}

func GenerateAuthRoleDataSource(
	resourceLabel string,
	name string,
	// Must explicitly use depends_on in terraform v0.13 when a data source references a resource
	// Fixed in v0.14 https://github.com/hashicorp/terraform/pull/26284
	dependsOnResource string) string {
	return fmt.Sprintf(`data "genesyscloud_auth_role" "%s" {
		name = %s
        depends_on=[%s]
	}
	`, resourceLabel, name, dependsOnResource)
}
