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
	"ctRestClient/testutil"
)

var _ = Describe("DynamicGroupsEndpoint", func() {

	var (
		httpClient   *httpclientfakes.FakeHTTPClient
		httpResponse *http.Response
	)

	BeforeEach(func() {
		httpClient = &httpclientfakes.FakeHTTPClient{}
	})

	var _ = Describe("GetGroupStatus", func() {

		BeforeEach(func(){
			httpResponse = &http.Response{
			StatusCode: 200,
			Body: io.NopCloser(testutil.JsonToBufferString(
				`{
					"dynamicGroupStatus": "active"
				}`)),
			}
		})

		It("returns a group status", func() {

			httpClient.DoReturns(httpResponse, nil)

			endpoint := rest.NewDynamicGroupsEndpoint(httpClient)
			group, err := endpoint.GetGroupStatus(1)

			Expect(err).NotTo(HaveOccurred())
			Expect(*group.Status).To(Equal("active"))
			request := httpClient.DoArgsForCall(0)
			Expect(request.URL.Path).To(Equal("/api/dynamicgroups/1/status"))
		})

		It("returns an error if the request cannot be send", func() {
			httpClient.DoReturns(nil, errors.New("request failed"))

			endpoint := rest.NewDynamicGroupsEndpoint(httpClient)
			_, err := endpoint.GetGroupStatus(1)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to send request, request failed"))
		})

		It("returns an error if the status code is wrong", func() {
			httpResponse := &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewBufferString(`{}`))}
			httpClient.DoReturns(httpResponse, nil)

			endpoint := rest.NewDynamicGroupsEndpoint(httpClient)
			_, err := endpoint.GetGroupStatus(1)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("received non-200 response code: 404"))
		})

		It("returns an error if the response body is an invalid json response", func() {
			httpResponse := &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(
					``))}
			httpClient.DoReturns(httpResponse, nil)

			endpoint := rest.NewDynamicGroupsEndpoint(httpClient)
			_, err := endpoint.GetGroupStatus(1)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("response body is not containing expected json"))
		})

		It("returns an error if the response body is missing the dynamicGroupStatus field", func() {
			httpResponse := &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(testutil.JsonToBufferString(
					`{
						"foo": []
					}`)),
			}
			httpClient.DoReturns(httpResponse, nil)

			endpoint := rest.NewDynamicGroupsEndpoint(httpClient)
			_, err := endpoint.GetGroupStatus(1)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("response body is missing dynamicGroupStatus field"))
		})
	})

	var _ = Describe("GetAllDynamicGroups", func() {

		BeforeEach(func(){
			httpResponse = &http.Response{
			StatusCode: 200,
			Body: io.NopCloser(testutil.JsonToBufferString(
				`{
					"data": [0,1,2,3]
				}`)),
			}
		})

		It("returns a list of dynamic groups", func() {

			httpClient.DoReturns(httpResponse, nil)

			endpoint := rest.NewDynamicGroupsEndpoint(httpClient)
			groups, err := endpoint.GetAllDynamicGroups()

			Expect(err).NotTo(HaveOccurred())
			Expect(groups.GroupIDs).To(Equal([]int{0, 1, 2, 3}))
			request := httpClient.DoArgsForCall(0)
			Expect(request.URL.Path).To(Equal("/api/dynamicgroups"))
		})

		It("returns an error if the request cannot be send", func() {
			httpClient.DoReturns(nil, errors.New("request failed"))

			endpoint := rest.NewDynamicGroupsEndpoint(httpClient)
			_, err := endpoint.GetAllDynamicGroups()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to send request, request failed"))
		})

		It("returns an error if the status code is wrong", func() {
			httpResponse := &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewBufferString(`{}`))}
			httpClient.DoReturns(httpResponse, nil)

			endpoint := rest.NewDynamicGroupsEndpoint(httpClient)
			_, err := endpoint.GetAllDynamicGroups()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("received non-200 response code: 404"))
		})

		It("returns an error if the response body is an invalid json response", func() {
			httpResponse := &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(
					``))}
			httpClient.DoReturns(httpResponse, nil)

			endpoint := rest.NewDynamicGroupsEndpoint(httpClient)
			_, err := endpoint.GetAllDynamicGroups()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("response body is not containing expected json"))
		})
	})
})
