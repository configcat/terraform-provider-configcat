package configcat

const (
	ID              = "id"
	Name            = "name"
	Description     = "description"
	Order           = "order"
	Color           = "color"
	NameFilterRegex = "name_filter_regex"

	OrganizationResourceName = "Organization"
	Organizations            = "organizations"
	OrganizationId           = "organization_id"

	ProductResourceName = "Product"
	Products            = "products"
	ProductId           = "product_id"

	ProductPreferencesResourceName                  = "Product Preferences"
	ProductPreferenceReasonRequired                 = "reason_required"
	ProductPreferenceReasonRequiredEnvironmentments = "reason_required_environments"
	ProductPreferenceKeyGenerationMode              = "key_generation_mode"
	ProductPreferenceShowVariationId                = "show_variation_id"
	ProductPreferenceMandatorySettingHint           = "mandatory_setting_hint"

	ConfigResourceName = "Config"
	Configs            = "configs"
	ConfigId           = "config_id"
	EvaluationVersion  = "evaluation_version"

	EnvironmentResourceName = "Environment"
	Environments            = "environments"
	EnvironmentId           = "environment_id"

	WebhookResourceName  = "Webhook"
	WebhookId            = "webhook_id"
	WebhookUrl           = "url"
	WebhookContent       = "content"
	WebhookHttpMethod    = "http_method"
	WebhookHeaders       = "webhook_headers"
	SecureWebhookHeaders = "secure_webhook_headers"
	WebhookHeaderKey     = "key"
	WebhookHeaderValue   = "value"

	WebhookSigningKeysResourceName = "Webhook Signing Keys"
	WebhookSigningKeyKey1          = "key1"
	WebhookSigningKeyKey2          = "key2"

	PermissionGroupResourceName                 = "Permission Group"
	PermissionGroups                            = "permission_groups"
	PermissionGroupId                           = "permission_group_id"
	PermissionGroupCanManageMembers             = "can_manage_members"
	PermissionGroupCanCreateOrUpdateConfig      = "can_createorupdate_config"
	PermissionGroupCanDeleteConfig              = "can_delete_config"
	PermissionGroupCanCreateOrUpdateEnvironment = "can_createorupdate_environment"
	PermissionGroupCanDeleteEnvironment         = "can_delete_environment"
	PermissionGroupCanCreateOrUpdateSetting     = "can_createorupdate_setting"
	PermissionGroupCanTagSetting                = "can_tag_setting"
	PermissionGroupCanDeleteSetting             = "can_delete_setting"
	PermissionGroupCanCreateOrUpdateTag         = "can_createorupdate_tag"
	PermissionGroupCanDeleteTag                 = "can_delete_tag"
	PermissionGroupCanManageWebhook             = "can_manage_webhook"
	PermissionGroupCanUseExportImport           = "can_use_exportimport"
	PermissionGroupCanManageProductPreferences  = "can_manage_product_preferences"
	PermissionGroupCanManageIntegrations        = "can_manage_integrations"
	PermissionGroupCanViewSdkKey                = "can_view_sdkkey"
	PermissionGroupCanRotateSdkKey              = "can_rotate_sdkkey"
	PermissionGroupCanCreateOrUpdateSegment     = "can_createorupdate_segment"
	PermissionGroupCanDeleteSegment             = "can_delete_segment"
	PermissionGroupCanViewProductAuditLogs      = "can_view_product_auditlog"
	PermissionGroupCanViewProductStatistics     = "can_view_product_statistics"
	PermissionGroupAccessType                   = "accesstype"
	PermissionGroupNewEnvironmentAccessType     = "new_environment_accesstype"
	PermissionGroupEnvironmentAccess            = "environment_accesses"

	SegmentResourceName        = "Segment"
	Segments                   = "segments"
	SegmentId                  = "segment_id"
	SegmentComparisonAttribute = "comparison_attribute"
	SegmentComparator          = "comparator"
	SegmentComparisonValue     = "comparison_value"

	TagResourceName = "Tag"
	Tags            = "tags"
	TagId           = "tag_id"

	SettingsResourceName  = "Feature Flags or Settings"
	SettingResourceName   = "Feature Flag or Setting"
	Settings              = "settings"
	SettingId             = "setting_id"
	SettingKey            = "key"
	SettingKeyFilterRegex = "key_filter_regex"
	SettingHint           = "hint"
	SettingType           = "setting_type"

	// V1
	SettingValueResourceName = "Feature Flag or Setting value"

	SettingValue   = "value"
	InitOnly       = "init_only"
	MandatoryNotes = "mandatory_notes"

	RolloutRules                   = "rollout_rules"
	RolloutRuleComparisonAttribute = "comparison_attribute"
	RolloutRuleComparator          = "comparator"
	RolloutRuleComparisonValue     = "comparison_value"
	RolloutRuleValue               = "value"
	RolloutRuleSegmentComparator   = "segment_comparator"
	RolloutRuleSegmentId           = "segment_id"

	RolloutPercentageItems          = "percentage_items"
	RolloutPercentageItemPercentage = "percentage"
	RolloutPercentageItemValue      = "value"

	// SDK KEY
	SdkKeyResourceName = "SDK Key"
	PrimarySdkKey      = "primary"
	SecondarySdkKey    = "secondary"

	// V2
	BoolValue   = "bool_value"
	StringValue = "string_value"
	IntValue    = "int_value"
	DoubleValue = "double_value"

	ListValues     = "list_values"
	ListValueValue = "value"
	ListValueHint  = "hint"

	PercentageEvaluationAttribute = "percentage_evaluation_attribute"

	DefaultValue = "value"

	TargetingRules = "targeting_rules"

	TargetingRuleValue = "value"

	TargetingRulePercentageOptions          = "percentage_options"
	TargetingRulePercentageOptionPercentage = "percentage"
	TargetingRulePercentageOptionValue      = "value"

	TargetingRuleConditions = "conditions"

	TargetingRuleUserCondition                    = "user_condition"
	TargetingRuleUserConditionComparisonAttribute = "comparison_attribute"
	TargetingRuleUserConditionComparator          = "comparator"
	TargetingRuleUserConditionComparisonValue     = "comparison_value"

	TargetingRuleSegmentCondition           = "segment_condition"
	TargetingRuleSegmentConditionSegmentId  = "segment_id"
	TargetingRuleSegmentConditionComparator = "comparator"

	TargetingRulePrerequisiteFlagCondition                = "prerequisite_flag_condition"
	TargetingRulePrerequisiteFlagConditionSettingId       = "prerequisite_setting_id"
	TargetingRulePrerequisiteFlagConditionComparator      = "comparator"
	TargetingRulePrerequisiteFlagConditionComparisonValue = "comparison_value"
)
