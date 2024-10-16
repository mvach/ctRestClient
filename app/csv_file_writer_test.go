package app_test

import (
	"ctRestClient/app"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CSVFileWriter", func() {

	var (
		csvHeader  []string
		csvRecords [][]string
	)

	BeforeEach(func() {
		csvHeader = []string{"FirstName", "LastName", "Email"}
		csvRecords = [][]string{
			{"John", "Doe", "john.doe@example.com"},
			{"Jane", "Smith", "jane.smith@example.com"},
		}
	})

	var _ = Describe("Write", func() {
		It("writes a csv file", func() {
			tmpfile, err := os.CreateTemp("", "test.csv")
			Expect(err).ToNot(HaveOccurred())

			defer os.Remove(tmpfile.Name())

			err = app.NewCSVFileWriter().Write(tmpfile.Name(), csvHeader, csvRecords)
			Expect(err).ToNot(HaveOccurred())

			content, err := os.ReadFile(tmpfile.Name())
			Expect(err).ToNot(HaveOccurred())

			// \ufeff is the UTF-8 BOM
			expectedOutput := "\ufeffFirstName;LastName;Email\nJohn;Doe;john.doe@example.com\nJane;Smith;jane.smith@example.com\n"
			Expect(string(content)).To(Equal(expectedOutput))
		})

		It("returns an error if the csv file cannot be created", func() {
			notAFile, err := os.MkdirTemp("", "testdir")
			Expect(err).ToNot(HaveOccurred())

			err = app.NewCSVFileWriter().Write(notAFile, csvHeader, csvRecords)
			Expect(err.Error()).To(ContainSubstring("failed to create csv file"))
		})
	})
})
