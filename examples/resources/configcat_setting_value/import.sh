# Feature Flag/Setting values can be imported using a combined EnvironmentID:SettingId ID.  
# Get the SettingId using e.g. the [List Flags API](https://api.configcat.com/docs/#tag/Feature-Flags-and-Settings/operation/get-settings).  
# Get the EnvironmentId using e.g. the [List Environments API](https://api.configcat.com/docs/#tag/Environments/operation/get-environments).  

terraform import configcat_setting_value.example 08d86d63-2726-47cd-8bfc-59608ecb91e2:1234