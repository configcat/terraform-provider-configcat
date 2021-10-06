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

	ROLLOUT_PERCENTAGE_ITEMS           = "percentage_items"
	ROLLOUT_PERCENTAGE_ITEM_PERCENTAGE = "percentage"
	ROLLOUT_PERCENTAGE_ITEM_VALUE      = "value"

	PRIMARY_SDK_KEY   = "primary"
	SECONDARY_SDK_KEY = "secondary"
)
