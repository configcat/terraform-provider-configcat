package client

import sw "github.com/configcat/configcat-publicapi-go-client/v2"

func (client *Client) GetEnvironments(productID string) ([]sw.EnvironmentModel, error) {
	model, response, err := client.apiClient.EnvironmentsAPI.GetEnvironments(client.GetAuthContext(), productID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) GetEnvironment(environmentID string) (*sw.EnvironmentModel, error) {
	model, response, err := client.apiClient.EnvironmentsAPI.GetEnvironment(client.GetAuthContext(), environmentID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) CreateEnvironment(productID string, body sw.CreateEnvironmentModel) (*sw.EnvironmentModel, error) {
	model, response, err := client.apiClient.EnvironmentsAPI.CreateEnvironment(
		client.GetAuthContext(),
		productID).CreateEnvironmentModel(body).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) UpdateEnvironment(environmentID string, body sw.UpdateEnvironmentModel) (*sw.EnvironmentModel, error) {
	model, response, err := client.apiClient.EnvironmentsAPI.UpdateEnvironment(
		client.GetAuthContext(),
		environmentID).UpdateEnvironmentModel(body).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) DeleteEnvironment(environmentID string) error {
	response, err := client.apiClient.EnvironmentsAPI.DeleteEnvironment(
		client.GetAuthContext(),
		environmentID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return error
}
