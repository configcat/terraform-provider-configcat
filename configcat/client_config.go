package configcat

import sw "github.com/configcat/configcat-publicapi-go-client"

func (client *Client) GetConfigs(productID string) ([]sw.ConfigModel, error) {
	model, response, err := client.apiClient.ConfigsApi.GetConfigs(client.GetAuthContext(), productID)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) GetConfig(configID string) (sw.ConfigModel, error) {
	model, response, err := client.apiClient.ConfigsApi.GetConfig(client.GetAuthContext(), configID)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) CreateConfig(productID string, body sw.CreateConfigRequest) (sw.ConfigModel, error) {
	model, response, err := client.apiClient.ConfigsApi.CreateConfig(
		client.GetAuthContext(),
		body,
		productID)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) UpdateConfig(configID string, body sw.UpdateConfigRequest) (sw.ConfigModel, error) {
	model, response, err := client.apiClient.ConfigsApi.UpdateConfig(
		client.GetAuthContext(),
		body,
		configID)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) DeleteConfig(configID string) error {
	response, err := client.apiClient.ConfigsApi.DeleteConfig(
		client.GetAuthContext(),
		configID)
	defer response.Body.Close()
	return handleAPIError(err)
}
