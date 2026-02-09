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

var _ = Describe("Group status check", func() {

	var _ = Describe("with dynamic groups", func() {
		BeforeEach(func() {
			configContent := testutil.YamlToString(fmt.Sprintf(`
			---
			instances:
			- hostname: %s
			  token_name: %s
			  groups:
			  - name: testGroup3
			    fields: [firstName, lastName, street, zip, city, sexId, birthday]
			`, ctInstanceUrl, ctInstanceUrl))

			configPath = filepath.Join(tempDir, "geburtstage-config.yml")
			err = os.WriteFile(configPath, []byte(configContent), 0644)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should not create the csv files", func() {
			runBinary()

			notExistingCsvFile := filepath.Join(filepath.Dir(logFilePath()), ctInstanceUrl, "testGroup1.csv")
			Expect(notExistingCsvFile).NotTo(BeAnExistingFile())
		})

		It("should log the expected output", func() {
			expectedLog := fmt.Sprintf(`%s
[INFO] +----------------------------------------------------------------------+
[INFO] | Processing instance '%s'                             |
[INFO] +----------------------------------------------------------------------+
[INFO]
[INFO]   processing group 'testGroup3'
[WARN]       skipping csv creation since the group is not active
`, LOG_HEADER, ctInstanceUrl)

			runBinary()

			normalizedLines := getNormalizedLogLines(logFilePath())
			expectedLogLines := strings.Split(expectedLog, "\n")

			for i, line := range normalizedLines {
				Expect(line).To(Equal(expectedLogLines[i]), fmt.Sprintf("Mismatch at line %d:\nExpected: %s\nActual:   %s", i+1, expectedLogLines[i], line))
			}
		})
	})

	var _ = Describe("with non dynamic groups", func() {
		BeforeEach(func() {
			configContent := testutil.YamlToString(fmt.Sprintf(`
			---
			instances:
			- hostname: %s
			  token_name: %s
			  groups:
			  - name: testGroup4
			    fields: [id, firstName, lastName, street, zip, city, sexId, birthday]
			`, ctInstanceUrl, ctInstanceUrl))

			configPath = filepath.Join(tempDir, "geburtstage-config.yml")
			err = os.WriteFile(configPath, []byte(configContent), 0644)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should log the expected output", func() {
			expectedLog := fmt.Sprintf(`%s
[INFO] +----------------------------------------------------------------------+
[INFO] | Processing instance '%s'                             |
[INFO] +----------------------------------------------------------------------+
[INFO]
[INFO]   processing group 'testGroup4'
[INFO]       the group has 1 persons
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
