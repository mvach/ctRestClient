package data_provider_test

import (
	"ctRestClient/data_provider"
	"ctRestClient/testutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FileDataProvider", func() {

	var (
		err            error
		tempDataDir    string
		dp             data_provider.FileDataProvider
		mappedFilePath string
	)

	BeforeEach(func() {
		tempDataDir, err = os.MkdirTemp("", "file_data_provider_test_")
		Expect(err).ToNot(HaveOccurred())

		mappedFilePath = filepath.Join(tempDataDir, "mappedField.yml")

		yamlContent := testutil.YamlToByteArray(`
			1: "number one"
			2: "number two"
			1.1: "number one point one"
			1.0: "number one point zero"
			"2": "string two"
			"A": "string A"
			"B": "string B"
			`)

		err = os.WriteFile(mappedFilePath, []byte(yamlContent), 0644)
		Expect(err).ToNot(HaveOccurred())
		dp = data_provider.NewFileDataProvider(tempDataDir)
	})

	AfterEach(func() {
		os.RemoveAll(tempDataDir)
	})

	var _ = Describe("GetData", func() {
		It("returns mapped data for int keys", func() {
			result, _ := dp.GetData("mappedField", []byte("1"))
			Expect(result).To(Equal("number one"))
			result, _ = dp.GetData("mappedField", []byte("2"))
			Expect(result).To(Equal("number two"))
		})

		It("returns mapped data for float keys", func() {
			result, _ := dp.GetData("mappedField", []byte("1.1"))
			Expect(result).To(Equal("number one point one"))

			result, _ = dp.GetData("mappedField", []byte("1.0"))
			Expect(result).To(Equal("number one point zero"))
		})

		It("returns mapped data for string keys", func() {
			result, _ := dp.GetData("mappedField", []byte("\"2\""))
			Expect(result).To(Equal("string two"))

			result, _ = dp.GetData("mappedField", []byte("A"))
			Expect(result).To(Equal("string A"))

			result, _ = dp.GetData("mappedField", []byte("\"B\""))
			Expect(result).To(Equal("string B"))
		})

		It("returns error for non-existing key", func() {
			_, err := dp.GetData("mappedField", []byte("999"))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("the value 999 is not in '" + mappedFilePath + "'"))
		})
	})
})
