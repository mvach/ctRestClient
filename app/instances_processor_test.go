package app_test

import (
	"ctRestClient/app"
	"ctRestClient/app/appfakes"
	"ctRestClient/config"
	"ctRestClient/csv/csvfakes"
	"ctRestClient/logger/loggerfakes"
	"encoding/json"
	"errors"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func ptr(s string) *string {
    return &s
}

var _ = Describe("InstanceProcessor", func() {

	var (
		groupExporter      *appfakes.FakeGroupExporter
		csvWriter          *csvfakes.FakeCSVFileWriter
		logger             *loggerfakes.FakeLogger
		keepassCli         *appfakes.FakeKeepassCli
		personDataProvider *csvfakes.FakeFileDataProvider
		cfg                config.Config
		instancesProcessor app.InstancesProcessor
		result             []json.RawMessage
	)

	BeforeEach(func() {
		groupExporter = &appfakes.FakeGroupExporter{}
		csvWriter = &csvfakes.FakeCSVFileWriter{}
		logger = &loggerfakes.FakeLogger{}
		keepassCli = &appfakes.FakeKeepassCli{}
		personDataProvider = &csvfakes.FakeFileDataProvider{}

		cfg = config.Config{
			Instances: []config.Instance{
				{
					Hostname:  "foo",
					TokenName: "THE_TOKEN",
					Groups: []config.Group{
						{
							Name:   "foo_group",
							Fields: []config.Field{{RawString: ptr("id")}, {RawString: ptr("firstName")}, {RawString: ptr("lastName")}},
						},
					},
				},
			},
		}

		instancesProcessor = app.NewInstancesProcessor(cfg, logger)

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

		result = []json.RawMessage{json.RawMessage(person1), json.RawMessage(person2)}

		keepassCli.GetPasswordReturns("the_token", nil)
	})

	var _ = Describe("Process", func() {
		It("writes a csv", func() {
			groupExporter.ExportGroupMembersReturns(result, nil)
			csvWriter.WriteReturns(nil)

			instancesProcessor.Process(groupExporter, csvWriter, os.TempDir(), personDataProvider, keepassCli)

			path, header, content := csvWriter.WriteArgsForCall(0)
			Expect(path).To(ContainSubstring("foo_group.csv"))
			Expect(header).To(Equal([]string{"id", "firstName", "lastName"}))
			Expect(content).To(Equal([][]string{{"1", "foo_firstname", "foo_lastname"}, {"2", "bar_firstname", "bar_lastname"}}))
		})

		It("logs a warning if a token is not in the environment", func() {
			cfg = config.Config{
				Instances: []config.Instance{
					{
						Hostname:  "foo",
						TokenName: "THE_UNKNOWN_TOKEN",
						Groups: []config.Group{
							{
								Name:   "foo_group",
								Fields: []config.Field{{RawString: ptr("id")}, {RawString: ptr("firstName")}, {RawString: ptr("lastName")}},
							},
						},
					},
				},
			}

			keepassCli.GetPasswordReturns("", errors.New("booom"))

			instancesProcessor = app.NewInstancesProcessor(cfg, logger)
			instancesProcessor.Process(groupExporter, csvWriter, os.TempDir(), personDataProvider, keepassCli)

			message := logger.WarnArgsForCall(0)
			Expect(message).To(Equal("  skipping export, failed to get token with name 'THE_UNKNOWN_TOKEN' from Keepass. Err: booom"))
		})

		It("logs empty groups", func() {
			emptyGroupResult := []json.RawMessage{}
			groupExporter.ExportGroupMembersReturns(emptyGroupResult, nil)
			csvWriter.WriteReturns(nil)

			instancesProcessor.Process(groupExporter, csvWriter, os.TempDir(), personDataProvider, keepassCli)

			Expect(logger.InfoArgsForCall(2)).To(ContainSubstring("Processing instance 'foo'"))
			Expect(logger.InfoArgsForCall(5)).To(Equal("  processing group 'foo_group'"))
			Expect(logger.InfoArgsForCall(6)).To(Equal("      the group is empty"))
		})

		It("logs the group size", func() {
			groupExporter.ExportGroupMembersReturns(result, nil)
			csvWriter.WriteReturns(nil)

			instancesProcessor.Process(groupExporter, csvWriter, os.TempDir(), personDataProvider, keepassCli)

			Expect(logger.InfoArgsForCall(2)).To(ContainSubstring("Processing instance 'foo'"))
			Expect(logger.InfoArgsForCall(5)).To(Equal("  processing group 'foo_group'"))
			Expect(logger.InfoArgsForCall(6)).To(Equal("      the group has 2 persons"))
		})

		It("returns an error if person data export fails", func() {
			groupExporter.ExportGroupMembersReturns(nil, errors.New("boom"))

			instancesProcessor.Process(groupExporter, csvWriter, os.TempDir(), personDataProvider, keepassCli)

			Expect(logger.InfoArgsForCall(2)).To(ContainSubstring("Processing instance 'foo'"))
			Expect(logger.InfoArgsForCall(5)).To(Equal("  processing group 'foo_group'"))
			Expect(logger.ErrorArgsForCall(0)).To(Equal("    failed to get person information: boom"))
		})

		It("logs an error if person data cannot be read", func() {
			result := []json.RawMessage{json.RawMessage(`[]`)}

			groupExporter.ExportGroupMembersReturns(result, nil)
			csvWriter.WriteReturns(nil)

			err := instancesProcessor.Process(groupExporter, csvWriter, os.TempDir(), personDataProvider, keepassCli)
			Expect(err).ToNot(HaveOccurred())

			Expect(logger.InfoArgsForCall(2)).To(ContainSubstring("Processing instance 'foo'"))
			Expect(logger.InfoArgsForCall(5)).To(Equal("  processing group 'foo_group'"))
			Expect(logger.ErrorArgsForCall(0)).To(ContainSubstring("    failed to extract persons:"))
		})
	})
})
