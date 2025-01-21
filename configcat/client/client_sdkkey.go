package client

import (
	sw "github.com/configcat/configcat-publicapi-go-client/v2"
)

func (client *Client) GetSdkKeys(configId string, environmentId string) (*sw.SdkKeysModel, error) {
	model, response, err := client.apiClient.SDKKeysAPI.GetSdkKeys(client.GetAuthContext(), configId, environmentId).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}
