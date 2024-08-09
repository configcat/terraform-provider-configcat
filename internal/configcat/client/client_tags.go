package client

import sw "github.com/configcat/configcat-publicapi-go-client/v2"

func (client *Client) GetTags(productID string) ([]sw.TagModel, error) {
	model, response, err := client.apiClient.TagsAPI.GetTags(client.GetAuthContext(), productID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) GetTag(tagID int64) (*sw.TagModel, error) {
	model, response, err := client.apiClient.TagsAPI.GetTag(client.GetAuthContext(), tagID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) CreateTag(productID string, body sw.CreateTagModel) (*sw.TagModel, error) {
	model, response, err := client.apiClient.TagsAPI.CreateTag(
		client.GetAuthContext(),
		productID).CreateTagModel(body).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) UpdateTag(tagID int64, body sw.UpdateTagModel) (*sw.TagModel, error) {
	model, response, err := client.apiClient.TagsAPI.UpdateTag(
		client.GetAuthContext(),
		tagID).UpdateTagModel(body).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) DeleteTag(tagID int64) error {
	response, err := client.apiClient.TagsAPI.DeleteTag(
		client.GetAuthContext(),
		tagID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return error
}
