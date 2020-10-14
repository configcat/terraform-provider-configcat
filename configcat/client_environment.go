package configcat

import sw "github.com/configcat/configcat-publicapi-go-client"

func (client *Client) GetEnvironments(productId string) ([]sw.EnvironmentModel, error) {
	model, response, err := client.apiClient.EnvironmentsApi.GetEnvironments(client.GetAuthContext(), productId)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) GetEnvironment(environmentId string) (sw.EnvironmentModel, error) {
	model, response, err := client.apiClient.EnvironmentsApi.GetEnvironment(client.GetAuthContext(), environmentId)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) CreateEnvironment(productId string, body sw.CreateEnvironmentModel) (sw.EnvironmentModel, error) {
	model, response, err := client.apiClient.EnvironmentsApi.CreateEnvironment(
		client.GetAuthContext(),
		body,
		productId)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) UpdateEnvironment(environmentID string, body sw.UpdateEnvironmentModel) (sw.EnvironmentModel, error) {
	model, response, err := client.apiClient.EnvironmentsApi.UpdateEnvironment(
		client.GetAuthContext(),
		body,
		environmentID)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) DeleteEnvironment(environmentID string) error {
	response, err := client.apiClient.EnvironmentsApi.DeleteEnvironment(
		client.GetAuthContext(),
		environmentID)
	defer response.Body.Close()
	return handleAPIError(err)
}
