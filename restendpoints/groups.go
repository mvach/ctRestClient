package restendpoints

import (
	"ctRestClient/httpclient"
	"encoding/json"
	"fmt"
	"io"

	"net/http"
	"net/url"
)

type groupsEndpoint struct {
    httpclient httpclient.HTTPClient
}

func NewGroupsEndpoint(httpclient httpclient.HTTPClient) *groupsEndpoint {
    return &groupsEndpoint{
        httpclient: httpclient,
    }
}

func (c *groupsEndpoint) GetGroupName(groupId int) (*GroupsResponseJson, error) {

    req, err := http.NewRequest("GET", "", nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request, %w", err)
    }

    params := url.Values{}
    params.Add("ids[]", fmt.Sprintf("%d", groupId))
    encodedQueryParam := params.Encode()

    req.URL.Path = fmt.Sprintf("/api/groups?%s", encodedQueryParam)


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

    return &response, nil
}
