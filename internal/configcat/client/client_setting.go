package client

import sw "github.com/configcat/configcat-publicapi-go-client/v2"

func (client *Client) GetSettings(configID string) ([]sw.SettingModel, error) {
	model, response, err := client.apiClient.FeatureFlagsSettingsAPI.GetSettings(client.GetAuthContext(), configID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) GetSetting(settingID int32) (*sw.SettingModel, error) {
	model, response, err := client.apiClient.FeatureFlagsSettingsAPI.GetSetting(client.GetAuthContext(), settingID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) CreateSetting(configID string, body sw.CreateSettingInitialValues) (*sw.SettingModel, error) {
	model, response, err := client.apiClient.FeatureFlagsSettingsAPI.CreateSetting(client.GetAuthContext(), configID).CreateSettingInitialValues(body).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) UpdateSetting(settingID int32, patchOperation []sw.JsonPatchOperation) (*sw.SettingModel, error) {
	model, response, err := client.apiClient.FeatureFlagsSettingsAPI.UpdateSetting(client.GetAuthContext(), settingID).JsonPatchOperation(patchOperation).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) DeleteSetting(settingID int32) error {
	response, err := client.apiClient.FeatureFlagsSettingsAPI.DeleteSetting(client.GetAuthContext(), settingID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return error
}
