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
