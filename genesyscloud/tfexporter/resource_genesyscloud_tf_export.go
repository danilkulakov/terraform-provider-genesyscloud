package tfexporter

import (
	"context"
	"fmt"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/validators"
	"os"
	"path/filepath"

	resourceExporter "github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/resource_exporter"

	registrar "github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/resource_register"

	tfExporterState "github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/tfexporter_state"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type fileMeta struct {
	Path  string
	IsDir bool
}

const ResourceType = "genesyscloud_tf_export"

func SetRegistrar(l registrar.Registrar) {
	l.RegisterResource(ResourceType, ResourceTfExport())
}

func ResourceTfExport() *schema.Resource {
	return &schema.Resource{
		Description: fmt.Sprintf(`
		Genesys Cloud Resource to export Terraform config and (optionally) tfstate files to a local directory.
		The config file is named '%s' or '%s', and the state file is named '%s'.
		`, defaultTfJSONFile, defaultTfHCLFile, defaultTfStateFile),

		CreateWithoutTimeout: createTfExport,
		ReadWithoutTimeout:   readTfExport,
		DeleteContext:        deleteTfExport,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"directory": {
				Description: "Directory where the config and state files will be exported.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "./genesyscloud",
				ForceNew:    true,
			},
			"resource_types": {
				Description: "Resource types to export, e.g. 'genesyscloud_user'. Defaults to all exportable types. NOTE: This field is deprecated and will be removed in future release.  Please use the include_filter_resources or exclude_filter_resources attribute.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validators.ValidateSubStringInSlice(resourceExporter.GetAvailableExporterTypes()),
				},
				ForceNew:      true,
				Deprecated:    "Use include_filter_resources attribute instead",
				ConflictsWith: []string{"include_filter_resources", "exclude_filter_resources"},
			},
			"include_filter_resources": {
				Description: "Include only resources that match either a resource type or a resource type::regular expression.  See export guide for additional information.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validators.ValidateSubStringInSlice(resourceExporter.GetAvailableExporterTypes()),
				},
				ForceNew:      true,
				ConflictsWith: []string{"resource_types", "exclude_filter_resources"},
			},
			"replace_with_datasource": {
				Description: "Include only resources that match either a resource type or a resource type::regular expression.  See export guide for additional information.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ForceNew: true,
			},
			"exclude_filter_resources": {
				Description: "Exclude resources that match either a resource type or a resource type::regular expression.  See export guide for additional information.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validators.ValidateSubStringInSlice(resourceExporter.GetAvailableExporterTypes()),
				},
				ForceNew:      true,
				ConflictsWith: []string{"resource_types", "include_filter_resources"},
			},
			"include_state_file": {
				Description: "Export a 'terraform.tfstate' file along with the config file. This can be used for orgs to begin managing existing resources with terraform. When `false`, GUID fields will be omitted from the config file unless a resource reference can be supplied. In this case, the resource type will need to be included in the `resource_types` array.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
			},
			"export_as_hcl": {
				Description:   "Export the config as HCL. Deprecated. Please use the export_format attribute instead",
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      true,
				Default:       false,
				ConflictsWith: []string{"export_format"},
			},
			"export_format": {
				Description: "Export the config as hcl or json or json_hcl.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "json",
				ForceNew:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"hcl",
					"json",
					"json_hcl",
					"hcl_json",
				}, true), // true enables case-insensitive matching
			},
			"split_files_by_resource": {
				Description: "Split export files by resource type. This will also split the terraform provider and variable declarations into their own files.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
			},
			"log_permission_errors": {
				Description: "Log permission/product issues rather than fail.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
			},
			"exclude_attributes": {
				Description: "Attributes to exclude from the config when exporting resources. Each value should be of the form {resource_type}.{attribute}, e.g. 'genesyscloud_user.skills'. Excluded attributes must be optional.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				ForceNew:    true,
			},
			"enable_dependency_resolution": {
				Description: "Adds a \"depends_on\" attribute to genesyscloud_flow resources with a list of resources that are referenced inside the flow configuration . This also resolves and exports all the dependent resources for any given resource. Resources mentioned in exclude_attributes will not be exported.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
			},
			"ignore_cyclic_deps": {
				Description: "Ignore Cyclic Dependencies when building the flows and do not throw an error.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    true,
			},
			"compress": {
				Description: "Compress exported results using zip format.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
			},
			"export_computed": {
				Description: "Export attributes that are marked as being Computed and Optional. Does not attempt to export attributes that are explicitly marked as read-only by the provider. Defaults to true to match existing functionality. This attribute's default value will likely switch to false in a future release.",
				Default:     true,
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
			},
			"use_legacy_architect_flow_exporter": {
				Description: "When set to `false`, architect flow configuration files will be downloaded as part of the flow export process.",
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				ForceNew:    true,
			},
		},
	}
}

func createTfExport(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	tfExporterState.ActivateExporterState()

	if _, ok := d.GetOk("include_filter_resources"); ok {
		gre, _ := NewGenesysCloudResourceExporter(ctx, d, meta, IncludeResources)
		diagErr := gre.Export()
		if diagErr.HasError() {
			return diagErr
		}

		d.SetId(gre.exportDirPath)
		return diagErr
	}

	if _, ok := d.GetOk("exclude_filter_resources"); ok {
		gre, _ := NewGenesysCloudResourceExporter(ctx, d, meta, ExcludeResources)
		diagErr := gre.Export()
		if diagErr.HasError() {
			return diagErr
		}

		d.SetId(gre.exportDirPath)
		return diagErr
	}

	//Dealing with the traditional resource
	gre, _ := NewGenesysCloudResourceExporter(ctx, d, meta, LegacyInclude)
	diagErr := gre.Export()
	if diagErr.HasError() {
		return diagErr
	}

	d.SetId(gre.exportDirPath)

	return diagErr
}

// If the output directory doesn't exist or empty, mark the resource for creation.
func readTfExport(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	path := d.Id()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		d.SetId("")
		return nil
	}
	if isEmpty, diagErr := isDirEmpty(path); isEmpty || diagErr != nil {
		d.SetId("")
		return diagErr
	}

	return nil
}

// Delete everything (files and subdirectories) inside the export directory
// not including the directory itself
func deleteTfExport(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	exportPath := d.Id()
	dir, err := os.ReadDir(exportPath)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, d := range dir {
		os.RemoveAll(filepath.Join(exportPath, d.Name()))
	}

	return nil
}
