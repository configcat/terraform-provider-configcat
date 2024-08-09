package client

import (
	sw "github.com/configcat/configcat-publicapi-go-client/v2"
)

func (client *Client) GetSettingValueV2(environmentID string, settingID int32) (*sw.SettingFormulaModel, error) {
	model, response, err := client.apiClient.FeatureFlagSettingValuesV2API.GetSettingValueV2(client.GetAuthContext(), environmentID, settingID).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}

func (client *Client) ReplaceSettingValueV2(environmentID string, settingID int32, body sw.UpdateEvaluationFormulaModel, reason string) (*sw.SettingFormulaModel, error) {
	model, response, err := client.apiClient.FeatureFlagSettingValuesV2API.ReplaceSettingValueV2(client.GetAuthContext(), environmentID, settingID).UpdateEvaluationFormulaModel(body).Reason(reason).Execute()
	error := handleAPIError(err)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return model, error
}
