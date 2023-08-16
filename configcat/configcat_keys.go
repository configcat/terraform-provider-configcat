package configcat

const (
	ORGANIZATIONS                  = "organizations"
	ORGANIZATION_ID                = "organization_id"
	ORGANIZATION_NAME              = "name"
	ORGANIZATION_NAME_FILTER_REGEX = "name_filter_regex"

	PRODUCTS                  = "products"
	PRODUCT_ID                = "product_id"
	PRODUCT_NAME              = "name"
	PRODUCT_DESCRIPTION       = "description"
	PRODUCT_NAME_FILTER_REGEX = "name_filter_regex"

	CONFIGS                  = "configs"
	CONFIG_ID                = "config_id"
	CONFIG_NAME              = "name"
	CONFIG_DESCRIPTION       = "description"
	CONFIG_NAME_FILTER_REGEX = "name_filter_regex"

	ENVIRONMENTS                  = "environments"
	ENVIRONMENT_ID                = "environment_id"
	ENVIRONMENT_NAME              = "name"
	ENVIRONMENT_DESCRIPTION       = "description"
	ENVIRONMENT_COLOR             = "color"
	ENVIRONMENT_NAME_FILTER_REGEX = "name_filter_regex"

	PERMISSION_GROUPS                                           = "permission_groups"
	PERMISSION_GROUP_ID                                         = "permission_group_id"
	PERMISSION_GROUP_NAME                                       = "name"
	PERMISSION_GROUP_NAME_FILTER_REGEX                          = "name_filter_regex"
	PERMISSION_GROUP_CAN_MANAGE_MEMBERS                         = "can_manage_members"
	PERMISSION_GROUP_CAN_CREATEORUPDATE_CONFIG                  = "can_createorupdate_config"
	PERMISSION_GROUP_CAN_DELETE_CONFIG                          = "can_delete_config"
	PERMISSION_GROUP_CAN_CREATEORUPDATE_ENVIRONMENT             = "can_createorupdate_environment"
	PERMISSION_GROUP_CAN_DELETE_ENVIRONMENT                     = "can_delete_environment"
	PERMISSION_GROUP_CAN_CREATEORUPDATE_SETTING                 = "can_createorupdate_setting"
	PERMISSION_GROUP_CAN_TAG_SETTING                            = "can_tag_setting"
	PERMISSION_GROUP_CAN_DELETE_SETTING                         = "can_delete_setting"
	PERMISSION_GROUP_CAN_CREATEORUPDATE_TAG                     = "can_createorupdate_tag"
	PERMISSION_GROUP_CAN_DELETE_TAG                             = "can_delete_tag"
	PERMISSION_GROUP_CAN_MANAGE_WEBHOOK                         = "can_manage_webhook"
	PERMISSION_GROUP_CAN_USE_EXPORTIMPORT                       = "can_use_exportimport"
	PERMISSION_GROUP_CAN_MANAGE_PRODUCT_PREFERENCES             = "can_manage_product_preferences"
	PERMISSION_GROUP_CAN_MANAGE_INTEGRATIONS                    = "can_manage_integrations"
	PERMISSION_GROUP_CAN_VIEW_SDKKEY                            = "can_view_sdkkey"
	PERMISSION_GROUP_CAN_ROTATE_SDKKEY                          = "can_rotate_sdkkey"
	PERMISSION_GROUP_CAN_CREATEORUPDATE_SEGMENT                 = "can_createorupdate_segment"
	PERMISSION_GROUP_CAN_DELETE_SEGMENT                         = "can_delete_segment"
	PERMISSION_GROUP_CAN_VIEW_PRODUCT_AUDITLOG                  = "can_view_product_auditlog"
	PERMISSION_GROUP_CAN_VIEW_PRODUCT_STATISTICS                = "can_view_product_statistics"
	PERMISSION_GROUP_ACCESSTYPE                                 = "accesstype"
	PERMISSION_GROUP_NEW_ENVIRONMENT_ACCESSTYPE                 = "new_environment_accesstype"
	PERMISSION_GROUP_ENVIRONMENT_ACCESSES                       = "environment_accesses"
	PERMISSION_GROUP_ENVIRONMENT_ACCESS_ENVIRONMENT_ID          = "environment_id"
	PERMISSION_GROUP_ENVIRONMENT_ACCESS_ENVIRONMENT_ACCESS_TYPE = "environment_access_type"

	SEGMENTS                     = "segments"
	SEGMENT_ID                   = "segment_id"
	SEGMENT_NAME                 = "name"
	SEGMENT_DESCRIPTION          = "description"
	SEGMENT_COMPARISON_ATTRIBUTE = "comparison_attribute"
	SEGMENT_COMPARATOR           = "comparator"
	SEGMENT_COMPARISON_VALUE     = "comparison_value"
	SEGMENT_NAME_FILTER_REGEX    = "name_filter_regex"

	TAGS                  = "tags"
	TAG_ID                = "tag_id"
	TAG_NAME              = "name"
	TAG_COLOR             = "color"
	TAG_NAME_FILTER_REGEX = "name_filter_regex"

	SETTINGS                 = "settings"
	SETTING_ID               = "setting_id"
	SETTING_KEY              = "key"
	SETTING_KEY_FILTER_REGEX = "key_filter_regex"
	SETTING_NAME             = "name"
	SETTING_HINT             = "hint"
	SETTING_TYPE             = "setting_type"

	SETTING_VALUE   = "value"
	INIT_ONLY       = "init_only"
	MANDATORY_NOTES = "mandatory_notes"

	ROLLOUT_RULES                     = "rollout_rules"
	ROLLOUT_RULE_COMPARISON_ATTRIBUTE = "comparison_attribute"
	ROLLOUT_RULE_COMPARATOR           = "comparator"
	ROLLOUT_RULE_COMPARISON_VALUE     = "comparison_value"
	ROLLOUT_RULE_VALUE                = "value"
	ROLLOUT_RULE_SEGMENT_COMPARATOR   = "segment_comparator"
	ROLLOUT_RULE_SEGMENT_ID           = "segment_id"

	ROLLOUT_PERCENTAGE_ITEMS           = "percentage_items"
	ROLLOUT_PERCENTAGE_ITEM_PERCENTAGE = "percentage"
	ROLLOUT_PERCENTAGE_ITEM_VALUE      = "value"

	PRIMARY_SDK_KEY   = "primary"
	SECONDARY_SDK_KEY = "secondary"
)
