package rest

import (
	"ctRestClient/httpclient"
	"encoding/json"
	"fmt"
	"io"

	"net/http"
	"net/url"
)

//counterfeiter:generate . GroupsEndpoint
type GroupsEndpoint interface {
    GetGroupNames(groupIDs []int) ([]GroupsResponse, error)

    GetGroupMembers(groupID int) ([]GroupsMembersResponse, error)
}

type groupsEndpoint struct {
    httpclient httpclient.HTTPClient
}

func NewGroupsEndpoint(httpclient httpclient.HTTPClient) GroupsEndpoint {
    return groupsEndpoint{
        httpclient: httpclient,
    }
}

func (c groupsEndpoint) GetGroupNames(groupIDs []int) ([]GroupsResponse, error) {

    req, err := http.NewRequest("GET", "", nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request, %w", err)
    }

    params := url.Values{}
    for _, groupID := range groupIDs {
        params.Add("ids[]", fmt.Sprintf("%d", groupID))
    }
    encodedQueryParam := params.Encode()

    req.URL.Path = "/api/groups"
    req.URL.RawQuery = encodedQueryParam

    resp, err := c.httpclient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send request, %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body, %w", err)
    }

    var response GroupsResponseJson
    if err := json.Unmarshal(body, &response); err != nil {
        return nil, fmt.Errorf("response body is not containing expected json, %w", err)
    }

    return response.Data, nil
}

func (c groupsEndpoint) GetGroupMembers(groupID int) ([]GroupsMembersResponse, error) {

    req, err := http.NewRequest("GET", "", nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request, %w", err)
    }

    params := url.Values{}
    params.Add("ids[]", fmt.Sprintf("%d", groupID))
    params.Add("with_deleted", "false")
    encodedQueryParam := params.Encode()

    req.URL.Path = "/api/groups/members"
    req.URL.RawQuery = encodedQueryParam

    resp, err := c.httpclient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send request, %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body, %w", err)
    }

    var response GroupsMembersResponseJson
    if err := json.Unmarshal(body, &response); err != nil {
        return nil, fmt.Errorf("response body is not containing expected json, %w", err)
    }

    return response.Data, nil
}
