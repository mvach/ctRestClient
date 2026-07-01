package config_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"ctRestClient/internal/config"
)


var _ = Describe("ExportGroupConfig", func() {
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

	var _ = Describe("Validate", func() {

		It("returns an error if mandatory group name field is missing", func() {
			subject := config.ExportGroupGroup{Name: ""}

			err := subject.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("property name is not set"))
		})

		It("returns an error if mandatory fields field is missing", func() {
			subject := config.ExportGroupGroup{Name: "foo"}

			err := subject.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("property fields is not set"))
		})

		It("returns an error if mandatory fields field is empty array", func() {
			subject := config.ExportGroupGroup{Name: "foo", Fields: []config.ExportGroupField{}}

			err := subject.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("property fields is not set"))
		})
	})

	// It("returns an error if fields array contains invalid object", func() {
	// 	yamlContent := testutil.YamlToByteArray(`
	// 		---
	// 		instances:
	// 		- hostname: foo
	// 		  token_name: foo
	// 		  groups:
	// 		  - name: foo_group_0
	// 		    fields: [{}]
	// 		`)
	// 	_, err := tempFile.Write([]byte(yamlContent))
	// 	Expect(err).ToNot(HaveOccurred())
	// 	tempFile.Close()

	// 	cfg, err := config.LoadConfig(tempFile.Name())
	// 	Expect(err).To(HaveOccurred())
	// 	Expect(err.Error()).To(ContainSubstring("failed to load invalid config file, both 'fieldname' and 'columnname' must be set"))
	// 	Expect(cfg).To(BeNil())
	// })
})