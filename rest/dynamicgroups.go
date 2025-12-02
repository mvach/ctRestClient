package rest

import (
	"ctRestClient/httpclient"
	"encoding/json"
	"fmt"
	"io"

	"net/http"
)

//counterfeiter:generate . DynamicGroupsEndpoint
type DynamicGroupsEndpoint interface {
	GetGroupStatus(groupID int) (DynamicGroupsStatusResponse, error)
}

type dynamicGroupsEndpoint struct {
	httpclient httpclient.HTTPClient
}

func NewDynamicGroupsEndpoint(httpclient httpclient.HTTPClient) DynamicGroupsEndpoint {
	return dynamicGroupsEndpoint{
		httpclient: httpclient,
	}
}

func (c dynamicGroupsEndpoint) GetGroupStatus(groupID int) (DynamicGroupsStatusResponse, error) {

	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		return DynamicGroupsStatusResponse{}, fmt.Errorf("failed to create request, %w", err)
	}

	req.URL.Path = fmt.Sprintf("/api/dynamicgroups/%d/status", groupID)

	resp, err := c.httpclient.Do(req)
	if err != nil {
		return DynamicGroupsStatusResponse{}, fmt.Errorf("failed to send request, %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return DynamicGroupsStatusResponse{}, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return DynamicGroupsStatusResponse{}, fmt.Errorf("failed to read response body, %w", err)
	}

	var response DynamicGroupsStatusResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return DynamicGroupsStatusResponse{}, fmt.Errorf("response body is not containing expected json, %w", err)
	}
	if response.Status == nil {
		return DynamicGroupsStatusResponse{}, fmt.Errorf("response body is missing dynamicGroupStatus field")
	}

	return response, nil
}
