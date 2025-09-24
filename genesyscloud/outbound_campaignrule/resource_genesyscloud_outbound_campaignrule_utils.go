package outbound_campaignrule

import (
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/util"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/util/resourcedata"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mypurecloud/platform-client-sdk-go/v165/platformclientv2"
)

func getCampaignruleFromResourceData(d *schema.ResourceData) platformclientv2.Campaignrule {
	matchAnyConditions := d.Get("match_any_conditions").(bool)

	campaignRule := platformclientv2.Campaignrule{
		Name:                   platformclientv2.String(d.Get("name").(string)),
		Enabled:                platformclientv2.Bool(false), // All campaign rules have to be created in an "off" state to start out with
		CampaignRuleEntities:   buildCampaignRuleEntities(d.Get("campaign_rule_entities").(*schema.Set)),
		CampaignRuleConditions: buildCampaignRuleConditions(d.Get("campaign_rule_conditions").([]interface{})),
		CampaignRuleActions:    buildCampaignRuleAction(d.Get("campaign_rule_actions").([]interface{})),
		MatchAnyConditions:     &matchAnyConditions,
		CampaignRuleProcessing: platformclientv2.String(d.Get("campaign_rule_processing").(string)),
		ConditionGroups:        buildCampaignRuleConditionGroups(d.Get("condition_groups").([]interface{})),
		ExecutionSettings:      buildCampaignRuleExecutionSettings(d.Get("execution_settings").(*schema.Set)),
	}
	return campaignRule
}

func buildCampaignRuleEntities(entities *schema.Set) *platformclientv2.Campaignruleentities {
	if entities == nil {
		return nil
	}
	var campaignRuleEntities platformclientv2.Campaignruleentities

	campaignRuleEntitiesList := entities.List()

	if len(campaignRuleEntitiesList) <= 0 {
		return &campaignRuleEntities
	}

	campaignRuleEntitiesMap := campaignRuleEntitiesList[0].(map[string]interface{})
	if campaigns, ok := campaignRuleEntitiesMap["campaign_ids"].([]interface{}); ok && campaigns != nil {
		campaignRuleEntities.Campaigns = util.BuildSdkDomainEntityRefArrFromArr(campaigns)
	}
	if sequences, ok := campaignRuleEntitiesMap["sequence_ids"].([]interface{}); ok && sequences != nil {
		campaignRuleEntities.Sequences = util.BuildSdkDomainEntityRefArrFromArr(sequences)
	}
	if smsCampaigns, ok := campaignRuleEntitiesMap["sms_campaign_ids"].([]interface{}); ok && smsCampaigns != nil {
		campaignRuleEntities.EmailCampaigns = util.BuildSdkDomainEntityRefArrFromArr(smsCampaigns)
	}
	if emailCampaigns, ok := campaignRuleEntitiesMap["email_campaign_ids"].([]interface{}); ok && emailCampaigns != nil {
		campaignRuleEntities.SmsCampaigns = util.BuildSdkDomainEntityRefArrFromArr(emailCampaigns)
	}
	return &campaignRuleEntities
}

func buildCampaignRuleConditions(campaignRuleConditions []interface{}) *[]platformclientv2.Campaignrulecondition {
	var campaignRuleConditionSlice []platformclientv2.Campaignrulecondition

	for _, campaignRuleCondition := range campaignRuleConditions {
		sdkCondition := platformclientv2.Campaignrulecondition{}
		conditionMap := campaignRuleCondition.(map[string]interface{})

		sdkCondition.Parameters = buildCampaignRuleParameters(conditionMap["parameters"].(*schema.Set))
		resourcedata.BuildSDKStringValueIfNotNil(&sdkCondition.Id, conditionMap, "id")
		resourcedata.BuildSDKStringValueIfNotNil(&sdkCondition.ConditionType, conditionMap, "condition_type")

		campaignRuleConditionSlice = append(campaignRuleConditionSlice, sdkCondition)
	}

	return &campaignRuleConditionSlice
}

func buildCampaignRuleAction(campaignRuleActions []interface{}) *[]platformclientv2.Campaignruleaction {
	var campaignRuleActionSlice []platformclientv2.Campaignruleaction

	for _, campaignRuleAction := range campaignRuleActions {
		var sdkCampaignRuleAction platformclientv2.Campaignruleaction
		actionMap := campaignRuleAction.(map[string]interface{})

		resourcedata.BuildSDKStringValueIfNotNil(&sdkCampaignRuleAction.Id, actionMap, "id")
		resourcedata.BuildSDKStringValueIfNotNil(&sdkCampaignRuleAction.ActionType, actionMap, "action_type")
		sdkCampaignRuleAction.Parameters = buildCampaignRuleParameters(actionMap["parameters"].(*schema.Set))
		sdkCampaignRuleAction.CampaignRuleActionEntities = buildCampaignRuleActionEntities(actionMap["campaign_rule_action_entities"].(*schema.Set))

		campaignRuleActionSlice = append(campaignRuleActionSlice, sdkCampaignRuleAction)
	}

	return &campaignRuleActionSlice
}

func buildCampaignRuleParameters(set *schema.Set) *platformclientv2.Campaignruleparameters {
	var sdkCampaignRuleParameters platformclientv2.Campaignruleparameters

	paramsList := set.List()

	if len(paramsList) <= 0 {
		return &sdkCampaignRuleParameters
	}

	paramsMap := paramsList[0].(map[string]interface{})

	resourcedata.BuildSDKStringValueIfNotNil(&sdkCampaignRuleParameters.Operator, paramsMap, "operator")
	resourcedata.BuildSDKStringValueIfNotNil(&sdkCampaignRuleParameters.Value, paramsMap, "value")
	resourcedata.BuildSDKStringValueIfNotNil(&sdkCampaignRuleParameters.Priority, paramsMap, "priority")
	resourcedata.BuildSDKStringValueIfNotNil(&sdkCampaignRuleParameters.DialingMode, paramsMap, "dialing_mode")

	if abandonRate, ok := paramsMap["abandon_rate"].(float32); ok {
		sdkCampaignRuleParameters.AbandonRate = platformclientv2.Float32(abandonRate)
	}
	if lineCount, ok := paramsMap["outbound_line_count"].(int); ok {
		sdkCampaignRuleParameters.OutboundLineCount = platformclientv2.Int(lineCount)
	}
	if weight, ok := paramsMap["relative_weight"].(int); ok {
		sdkCampaignRuleParameters.RelativeWeight = platformclientv2.Int(weight)
	}
	if maxCpa, ok := paramsMap["max_calls_per_agent"].(float32); ok {
		sdkCampaignRuleParameters.MaxCallsPerAgent = platformclientv2.Float32(maxCpa)
	}
	if queueId, ok := paramsMap["queue_id"].(string); ok {
		sdkCampaignRuleParameters.Queue = &platformclientv2.Domainentityref{Id: &queueId}
	}
	if messagesPerMinute, ok := paramsMap["messages_per_minute"].(int); ok {
		sdkCampaignRuleParameters.MessagesPerMinute = platformclientv2.Int(messagesPerMinute)
	}
	if smsMessagesPerMinute, ok := paramsMap["sms_messages_per_minute"].(int); ok {
		sdkCampaignRuleParameters.SmsMessagesPerMinute = platformclientv2.Int(smsMessagesPerMinute)
	}
	if emailMessagesPerMinute, ok := paramsMap["email_messages_per_minute"].(int); ok {
		sdkCampaignRuleParameters.EmailMessagesPerMinute = platformclientv2.Int(emailMessagesPerMinute)
	}
	if emailTemplateId, ok := paramsMap["email_content_template"].(string); ok {
		sdkCampaignRuleParameters.EmailContentTemplate = &platformclientv2.Domainentityref{Id: &emailTemplateId}
	}
	if smsTemplateId, ok := paramsMap["sms_content_template"].(string); ok {
		sdkCampaignRuleParameters.SmsContentTemplate = &platformclientv2.Domainentityref{Id: &smsTemplateId}
	}

	return &sdkCampaignRuleParameters
}

func buildCampaignRuleActionEntities(set *schema.Set) *platformclientv2.Campaignruleactionentities {
	var (
		sdkCampaignRuleActionEntities platformclientv2.Campaignruleactionentities
		entities                      = set.List()
	)

	if len(entities) <= 0 {
		return &sdkCampaignRuleActionEntities
	}

	entitiesMap := entities[0].(map[string]interface{})

	sdkCampaignRuleActionEntities.UseTriggeringEntity = platformclientv2.Bool(entitiesMap["use_triggering_entity"].(bool))

	if campaignIds, ok := entitiesMap["campaign_ids"].([]interface{}); ok && campaignIds != nil {
		sdkCampaignRuleActionEntities.Campaigns = util.BuildSdkDomainEntityRefArrFromArr(campaignIds)
	}
	if sequenceIds, ok := entitiesMap["sequence_ids"].([]interface{}); ok && sequenceIds != nil {
		sdkCampaignRuleActionEntities.Sequences = util.BuildSdkDomainEntityRefArrFromArr(sequenceIds)
	}
	if smsCampaignIds, ok := entitiesMap["sms_campaign_ids"].([]interface{}); ok && smsCampaignIds != nil {
		sdkCampaignRuleActionEntities.EmailCampaigns = util.BuildSdkDomainEntityRefArrFromArr(smsCampaignIds)
	}
	if emailCampaignIds, ok := entitiesMap["email_campaign_ids"].([]interface{}); ok && emailCampaignIds != nil {
		sdkCampaignRuleActionEntities.SmsCampaigns = util.BuildSdkDomainEntityRefArrFromArr(emailCampaignIds)
	}

	return &sdkCampaignRuleActionEntities
}

func buildCampaignRuleConditionGroups(campaignRuleConditionGroups []interface{}) *[]platformclientv2.Campaignruleconditiongroup {
	var campaignRuleConditionGroupSlice []platformclientv2.Campaignruleconditiongroup

	for _, campaignRuleConditionGroup := range campaignRuleConditionGroups {
		var sdkCampaignRuleConditionGroup platformclientv2.Campaignruleconditiongroup
		conditionGroupMap := campaignRuleConditionGroup.(map[string]interface{})

		if matchAnyConditions, ok := conditionGroupMap["match_any_conditions"].(bool); ok {
			sdkCampaignRuleConditionGroup.MatchAnyConditions = platformclientv2.Bool(matchAnyConditions)
		}
		sdkCampaignRuleConditionGroup.Conditions = buildCampaignRuleConditions(conditionGroupMap["conditions"].([]interface{}))

		campaignRuleConditionGroupSlice = append(campaignRuleConditionGroupSlice, sdkCampaignRuleConditionGroup)
	}

	return &campaignRuleConditionGroupSlice
}

func buildCampaignRuleExecutionSettings(set *schema.Set) *platformclientv2.Campaignruleexecutionsettings {
	var (
		sdkCampaignRuleExecutionSettings platformclientv2.Campaignruleexecutionsettings
		settings                         = set.List()
	)

	if len(settings) <= 0 {
		return &sdkCampaignRuleExecutionSettings
	}

	settingsMap := settings[0].(map[string]interface{})

	resourcedata.BuildSDKStringValueIfNotNil(&sdkCampaignRuleExecutionSettings.Frequency, settingsMap, "frequency")
	resourcedata.BuildSDKStringValueIfNotNil(&sdkCampaignRuleExecutionSettings.TimeZoneId, settingsMap, "time_zone_id")

	return &sdkCampaignRuleExecutionSettings
}

func flattenCampaignRuleEntities(campaignRuleEntities *platformclientv2.Campaignruleentities) *schema.Set {
	var (
		campaignRuleEntitiesSet = schema.NewSet(schema.HashResource(outboundCampaignRuleEntities), []interface{}{})
		campaignRuleEntitiesMap = make(map[string]interface{})

		// had to change from []string to []interface{}
		campaigns      []interface{}
		sequences      []interface{}
		smsCampaigns   []interface{}
		emailCampaigns []interface{}
	)

	if campaignRuleEntities == nil {
		return nil
	}

	if campaignRuleEntities.Campaigns != nil {
		for _, v := range *campaignRuleEntities.Campaigns {
			campaigns = append(campaigns, *v.Id)
		}
	}

	if campaignRuleEntities.Sequences != nil {
		for _, v := range *campaignRuleEntities.Sequences {
			sequences = append(sequences, *v.Id)
		}
	}

	if campaignRuleEntities.SmsCampaigns != nil {
		for _, v := range *campaignRuleEntities.SmsCampaigns {
			smsCampaigns = append(smsCampaigns, *v.Id)
		}
	}

	if campaignRuleEntities.EmailCampaigns != nil {
		for _, v := range *campaignRuleEntities.EmailCampaigns {
			emailCampaigns = append(emailCampaigns, *v.Id)
		}
	}

	campaignRuleEntitiesMap["campaign_ids"] = campaigns
	campaignRuleEntitiesMap["sequence_ids"] = sequences
	campaignRuleEntitiesMap["sms_campaigns_ids"] = smsCampaigns
	campaignRuleEntitiesMap["email_campaigns_ids"] = emailCampaigns

	campaignRuleEntitiesSet.Add(campaignRuleEntitiesMap)
	return campaignRuleEntitiesSet
}

func flattenCampaignRuleConditions(campaignRuleConditions *[]platformclientv2.Campaignrulecondition) []interface{} {
	if campaignRuleConditions == nil || len(*campaignRuleConditions) == 0 {
		return nil
	}

	var ruleConditionList []interface{}

	for _, currentSdkCondition := range *campaignRuleConditions {
		campaignRuleConditionsMap := make(map[string]interface{})

		resourcedata.SetMapValueIfNotNil(campaignRuleConditionsMap, "id", currentSdkCondition.Id)
		resourcedata.SetMapValueIfNotNil(campaignRuleConditionsMap, "condition_type", currentSdkCondition.ConditionType)
		resourcedata.SetMapInterfaceArrayWithFuncIfNotNil(campaignRuleConditionsMap, "parameters", currentSdkCondition.Parameters, flattenRuleParameters)

		ruleConditionList = append(ruleConditionList, campaignRuleConditionsMap)
	}
	return ruleConditionList
}

func flattenCampaignRuleAction[T any](campaignRuleActions *[]platformclientv2.Campaignruleaction, actionEntitiesFunc func(*platformclientv2.Campaignruleactionentities) T) []interface{} {
	if campaignRuleActions == nil {
		return nil
	}

	var ruleActionsList []interface{}

	for _, currentAction := range *campaignRuleActions {
		actionMap := make(map[string]interface{})

		resourcedata.SetMapValueIfNotNil(actionMap, "id", currentAction.Id)
		resourcedata.SetMapValueIfNotNil(actionMap, "action_type", currentAction.ActionType)
		resourcedata.SetMapInterfaceArrayWithFuncIfNotNil(actionMap, "parameters", currentAction.Parameters, flattenRuleParameters)
		if currentAction.CampaignRuleActionEntities != nil {
			actionMap["campaign_rule_action_entities"] = actionEntitiesFunc(currentAction.CampaignRuleActionEntities)
		}

		ruleActionsList = append(ruleActionsList, actionMap)
	}

	return ruleActionsList
}

func flattenCampaignRuleActionEntities(sdkActionEntity *platformclientv2.Campaignruleactionentities) *schema.Set {
	var (
		campaigns      []interface{}
		sequences      []interface{}
		smsCampaigns   []interface{}
		emailCampaigns []interface{}
		entitiesSet    = schema.NewSet(schema.HashResource(outboundCampaignRuleActionEntities), []interface{}{})
		entitiesMap    = make(map[string]interface{})
	)

	if sdkActionEntity == nil {
		return nil
	}

	if sdkActionEntity.Campaigns != nil {
		for _, campaign := range *sdkActionEntity.Campaigns {
			campaigns = append(campaigns, *campaign.Id)
		}
	}

	if sdkActionEntity.Sequences != nil {
		for _, sequence := range *sdkActionEntity.Sequences {
			sequences = append(sequences, *sequence.Id)
		}
	}

	if sdkActionEntity.SmsCampaigns != nil {
		for _, campaign := range *sdkActionEntity.SmsCampaigns {
			smsCampaigns = append(smsCampaigns, *campaign.Id)
		}
	}

	if sdkActionEntity.EmailCampaigns != nil {
		for _, campaign := range *sdkActionEntity.EmailCampaigns {
			emailCampaigns = append(emailCampaigns, *campaign.Id)
		}
	}

	entitiesMap["campaign_ids"] = campaigns
	entitiesMap["sequence_ids"] = sequences
	entitiesMap["sms_campaign_ids"] = smsCampaigns
	entitiesMap["email_campaign_ids"] = emailCampaigns
	entitiesMap["use_triggering_entity"] = *sdkActionEntity.UseTriggeringEntity

	entitiesSet.Add(entitiesMap)
	return entitiesSet
}

func flattenRuleParameters(params *platformclientv2.Campaignruleparameters) []interface{} {
	paramsMap := make(map[string]interface{})

	resourcedata.SetMapValueIfNotNil(paramsMap, "operator", params.Operator)
	resourcedata.SetMapValueIfNotNil(paramsMap, "value", params.Value)
	resourcedata.SetMapValueIfNotNil(paramsMap, "priority", params.Priority)
	resourcedata.SetMapValueIfNotNil(paramsMap, "dialing_mode", params.DialingMode)
	resourcedata.SetMapValueIfNotNil(paramsMap, "abandon_rate", params.AbandonRate)
	resourcedata.SetMapValueIfNotNil(paramsMap, "outbound_line_count", params.OutboundLineCount)
	resourcedata.SetMapValueIfNotNil(paramsMap, "relative_weight", params.RelativeWeight)
	resourcedata.SetMapValueIfNotNil(paramsMap, "max_calls_per_agent", params.MaxCallsPerAgent)
	resourcedata.SetMapValueIfNotNil(paramsMap, "queue_id", params.Queue)
	resourcedata.SetMapValueIfNotNil(paramsMap, "messages_per_minute", params.MessagesPerMinute)
	resourcedata.SetMapValueIfNotNil(paramsMap, "sms_messages_per_minute", params.SmsMessagesPerMinute)
	resourcedata.SetMapValueIfNotNil(paramsMap, "email_messages_per_minute", params.EmailMessagesPerMinute)
	resourcedata.SetMapValueIfNotNil(paramsMap, "sms_content_template", params.SmsContentTemplate)
	resourcedata.SetMapValueIfNotNil(paramsMap, "email_content_template", params.EmailContentTemplate)

	return []interface{}{paramsMap}
}
