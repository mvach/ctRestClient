package data_provider_test

import (
	"ctRestClient/config"
	"ctRestClient/data_provider"
	"ctRestClient/logger/loggerfakes"
	"ctRestClient/testutil"
	"encoding/json"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BlocklistDataProvider", func() {

	var (
		err         error
		tempDataDir string
		dp          data_provider.BlockListDataProvider
		logger      *loggerfakes.FakeLogger
		personJson  map[string]json.RawMessage
		group       config.Group
	)

	BeforeEach(func() {
		tempDataDir, err = os.MkdirTemp("", "blocklist_data_provider_test_")
		Expect(err).ToNot(HaveOccurred())

		personJson = map[string]json.RawMessage{
			"street": json.RawMessage(`"Mainstreet"`),
			"zip":    json.RawMessage(`"12345"`),
			"city":   json.RawMessage(`"Anytown"`),
			"age":    json.RawMessage(`30`),
			"isDead": json.RawMessage(`false`),
			"sexId":  json.RawMessage(`1`),
			"weddingDate":  json.RawMessage(`null`),
		}
		logger = &loggerfakes.FakeLogger{}
		dp = data_provider.NewBlockListDataProvider(tempDataDir, logger)

		group = config.Group{Name: "mappedField"}
	})

	AfterEach(func() {
		os.RemoveAll(tempDataDir)
	})

	var _ = Describe("IsBlocked", func() {
		It("returns false if blocklist is not existing", func() {

			result, err := dp.IsBlocked(personJson, config.Group{Name: "not_existing_blocklist"})
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(false))
		})

		It("returns false if blocklist is empty", func() {
			blocklistFilePath := filepath.Join(tempDataDir, "mappedField.yml")
			err = os.WriteFile(blocklistFilePath, []byte(``), 0644)
			Expect(err).ToNot(HaveOccurred())

			result, err := dp.IsBlocked(personJson, group)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(false))
		})

		It("returns false if blocklist has no blocked addresses", func() {
			blocklistFilePath := filepath.Join(tempDataDir, "mappedField.yml")
			err = os.WriteFile(blocklistFilePath, []byte(`---`), 0644)
			Expect(err).ToNot(HaveOccurred())

			result, err := dp.IsBlocked(personJson, group)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(false))
		})

		It("returns false if blocklist is not matching the person json - zipcode", func() {
			blocklistFilePath := filepath.Join(tempDataDir, "mappedField.yml")
			yamlContent := testutil.YamlToByteArray(`
				---
				- zip: ""
				`)
			err = os.WriteFile(blocklistFilePath, []byte(yamlContent), 0644)
			Expect(err).ToNot(HaveOccurred())

			result, err := dp.IsBlocked(personJson, group)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(false))
		})

		It("returns false if blocklist is not matching the person json - age", func() {
			blocklistFilePath := filepath.Join(tempDataDir, "mappedField.yml")
			yamlContent := testutil.YamlToByteArray(`
				---
				- age: ""
				`)
			err = os.WriteFile(blocklistFilePath, []byte(yamlContent), 0644)
			Expect(err).ToNot(HaveOccurred())

			result, err := dp.IsBlocked(personJson, group)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(false))
		})

		It("returns false if blocklist is not matching the person json - isDead", func() {
			blocklistFilePath := filepath.Join(tempDataDir, "mappedField.yml")
			yamlContent := testutil.YamlToByteArray(`
				---
				- isDead: ""
				`)
			err = os.WriteFile(blocklistFilePath, []byte(yamlContent), 0644)
			Expect(err).ToNot(HaveOccurred())

			result, err := dp.IsBlocked(personJson, group)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(false))
		})

		It("returns false if blocklist is not matching the person json - sexId", func() {
			blocklistFilePath := filepath.Join(tempDataDir, "mappedField.yml")
			yamlContent := testutil.YamlToByteArray(`
				---
				- sexId: ""
				`)
			err = os.WriteFile(blocklistFilePath, []byte(yamlContent), 0644)
			Expect(err).ToNot(HaveOccurred())

			result, err := dp.IsBlocked(personJson, group)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(false))
		})

		It("returns false if blocklist is not matching the person json - weddingDate", func() {
			blocklistFilePath := filepath.Join(tempDataDir, "mappedField.yml")
			yamlContent := testutil.YamlToByteArray(`
				---
				- weddingDate: ""
				`)
			err = os.WriteFile(blocklistFilePath, []byte(yamlContent), 0644)
			Expect(err).ToNot(HaveOccurred())

			result, err := dp.IsBlocked(personJson, group)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(false))
		})

		It("returns false if blocklist is not fully matching the person json", func() {
			blocklistFilePath := filepath.Join(tempDataDir, "mappedField.yml")
			yamlContent := testutil.YamlToByteArray(`
				---
				- street: "Mainstreet"
				  city: "Anytown"
				  zip: "9999"
				`)
			err = os.WriteFile(blocklistFilePath, []byte(yamlContent), 0644)
			Expect(err).ToNot(HaveOccurred())

			result, err := dp.IsBlocked(personJson, group)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(false))
		})

		It("returns true if blocklist is matching the person json", func() {
			blocklistFilePath := filepath.Join(tempDataDir, "mappedField.yml")
			yamlContent := testutil.YamlToByteArray(`
				---
				- street: "Mainstreet"
				city: "Anytown"
				zip: "12345"
				`)
			err = os.WriteFile(blocklistFilePath, []byte(yamlContent), 0644)
			Expect(err).ToNot(HaveOccurred())

			result, err := dp.IsBlocked(personJson, group)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(true))
		})

		It("returns true if blocklist is matching the person json (comparing all data)", func() {
			blocklistFilePath := filepath.Join(tempDataDir, "mappedField.yml")
			yamlContent := testutil.YamlToByteArray(`
				---
				- street: "Mainstreet"
				city: "Anytown"
				zip: "12345"
				age: 30
				isDead: false
				sexId: 1
				weddingDate: null
				`)
			err = os.WriteFile(blocklistFilePath, []byte(yamlContent), 0644)
			Expect(err).ToNot(HaveOccurred())

			result, err := dp.IsBlocked(personJson, group)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(true))
		})
	})
})
