package configcat

import (
	sw "github.com/configcat/configcat-publicapi-go-client"
)

func (client *Client) GetSettingValue(environmentId string, settingId int32) (sw.SettingValueModel, error) {
	model, response, err := client.apiClient.FeatureFlagSettingValuesApi.GetSettingValue(client.GetAuthContext(), environmentId, settingId)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) ReplaceSettingValue(environmentId string, settingId int32, body sw.UpdateSettingValueModel) (sw.SettingValueModel, error) {
	model, response, err := client.apiClient.FeatureFlagSettingValuesApi.ReplaceSettingValue(client.GetAuthContext(), body, environmentId, settingId, nil)
	defer response.Body.Close()
	return model, handleAPIError(err)
}
