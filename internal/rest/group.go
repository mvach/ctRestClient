package rest

import (
	"ctRestClient/internal/httpclient"
	"encoding/json"
	"fmt"
	"io"

	"net/http"
)

//counterfeiter:generate . GroupEndpoint
type GroupEndpoint interface {
	GetGroupType(groupTypeID int) (GroupTypeResponse, error)
}

type groupEndpoint struct {
	httpclient httpclient.HTTPClient
}

func NewGroupEndpoint(httpclient httpclient.HTTPClient) GroupEndpoint {
	return groupEndpoint{
		httpclient: httpclient,
	}
}

func (c groupEndpoint) GetGroupType(groupTypeID int) (GroupTypeResponse, error) {

	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		return GroupTypeResponse{}, fmt.Errorf("failed to create request, %w", err)
	}

	req.URL.Path = fmt.Sprintf("/api/group/grouptypes/%d", groupTypeID)

	resp, err := c.httpclient.Do(req)
	if err != nil {
		return GroupTypeResponse{}, fmt.Errorf("failed to send request, %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GroupTypeResponse{}, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GroupTypeResponse{}, fmt.Errorf("failed to read response body, %w", err)
	}

	var response GroupTypeResponseJson
	if err := json.Unmarshal(body, &response); err != nil {
		return GroupTypeResponse{}, fmt.Errorf("response body is not containing expected json, %w", err)
	}

	return response.Data, nil
}
