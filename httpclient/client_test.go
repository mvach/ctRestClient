package httpclient_test

import (
    "net/http"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    "ctRestClient/httpclient"
)

var _ = Describe("HTTPClient", func() {

    var _ = Describe("Do", func() {

        It("sets accept headers", func() {
            request, err := http.NewRequest("GET", "", nil)
            Expect(err).NotTo(HaveOccurred())

            client := httpclient.NewHTTPClient("hostname", "token")
            client.Do(request)

            Expect(request.Header["Accept"][0]).To(Equal("application/json"))
        })

        It("sets authorization headers", func() {
            request, err := http.NewRequest("GET", "", nil)
            Expect(err).NotTo(HaveOccurred())

            client := httpclient.NewHTTPClient("hostname", "token")
            client.Do(request)

            Expect(request.Header["Authorization"][0]).To(Equal("Login token"))
        })

        It("uses https", func() {
            request, err := http.NewRequest("GET", "", nil)
            Expect(err).NotTo(HaveOccurred())

            client := httpclient.NewHTTPClient("hostname", "token")
            client.Do(request)

            url := *request.URL
            Expect(url.Scheme).To(Equal("https"))
        })

        It("uses https", func() {
            request, err := http.NewRequest("GET", "", nil)
            Expect(err).NotTo(HaveOccurred())

            client := httpclient.NewHTTPClient("hostname", "token")
            client.Do(request)

            url := *request.URL
            Expect(url.Hostname()).To(Equal("hostname"))
        })
    })
})
