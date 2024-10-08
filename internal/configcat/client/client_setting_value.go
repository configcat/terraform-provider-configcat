package client

import (
	sw "github.com/configcat/configcat-publicapi-go-client/v2"
)

func (client *Client) GetSettingValue(environmentID string, settingID int32) (*sw.SettingValueModel, error) {
	model, response, err := client.apiClient.FeatureFlagSettingValuesAPI.GetSettingValue(client.GetAuthContext(), environmentID, settingID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) ReplaceSettingValue(environmentID string, settingID int32, body sw.UpdateSettingValueModel, reason string) (*sw.SettingValueModel, error) {
	model, response, err := client.apiClient.FeatureFlagSettingValuesAPI.ReplaceSettingValue(client.GetAuthContext(), environmentID, settingID).UpdateSettingValueModel(body).Reason(reason).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}
