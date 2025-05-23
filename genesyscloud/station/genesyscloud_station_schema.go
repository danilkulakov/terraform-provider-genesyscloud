package station

import (
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/provider"
	registrar "github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/resource_register"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const ResourceType = "genesyscloud_station"

// SetRegistrar registers all of the resources, datasources and exporters in the package
func SetRegistrar(l registrar.Registrar) {
	l.RegisterDataSource(ResourceType, DataSourceStation())
}

// DataSourceStation registers the genesyscloud_station data source
func DataSourceStation() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for Genesys Cloud Stations. Select a station by name.",
		ReadContext: provider.ReadWithPooledClient(dataSourceStationRead),
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Station name.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}
