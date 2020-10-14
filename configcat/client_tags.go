package configcat

import sw "github.com/configcat/configcat-publicapi-go-client"

func (client *Client) GetTags(productId string) ([]sw.TagModel, error) {
	model, response, err := client.apiClient.TagsApi.GetTags(client.GetAuthContext(), productId)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) GetTag(tagId int64) (sw.TagModel, error) {
	model, response, err := client.apiClient.TagsApi.GetTag(client.GetAuthContext(), tagId)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) CreateTag(productId string, body sw.CreateTagModel) (sw.TagModel, error) {
	model, response, err := client.apiClient.TagsApi.CreateTag(
		client.GetAuthContext(),
		body,
		productId)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) UpdateTag(tagID int64, body sw.UpdateTagModel) (sw.TagModel, error) {
	model, response, err := client.apiClient.TagsApi.UpdateTag(
		client.GetAuthContext(),
		body,
		tagID)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) DeleteTag(tagID int64) error {
	response, err := client.apiClient.TagsApi.DeleteTag(
		client.GetAuthContext(),
		tagID)
	defer response.Body.Close()
	return handleAPIError(err)
}
