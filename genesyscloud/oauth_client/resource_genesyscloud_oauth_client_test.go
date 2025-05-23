package oauth_client

import (
	"fmt"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/provider"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/util"
	"strconv"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceOAuthClient(t *testing.T) {
	var (
		clientResourceLabel1 = "test-client"
		clientName1          = "terraform1-" + uuid.NewString()
		clientName2          = "terraform2-" + uuid.NewString()
		clientDesc1          = "terraform test client1"
		clientDesc2          = "terraform test client2"
		tokenSec1            = "300"
		tokenSec2            = "172800"
		redirectURI1         = "https://example.com/auth1"
		redirectURI2         = "https://example.com/auth2"
		grantTypeClientCreds = "CLIENT-CREDENTIALS"
		grantTypeCode        = "CODE"
		scope1               = "oauth"
		stateActive          = "active"
		stateInactive        = "inactive"
		credentialName1      = "terraform3" + uuid.NewString()

		roleResourceLabel1 = "admin-role"
		roleName1          = "admin" // Must use a role already assigned to the TF OAuth client
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { util.TestAccPreCheck(t) },
		ProviderFactories: provider.GetProviderFactories(providerResources, providerDataSources),
		Steps: []resource.TestStep{
			{
				// Create client cred client with 1 role in default division
				Config: generateAuthRoleDataSource(
					roleResourceLabel1,
					strconv.Quote(roleName1),
					"",
				) + generateOauthClientWithCredential(
					clientResourceLabel1,
					clientName1,
					clientDesc1,
					grantTypeClientCreds,
					tokenSec1,
					util.NullValue, // Default state
					util.GenerateStringArray(strconv.Quote(redirectURI1)),
					util.NullValue, // No scopes for client creds
					credentialName1,
					generateOauthClientRoles("data.genesyscloud_auth_role."+roleResourceLabel1+".id", util.NullValue),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("genesyscloud_oauth_client."+clientResourceLabel1, "name", clientName1),
					resource.TestCheckResourceAttr("genesyscloud_oauth_client."+clientResourceLabel1, "description", clientDesc1),
					resource.TestCheckResourceAttr("genesyscloud_oauth_client."+clientResourceLabel1, "authorized_grant_type", grantTypeClientCreds),
					resource.TestCheckResourceAttr("genesyscloud_oauth_client."+clientResourceLabel1, "access_token_validity_seconds", tokenSec1),
					resource.TestCheckResourceAttr("genesyscloud_oauth_client."+clientResourceLabel1, "state", stateActive),
					resource.TestCheckResourceAttr("genesyscloud_oauth_client."+clientResourceLabel1, "integration_credential_name", credentialName1),
					resource.TestCheckNoResourceAttr("genesyscloud_oauth_client."+clientResourceLabel1, "scopes.%"),
					util.ValidateStringInArray("genesyscloud_oauth_client."+clientResourceLabel1, "registered_redirect_uris", redirectURI1),
					validateOauthRole("genesyscloud_oauth_client."+clientResourceLabel1, "data.genesyscloud_auth_role."+roleResourceLabel1, ""),
				),
			},
			{
				// Update client cred client attributes
				Config: generateAuthRoleDataSource(
					roleResourceLabel1,
					strconv.Quote(roleName1),
					"",
				) + generateOauthClient(
					clientResourceLabel1,
					clientName2,
					clientDesc2,
					grantTypeClientCreds,
					tokenSec2,
					strconv.Quote(stateInactive),
					util.GenerateStringArray(strconv.Quote(redirectURI2)),
					util.NullValue, // No scopes for client creds
					generateOauthClientRoles("data.genesyscloud_auth_role."+roleResourceLabel1+".id", util.NullValue),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("genesyscloud_oauth_client."+clientResourceLabel1, "name", clientName2),
					resource.TestCheckResourceAttr("genesyscloud_oauth_client."+clientResourceLabel1, "description", clientDesc2),
					resource.TestCheckResourceAttr("genesyscloud_oauth_client."+clientResourceLabel1, "authorized_grant_type", grantTypeClientCreds),
					resource.TestCheckResourceAttr("genesyscloud_oauth_client."+clientResourceLabel1, "access_token_validity_seconds", tokenSec2),
					resource.TestCheckResourceAttr("genesyscloud_oauth_client."+clientResourceLabel1, "state", stateInactive),
					resource.TestCheckNoResourceAttr("genesyscloud_oauth_client."+clientResourceLabel1, "scopes.%"),
					util.ValidateStringInArray("genesyscloud_oauth_client."+clientResourceLabel1, "registered_redirect_uris", redirectURI2),
				),
			},
			{
				// Change to a CODE grant type with scopes instead of a role
				Config: generateOauthClient(
					clientResourceLabel1,
					clientName1,
					clientDesc1,
					grantTypeCode,
					tokenSec1,
					strconv.Quote(stateActive),
					util.GenerateStringArray(strconv.Quote(redirectURI1)),
					util.GenerateStringArray(strconv.Quote(scope1)),
					// No roles for CODE type
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("genesyscloud_oauth_client."+clientResourceLabel1, "name", clientName1),
					resource.TestCheckResourceAttr("genesyscloud_oauth_client."+clientResourceLabel1, "description", clientDesc1),
					resource.TestCheckResourceAttr("genesyscloud_oauth_client."+clientResourceLabel1, "authorized_grant_type", grantTypeCode),
					resource.TestCheckResourceAttr("genesyscloud_oauth_client."+clientResourceLabel1, "access_token_validity_seconds", tokenSec1),
					resource.TestCheckResourceAttr("genesyscloud_oauth_client."+clientResourceLabel1, "state", stateActive),
					resource.TestCheckNoResourceAttr("genesyscloud_oauth_client."+clientResourceLabel1, "roles.%"),
					util.ValidateStringInArray("genesyscloud_oauth_client."+clientResourceLabel1, "registered_redirect_uris", redirectURI1),
					util.ValidateStringInArray("genesyscloud_oauth_client."+clientResourceLabel1, "scopes", scope1),
				),
			},
			{
				// Import/Read
				ResourceName:            "genesyscloud_oauth_client." + clientResourceLabel1,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"integration_credential_id", "integration_credential_name", "client_id", "client_secret"},
			},
		},
	})
}

func generateOauthClient(resourceLabel, name, description, grantType, tokenSec, state, uris, scopes string, blocks ...string) string {
	return fmt.Sprintf(`resource "genesyscloud_oauth_client" "%s" {
		name = "%s"
		description = "%s"
        authorized_grant_type = "%s"
        access_token_validity_seconds = %s
        state = %s
        registered_redirect_uris = %s
        scopes = %s
        %s
	}
	`, resourceLabel, name, description, grantType, tokenSec, state, uris, scopes, strings.Join(blocks, "\n"))
}

func generateOauthClientWithCredential(resourceLabel, name, description, grantType, tokenSec, state, uris, scopes string, credentialName string, blocks ...string) string {
	return fmt.Sprintf(`resource "genesyscloud_oauth_client" "%s" {
		name = "%s"
		description = "%s"
        authorized_grant_type = "%s"
        access_token_validity_seconds = %s
        state = %s
        registered_redirect_uris = %s
        scopes = %s
        integration_credential_name = "%s"
        %s
	}
	`, resourceLabel, name, description, grantType, tokenSec, state, uris, scopes, credentialName, strings.Join(blocks, "\n"))
}

func generateOauthClientRoles(roleID string, divisionId string) string {
	return fmt.Sprintf(`roles {
		role_id = %s
		division_id = %s
	}
	`, roleID, divisionId)
}

func validateOauthRole(resourcePath string, roleResourcePath string, division string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		resourceState, ok := state.RootModule().Resources[resourcePath]
		if !ok {
			return fmt.Errorf("Failed to find %s in state", resourcePath)
		}
		resourceLabel := resourceState.Primary.ID

		roleResource, ok := state.RootModule().Resources[roleResourcePath]
		if !ok {
			return fmt.Errorf("Failed to find role %s in state", roleResourcePath)
		}
		roleID := roleResource.Primary.ID

		if division == "" {
			// If no division specified, role should be in the home division
			homeDiv, err := util.GetHomeDivisionID()
			if err != nil {
				return fmt.Errorf("Failed to query home div: %v", err)
			}
			division = homeDiv
		} else if division != "*" {
			// Get the division ID from state
			divResource, ok := state.RootModule().Resources[division]
			if !ok {
				return fmt.Errorf("Failed to find %s in state", division)
			}
			division = divResource.Primary.ID
		}

		resourceAttrs := resourceState.Primary.Attributes
		numRolesAttr := resourceAttrs["roles.#"]
		numRoles, _ := strconv.Atoi(numRolesAttr)
		for i := 0; i < numRoles; i++ {
			if resourceAttrs["roles."+strconv.Itoa(i)+".role_id"] == roleID {
				divId := resourceAttrs["roles."+strconv.Itoa(i)+".division_id"]
				if divId == division {
					// Found expected role and division
					return nil
				}
			}
		}
		return fmt.Errorf("Missing expected role/division for oauth client %s in state: %s/%s", resourceLabel, roleID, division)
	}
}
