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

		It("returns an error if mandatory hostname field is missing", func() {
			subject := config.ExportGroupInstance{Hostname: ""}

			err := subject.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("property hostname is not set"))
		})

		It("returns an error if mandatory token_name field is missing", func() {
			subject := config.ExportGroupInstance{Hostname: "foo"}

			err := subject.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("property token_name is not set"))
		})

		It("returns an error if mandatory groups field is missing", func() {
				subject := config.ExportGroupInstance{Hostname: "foo", TokenName: "bar"}

				err := subject.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to validate the config file, property groups is not set"))
			})

			It("returns an error if mandatory groups field is nil", func() {
				subject := config.ExportGroupInstance{Hostname: "foo", TokenName: "bar", Groups: nil}
				
				err := subject.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to validate the config file, property groups is not set"))
			})

			It("returns an error if mandatory groups field is empty array", func() {
				subject := config.ExportGroupInstance{Hostname: "foo", TokenName: "bar", Groups: []config.ExportGroupGroup{}}
				
				err := subject.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to validate the config file, property groups is not set"))
			})
	})

		// var _ = Describe("group name property errors", func() {
		// 	It("returns an error if mandatory group name field is missing", func() {
		// 		yamlContent := testutil.YamlToByteArray(`
		// 			---
		// 			instances:
		// 			- hostname: foo
		// 			  token_name: foo
		// 			  groups:
		// 			  - foo: bar
		// 			`)
		// 		_, err := tempFile.Write([]byte(yamlContent))
		// 		Expect(err).ToNot(HaveOccurred())
		// 		tempFile.Close()

		// 		cfg, err := config.LoadConfig(tempFile.Name())
		// 		Expect(err).To(HaveOccurred())
		// 		Expect(err.Error()).To(Equal("failed to validate the config file, property name is not set"))
		// 		Expect(cfg).To(BeNil())
		// 	})

		// 	It("returns an error if mandatory group name field is nil", func() {
		// 		yamlContent := testutil.YamlToByteArray(`
		// 			---
		// 			instances:
		// 			- hostname: foo
		// 			  token_name: foo
		// 			  groups:
		// 			  - name:
		// 			`)
		// 		_, err := tempFile.Write([]byte(yamlContent))
		// 		Expect(err).ToNot(HaveOccurred())
		// 		tempFile.Close()

		// 		cfg, err := config.LoadConfig(tempFile.Name())
		// 		Expect(err).To(HaveOccurred())
		// 		Expect(err.Error()).To(Equal("failed to validate the config file, property name is not set"))
		// 		Expect(cfg).To(BeNil())
		// 	})
		// })

		// var _ = Describe("fields property errors", func() {

		// 	It("returns an error if mandatory fields field is missing", func() {
		// 		yamlContent := testutil.YamlToByteArray(`
		// 			---
		// 			instances:
		// 			- hostname: foo
		// 			  token_name: foo
		// 			  groups:
		// 			  - name: foo_group_0
		// 			`)
		// 		_, err := tempFile.Write([]byte(yamlContent))
		// 		Expect(err).ToNot(HaveOccurred())
		// 		tempFile.Close()

		// 		cfg, err := config.LoadConfig(tempFile.Name())
		// 		Expect(err).To(HaveOccurred())
		// 		Expect(err.Error()).To(Equal("failed to validate the config file, property fields is not set"))
		// 		Expect(cfg).To(BeNil())
		// 	})

		// 	It("returns an error if mandatory fields field is nil", func() {
		// 		yamlContent := testutil.YamlToByteArray(`
		// 			---
		// 			instances:
		// 			- hostname: foo
		// 			  token_name: foo
		// 			  groups:
		// 			  - name: foo_group_0
		// 			    fields:
		// 			`)
		// 		_, err := tempFile.Write([]byte(yamlContent))
		// 		Expect(err).ToNot(HaveOccurred())
		// 		tempFile.Close()

		// 		cfg, err := config.LoadConfig(tempFile.Name())
		// 		Expect(err).To(HaveOccurred())
		// 		Expect(err.Error()).To(Equal("failed to validate the config file, property fields is not set"))
		// 		Expect(cfg).To(BeNil())
		// 	})

		// 	It("returns an error if mandatory fields field has wrong type", func() {
		// 		yamlContent := testutil.YamlToByteArray(`
		// 			---
		// 			instances:
		// 			- hostname: foo
		// 			  token_name: foo
		// 			  groups:
		// 			  - name: foo_group_0
		// 			    fields: "not_array"
		// 			`)
		// 		_, err := tempFile.Write([]byte(yamlContent))
		// 		Expect(err).ToNot(HaveOccurred())
		// 		tempFile.Close()

		// 		cfg, err := config.LoadConfig(tempFile.Name())
		// 		Expect(err).To(HaveOccurred())
		// 		Expect(err.Error()).To(ContainSubstring("failed to load invalid config file"))
		// 		Expect(cfg).To(BeNil())
		// 	})

		// 	It("returns an error if mandatory fields field is empty array", func() {
		// 		yamlContent := testutil.YamlToByteArray(`
		// 			---
		// 			instances:
		// 			- hostname: foo
		// 			  token_name: foo
		// 			  groups:
		// 			  - name: foo_group_0
		// 			    fields: []
		// 			`)
		// 		_, err := tempFile.Write([]byte(yamlContent))
		// 		Expect(err).ToNot(HaveOccurred())
		// 		tempFile.Close()

		// 		cfg, err := config.LoadConfig(tempFile.Name())
		// 		Expect(err).To(HaveOccurred())
		// 		Expect(err.Error()).To(ContainSubstring("failed to validate the config file, property fields is not set"))
		// 		Expect(cfg).To(BeNil())
		// 	})

		// 	It("returns an error if fields array contains invalid object", func() {
		// 		yamlContent := testutil.YamlToByteArray(`
		// 			---
		// 			instances:
		// 			- hostname: foo
		// 			  token_name: foo
		// 			  groups:
		// 			  - name: foo_group_0
		// 			    fields: [{}]
		// 			`)
		// 		_, err := tempFile.Write([]byte(yamlContent))
		// 		Expect(err).ToNot(HaveOccurred())
		// 		tempFile.Close()

		// 		cfg, err := config.LoadConfig(tempFile.Name())
		// 		Expect(err).To(HaveOccurred())
		// 		Expect(err.Error()).To(ContainSubstring("failed to load invalid config file, both 'fieldname' and 'columnname' must be set"))
		// 		Expect(cfg).To(BeNil())
		// 	})
		// })
})