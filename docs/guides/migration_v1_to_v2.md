---
page_title: "Migration guide from Config V1 to V2"
---

# Migration guide from Config V1 to V2

If you start the [migration process to Config V2](https://configcat.com/docs/advanced/config-v2-migration-guide/), the V2 Config and the Feature Flags or Settings in the V2 Config will be brand new resources in ConfigCat.
If you manage your V1 Configs and Feature Flags with Terraform, and you would like to manage the new V2 resources with Terraform too, you'll have to import the new resources in your Terraform configuration file.

To help this migration process we prepared a small bash script that outputs the terraform resources and the import statements required to complete this process.  
There are 2 variables in the beginning of the script that you have to replace with your Public Api Access Token and the migrated, V2 Config's ConfigId.

```bash
#!/bin/bash

# Get your Public Api Access Token at https://app.configcat.com/my-account/public-api-credentials and use the "Base64 encoded authorization header" as the auth_header
auth_header="#YOUR_AUTH_HEADER"

# Use the V2 config's ConfigId as config_id
config_id="#V2_CONFIG_ID#"


config=$(curl -s -L -X GET "https://api.configcat.com/v1/configs/${config_id}" -H "Accept: application/json" -H "Authorization: ${auth_header}")

echo "Importing config\n"

echo "1. Insert these terraform resources to your terraform configuration file (.tf): \n" 
config_resource="resource \"configcat_config\" \"my_config_$config_id\" {\n"
config_resource="$config_resource    product_id = $(echo $config | jq -c '.product.productId')\n"
config_resource="$config_resource    name = $(echo $config | jq -c '.name')\n"
config_resource="$config_resource    evaluation_version = $(echo $config | jq -c '.evaluationVersion')\n"
config_resource="$config_resource    order = $(echo $config | jq -c '.order')\n"

description=$(echo $config | jq -c '.description')
if [ $description != "\"\"" ]
then
    config_resource="$config_resource  description = $description\n"
fi
config_resource="$config_resource}\n"
echo $config_resource

settings=$(curl -s -L -X GET "https://api.configcat.com/v1/configs/${config_id}/settings" -H "Accept: application/json" -H "Authorization: ${auth_header}")

echo $settings | jq -c '.[]' | while read setting; do
    setting_resource="resource \"configcat_setting\" \"my_setting_$(echo $setting | jq -c '.settingId')\" {\n"
    setting_resource="$setting_resource    config_id = \"$config_id\"\n"
    setting_resource="$setting_resource    key = $(echo $setting | jq -c '.key')\n"
    setting_resource="$setting_resource    name = $(echo $setting | jq -c '.name')\n"
    setting_resource="$setting_resource    setting_type = $(echo $setting | jq -c '.settingType')\n"
    setting_resource="$setting_resource    order = $(echo $setting | jq -c '.order')\n"
    
    hint=$(echo $setting | jq -c '.hint')
    if [ $hint != "\"\"" ]
    then
        setting_resource="$setting_resource  hint = $hint\n"
    fi
    setting_resource="$setting_resource}\n"
    echo $setting_resource
done


echo "\n2. Then import these resources into the terraform state with these statements: \n"
echo "terraform import configcat_config.my_config_$config_id $config_id"
echo $settings | jq -c '.[]' | while read setting; do
    setting_id=$(echo $setting | jq -c '.settingId')
    echo "terraform import configcat_setting.my_setting_$setting_id $setting_id"
done
```

After executing this script, the output will contain the resource definitions that you have to insert into your terraform files and the terraform import statement that will import the resources into Terraform's state.
Please note that this is just an example configuration, you may need to modify the resource names based on your preferences.

Example output:
```text
1. Insert this terraform resource to your terraform configuration file:

resource "configcat_config" "my_config_CONFIG_ID" {
 product_id = "PRODUCT_ID"
 name = "My Config"
 evaluation_version = "v2"
 order = 0
}

resource "configcat_setting" "my_setting_SETTING_ID_1" {
 config_id = "CONFIG_ID"
 key = "isAwesomeFeatureEnabled"
 name = "Is awesome feature enabled"
 setting_type = "boolean"
 order = 0
 hint = "hint"
}

resource "configcat_setting" "my_setting_SETTING_ID_2" {
 config_id = "CONFIG_ID"
 key = "myTextFeatureFlag"
 name = "My text feature flag"
 setting_type = "string"
 order = 1
}

2. Then import the resources into the terraform state with these statements:

terraform import configcat_config.my_config_CONFIG_ID CONFIG_ID
terraform import configcat_setting.my_setting_SETTING_ID_1 SETTING_ID_1
terraform import configcat_setting.my_setting_SETTING_ID_2 SETTING_ID_2
```
