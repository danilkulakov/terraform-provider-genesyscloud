package journey_segment

import (
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/provider"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/util"
	"testing"

	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/util/testrunner"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceJourneySegment(t *testing.T) {
	runDataJourneySegmentTestCase(t, "find_by_name")
}

func runDataJourneySegmentTestCase(t *testing.T, testCaseName string) {
	testObjectName := testrunner.TestObjectIdPrefix + testCaseName
	testObjectFullName := ResourceType + "." + testObjectName
	setupJourneySegment(t, testCaseName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { util.TestAccPreCheck(t) },
		ProviderFactories: provider.GetProviderFactories(providerResources, providerDataSources),
		Steps: testrunner.GenerateDataJourneySourceTestSteps(ResourceType, testCaseName, []resource.TestCheckFunc{
			resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrPair("data."+testObjectFullName, "id", testObjectFullName, "id"),
				resource.TestCheckResourceAttr(testObjectFullName, "display_name", testObjectName+"_to_find"),
			),
		}),
	})
}
