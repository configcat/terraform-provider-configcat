# Setting Tags can be imported using a combined SettingId:TagId ID.  
# Get the SettingId using e.g. the [List Flags API](https://api.configcat.com/docs/#tag/Feature-Flags-and-Settings/operation/get-settings).  
# Get the TagId using e.g. the [List Tags API](https://api.configcat.com/docs/#tag/Tags/operation/get-tags).  

terraform import configcat_setting_tag.example 1234:5678