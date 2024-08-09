package client

import sw "github.com/configcat/configcat-publicapi-go-client/v2"

func (client *Client) GetPermissionGroups(productID string) ([]sw.PermissionGroupModel, error) {
	model, response, err := client.apiClient.PermissionGroupsAPI.GetPermissionGroups(client.GetAuthContext(), productID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) GetPermissionGroup(permissionGroupID int64) (*sw.PermissionGroupModel, error) {
	model, response, err := client.apiClient.PermissionGroupsAPI.GetPermissionGroup(client.GetAuthContext(), permissionGroupID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) CreatePermissionGroup(productID string, body sw.CreatePermissionGroupRequest) (*sw.PermissionGroupModel, error) {
	model, response, err := client.apiClient.PermissionGroupsAPI.CreatePermissionGroup(
		client.GetAuthContext(),
		productID).CreatePermissionGroupRequest(body).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) UpdatePermissionGroup(permissionGroupID int64, body sw.UpdatePermissionGroupRequest) (*sw.PermissionGroupModel, error) {
	model, response, err := client.apiClient.PermissionGroupsAPI.UpdatePermissionGroup(
		client.GetAuthContext(),
		permissionGroupID).UpdatePermissionGroupRequest(body).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) DeletePermissionGroup(permissionGroupID int64) error {
	response, err := client.apiClient.PermissionGroupsAPI.DeletePermissionGroup(
		client.GetAuthContext(),
		permissionGroupID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return error
}
