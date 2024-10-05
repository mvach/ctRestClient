package restendpoints_test

import (
    "bytes"
    "errors"
    "io"
    "net/http"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    "ctRestClient/httpclient/httpclientfakes"
    "ctRestClient/restendpoints"
)

var _ = Describe("DynamicGroups", func() {

    var (
        httpClient *httpclientfakes.FakeHTTPClient
    )

    BeforeEach(func() {
        httpClient = &httpclientfakes.FakeHTTPClient{}
    })

    var _ = Describe("GetDynamicGroupIds", func() {

        It("returns dynamic groups", func() {
            httpResponse := &http.Response{
                StatusCode: 200,
                Body: io.NopCloser(bytes.NewBufferString(
                    `{
                        "data": [
                            10,
                            11,
                            12
                        ],
                        "meta": {
                            "count": 3
                        }
                    }`))}
            httpClient.DoReturns(httpResponse, nil)

            dynamicGroupClient := restendpoints.NewDynamicGroups(httpClient)
            resp, err := dynamicGroupClient.GetDynamicGroupIds()

            Expect(err).NotTo(HaveOccurred())
            Expect(resp.Data).To(Equal([]int{10, 11, 12}))
        })

        It("returns an error if the request cannot be send", func() {
            httpClient.DoReturns(nil, errors.New("request failed"))

            dynamicGroupClient := restendpoints.NewDynamicGroups(httpClient)
            _, err := dynamicGroupClient.GetDynamicGroupIds()

            Expect(err).To(HaveOccurred())
            Expect(err.Error()).To(Equal("failed to send request, request failed"))
        })

        It("returns an error if the status code is wrong", func() {
            httpResponse := &http.Response{
                StatusCode: 404,
                Body: io.NopCloser(bytes.NewBufferString(
                    `{
                        "data": [],
                        "meta": { "count": 0 }
                    }`))}
            httpClient.DoReturns(httpResponse, nil)

            dynamicGroupClient := restendpoints.NewDynamicGroups(httpClient)
            _, err := dynamicGroupClient.GetDynamicGroupIds()

            Expect(err).To(HaveOccurred())
            Expect(err.Error()).To(Equal("received non-200 response code: 404"))
        })

        It("returns an error if the response body is not a church tools json response", func() {
            httpResponse := &http.Response{
                StatusCode: 200,
                Body: io.NopCloser(bytes.NewBufferString(
                    `{
                        "foo": [],
                    }`))}
            httpClient.DoReturns(httpResponse, nil)

            dynamicGroupClient := restendpoints.NewDynamicGroups(httpClient)
            _, err := dynamicGroupClient.GetDynamicGroupIds()

            Expect(err).To(HaveOccurred())
            Expect(err.Error()).To(ContainSubstring("response body is not containing church tools json"))
        })
    })
})
