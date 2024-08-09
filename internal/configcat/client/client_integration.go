package client

import sw "github.com/configcat/configcat-publicapi-go-client/v2"

func (client *Client) GetIntegrations(productID string) (*sw.IntegrationsModel, error) {
	model, response, err := client.apiClient.IntegrationsAPI.GetIntegrations(client.GetAuthContext(), productID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) GetIntegration(integrationID string) (*sw.IntegrationModel, error) {
	model, response, err := client.apiClient.IntegrationsAPI.GetIntegration(client.GetAuthContext(), integrationID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) CreateIntegration(productID string, body sw.CreateIntegrationModel) (*sw.IntegrationModel, error) {
	model, response, err := client.apiClient.IntegrationsAPI.CreateIntegration(
		client.GetAuthContext(),
		productID).CreateIntegrationModel(body).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) UpdateIntegration(integrationID string, body sw.ModifyIntegrationRequest) (*sw.IntegrationModel, error) {
	model, response, err := client.apiClient.IntegrationsAPI.UpdateIntegration(
		client.GetAuthContext(),
		integrationID).ModifyIntegrationRequest(body).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) DeleteIntegration(integrationID string) error {
	response, err := client.apiClient.IntegrationsAPI.DeleteIntegration(
		client.GetAuthContext(),
		integrationID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return error
}
