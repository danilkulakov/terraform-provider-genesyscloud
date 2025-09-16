package outbound_campaignrule

import (
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/provider"
	resourceExporter "github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/resource_exporter"
	registrar "github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/resource_register"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"
)

/*
resource_genesycloud_outbound_campaignrule_schema.go holds four functions within it:

1.  The registration code that registers the Datasource, Resource and Exporter for the package.
2.  The resource schema definitions for the outbound_campaignrule resource.
3.  The datasource schema definitions for the outbound_campaignrule datasource.
4.  The resource exporter configuration for the outbound_campaignrule exporter.
*/
const ResourceType = "genesyscloud_outbound_campaignrule"

func getAllowedActions() []string {
	return []string{
		"turnOnCampaign",
		"turnOffCampaign",
		"turnOnSequence",
		"turnOffSequence",
		"setCampaignPriority",
		"recycleCampaign",
		"setCampaignDialingMode",
		"setCampaignAbandonRate",
		"setCampaignNumberOfLines",
		"setCampaignWeight",
		"setCampaignMaxCallsPerAgent",
		"changeCampaignQueue",
	}
}

func getAllowedConditions() []string {
	return []string{
		"campaignProgress",
		"campaignAgents",
		"campaignRecordsAttempted",
		"campaignContactsMessaged",
		"campaignBusinessSuccess",
		"campaignBusinessNeutral",
		"campaignBusinessFailure",
		"campaignValidAttempts",
		"campaignRightPartyContacts",
	}
}

var (
	outboundCampaignRuleEntities = &schema.Resource{
		Schema: map[string]*schema.Schema{
			`campaign_ids`:        outboundCampaignRuleEntityCampaignRuleId,
			`sequence_ids`:        outboundCampaignRuleEntitySequenceRuleId,
			`email_campaigns_ids`: outboundCampaignRuleEntityEmailCampaignRuleId,
			`sms_campaigns_ids`:   outboundCampaignRuleEntitySmsCampaignRuleId,
		},
	}

	outboundCampaignRuleExecutionSettings = &schema.Resource{
		Schema: map[string]*schema.Schema{
			`frequency`: {
				Description:  `Execution control frequency.`,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{`oncePerDay`, `onEachTrigger`}, true),
			},
			`time_zone_id`: {
				Description: `The time zone for the execution control frequency="oncePerDay"; for example, Africa/Abidjan. This property is ignored when frequency is not "oncePerDay".`,
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}

	outboundCampaignRuleEntityCampaignRuleId = &schema.Schema{
		Description: `The list of campaigns for a CampaignRule to monitor. Required if the CampaignRule has any conditions that run on a campaign. Changing the outboundCampaignRuleEntityCampaignRuleId attribute will cause the outbound_campaignrule object to be dropped and recreated with a new ID.`,
		Optional:    true,
		ForceNew:    true,
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
	}

	outboundCampaignRuleEntitySequenceRuleId = &schema.Schema{
		Description: `The list of sequences for a CampaignRule to monitor. Required if the CampaignRule has any conditions that run on a sequence. Changing the outboundCampaignRuleEntitySequenceRuleId attribute will cause the outbound_campaignrule object to be dropped and recreated with a new ID.`,
		Optional:    true,
		ForceNew:    true,
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
	}

	outboundCampaignRuleEntityEmailCampaignRuleId = &schema.Schema{
		Description: `The list of Email campaigns for a CampaignRule to monitor. Required if the CampaignRule has any conditions that run on a Email campaign. Changing the outboundCampaignRuleEntityCampaignRuleId attribute will cause the outbound_campaignrule object to be dropped and recreated with a new ID.`,
		Optional:    true,
		ForceNew:    true,
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
	}

	outboundCampaignRuleEntitySmsCampaignRuleId = &schema.Schema{
		Description: `The list of SMS campaigns for a CampaignRule to monitor. Required if the CampaignRule has any conditions that run on a SMS campaign. Changing the outboundCampaignRuleEntityCampaignRuleId attribute will cause the outbound_campaignrule object to be dropped and recreated with a new ID.`,
		Optional:    true,
		ForceNew:    true,
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
	}

	outboundCampaignRuleActionEntities = &schema.Resource{
		Schema: map[string]*schema.Schema{
			`campaign_ids`:        outboundCampaignRuleEntityCampaignRuleId,
			`sequence_ids`:        outboundCampaignRuleEntitySequenceRuleId,
			`email_campaigns_ids`: outboundCampaignRuleEntityEmailCampaignRuleId,
			`sms_campaigns_ids`:   outboundCampaignRuleEntitySmsCampaignRuleId,
			`use_triggering_entity`: {
				Description: `If true, the CampaignRuleAction will apply to the same entity that triggered the CampaignRuleCondition.`,
				Optional:    true,
				Type:        schema.TypeBool,
				Default:     false,
			},
		},
	}
)

// SetRegistrar registers all of the resources, datasources and exporters in the package
func SetRegistrar(regInstance registrar.Registrar) {
	regInstance.RegisterResource(ResourceType, ResourceOutboundCampaignrule())
	regInstance.RegisterDataSource(ResourceType, DataSourceOutboundCampaignrule())
	regInstance.RegisterExporter(ResourceType, OutboundCampaignruleExporter())
}

// ResourceOutboundCampaignrule registers the genesyscloud_outbound_campaignrule resource with Terraform
func ResourceOutboundCampaignrule() *schema.Resource {
	campaignRuleParameters := &schema.Resource{
		Schema: map[string]*schema.Schema{
			`operator`: {
				Description:  `The operator for comparison. Required for a CampaignRuleCondition.`,
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"equals", "greaterThan", "greaterThanEqualTo", "lessThan", "lessThanEqualTo"}, true),
			},
			`value`: {
				Description: `The value for comparison. Required for a CampaignRuleCondition.`,
				Optional:    true,
				Type:        schema.TypeString,
			},
			`priority`: {
				Description:  `The priority to set a campaign to (1 | 2 | 3 | 4 | 5). Required for the 'setCampaignPriority' action.`,
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"1", "2", "3", "4", "5"}, true),
			},
			`dialing_mode`: {
				Description:  `The dialing mode to set a campaign to. Required for the 'setCampaignDialingMode' action (agentless | preview | power | predictive | progressive | external).`,
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"agentless", "preview", "power", "predictive", "progressive", "external"}, true),
			},
			`abandon_rate`: {
				Description:  `Compliance Abandon Rate. Required for 'setCampaignAbandonRate' action`,
				Optional:     true,
				Type:         schema.TypeFloat,
				ValidateFunc: validation.FloatAtLeast(0.1),
			},
			`outbound_line_count`: {
				Description:  `Number of Outbound lines. Required for 'setCampaignNumberOfLines' action`,
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntAtLeast(0),
			},
			`relative_weight`: {
				Description:  `Relative weight. Required for 'setCampaignWeight' action`,
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(0, 100),
			},
			`max_calls_per_agent`: {
				Description:  `Max calls per agent. Optional parameter for 'setCampaignMaxCallsPerAgent' action`,
				Optional:     true,
				Type:         schema.TypeFloat,
				ValidateFunc: validation.FloatAtLeast(1.0),
			},
			`queue_id`: {
				Description: `The ID of the Queue. Required for 'changeCampaignQueue' action`,
				Optional:    true,
				Type:        schema.TypeString,
			},
			`messages_per_minute`: {
				Description:  `The number of messages per minute to set a messaging campaign to.`,
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntAtLeast(0),
			},
			`sms_messages_per_minute`: {
				Description:  `The number of messages per minute to set a SMS messaging campaign to.`,
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntAtLeast(0),
			},
			`email_messages_per_minute`: {
				Description:  `The number of messages per minute to set an Email messaging campaign to.`,
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntAtLeast(0),
			},
			`sms_content_template`: {
				Description: `The content template to set a SMS campaign to.`,
				Optional:    true,
				Type:        schema.TypeString,
			},
			`email_content_template`: {
				Description: `The content template to set an Email campaign to.`,
				Optional:    true,
				Type:        schema.TypeString,
			},
		},
	}

	outboundCampaignRuleCondition := &schema.Resource{
		Schema: map[string]*schema.Schema{
			`id`: {
				Description: `The ID of the CampaignRuleCondition.`,
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
			},
			`parameters`: {
				Description: `The parameters for the CampaignRuleCondition.`,
				Required:    true,
				Type:        schema.TypeSet,
				Elem:        campaignRuleParameters,
			},
			`condition_type`: {
				Description:  `The type of condition to evaluate (` + strings.Join(getAllowedConditions(), ` | `) + `)`,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(getAllowedConditions(), true),
			},
		},
	}

	outboundCampaignRuleConditionGroup := &schema.Resource{
		Schema: map[string]*schema.Schema{
			`match_any_conditions`: {
				Description: `Whether actions are executed if any condition is met, or only when all conditions are met.`,
				Optional:    true,
				Default:     false,
				Type:        schema.TypeBool,
			},
			`conditions`: {
				Description: `The list of conditions that are evaluated on the entities.`,
				Optional:    true,
				MinItems:    1,
				Type:        schema.TypeList,
				Elem:        outboundCampaignRuleCondition,
			},
		},
	}

	outboundCampaignRuleAction := &schema.Resource{
		Schema: map[string]*schema.Schema{
			`id`: {
				Description: `The ID of the CampaignRuleAction.`,
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
			},
			`parameters`: {
				Description: `The parameters for the CampaignRuleAction. Required for certain actionTypes.`,
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        campaignRuleParameters,
			},
			`action_type`: {
				Description:  `The action to take on the campaignRuleActionEntities (` + strings.Join(getAllowedActions(), ` | `) + `)`,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(getAllowedActions(), true),
			},
			`campaign_rule_action_entities`: {
				Description: `The list of entities that this action will apply to.`,
				Required:    true,
				Type:        schema.TypeSet,
				Elem:        outboundCampaignRuleActionEntities,
			},
		},
	}

	return &schema.Resource{
		Description: `Genesys Cloud outbound campaign rule`,

		CreateContext: provider.CreateWithPooledClient(createOutboundCampaignRule),
		ReadContext:   provider.ReadWithPooledClient(readOutboundCampaignRule),
		UpdateContext: provider.UpdateWithPooledClient(updateOutboundCampaignRule),
		DeleteContext: provider.DeleteWithPooledClient(deleteOutboundCampaignRule),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			`name`: {
				Description: `The name of the campaign rule.`,
				Required:    true,
				Type:        schema.TypeString,
			},
			`campaign_rule_entities`: {
				Description: `The list of entities that this campaign rule monitors.`,
				Required:    true,
				MaxItems:    1,
				Type:        schema.TypeSet,
				Elem:        outboundCampaignRuleEntities,
			},
			`campaign_rule_conditions`: {
				Description: `The list of conditions that are evaluated on the entities.`,
				Required:    true,
				Type:        schema.TypeList,
				Elem:        outboundCampaignRuleCondition,
			},
			`campaign_rule_actions`: {
				Description: `The list of actions that are executed if the conditions are satisfied.`,
				Required:    true,
				Type:        schema.TypeList,
				Elem:        outboundCampaignRuleAction,
			},
			`match_any_conditions`: {
				Description: `Whether actions are executed if any condition is met, or only when all conditions are met.`,
				Optional:    true,
				Default:     false,
				Type:        schema.TypeBool,
			},
			`enabled`: {
				Description: `Whether or not this campaign rule is currently enabled.`,
				Optional:    true,
				Default:     false,
				Type:        schema.TypeBool,
			},
			`campaign_rule_processing`: {
				Description:  `CampaignRule processing algorithm.`,
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{`V2`}, true),
			},
			`condition_groups`: {
				Description: `List of condition groups that are evaluated, used only with campaignRuleProcessing="v2"`,
				Optional:    true,
				MinItems:    1,
				Type:        schema.TypeList,
				Elem:        outboundCampaignRuleConditionGroup,
			},
			`execution_settings`: {
				Description: `Campaign rule execution settings.`,
				Optional:    true,
				MaxItems:    1,
				Type:        schema.TypeList,
				Elem:        outboundCampaignRuleExecutionSettings,
			},
		},
	}
}

// OutboundCampaignruleExporter returns the resourceExporter object used to hold the genesyscloud_outbound_campaignrule exporter's config
func OutboundCampaignruleExporter() *resourceExporter.ResourceExporter {
	return &resourceExporter.ResourceExporter{
		GetResourcesFunc: provider.GetAllWithPooledClient(getAllAuthCampaignRules),
		RefAttrs: map[string]*resourceExporter.RefAttrSettings{
			`campaign_rule_actions.campaign_rule_action_entities.campaign_ids`: {
				RefType: "genesyscloud_outbound_campaign",
			},
			`campaign_rule_actions.campaign_rule_action_entities.sequence_ids`: {
				RefType: "genesyscloud_outbound_sequence",
			},
			`campaign_rule_entities.campaign_ids`: {
				RefType: "genesyscloud_outbound_campaign",
			},
			`campaign_rule_entities.sequence_ids`: {
				RefType: "genesyscloud_outbound_sequence",
			},
		},
	}
}

// DataSourceOutboundCampaignrule registers the genesyscloud_outbound_campaignrule data source
func DataSourceOutboundCampaignrule() *schema.Resource {
	return &schema.Resource{
		Description: `Genesys Cloud outbound campaign rule data source. Select a campaign rule by name`,
		ReadContext: provider.ReadWithPooledClient(dataSourceOutboundCampaignruleRead),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Campaign Rule name.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}
