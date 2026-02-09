package endtoend_test

import (
	"ctRestClient/internal/testutil"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var blocklistsDir string

var _ = Describe("running ctRestClient with configured blocklist", func() {

	BeforeEach(func() {
		configContent := testutil.YamlToString(fmt.Sprintf(`
			---
			instances:
			- hostname: %s
			  token_name: %s
			  groups:
			  - name: testGroup1
			    fields: [firstName, lastName, street, zip, city, sexId, birthday]
		`, ctInstanceUrl, ctInstanceUrl))

		configPath = filepath.Join(tempDir, "geburtstage-config.yml")
		err = os.WriteFile(configPath, []byte(configContent), 0644)
		Expect(err).ToNot(HaveOccurred())

		blocklistsDir = filepath.Join(dataDir, "blocklists")
		err := os.MkdirAll(blocklistsDir, 0755)
		Expect(err).NotTo(HaveOccurred())
	})

	var _ = Describe("with an empty blocklist", func() {
		BeforeEach(func() {
			emptyBlocklist := testutil.YamlToString(``)
			err = os.WriteFile(filepath.Join(blocklistsDir, "testGroup1.yml"), []byte(emptyBlocklist), 0644)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should create the expected csv files", func() {
			runBinary()

			csvString := csvFileToString("testGroup1.csv")

			Expect(csvString).To(Equal(`firstName;lastName;street;zip;city;sexId;birthday
Anna;Mustermann;Musterstraße 1;11111;Musterstadt;2;2008-09-12
Max;Beispiel;Beispielweg 2;22222;Musterdorf;1;2008-09-24
Lisa;Tester;Testallee 3;33333;Mustertal;2;2008-09-29
Tom;Probe;Probestraße 4/1;44444;Musterhausen;1;2008-09-18
`))
		})

		It("should log the expected output", func() {
			expectedLog := fmt.Sprintf(`%s
[INFO] +----------------------------------------------------------------------+
[INFO] | Processing instance '%s'                             |
[INFO] +----------------------------------------------------------------------+
[INFO]
[INFO]   processing group 'testGroup1'
[INFO]       the group has 4 persons
[INFO]       using blocklist 'testGroup1.yml'
[INFO]       blocked 0 persons
`, LOG_HEADER, ctInstanceUrl)

			runBinary()

			normalizedLines := getNormalizedLogLines(logFilePath())
			expectedLogLines := strings.Split(expectedLog, "\n")

			for i, line := range normalizedLines {
				Expect(line).To(Equal(expectedLogLines[i]), fmt.Sprintf("Mismatch at line %d:\nExpected: %s\nActual:   %s", i+1, expectedLogLines[i], line))
			}
		})
	})

	var _ = Describe("with a not matching blocklist", func() {
		BeforeEach(func() {
			emptyBlocklist := testutil.YamlToString(`
			- street: "Non Matching Street"
			`)
			err = os.WriteFile(filepath.Join(blocklistsDir, "testGroup1.yml"), []byte(emptyBlocklist), 0644)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should create the expected csv files", func() {
			runBinary()

			csvString := csvFileToString("testGroup1.csv")

			Expect(csvString).To(Equal(`firstName;lastName;street;zip;city;sexId;birthday
Anna;Mustermann;Musterstraße 1;11111;Musterstadt;2;2008-09-12
Max;Beispiel;Beispielweg 2;22222;Musterdorf;1;2008-09-24
Lisa;Tester;Testallee 3;33333;Mustertal;2;2008-09-29
Tom;Probe;Probestraße 4/1;44444;Musterhausen;1;2008-09-18
`))
		})

		It("should log the expected output", func() {
			expectedLog := fmt.Sprintf(`%s
[INFO] +----------------------------------------------------------------------+
[INFO] | Processing instance '%s'                             |
[INFO] +----------------------------------------------------------------------+
[INFO]
[INFO]   processing group 'testGroup1'
[INFO]       the group has 4 persons
[INFO]       using blocklist 'testGroup1.yml'
[INFO]       blocked 0 persons
`, LOG_HEADER, ctInstanceUrl)

			runBinary()

			normalizedLines := getNormalizedLogLines(logFilePath())
			expectedLogLines := strings.Split(expectedLog, "\n")

			for i, line := range normalizedLines {
				Expect(line).To(Equal(expectedLogLines[i]), fmt.Sprintf("Mismatch at line %d:\nExpected: %s\nActual:   %s", i+1, expectedLogLines[i], line))
			}
		})
	})

	var _ = Describe("with matching blocklist", func() {
		BeforeEach(func() {
			emptyBlocklist := testutil.YamlToString(`
			- street: "Testallee 3"
			`)
			err = os.WriteFile(filepath.Join(blocklistsDir, "testGroup1.yml"), []byte(emptyBlocklist), 0644)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should create the expected csv files", func() {
			runBinary()

			csvString := csvFileToString("testGroup1.csv")

			Expect(csvString).To(Equal(`firstName;lastName;street;zip;city;sexId;birthday
Anna;Mustermann;Musterstraße 1;11111;Musterstadt;2;2008-09-12
Max;Beispiel;Beispielweg 2;22222;Musterdorf;1;2008-09-24
Tom;Probe;Probestraße 4/1;44444;Musterhausen;1;2008-09-18
`))
		})

		It("should log the expected output", func() {
			expectedLog := fmt.Sprintf(`%s
[INFO] +----------------------------------------------------------------------+
[INFO] | Processing instance '%s'                             |
[INFO] +----------------------------------------------------------------------+
[INFO]
[INFO]   processing group 'testGroup1'
[INFO]       the group has 4 persons
[INFO]       using blocklist 'testGroup1.yml'
[INFO]       -> "Lisa" "Tester" will not be added to csv file
[INFO]       blocked 1 persons
`, LOG_HEADER, ctInstanceUrl)

			runBinary()

			normalizedLines := getNormalizedLogLines(logFilePath())
			expectedLogLines := strings.Split(expectedLog, "\n")

			for i, line := range normalizedLines {
				Expect(line).To(Equal(expectedLogLines[i]), fmt.Sprintf("Mismatch at line %d:\nExpected: %s\nActual:   %s", i+1, expectedLogLines[i], line))
			}
		})
	})
})
