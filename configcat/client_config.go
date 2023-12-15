package configcat

import sw "github.com/configcat/configcat-publicapi-go-client"

func (client *Client) GetConfigs(productID string) ([]sw.ConfigModel, error) {
	model, response, err := client.apiClient.ConfigsApi.GetConfigs(client.GetAuthContext(), productID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) GetConfig(configID string) (*sw.ConfigModel, error) {
	model, response, err := client.apiClient.ConfigsApi.GetConfig(client.GetAuthContext(), configID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) CreateConfig(productID string, body sw.CreateConfigRequest) (*sw.ConfigModel, error) {
	model, response, err := client.apiClient.ConfigsApi.CreateConfig(
		client.GetAuthContext(),
		productID).CreateConfigRequest(body).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) UpdateConfig(configID string, body sw.UpdateConfigRequest) (*sw.ConfigModel, error) {
	model, response, err := client.apiClient.ConfigsApi.UpdateConfig(
		client.GetAuthContext(),
		configID).UpdateConfigRequest(body).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) DeleteConfig(configID string) error {
	response, err := client.apiClient.ConfigsApi.DeleteConfig(
		client.GetAuthContext(),
		configID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return error
}
