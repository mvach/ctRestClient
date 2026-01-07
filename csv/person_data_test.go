package csv_test

import (
	"ctRestClient/config"
	"ctRestClient/csv"
	"ctRestClient/csv/csvfakes"
	"ctRestClient/logger/loggerfakes"
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func ptr(s string) *string {
	return &s
}

var _ = Describe("PersonData", func() {

	var (
		persons                []json.RawMessage
		fileDataProvider       *csvfakes.FakeFileDataProvider
		blocklistsDataProvider *csvfakes.FakeBlockListDataProvider
		logger                 *loggerfakes.FakeLogger
	)

	BeforeEach(func() {
		person1 := `{	
            "id": 1,
            "firstName": "foo_firstname",
            "lastName": "foo_lastname",
			"height": 2.75
        }`
		person2 := `{	
            "id": 2,
            "firstName": "bar_firstname",
            "lastName": "bar_lastname",
			"height": 1.0
        }`
		persons = []json.RawMessage{json.RawMessage(person1), json.RawMessage(person2)}
		fileDataProvider = &csvfakes.FakeFileDataProvider{}
		blocklistsDataProvider = &csvfakes.FakeBlockListDataProvider{}
		logger = &loggerfakes.FakeLogger{}
	})

	var _ = Describe("NewPersonData", func() {
		It("returns persons", func() {
			group := config.Group{Fields: []config.Field{{FieldName: ptr("id")}, {FieldName: ptr("firstName")}, {FieldName: ptr("lastName")}, {FieldName: ptr("height")}}}
			data, err := csv.NewPersonData(persons, group, fileDataProvider, blocklistsDataProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(data.Header()).To(Equal([]string{"id", "firstName", "lastName", "height"}))
			Expect(data.Records()).To(HaveLen(2))
			Expect(data.Records()[0]).To(Equal([]string{"1", "foo_firstname", "foo_lastname", "2.75"}))
			Expect(data.Records()[1]).To(Equal([]string{"2", "bar_firstname", "bar_lastname", "1.0"}))
		})

		It("returns an error if json cannot be read", func() {
			persons := []json.RawMessage{json.RawMessage(`[]`)}

			group := config.Group{Fields: []config.Field{{FieldName: ptr("id")}, {FieldName: ptr("firstName")}, {FieldName: ptr("lastName")}}}
			data, err := csv.NewPersonData(persons, group, fileDataProvider, blocklistsDataProvider, logger)
			Expect(data).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("failed to read person information raw json"))
		})

		It("skips blocked persons", func() {
			blocklistsDataProvider.IsBlockedReturnsOnCall(0, true, nil)
			blocklistsDataProvider.IsBlockedReturnsOnCall(1, false, nil)

			group := config.Group{Name: "test", Fields: []config.Field{{FieldName: ptr("id")}, {FieldName: ptr("unknown")}, {FieldName: ptr("lastName")}}}
			data, err := csv.NewPersonData(persons, group, fileDataProvider, blocklistsDataProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(logger.InfoArgsForCall(0)).To(Equal("      -> \"foo_firstname\" \"foo_lastname\" will not be added to csv file"))
			Expect(logger.WarnArgsForCall(0)).To(Equal("      Field 'unknown' does not exist"))

			Expect(data.Header()).To(Equal([]string{"id", "unknown", "lastName"}))
			Expect(data.Records()).To(HaveLen(1))
			Expect(data.Records()[0]).To(Equal([]string{"2", "", "bar_lastname"}))
		})

		It("does not skip persons if an error occurs while checking blocklists", func() {
			blocklistsDataProvider.IsBlockedReturnsOnCall(0, false, nil)
			blocklistsDataProvider.IsBlockedReturnsOnCall(1, false, errors.New("boom"))

			group := config.Group{Name: "test", Fields: []config.Field{{FieldName: ptr("id")}, {FieldName: ptr("unknown")}, {FieldName: ptr("lastName")}}}
			data, err := csv.NewPersonData(persons, group, fileDataProvider, blocklistsDataProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(logger.ErrorArgsForCall(0)).To(Equal("      failed to check if person is blocked: 'boom'"))
			Expect(logger.WarnArgsForCall(0)).To(Equal("      Field 'unknown' does not exist"))

			Expect(data.Header()).To(Equal([]string{"id", "unknown", "lastName"}))
			Expect(data.Records()).To(HaveLen(2))
			Expect(data.Records()[0]).To(Equal([]string{"1", "", "foo_lastname"}))
			Expect(data.Records()[1]).To(Equal([]string{"2", "", "bar_lastname"}))
		})

		It("sets unknown fields to empty string", func() {
			group := config.Group{Fields: []config.Field{{FieldName: ptr("id")}, {FieldName: ptr("unknown")}, {FieldName: ptr("lastName")}}}
			data, err := csv.NewPersonData(persons, group, fileDataProvider, blocklistsDataProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(logger.WarnArgsForCall(0)).To(Equal("      Field 'unknown' does not exist"))

			Expect(data.Header()).To(Equal([]string{"id", "unknown", "lastName"}))
			Expect(data.Records()).To(HaveLen(2))
			Expect(data.Records()[0]).To(Equal([]string{"1", "", "foo_lastname"}))
			Expect(data.Records()[1]).To(Equal([]string{"2", "", "bar_lastname"}))
		})

		It("sets null values to empty string", func() {
			person := `{
				"id": 1,
				"date": null
        	}`
			persons = []json.RawMessage{json.RawMessage(person)}
			group := config.Group{Fields: []config.Field{{FieldName: ptr("id")}, {FieldName: ptr("date")}}}
			data, err := csv.NewPersonData(persons, group, fileDataProvider, blocklistsDataProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(data.Header()).To(Equal([]string{"id", "date"}))
			Expect(data.Records()).To(HaveLen(1))
			Expect(data.Records()[0]).To(Equal([]string{"1", ""}))
		})

		It("preserves trailing zeros in float values", func() {
			person := `{
				"id": 1,
				"height": 2.7500
			}`
			persons = []json.RawMessage{json.RawMessage(person)}
			group := config.Group{Fields: []config.Field{{FieldName: ptr("id")}, {FieldName: ptr("height")}}}
			data, err := csv.NewPersonData(persons, group, fileDataProvider, blocklistsDataProvider, logger)

			Expect(err).NotTo(HaveOccurred())
			Expect(data.Header()).To(Equal([]string{"id", "height"}))
			Expect(data.Records()).To(HaveLen(1))
			Expect(data.Records()[0]).To(Equal([]string{"1", "2.7500"}))
		})

		It("sets booleans as strings", func() {
			person := `{
				"id": 1,
				"isSet": true
        	}`
			persons = []json.RawMessage{json.RawMessage(person)}
			group := config.Group{Fields: []config.Field{{FieldName: ptr("id")}, {FieldName: ptr("isSet")}}}
			data, err := csv.NewPersonData(persons, group, fileDataProvider, blocklistsDataProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(data.Header()).To(Equal([]string{"id", "isSet"}))
			Expect(data.Records()).To(HaveLen(1))
			Expect(data.Records()[0]).To(Equal([]string{"1", "true"}))
		})

		It("sets unknown fields to empty string", func() {
			group := config.Group{Fields: []config.Field{{FieldName: ptr("id")}, {FieldName: ptr("unknown")}}}
			_, err := csv.NewPersonData(persons, group, fileDataProvider, blocklistsDataProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(logger.WarnArgsForCall(0)).To(Equal("      Field 'unknown' does not exist"))
		})

		It("returns mapped data for string key fields", func() {
			person1 := `{	
				"id": 1,
				"key": "value"
			}`
			persons = []json.RawMessage{json.RawMessage(person1)}

			fileDataProvider.GetDataReturnsOnCall(0, "mapped_value", nil)

			group := config.Group{Fields: []config.Field{{FieldName: ptr("id")}, {Object: &config.FieldInformation{FieldName: "key", ColumnName: "mappedColumn"}}}}
			data, err := csv.NewPersonData(persons, group, fileDataProvider, blocklistsDataProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(fileDataProvider.GetDataCallCount()).To(Equal(1))

			fieldName, fieldValue := fileDataProvider.GetDataArgsForCall(0)
			Expect(fieldName).To(Equal("key"))
			Expect(fieldValue).To(Equal(json.RawMessage(`"value"`)))

			Expect(data.Header()).To(Equal([]string{"id", "mappedColumn"}))
			Expect(data.Records()).To(HaveLen(1))
			Expect(data.Records()[0]).To(Equal([]string{"1", "mapped_value"}))

		})

		It("returns empty string for string key fields that are not mapped", func() {
			person1 := `{	
				"id": 1,
				"key": "value"
			}`
			persons = []json.RawMessage{json.RawMessage(person1)}

			fileDataProvider.GetDataReturnsOnCall(0, "", errors.New("not found"))

			group := config.Group{Fields: []config.Field{{FieldName: ptr("id")}, {Object: &config.FieldInformation{FieldName: "key", ColumnName: "mappedColumn"}}}}
			data, err := csv.NewPersonData(persons, group, fileDataProvider, blocklistsDataProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(fileDataProvider.GetDataCallCount()).To(Equal(1))

			fieldName, fieldValue := fileDataProvider.GetDataArgsForCall(0)
			Expect(fieldName).To(Equal("key"))
			Expect(fieldValue).To(Equal(json.RawMessage(`"value"`)))

			Expect(data.Header()).To(Equal([]string{"id", "mappedColumn"}))
			Expect(data.Records()).To(HaveLen(1))
			Expect(data.Records()[0]).To(Equal([]string{"1", ""}))

		})

		It("returns mapped data for float64 key fields", func() {
			person1 := `{	
				"id": 1,
				"key": 1.2
			}`
			persons = []json.RawMessage{json.RawMessage(person1)}

			fileDataProvider.GetDataReturnsOnCall(0, "mapped_value", nil)

			group := config.Group{Fields: []config.Field{{FieldName: ptr("id")}, {Object: &config.FieldInformation{FieldName: "key", ColumnName: "mappedColumn"}}}}
			data, err := csv.NewPersonData(persons, group, fileDataProvider, blocklistsDataProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(fileDataProvider.GetDataCallCount()).To(Equal(1))

			fieldName, fieldValue := fileDataProvider.GetDataArgsForCall(0)
			Expect(fieldName).To(Equal("key"))
			Expect(fieldValue).To(Equal(json.RawMessage(`1.2`)))

			Expect(data.Header()).To(Equal([]string{"id", "mappedColumn"}))
			Expect(data.Records()).To(HaveLen(1))
			Expect(data.Records()[0]).To(Equal([]string{"1", "mapped_value"}))

		})

		It("returns empty string for float64 key fields that are not mapped", func() {
			person1 := `{	
				"id": 1,
				"key": 1.2
			}`
			persons = []json.RawMessage{json.RawMessage(person1)}

			fileDataProvider.GetDataReturnsOnCall(0, "", errors.New("not found"))

			group := config.Group{Fields: []config.Field{{FieldName: ptr("id")}, {Object: &config.FieldInformation{FieldName: "key", ColumnName: "mappedColumn"}}}}
			data, err := csv.NewPersonData(persons, group, fileDataProvider, blocklistsDataProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(fileDataProvider.GetDataCallCount()).To(Equal(1))

			fieldName, fieldValue := fileDataProvider.GetDataArgsForCall(0)
			Expect(fieldName).To(Equal("key"))
			Expect(fieldValue).To(Equal(json.RawMessage(`1.2`)))

			Expect(data.Header()).To(Equal([]string{"id", "mappedColumn"}))
			Expect(data.Records()).To(HaveLen(1))
			Expect(data.Records()[0]).To(Equal([]string{"1", ""}))

		})
	})
})
