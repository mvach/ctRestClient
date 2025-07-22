package config_test

import (
    "os"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    "ctRestClient/config"
)

func ptr(s string) *string {
    return &s
}

var _ = Describe("Config", func() {
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

    var _ = Describe("LoadConfig", func() {

        It("should load the configuration", func() {
            yamlContent := `
---
instances:
  - hostname: foo.com
    token_name: foo
    groups:
    - name: foo_group_0
      fields:
      - foo_field_1
      - foo_field_2
    - name: foo_group_1
      fields:
      - {fieldname: foo_field_3, columnname: foo_column_3}
  - hostname: bar.com
    token_name: bar
    groups:
    - name: bar_group_0
      fields:
      - bar_field_1
`
            _, err := tempFile.Write([]byte(yamlContent))
            Expect(err).ToNot(HaveOccurred())
            tempFile.Close()

            cfg, err := config.LoadConfig(tempFile.Name())
            Expect(err).ToNot(HaveOccurred())
            Expect(cfg.Instances).To(HaveLen(2))

            Expect(cfg.Instances[0].Hostname).To(Equal("foo.com"))
            Expect(cfg.Instances[0].TokenName).To(Equal("foo"))
            Expect(cfg.Instances[0].Groups).To(HaveLen(2))
            Expect(cfg.Instances[0].Groups[0].Name).To(Equal("foo_group_0"))
            Expect(cfg.Instances[0].Groups[0].Fields).To(Equal([]config.Field{{RawString:ptr("foo_field_1")}, {RawString:ptr("foo_field_2")}}))
            Expect(cfg.Instances[0].Groups[1].Name).To(Equal("foo_group_1"))
            Expect(cfg.Instances[0].Groups[1].Fields).To(Equal([]config.Field{{RawObject: &config.FieldInformation{FieldName: "foo_field_3", ColumnName: "foo_column_3"}}})) 

            Expect(cfg.Instances[1].Hostname).To(Equal("bar.com"))
            Expect(cfg.Instances[1].TokenName).To(Equal("bar"))
            Expect(cfg.Instances[1].Groups).To(HaveLen(1))
            Expect(cfg.Instances[1].Groups[0].Name).To(Equal("bar_group_0"))
            Expect(cfg.Instances[1].Groups[0].Fields).To(Equal([]config.Field{{RawString:ptr("bar_field_1")}}))
        })

        var _ = Describe("instances property errors", func() {
            It("returns an error if mandatory instances field is missing", func() {
                yamlContent := `
--- 
`
                _, err := tempFile.Write([]byte(yamlContent))
                Expect(err).ToNot(HaveOccurred())
                tempFile.Close()

                cfg, err := config.LoadConfig(tempFile.Name())
                Expect(err).To(HaveOccurred())
                Expect(err.Error()).To(Equal("failed to validate the config file, property instances is not set"))
                Expect(cfg).To(BeNil())
            })

            It("returns an error if mandatory instances field is nil", func() {
                yamlContent := `
---
instances:
`
                _, err := tempFile.Write([]byte(yamlContent))
                Expect(err).ToNot(HaveOccurred())
                tempFile.Close()

                cfg, err := config.LoadConfig(tempFile.Name())
                Expect(err).To(HaveOccurred())
                Expect(err.Error()).To(Equal("failed to validate the config file, property instances is not set"))
                Expect(cfg).To(BeNil())
            })

            It("returns an error if mandatory instances field has wrong type", func() {
                yamlContent := `
---
instances: "not_array"
`
                _, err := tempFile.Write([]byte(yamlContent))
                Expect(err).ToNot(HaveOccurred())
                tempFile.Close()

                cfg, err := config.LoadConfig(tempFile.Name())
                Expect(err).To(HaveOccurred())
                Expect(err.Error()).To(ContainSubstring("failed to load invalid config file"))
                Expect(cfg).To(BeNil())
            })

            It("returns an error if mandatory instances field is empty array", func() {
                yamlContent := `
---
instances: []
`
                _, err := tempFile.Write([]byte(yamlContent))
                Expect(err).ToNot(HaveOccurred())
                tempFile.Close()

                cfg, err := config.LoadConfig(tempFile.Name())
                Expect(err).To(HaveOccurred())
                Expect(err.Error()).To(ContainSubstring("failed to validate the config file, property instances is not set"))
                Expect(cfg).To(BeNil())
            })
        })

        var _ = Describe("hostname property errors", func() {
            It("returns an error if mandatory hostname field is missing", func() {
                yamlContent := `
---
instances:
  - foo: bar
`
                _, err := tempFile.Write([]byte(yamlContent))
                Expect(err).ToNot(HaveOccurred())
                tempFile.Close()

                cfg, err := config.LoadConfig(tempFile.Name())
                Expect(err).To(HaveOccurred())
                Expect(err.Error()).To(Equal("failed to validate the config file, property hostname is not set"))
                Expect(cfg).To(BeNil())
            })

            It("returns an error if mandatory hostname field is nil", func() {
                yamlContent := `
---
instances:
  - hostname:
`
                _, err := tempFile.Write([]byte(yamlContent))
                Expect(err).ToNot(HaveOccurred())
                tempFile.Close()

                cfg, err := config.LoadConfig(tempFile.Name())
                Expect(err).To(HaveOccurred())
                Expect(err.Error()).To(Equal("failed to validate the config file, property hostname is not set"))
                Expect(cfg).To(BeNil())
            })
        })

        var _ = Describe("token_name property errors", func() {
            It("returns an error if mandatory token_name field is missing", func() {
                yamlContent := `
---
instances:
  - hostname: foo
    foo: bar
`
                _, err := tempFile.Write([]byte(yamlContent))
                Expect(err).ToNot(HaveOccurred())
                tempFile.Close()

                cfg, err := config.LoadConfig(tempFile.Name())
                Expect(err).To(HaveOccurred())
                Expect(err.Error()).To(Equal("failed to validate the config file, property token_name is not set"))
                Expect(cfg).To(BeNil())
            })

            It("returns an error if mandatory token_name field is nil", func() {
                yamlContent := `
---
instances:
  - hostname: foo
    token_name:
`
                _, err := tempFile.Write([]byte(yamlContent))
                Expect(err).ToNot(HaveOccurred())
                tempFile.Close()

                cfg, err := config.LoadConfig(tempFile.Name())
                Expect(err).To(HaveOccurred())
                Expect(err.Error()).To(Equal("failed to validate the config file, property token_name is not set"))
                Expect(cfg).To(BeNil())
            })
        })

        var _ = Describe("groups property errors", func() {
            It("returns an error if mandatory groups field is missing", func() {
                yamlContent := `
---
instances:
  - hostname: foo
    token_name: foo
`
                _, err := tempFile.Write([]byte(yamlContent))
                Expect(err).ToNot(HaveOccurred())
                tempFile.Close()

                cfg, err := config.LoadConfig(tempFile.Name())
                Expect(err).To(HaveOccurred())
                Expect(err.Error()).To(Equal("failed to validate the config file, property groups is not set"))
                Expect(cfg).To(BeNil())
            })

            It("returns an error if mandatory groups field is nil", func() {
                yamlContent := `
---
instances:
  - hostname: foo
    token_name: foo
    groups:
`
                _, err := tempFile.Write([]byte(yamlContent))
                Expect(err).ToNot(HaveOccurred())
                tempFile.Close()

                cfg, err := config.LoadConfig(tempFile.Name())
                Expect(err).To(HaveOccurred())
                Expect(err.Error()).To(Equal("failed to validate the config file, property groups is not set"))
                Expect(cfg).To(BeNil())
            })

            It("returns an error if mandatory groups field has wrong type", func() {
                yamlContent := `
---
instances:
  - hostname: foo
    token_name: foo
    groups: "not_array"
`
                _, err := tempFile.Write([]byte(yamlContent))
                Expect(err).ToNot(HaveOccurred())
                tempFile.Close()

                cfg, err := config.LoadConfig(tempFile.Name())
                Expect(err).To(HaveOccurred())
                Expect(err.Error()).To(ContainSubstring("failed to load invalid config file"))
                Expect(cfg).To(BeNil())
            })

            It("returns an error if mandatory groups field is empty array", func() {
                yamlContent := `
---
instances:
  - hostname: foo
    token_name: foo
    groups: []
`
                _, err := tempFile.Write([]byte(yamlContent))
                Expect(err).ToNot(HaveOccurred())
                tempFile.Close()

                cfg, err := config.LoadConfig(tempFile.Name())
                Expect(err).To(HaveOccurred())
                Expect(err.Error()).To(ContainSubstring("failed to validate the config file, property groups is not set"))
                Expect(cfg).To(BeNil())
            })
        })

        var _ = Describe("group name property errors", func() {
            It("returns an error if mandatory group name field is missing", func() {
                yamlContent := `
---
instances:
  - hostname: foo
    token_name: foo
    groups: 
      - foo: bar
`
                _, err := tempFile.Write([]byte(yamlContent))
                Expect(err).ToNot(HaveOccurred())
                tempFile.Close()

                cfg, err := config.LoadConfig(tempFile.Name())
                Expect(err).To(HaveOccurred())
                Expect(err.Error()).To(Equal("failed to validate the config file, property name is not set"))
                Expect(cfg).To(BeNil())
            })

            It("returns an error if mandatory group name field is nil", func() {
                yamlContent := `
---
instances:
  - hostname: foo
    token_name: foo
    groups: 
      - name:
`
                _, err := tempFile.Write([]byte(yamlContent))
                Expect(err).ToNot(HaveOccurred())
                tempFile.Close()

                cfg, err := config.LoadConfig(tempFile.Name())
                Expect(err).To(HaveOccurred())
                Expect(err.Error()).To(Equal("failed to validate the config file, property name is not set"))
                Expect(cfg).To(BeNil())
            })
        })

        var _ = Describe("fields property errors", func() {
            It("returns an error if mandatory fields field is missing", func() {
                yamlContent := `
---
instances:
  - hostname: foo
    token_name: foo
    groups:
      - name: foo_group_0
`
                _, err := tempFile.Write([]byte(yamlContent))
                Expect(err).ToNot(HaveOccurred())
                tempFile.Close()

                cfg, err := config.LoadConfig(tempFile.Name())
                Expect(err).To(HaveOccurred())
                Expect(err.Error()).To(Equal("failed to validate the config file, property fields is not set"))
                Expect(cfg).To(BeNil())
            })

            It("returns an error if mandatory fields field is nil", func() {
                yamlContent := `
---
instances:
  - hostname: foo
    token_name: foo
    groups:
      - name: foo_group_0
        fields:
`
                _, err := tempFile.Write([]byte(yamlContent))
                Expect(err).ToNot(HaveOccurred())
                tempFile.Close()

                cfg, err := config.LoadConfig(tempFile.Name())
                Expect(err).To(HaveOccurred())
                Expect(err.Error()).To(Equal("failed to validate the config file, property fields is not set"))
                Expect(cfg).To(BeNil())
            })

            It("returns an error if mandatory fields field has wrong type", func() {
                yamlContent := `
---
instances:
  - hostname: foo
    token_name: foo
    groups:
      - name: foo_group_0
        fields: "not_array"
`
                _, err := tempFile.Write([]byte(yamlContent))
                Expect(err).ToNot(HaveOccurred())
                tempFile.Close()

                cfg, err := config.LoadConfig(tempFile.Name())
                Expect(err).To(HaveOccurred())
                Expect(err.Error()).To(ContainSubstring("failed to load invalid config file"))
                Expect(cfg).To(BeNil())
            })

            It("returns an error if mandatory fields field is empty array", func() {
                yamlContent := `
---
instances:
  - hostname: foo
    token_name: foo
    groups:
      - name: foo_group_0
        fields: []
`
                _, err := tempFile.Write([]byte(yamlContent))
                Expect(err).ToNot(HaveOccurred())
                tempFile.Close()

                cfg, err := config.LoadConfig(tempFile.Name())
                Expect(err).To(HaveOccurred())
                Expect(err.Error()).To(ContainSubstring("failed to validate the config file, property fields is not set"))
                Expect(cfg).To(BeNil())
            })

            It("returns an error if fields array contains invalid object", func() {
                yamlContent := `
---
instances:
  - hostname: foo
    token_name: foo
    groups:
      - name: foo_group_0
        fields: [{}]
`
                _, err := tempFile.Write([]byte(yamlContent))
                Expect(err).ToNot(HaveOccurred())
                tempFile.Close()

                cfg, err := config.LoadConfig(tempFile.Name())
                Expect(err).To(HaveOccurred())
                Expect(err.Error()).To(ContainSubstring("failed to load invalid config file, both 'fieldname' and 'columnname' must be set"))
                Expect(cfg).To(BeNil())
            })
        })
    })
})
