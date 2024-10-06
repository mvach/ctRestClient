package restendpoints

import (
	"ctRestClient/httpclient"
	"encoding/json"
	"fmt"
	"io"

	"net/http"
	"net/url"
)

type personsEndpoint struct {
    httpclient httpclient.HTTPClient
}

func NewPersonsEndpoint(httpclient httpclient.HTTPClient) *personsEndpoint {
    return &personsEndpoint{
        httpclient: httpclient,
    }
}

func (c *personsEndpoint) GetPerson(personId int) (*PersonResponseJson, error) {

    req, err := http.NewRequest("GET", "", nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request, %w", err)
    }

    params := url.Values{}
    params.Add("ids[]", fmt.Sprintf("%d", personId))
    encodedQueryParam := params.Encode()

    req.URL.Path = fmt.Sprintf("api/persons?%s", encodedQueryParam)


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

    var response PersonResponseJson
    if err := json.Unmarshal(body, &response); err != nil {
        return nil, fmt.Errorf("response body is not containing expected json, %w", err)
    }

    return &response, nil
}
