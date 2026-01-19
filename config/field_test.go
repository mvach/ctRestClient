package config_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"ctRestClient/config"
	"ctRestClient/testutil"
)

var _ = Describe("Field", func() {
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

	var _ = Describe("GetFieldName", func() {
		It("return the field name from simple field name", func() {
			yamlContent := testutil.YamlToByteArray(`
				---
				instances:
				- hostname: foo
					token_name: foo
					groups:
					- name: foo
						fields:
						- foo_field_1
				`)

			_, err := tempFile.Write(yamlContent)
			Expect(err).ToNot(HaveOccurred())
			tempFile.Close()

			cfg, err := config.LoadConfig(tempFile.Name())
			Expect(err).ToNot(HaveOccurred())
			Expect(cfg).ToNot(BeNil())

			Expect(cfg.Instances[0].Groups[0].Fields[0].GetFieldName()).To(Equal("foo_field_1"))
		})

		It("return the field name from object", func() {
			yamlContent := testutil.YamlToByteArray(`
				---
				instances:
				- hostname: foo
					token_name: foo
					groups:
					- name: foo
						fields:
						- {fieldname: foo_field_1, columnname: foo_column_1}
				`)

			_, err := tempFile.Write(yamlContent)
			Expect(err).ToNot(HaveOccurred())
			tempFile.Close()

			cfg, err := config.LoadConfig(tempFile.Name())
			Expect(err).ToNot(HaveOccurred())
			Expect(cfg).ToNot(BeNil())

			Expect(cfg.Instances[0].Groups[0].Fields[0].GetFieldName()).To(Equal("foo_field_1"))
		})
	})

	var _ = Describe("GetColumnName", func() {
		It("return the simple field name as column name", func() {
			yamlContent := testutil.YamlToByteArray(`
				---
				instances:
				- hostname: foo
					token_name: foo
					groups:
					- name: foo
						fields:
						- foo_field_1
				`)

			_, err := tempFile.Write(yamlContent)
			Expect(err).ToNot(HaveOccurred())
			tempFile.Close()

			cfg, err := config.LoadConfig(tempFile.Name())
			Expect(err).ToNot(HaveOccurred())
			Expect(cfg).ToNot(BeNil())

			Expect(cfg.Instances[0].Groups[0].Fields[0].GetColumnName()).To(Equal("foo_field_1"))
		})

		It("return the columnname from object", func() {
			yamlContent := testutil.YamlToByteArray(`
				---
				instances:
				- hostname: foo
					token_name: foo
					groups:
					- name: foo
						fields:
						- {fieldname: foo_field_1, columnname: foo_column_1}
				`)

			_, err := tempFile.Write(yamlContent)
			Expect(err).ToNot(HaveOccurred())
			tempFile.Close()

			cfg, err := config.LoadConfig(tempFile.Name())
			Expect(err).ToNot(HaveOccurred())
			Expect(cfg).ToNot(BeNil())

			Expect(cfg.Instances[0].Groups[0].Fields[0].GetColumnName()).To(Equal("foo_column_1"))
		})
	})

	var _ = Describe("IsMappedData", func() {
		It("returns false if field name is present", func() {
			yamlContent := testutil.YamlToByteArray(`
				---
				instances:
				- hostname: foo
					token_name: foo
					groups:
					- name: foo
						fields:
						- foo_field_1
				`)

			_, err := tempFile.Write(yamlContent)
			Expect(err).ToNot(HaveOccurred())
			tempFile.Close()

			cfg, err := config.LoadConfig(tempFile.Name())
			Expect(err).ToNot(HaveOccurred())
			Expect(cfg).ToNot(BeNil())

			Expect(cfg.Instances[0].Groups[0].Fields[0].IsMappedData()).To(Equal(false))
		})

		It("returns true if field object is present", func() {
			yamlContent := testutil.YamlToByteArray(`
				---
				instances:
				- hostname: foo
					token_name: foo
					groups:
					- name: foo
						fields:
						- {fieldname: foo_field_1, columnname: foo_column_1}
				`)

			_, err := tempFile.Write([]byte(yamlContent))
			Expect(err).ToNot(HaveOccurred())
			tempFile.Close()

			cfg, err := config.LoadConfig(tempFile.Name())
			Expect(err).ToNot(HaveOccurred())
			Expect(cfg).ToNot(BeNil())

			Expect(cfg.Instances[0].Groups[0].Fields[0].IsMappedData()).To(Equal(true))
		})
	})
})
