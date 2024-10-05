package restendpoints

import (
    "ctRestClient/httpclient"
    "encoding/json"
    "fmt"
    "io"

    "net/http"
)


type dynamicGroups struct {
    httpclient httpclient.HTTPClient
}

func NewDynamicGroups(httpclient httpclient.HTTPClient) *dynamicGroups {
    return &dynamicGroups{
        httpclient: httpclient,
    }
}

func (c *dynamicGroups) GetDynamicGroupIds() (*DynamicGroupIdsResponseJson, error) {

    req, err := http.NewRequest("GET", "", nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request, %w", err)
    }

    req.URL.Path = "/api/dynamicgroups"

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

    var response DynamicGroupIdsResponseJson
    if err := json.Unmarshal(body, &response); err != nil {
        return nil, fmt.Errorf("response body is not containing church tools json, %w", err)
    }

    return &response, nil
}
