package configcat

import sw "github.com/configcat/configcat-publicapi-go-client"

func (client *Client) GetPermissionGroups(productID string) ([]sw.PermissionGroupModel, error) {
	model, response, err := client.apiClient.PermissionGroupsApi.GetPermissionGroups(client.GetAuthContext(), productID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) GetPermissionGroup(permissionGroupID int64) (*sw.PermissionGroupModel, error) {
	model, response, err := client.apiClient.PermissionGroupsApi.GetPermissionGroup(client.GetAuthContext(), permissionGroupID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) CreatePermissionGroup(productID string, body sw.CreatePermissionGroupRequest) (*sw.PermissionGroupModel, error) {
	model, response, err := client.apiClient.PermissionGroupsApi.CreatePermissionGroup(
		client.GetAuthContext(),
		productID).CreatePermissionGroupRequest(body).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) UpdatePermissionGroup(permissionGroupID int64, body sw.UpdatePermissionGroupRequest) (*sw.PermissionGroupModel, error) {
	model, response, err := client.apiClient.PermissionGroupsApi.UpdatePermissionGroup(
		client.GetAuthContext(),
		permissionGroupID).UpdatePermissionGroupRequest(body).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) DeletePermissionGroup(permissionGroupID int64) error {
	response, err := client.apiClient.PermissionGroupsApi.DeletePermissionGroup(
		client.GetAuthContext(),
		permissionGroupID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return error
}
