package configcat

import sw "github.com/configcat/configcat-publicapi-go-client"

func (client *Client) GetTags(productID string) ([]sw.TagModel, error) {
	model, response, err := client.apiClient.TagsApi.GetTags(client.GetAuthContext(), productID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) GetTag(tagID int64) (*sw.TagModel, error) {
	model, response, err := client.apiClient.TagsApi.GetTag(client.GetAuthContext(), tagID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) CreateTag(productID string, body sw.CreateTagModel) (*sw.TagModel, error) {
	model, response, err := client.apiClient.TagsApi.CreateTag(
		client.GetAuthContext(),
		productID).CreateTagModel(body).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) UpdateTag(tagID int64, body sw.UpdateTagModel) (*sw.TagModel, error) {
	model, response, err := client.apiClient.TagsApi.UpdateTag(
		client.GetAuthContext(),
		tagID).UpdateTagModel(body).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) DeleteTag(tagID int64) error {
	response, err := client.apiClient.TagsApi.DeleteTag(
		client.GetAuthContext(),
		tagID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return error
}
