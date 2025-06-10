package client

import sw "github.com/configcat/configcat-publicapi-go-client/v3"

func (client *Client) GetWebhooks(productID string) ([]sw.WebhookResponseModel, error) {
	model, response, err := client.apiClient.WebhooksAPI.GetWebhooks(client.GetAuthContext(), productID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) GetWebhook(webhookId int32) (*sw.WebhookResponseModel, error) {
	model, response, err := client.apiClient.WebhooksAPI.GetWebhook(client.GetAuthContext(), webhookId).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) GetWebhookSigningKeys(webhookId int32) (*sw.WebhookSigningKeysModel, error) {
	model, response, err := client.apiClient.WebhooksAPI.GetWebhookSigningKeys(client.GetAuthContext(), webhookId).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) CreateWebhook(configId string, environmentId string, body sw.WebHookRequestModel) (*sw.WebhookResponseModel, error) {
	model, response, err := client.apiClient.WebhooksAPI.CreateWebhook(
		client.GetAuthContext(),
		configId, environmentId).WebHookRequestModel(body).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) UpdateWebhook(webhookId int32, body sw.WebHookRequestModel) (*sw.WebhookResponseModel, error) {
	model, response, err := client.apiClient.WebhooksAPI.ReplaceWebhook(
		client.GetAuthContext(),
		webhookId).WebHookRequestModel(body).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) DeleteWebhook(webhookId int32) error {
	response, err := client.apiClient.WebhooksAPI.DeleteWebhook(
		client.GetAuthContext(),
		webhookId).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return error
}
