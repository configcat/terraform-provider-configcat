package configcat

import sw "github.com/configcat/configcat-publicapi-go-client"

func (client *Client) GetSegments(productID string) ([]sw.SegmentListModel, error) {
	model, response, err := client.apiClient.SegmentsApi.GetSegments(client.GetAuthContext(), productID).Execute()
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) GetSegment(segmentID string) (*sw.SegmentModel, error) {
	model, response, err := client.apiClient.SegmentsApi.GetSegment(client.GetAuthContext(), segmentID).Execute()
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) CreateSegment(productID string, body sw.CreateSegmentModel) (*sw.SegmentModel, error) {
	model, response, err := client.apiClient.SegmentsApi.CreateSegment(
		client.GetAuthContext(),
		productID).CreateSegmentModel(body).Execute()
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) UpdateSegment(segmentId string, body sw.UpdateSegmentModel) (*sw.SegmentModel, error) {
	model, response, err := client.apiClient.SegmentsApi.UpdateSegment(
		client.GetAuthContext(),
		segmentId).UpdateSegmentModel(body).Execute()
	defer response.Body.Close()
	return model, handleAPIError(err)
}

func (client *Client) DeleteSegment(segmentId string) error {
	response, err := client.apiClient.SegmentsApi.DeleteSegment(
		client.GetAuthContext(),
		segmentId).Execute()
	defer response.Body.Close()
	return handleAPIError(err)
}
