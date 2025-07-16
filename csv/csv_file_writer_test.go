package csv_test

import (
	"ctRestClient/csv"
	"os"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"

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
		It("writes a UTF-16 csv file", func() {
			tmpfile, err := os.CreateTemp("", "test.csv")
			Expect(err).ToNot(HaveOccurred())

			defer os.Remove(tmpfile.Name())

			err = csv.NewCSVFileWriter().Write(tmpfile.Name(), csvHeader, csvRecords)
			Expect(err).ToNot(HaveOccurred())

			content, err := os.ReadFile(tmpfile.Name())
			Expect(err).ToNot(HaveOccurred())

			encoder := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewEncoder()

			// Transform the string to UTF-16
			expectedUTF16Output, _, err := transform.String(encoder, "FirstName;LastName;Email\nJohn;Doe;john.doe@example.com\nJane;Smith;jane.smith@example.com\n")
			Expect(err).ToNot(HaveOccurred())

			Expect(string(content)).To(Equal(expectedUTF16Output))
		})

		It("returns an error if the csv file cannot be created", func() {
			notAFile, err := os.MkdirTemp("", "testdir")
			Expect(err).ToNot(HaveOccurred())

			err = csv.NewCSVFileWriter().Write(notAFile, csvHeader, csvRecords)
			Expect(err.Error()).To(ContainSubstring("failed to create csv file"))
		})
	})
})
