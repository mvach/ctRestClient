package config_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"ctRestClient/config"
)

var _ = Describe("Group", func() {
	var (
		tempFile *os.File
		err      error
	)

	BeforeEach(func() {
		tempFile, err = os.CreateTemp("", "config.yml")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		os.Remove(tempFile.Name())
	})

	var _ = Describe("SanitizedGroupName", func() {
		It("sanitizes group names correctly", func() {
			yamlContent := `
---
instances:
  - hostname: foo
    token_name: foo
    groups:
      - name: foo- ,äöüÄÖÜ-group
        fields:
        - foo_field_1
`
			_, err := tempFile.Write([]byte(yamlContent))
			Expect(err).ToNot(HaveOccurred())
			tempFile.Close()

			cfg, err := config.LoadConfig(tempFile.Name())
			Expect(err).ToNot(HaveOccurred())
			Expect(cfg).ToNot(BeNil())

			Expect(cfg.Instances[0].Groups[0].SanitizedGroupName()).To(Equal("foo-_.aeoeueAeOeUe-group"))
		})
	})
})
