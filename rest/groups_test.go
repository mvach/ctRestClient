package rest_test

import (
    "bytes"
    "errors"
    "io"
    "net/http"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    "ctRestClient/httpclient/httpclientfakes"
    "ctRestClient/rest"
)

var _ = Describe("GroupsEndpoint", func() {

    var (
        httpClient *httpclientfakes.FakeHTTPClient
    )

    BeforeEach(func() {
        httpClient = &httpclientfakes.FakeHTTPClient{}
    })

    var _ = Describe("GetGroupName", func() {

        It("returns groupname", func() {
            httpResponse := &http.Response{
                StatusCode: 200,
                Body: io.NopCloser(bytes.NewBufferString(
                    `{
                        "data": [
                            {
                                "id": 10,
                                "guid": "1234",
                                "name": "group1"
                            }
                        ],
                        "meta": {
                            "count": 1
                        }
                    }`))}
            httpClient.DoReturns(httpResponse, nil)

            groupsEndpoint := rest.NewGroupsEndpoint(httpClient)
            resp, err := groupsEndpoint.GetGroupName(10)

            Expect(err).NotTo(HaveOccurred())
            Expect(resp[0].Name).To(Equal("group1"))
        })

        It("returns an error if the request cannot be send", func() {
            httpClient.DoReturns(nil, errors.New("request failed"))

            groupsEndpoint := rest.NewGroupsEndpoint(httpClient)
            _, err := groupsEndpoint.GetGroupName(10)

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

            groupsEndpoint := rest.NewGroupsEndpoint(httpClient)
            _, err := groupsEndpoint.GetGroupName(10)

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

            groupsEndpoint := rest.NewGroupsEndpoint(httpClient)
            _, err := groupsEndpoint.GetGroupName(10)

            Expect(err).To(HaveOccurred())
            Expect(err.Error()).To(ContainSubstring("response body is not containing expected json"))
        })
    })

    var _ = Describe("GetGroupMembers", func() {

        It("returns group members", func() {
            httpResponse := &http.Response{
                StatusCode: 200,
                Body: io.NopCloser(bytes.NewBufferString(
                    `{
                        "data": [
                            {
                                "personId": 1,
                                "groupId": 71
                            },{
                                "personId": 2,
                                "groupId": 71
                            }
                        ],
                        "meta": {
                            "count": 1
                        }
                    }`))}
            httpClient.DoReturns(httpResponse, nil)

            groupsEndpoint := rest.NewGroupsEndpoint(httpClient)
            resp, err := groupsEndpoint.GetGroupMembers(10)

            Expect(err).NotTo(HaveOccurred())
            Expect(resp[0].PersonId).To(Equal(1))
            Expect(resp[1].PersonId).To(Equal(2))
        })

        It("returns an error if the request cannot be send", func() {
            httpClient.DoReturns(nil, errors.New("request failed"))

            groupsEndpoint := rest.NewGroupsEndpoint(httpClient)
            _, err := groupsEndpoint.GetGroupMembers(10)

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

            groupsEndpoint := rest.NewGroupsEndpoint(httpClient)
            _, err := groupsEndpoint.GetGroupMembers(10)

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

            groupsEndpoint := rest.NewGroupsEndpoint(httpClient)
            _, err := groupsEndpoint.GetGroupMembers(10)

            Expect(err).To(HaveOccurred())
            Expect(err.Error()).To(ContainSubstring("response body is not containing expected json"))
        })
    })

})
