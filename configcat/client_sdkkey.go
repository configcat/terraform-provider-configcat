package configcat

import (
	sw "github.com/configcat/configcat-publicapi-go-client"
)

func (client *Client) GetSdkKeys(configId string, environmentId string) (sw.SdkKeysModel, error) {
	model, response, err := client.apiClient.SDKKeysApi.GetSdkKeys(client.GetAuthContext(), configId, environmentId)
	defer response.Body.Close()
	return model, handleAPIError(err)
}
