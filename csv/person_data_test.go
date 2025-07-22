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
		persons            []json.RawMessage
		personDataProvider *csvfakes.FakeFileDataProvider
		logger             *loggerfakes.FakeLogger
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
		personDataProvider = &csvfakes.FakeFileDataProvider{}
		logger = &loggerfakes.FakeLogger{}
	})

	var _ = Describe("NewPersonData", func() {
		It("returns persons", func() {
			data, err := csv.NewPersonData(persons, []config.Field{{RawString: ptr("id")}, {RawString: ptr("firstName")}, {RawString: ptr("lastName")}, {RawString: ptr("height")}}, personDataProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(data.Header()).To(Equal([]string{"id", "firstName", "lastName", "height"}))
			Expect(data.Records()).To(HaveLen(2))
			Expect(data.Records()[0]).To(Equal([]string{"1", "foo_firstname", "foo_lastname", "2.75"}))
			Expect(data.Records()[1]).To(Equal([]string{"2", "bar_firstname", "bar_lastname", "1.0"}))
		})

		It("returns an error if json cannot be read", func() {
			persons := []json.RawMessage{json.RawMessage(`[]`)}

			data, err := csv.NewPersonData(persons, []config.Field{{RawString: ptr("id")}, {RawString: ptr("firstName")}, {RawString: ptr("lastName")}}, personDataProvider, logger)
			Expect(data).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("failed to read person information raw json"))
		})

		It("sets unknown fields to empty string", func() {
			data, err := csv.NewPersonData(persons, []config.Field{{RawString: ptr("id")}, {RawString: ptr("unknown")}, {RawString: ptr("lastName")}}, personDataProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(logger.WarnArgsForCall(0)).To(Equal("    Field 'unknown' does not exist"))

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
			fields := []config.Field{{RawString: ptr("id")}, {RawString: ptr("date")}}
			data, err := csv.NewPersonData(persons, fields, personDataProvider, logger)
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
			fields := []config.Field{{RawString: ptr("id")}, {RawString: ptr("height")}}
			data, err := csv.NewPersonData(persons, fields, personDataProvider, logger)
			
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
			fields := []config.Field{{RawString: ptr("id")}, {RawString: ptr("isSet")}}
			data, err := csv.NewPersonData(persons, fields, personDataProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(data.Header()).To(Equal([]string{"id", "isSet"}))
			Expect(data.Records()).To(HaveLen(1))
			Expect(data.Records()[0]).To(Equal([]string{"1", "true"}))
		})

		It("sets unknown fields to empty string", func() {
			fields := []config.Field{{RawString: ptr("id")}, {RawString: ptr("unknown")}}
			_, err := csv.NewPersonData(persons, fields, personDataProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(logger.WarnArgsForCall(0)).To(Equal("    Field 'unknown' does not exist"))
		})

		It("returns mapped data for string key fields", func() {
			person1 := `{	
				"id": 1,
				"key": "value"
			}`
			persons = []json.RawMessage{json.RawMessage(person1)}

			personDataProvider.GetDataReturnsOnCall(0, "mapped_value", nil)

			fields := []config.Field{{RawString: ptr("id")}, {RawObject: &config.FieldInformation{FieldName: "key", ColumnName: "mappedColumn"}}}
			data, err := csv.NewPersonData(persons, fields, personDataProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(personDataProvider.GetDataCallCount()).To(Equal(1))

			fieldName, fieldValue := personDataProvider.GetDataArgsForCall(0)
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

			personDataProvider.GetDataReturnsOnCall(0, "", errors.New("not found"))

			fields := []config.Field{{RawString: ptr("id")}, {RawObject: &config.FieldInformation{FieldName: "key", ColumnName: "mappedColumn"}}}
			data, err := csv.NewPersonData(persons, fields, personDataProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(personDataProvider.GetDataCallCount()).To(Equal(1))

			fieldName, fieldValue := personDataProvider.GetDataArgsForCall(0)
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

			personDataProvider.GetDataReturnsOnCall(0, "mapped_value", nil)

			fields := []config.Field{{RawString: ptr("id")}, {RawObject: &config.FieldInformation{FieldName: "key", ColumnName: "mappedColumn"}}}
			data, err := csv.NewPersonData(persons, fields, personDataProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(personDataProvider.GetDataCallCount()).To(Equal(1))

			fieldName, fieldValue := personDataProvider.GetDataArgsForCall(0)
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

			personDataProvider.GetDataReturnsOnCall(0, "", errors.New("not found"))

			fields := []config.Field{{RawString: ptr("id")}, {RawObject: &config.FieldInformation{FieldName: "key", ColumnName: "mappedColumn"}}}
			data, err := csv.NewPersonData(persons, fields, personDataProvider, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(personDataProvider.GetDataCallCount()).To(Equal(1))

			fieldName, fieldValue := personDataProvider.GetDataArgsForCall(0)
			Expect(fieldName).To(Equal("key"))
			Expect(fieldValue).To(Equal(json.RawMessage(`1.2`)))

			Expect(data.Header()).To(Equal([]string{"id", "mappedColumn"}))
			Expect(data.Records()).To(HaveLen(1))
			Expect(data.Records()[0]).To(Equal([]string{"1", ""}))

		})
	})
})
