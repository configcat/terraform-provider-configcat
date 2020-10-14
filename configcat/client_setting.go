package configcat

import sw "github.com/configcat/configcat-publicapi-go-client"

func (client *Client) GetSettings(configId string) ([]sw.SettingModel, error) {
	model, response, err := client.apiClient.FeatureFlagsSettingsApi.GetSettings(client.GetAuthContext(), configId)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) GetSetting(settingId int32) (sw.SettingModel, error) {
	model, response, err := client.apiClient.FeatureFlagsSettingsApi.GetSetting(client.GetAuthContext(), settingId)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) CreateSetting(configId string, body sw.CreateSettingModel) (sw.SettingModel, error) {
	model, response, err := client.apiClient.FeatureFlagsSettingsApi.CreateSetting(client.GetAuthContext(), body, configId)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) UpdateSetting(settingId int32, body []sw.Operation) (sw.SettingModel, error) {
	model, response, err := client.apiClient.FeatureFlagsSettingsApi.UpdateSetting(client.GetAuthContext(), body, settingId)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) DeleteSetting(settingId int32) error {
	response, err := client.apiClient.FeatureFlagsSettingsApi.DeleteSetting(client.GetAuthContext(), settingId)
	defer response.Body.Close()
	return handleAPIError(err)
}
