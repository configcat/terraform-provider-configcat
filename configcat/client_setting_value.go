package configcat

import (
	sw "github.com/configcat/configcat-publicapi-go-client"
)

func (client *Client) GetSettingValue(environmentID string, settingID int32) (*sw.SettingValueModel, error) {
	model, response, err := client.apiClient.FeatureFlagSettingValuesApi.GetSettingValue(client.GetAuthContext(), environmentID, settingID).Execute()
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) ReplaceSettingValue(environmentID string, settingID int32, body sw.UpdateSettingValueModel, reason string) (*sw.SettingValueModel, error) {
	model, response, err := client.apiClient.FeatureFlagSettingValuesApi.ReplaceSettingValue(client.GetAuthContext(), environmentID, settingID).UpdateSettingValueModel(body).Reason(reason).Execute()
	defer response.Body.Close()
	return model, handleAPIError(err)
}
