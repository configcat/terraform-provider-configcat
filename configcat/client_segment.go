package configcat

import sw "github.com/configcat/configcat-publicapi-go-client"

func (client *Client) GetSegments(productID string) ([]sw.SegmentListModel, error) {
	model, response, err := client.apiClient.SegmentsApi.GetSegments(client.GetAuthContext(), productID)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) GetSegment(segmentID string) (sw.SegmentModel, error) {
	model, response, err := client.apiClient.SegmentsApi.GetSegment(client.GetAuthContext(), segmentID)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) CreateSegment(productID string, body sw.CreateSegmentModel) (sw.SegmentModel, error) {
	model, response, err := client.apiClient.SegmentsApi.CreateSegment(
		client.GetAuthContext(),
		body,
		productID)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) UpdateSegment(segmentId string, body sw.UpdateSegmentModel) (sw.SegmentModel, error) {
	model, response, err := client.apiClient.SegmentsApi.UpdateSegment(
		client.GetAuthContext(),
		body,
		segmentId)
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) DeleteSegment(segmentId string) error {
	response, err := client.apiClient.SegmentsApi.DeleteSegment(
		client.GetAuthContext(),
		segmentId)
	defer response.Body.Close()
	return handleAPIError(err)
}
