package task_management_worktype_status

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/provider"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/util"
)

/*
   The data_source_genesyscloud_task_management_worktype_status.go contains the data source implementation
   for the resource.
*/

// dataSourceTaskManagementWorktypeStatusRead retrieves by name the id in question
func dataSourceTaskManagementWorktypeStatusRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sdkConfig := meta.(*provider.ProviderMeta).ClientConfig
	proxy := getTaskManagementWorktypeStatusProxy(sdkConfig)

	worktypeId := d.Get("worktype_id").(string)
	name := d.Get("name").(string)

	return util.WithRetries(ctx, 15*time.Second, func() *retry.RetryError {
		worktypeStatusId, resp, retryable, err := proxy.getTaskManagementWorktypeStatusIdByName(ctx, worktypeId, name)

		if err != nil && !retryable {
			return retry.NonRetryableError(util.BuildWithRetriesApiDiagnosticError(ResourceType, fmt.Sprintf("Error searching task management worktype %s status %s | error: %s", worktypeId, name, err), resp))
		}

		if retryable {
			return retry.RetryableError(util.BuildWithRetriesApiDiagnosticError(ResourceType, fmt.Sprintf("No task management worktype %s status found with name %s", worktypeId, name), resp))
		}

		d.SetId(worktypeId + "/" + worktypeStatusId)
		return nil
	})
}
