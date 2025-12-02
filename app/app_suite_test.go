package app_test

import (
    "testing"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

func TestClient(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "App Suite")
}


func ptr(s string) *string {
	return &s
}