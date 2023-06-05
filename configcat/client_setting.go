package configcat

import sw "github.com/configcat/configcat-publicapi-go-client"

func (client *Client) GetSettings(configID string) ([]sw.SettingModel, error) {
	model, response, err := client.apiClient.FeatureFlagsSettingsApi.GetSettings(client.GetAuthContext(), configID).Execute()
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) GetSetting(settingID int32) (*sw.SettingModel, error) {
	model, response, err := client.apiClient.FeatureFlagsSettingsApi.GetSetting(client.GetAuthContext(), settingID).Execute()
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) CreateSetting(configID string, body sw.CreateSettingInitialValues) (*sw.SettingModel, error) {
	model, response, err := client.apiClient.FeatureFlagsSettingsApi.CreateSetting(client.GetAuthContext(), configID).CreateSettingInitialValues(body).Execute()
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) UpdateSetting(settingID int32, jsonPatch sw.JsonPatch) (*sw.SettingModel, error) {
	model, response, err := client.apiClient.FeatureFlagsSettingsApi.UpdateSetting(client.GetAuthContext(), settingID).JsonPatch(jsonPatch).Execute()
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) DeleteSetting(settingID int32) error {
	response, err := client.apiClient.FeatureFlagsSettingsApi.DeleteSetting(client.GetAuthContext(), settingID).Execute()
	defer response.Body.Close()
	return handleAPIError(err)
}
