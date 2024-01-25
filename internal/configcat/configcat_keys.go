package configcat

const (
	ID              = "id"
	Name            = "name"
	Description     = "description"
	Order           = "order"
	Color           = "color"
	NameFilterRegex = "name_filter_regex"

	Organizations  = "organizations"
	OrganizationId = "organization_id"

	Products  = "products"
	ProductId = "product_id"

	ConfigResourceName = "Config"
	Configs            = "configs"
	ConfigId           = "config_id"

	Environments  = "environments"
	EnvironmentId = "environment_id"

	PermissionGroups                            = "permission_groups"
	PermissionGroupId                           = "permission_group_id"
	PermissionGroupCanManageMembers             = "can_manage_members"
	PermissionGroupCanCreateOrUpdatConfig       = "can_createorupdate_config"
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

	Segments                   = "segments"
	SegmentId                  = "segment_id"
	SegmentComparisonAttribute = "comparison_attribute"
	SegmentComparator          = "comparator"
	SegmentComparisonValue     = "comparison_value"

	Tags  = "tags"
	TagId = "tag_id"

	Settings              = "settings"
	SettingId             = "setting_id"
	SettingKey            = "key"
	SettingKeyFilterRegex = "key_filter_regex"
	SettingHint           = "hint"
	SettingType           = "setting_type"

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

	PrimarySdkKey   = "primary"
	SecondarySdkKey = "secondary"
)
