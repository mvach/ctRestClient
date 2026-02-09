package endtoend_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var LOG_HEADER = `[INFO]
[INFO] +----------------------------------------------------------------------+
[INFO] | User: 'matthias'                                                     |
[INFO] | OS:   'Linux'                                                        |
[INFO] | Date: 'xxxx-xx-xx xx:xx:xx'                                          |
[INFO] +----------------------------------------------------------------------+
[INFO]
[INFO]`

var (
	executable    string
	tempDir       string
	outputDir     string
	dataDir       string
	kdbxPath      string
	configPath    string
	ctInstanceUrl string
	err           error
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Datamapping Suite")
}

var _ = BeforeSuite(func() {
	// List of required environment variables
	requiredVars := []string{"KEY_PASS_PASSWORD", "CT_INSTANCE_URL"}

	for _, varName := range requiredVars {
		value := os.Getenv(varName)
		Expect(value).ToNot(BeEmpty(), "Environment variable %s must be set", varName)
	}
})

var _ = BeforeEach(func() {
	executable = filepath.Join(ModuleRoot(), "dist", "ctRestClient-linux-amd64")
	tempDir = GinkgoT().TempDir()
	outputDir = filepath.Join(tempDir, "exports")
	dataDir = filepath.Join(tempDir, "data")
	kdbxPath = filepath.Join(ModuleRoot(), "churchtools-tokens.kdbx")
	ctInstanceUrl = os.Getenv("CT_INSTANCE_URL")
})

func ModuleRoot() string {
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}")
	out, err := cmd.Output()
	Expect(err).NotTo(HaveOccurred())
	return strings.TrimSpace(string(out))
}

func logFilePath() string {
	fileName := "ctRestClient.log"
	filePath := ""

	err := filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == fileName {
			filePath = path
			return filepath.SkipDir // Stop walking after finding the file
		}
		return nil
	})
	Expect(err).ToNot(HaveOccurred())
	Expect(filePath == "").To(BeFalse())

	return filePath
}

func csvFileToString(fileName string) string {
	csvContent, err := os.ReadFile(filepath.Join(filepath.Dir(logFilePath()), ctInstanceUrl, fileName))
	Expect(err).ToNot(HaveOccurred())

	utf16Decoder := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()
	decodedBytes, _, err := transform.Bytes(utf16Decoder, csvContent)
	Expect(err).ToNot(HaveOccurred())

	return string(decodedBytes)
}

func getNormalizedLogLines(logFilePath string) []string {
	logContent, err := os.ReadFile(logFilePath)
	Expect(err).NotTo(HaveOccurred())

	log := string(logContent)

	lines := strings.Split(log, "\n")
	var normalizedLines []string

	for _, line := range lines {
		// Remove leading timestamp (e.g., "2026/02/07 20:22:55 ")
		reLeading := regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} `)
		line = reLeading.ReplaceAllString(line, "")

		// Replace timestamp in the Date field (e.g., "2026-02-07 20:22:55") with <TIMESTAMP>
		reDate := regexp.MustCompile(`'\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}'`)
		line = reDate.ReplaceAllString(line, "'xxxx-xx-xx xx:xx:xx'")

		// Drop trailing spaces at the end of the line
		line = strings.TrimRight(line, " ")

		normalizedLines = append(normalizedLines, line)
	}

	return normalizedLines
}

func runBinary() {
	cmd := exec.Command(executable, "-k", kdbxPath, "-c", configPath, "-o", outputDir, "-d", dataDir)
	cmd.Env = append(os.Environ(), "KEY_PASS_PASSWORD="+os.Getenv("KEY_PASS_PASSWORD"))

	// // Log command
	// args := append([]string{executable}, cmd.Args[1:]...)
	// commandString := strings.Join(args, " ")
	// fmt.Printf("Running command: %s\n\n", commandString)

	err = cmd.Run()
	Expect(err).NotTo(HaveOccurred())
}
