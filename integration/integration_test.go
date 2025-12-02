package integration

import (
	"ctRestClient/config"
	"ctRestClient/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ChurchTools Integration Test with HTTP Server", func() {
	var (
		server        *httptest.Server
		tempDir       string
		configPath    string
		dataDir       string
		outputDir     string
		logFile       string
		appLogger     logger.Logger
		keepassDbPath string
	)

	BeforeEach(func() {
		// Create temporary directories for test
		var err error
		tempDir, err = os.MkdirTemp("", "ctRestClient_http_integration_test_")
		Expect(err).ToNot(HaveOccurred())

		configPath = filepath.Join(tempDir, "config.yml")
		dataDir = filepath.Join(tempDir, "data")
		outputDir = filepath.Join(tempDir, "output")

		// Create necessary directories
		err = os.MkdirAll(filepath.Join(dataDir, "persons"), 0755)
		Expect(err).ToNot(HaveOccurred())
		err = os.MkdirAll(outputDir, 0755)
		Expect(err).ToNot(HaveOccurred())

		// Create application logger
		logFile = filepath.Join(outputDir, "integration_test.log")
		appLogger = logger.NewLogger(logFile)

		// Copy the real Keepass database file from integration directory
		sourceKeepassPath := filepath.Join(".", "churchtools-tokens.kdbx")
		keepassDbPath = filepath.Join(tempDir, "churchtools-tokens.kdbx")

		sourceData, err := os.ReadFile(sourceKeepassPath)
		Expect(err).ToNot(HaveOccurred())

		err = os.WriteFile(keepassDbPath, sourceData, 0644)
		Expect(err).ToNot(HaveOccurred())

		// Create mock server
		server = httptest.NewTLSServer(createChurchToolsHandler())

		os.Setenv("ALLOW_SELF_SIGNED_CERTS", "true")
	})

	AfterEach(func() {
		if server != nil {
			server.Close()
		}

		os.Unsetenv("ALLOW_SELF_SIGNED_CERTS")
		os.RemoveAll(tempDir)
	})

	Context("Full end-to-end integration test", func() {
		It("exports group members into a CSV file", func() {

			configContent := fmt.Sprintf(`---
instances:
  - hostname: %s
    token_name: TEST_TOKEN
    groups:
    - name: youth_group
      fields:
      - id
      - firstName  
      - lastName
      - email
      - sexId`, strings.TrimPrefix(server.URL, "https://"))

			err := os.WriteFile(configPath, []byte(configContent), 0644)
			Expect(err).ToNot(HaveOccurred())

			// Load configuration
			cfg, err := config.LoadConfig(configPath)
			Expect(err).ToNot(HaveOccurred())

			err = RunApplicationWrapper(cfg, outputDir, dataDir, keepassDbPath, "abcd1234", appLogger)
			Expect(err).ToNot(HaveOccurred())

			// Check if log file was created
			_, err = os.Stat(logFile)
			Expect(err).ToNot(HaveOccurred())

			// Check if output CSV file was created
			csvFilePath := getCsvPath(outputDir, "youth_group")

			// Validate CSV content
			csvContent, _ := os.ReadFile(csvFilePath)
			csvString := getUTF8String(csvContent)

			Expect(csvString).To(Equal("id;firstName;lastName;email;sexId\n101;John;Doe;john.doe@example.com;1\n102;Jane;Smith;jane.smith@example.com;2\n"))
		})

		It("exports replaces field of group members by data files", func() {

			configContent := fmt.Sprintf(`---
instances:
  - hostname: %s
    token_name: TEST_TOKEN
    groups:
    - name: youth_group
      fields:
      - id
      - firstName  
      - lastName
      - email
      - {fieldname: sexId, columnname: "sex"}`, strings.TrimPrefix(server.URL, "https://"))

			err := os.WriteFile(configPath, []byte(configContent), 0644)
			Expect(err).ToNot(HaveOccurred())

			// Create mapping data for sexId field
			sexMappingContent := `---
1: "Male"
2: "Female"`
			err = os.WriteFile(filepath.Join(dataDir, "persons", "sexId.yml"), []byte(sexMappingContent), 0644)
			Expect(err).ToNot(HaveOccurred())

			// Load configuration
			cfg, err := config.LoadConfig(configPath)
			Expect(err).ToNot(HaveOccurred())

			err = RunApplicationWrapper(cfg, outputDir, dataDir, keepassDbPath, "abcd1234", appLogger)
			Expect(err).ToNot(HaveOccurred())

			// Check if log file was created
			_, err = os.Stat(logFile)
			Expect(err).ToNot(HaveOccurred())

			// Check if output CSV file was created
			csvFilePath := getCsvPath(outputDir, "youth_group")

			// Validate CSV content
			csvContent, _ := os.ReadFile(csvFilePath)
			csvString := getUTF8String(csvContent)

			Expect(csvString).To(Equal("id;firstName;lastName;email;sex\n101;John;Doe;john.doe@example.com;Male\n102;Jane;Smith;jane.smith@example.com;Female\n"))
		})
	})
})

func getUTF8String(content []byte) string {
	utf16Decoder := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()
	decodedBytes, _, err := transform.Bytes(utf16Decoder, content)
	if err != nil {
		return ""
	}
	return string(decodedBytes)
}

func getCsvPath(outputDir, groupName string) string {
	csvFileName := fmt.Sprintf("%s.csv", groupName)
	var csvFilePath string

	err := filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == csvFileName {
			csvFilePath = path
			return filepath.SkipDir // Stop searching once found
		}
		return nil
	})

	if err != nil {
		return ""
	}

	return csvFilePath
}

// createChurchToolsHandler creates an HTTP handler that simulates ChurchTools API responses
func createChurchToolsHandler() http.Handler {
	mux := http.NewServeMux()

	// Handle authentication/login endpoint
	mux.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"data": map[string]string{
				"token": "mock-jwt-token",
			},
		}
		err := json.NewEncoder(w).Encode(response)
		Expect(err).NotTo(HaveOccurred())
	})

	// Handle groups endpoint - return groups list
	mux.HandleFunc("/api/groups", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Get the query parameter to filter groups
		queryParam := r.URL.Query().Get("query")

		var responseData []map[string]interface{}

		// All available groups
		allGroups := []map[string]interface{}{
			{
				"id":   1,
				"guid": "youth-group-guid",
				"name": "youth_group",
			},
			{
				"id":   2,
				"guid": "adult-group-guid",
				"name": "adult_group",
			},
		}

		// Filter groups based on query parameter
		if queryParam != "" {
			for _, group := range allGroups {
				if group["name"] == queryParam {
					responseData = append(responseData, group)
				}
			}
		} else {
			// Return all groups if no query parameter
			responseData = allGroups
		}

		response := map[string]interface{}{
			"data": responseData,
		}
		err := json.NewEncoder(w).Encode(response)
		Expect(err).NotTo(HaveOccurred())
	})

	// Handle group members endpoint
	mux.HandleFunc("/api/groups/members", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Check if requesting youth group (id=1)
		if r.URL.Query().Get("ids[]") == "1" {
			response := map[string]interface{}{
				"data": []map[string]interface{}{
					{
						"personId":          101,
						"groupId":           1,
						"groupTypeRoleId":   1,
						"groupMemberStatus": "active",
						"deleted":           false,
					},
					{
						"personId":          102,
						"groupId":           1,
						"groupTypeRoleId":   1,
						"groupMemberStatus": "active",
						"deleted":           false,
					},
				},
			}
			err := json.NewEncoder(w).Encode(response)
			Expect(err).NotTo(HaveOccurred())
		} else {
			// Return empty for other groups
			response := map[string]interface{}{
				"data": []map[string]interface{}{},
			}
			err := json.NewEncoder(w).Encode(response)
			Expect(err).NotTo(HaveOccurred())
		}
	})

	// Handle dynamic group status endpoint for youth_group (id=1)
	mux.HandleFunc("/api/dynamicgroups/1/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := map[string]interface{}{
			"dynamicGroupStatus": "active",
		}
		err := json.NewEncoder(w).Encode(response)
		Expect(err).NotTo(HaveOccurred())
	})

	// Handle persons endpoint - return person details
	mux.HandleFunc("/api/persons", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Parse requested person IDs
		personIds := r.URL.Query()["ids[]"]

		var persons []json.RawMessage

		for _, idStr := range personIds {
			if idStr == "101" {
				person := json.RawMessage(`{
					"id": 101,
					"firstName": "John",
					"lastName": "Doe", 
					"email": "john.doe@example.com",
					"sexId": 1
				}`)
				persons = append(persons, person)
			} else if idStr == "102" {
				person := json.RawMessage(`{
					"id": 102,
					"firstName": "Jane",
					"lastName": "Smith",
					"email": "jane.smith@example.com", 
					"sexId": 2
				}`)
				persons = append(persons, person)
			}
		}

		response := map[string]interface{}{
			"data": persons,
		}
		err := json.NewEncoder(w).Encode(response)
		Expect(err).NotTo(HaveOccurred())
	})

	// Handle any other API calls with a generic response
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received request: %s %s\n", r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"data": []interface{}{},
		}
		err := json.NewEncoder(w).Encode(response)
		Expect(err).NotTo(HaveOccurred())
	})

	return mux
}
