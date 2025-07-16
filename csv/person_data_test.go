package csv_test

import (
	"encoding/json"
	"ctRestClient/logger/loggerfakes"
	"ctRestClient/csv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PersonData", func() {

	var (
		persons             []json.RawMessage
		logger             *loggerfakes.FakeLogger
	)

	BeforeEach(func() {
		person1 := `{	
            "id": 1,
            "firstName": "foo_firstname",
            "lastName": "foo_lastname"
        }`
        person2 := `{	
            "id": 2,
            "firstName": "bar_firstname",
            "lastName": "bar_lastname"
        }`
        persons = []json.RawMessage{json.RawMessage(person1), json.RawMessage(person2)}

		logger = &loggerfakes.FakeLogger{}
	})
		

	var _ = Describe("NewPersonData", func() {
		It("returns persons", func() {
			data, err := csv.NewPersonData(persons, []string{"id", "firstName", "lastName"}, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(data.Header()).To(Equal([]string{"id", "firstName", "lastName"}))
			Expect(data.Records()).To(HaveLen(2))
			Expect(data.Records()[0]).To(Equal([]string{"1", "foo_firstname", "foo_lastname"}))
			Expect(data.Records()[1]).To(Equal([]string{"2", "bar_firstname", "bar_lastname"}))
		})

		It("returns an error if json cannot be read", func() {
            persons := []json.RawMessage{json.RawMessage(`[]`)}

            _, err := csv.NewPersonData(persons, []string{"id", "firstName", "lastName"}, logger)

            Expect(err.Error()).To(ContainSubstring("failed to read person information raw json"))
        })

		It("sets unknown fields to empty string", func() {
            data, err := csv.NewPersonData(persons, []string{"id", "unknown", "lastName"}, logger)
			Expect(err).NotTo(HaveOccurred())

            Expect(logger.WarnArgsForCall(0)).To(Equal("    Field 'unknown' does not exist"))

			Expect(data.Header()).To(Equal([]string{"id", "unknown", "lastName"}))
			Expect(data.Records()).To(HaveLen(2))
			Expect(data.Records()[0]).To(Equal([]string{"1","", "foo_lastname"}))
			Expect(data.Records()[1]).To(Equal([]string{"2","", "bar_lastname"}))
        })


		It("sets null values to empty string", func() {
            person := `{
				"id": 1,
				"date": null
        	}`
        	persons = []json.RawMessage{json.RawMessage(person)}
			data, err := csv.NewPersonData(persons, []string{"id", "date"}, logger)
			Expect(err).NotTo(HaveOccurred())

            Expect(data.Header()).To(Equal([]string{"id", "date"}))
			Expect(data.Records()).To(HaveLen(1))
			Expect(data.Records()[0]).To(Equal([]string{"1",""}))
        })

		It("sets unknown data types to empty string", func() {
            person := `{
				"id": 1,
				"isSet": true
        	}`
        	persons = []json.RawMessage{json.RawMessage(person)}
			data, err := csv.NewPersonData(persons, []string{"id", "isSet"}, logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(logger.WarnArgsForCall(0)).To(Equal("    Field 'isSet' is not a string or int"))

            Expect(data.Header()).To(Equal([]string{"id", "isSet"}))
			Expect(data.Records()).To(HaveLen(1))
			Expect(data.Records()[0]).To(Equal([]string{"1",""}))
        })

		It("sets unknown fields to empty string", func() {
            _, err := csv.NewPersonData(persons, []string{"id", "unknown"}, logger)
			Expect(err).NotTo(HaveOccurred())

            Expect(logger.WarnArgsForCall(0)).To(Equal("    Field 'unknown' does not exist"))
        })
	})
})
