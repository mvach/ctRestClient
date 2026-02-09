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

var mappingsDir string

var _ = Describe("running ctRestClient with configured mapping for sexId", func() {

	BeforeEach(func() {
		mappingsDir = filepath.Join(dataDir, "mappings", "persons")
		err := os.MkdirAll(mappingsDir, 0755)
		Expect(err).NotTo(HaveOccurred())

		sexIdMappings := testutil.YamlToString(`
			---
			1: männlich
			2: weiblich
			3: divers
			`)
		err = os.WriteFile(filepath.Join(mappingsDir, "sexId.yml"), []byte(sexIdMappings), 0644)
		Expect(err).NotTo(HaveOccurred())

		configContent := testutil.YamlToString(fmt.Sprintf(`
			---
			instances:
			- hostname: %s
			  token_name: %s
			  groups:
			  - name: testGroup1
			    fields: [firstName, lastName, street, zip, city, {fieldname: sexId, columnname: "geschlecht"}, birthday]
			  - name: testGroup2
			    fields: [firstName, lastName, street, zip, city, {fieldname: sexId, columnname: "geschlecht"}, birthday]
			`, ctInstanceUrl, ctInstanceUrl))

		configPath = filepath.Join(tempDir, "geburtstage-config.yml")
		err = os.WriteFile(configPath, []byte(configContent), 0644)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should create the expected csv files with sex field", func() {
		runBinary()

		csvString := csvFileToString("testGroup1.csv")

		Expect(csvString).To(Equal(`firstName;lastName;street;zip;city;geschlecht;birthday
Anna;Mustermann;Musterstraße 1;11111;Musterstadt;weiblich;2008-09-12
Max;Beispiel;Beispielweg 2;22222;Musterdorf;männlich;2008-09-24
Lisa;Tester;Testallee 3;33333;Mustertal;weiblich;2008-09-29
Tom;Probe;Probestraße 4/1;44444;Musterhausen;männlich;2008-09-18
`))

		csvString = csvFileToString("testGroup2.csv")

		Expect(csvString).To(Equal(`firstName;lastName;street;zip;city;geschlecht;birthday
Kerstin;Demonstration;Demoweg 5;55555;Musterberg;weiblich;1982-03-25
Markus;Vorschau;Vorschauplatz 6;66666;Musterfeld;männlich;1982-03-04
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
[INFO]
[INFO]   processing group 'testGroup2'
[INFO]       the group has 2 persons
`, LOG_HEADER, ctInstanceUrl)

		runBinary()

		normalizedLines := getNormalizedLogLines(logFilePath())
		expectedLogLines := strings.Split(expectedLog, "\n")

		for i, line := range normalizedLines {
			Expect(line).To(Equal(expectedLogLines[i]), fmt.Sprintf("Mismatch at line %d:\nExpected: %s\nActual:   %s", i+1, expectedLogLines[i], line))
		}
	})
})
