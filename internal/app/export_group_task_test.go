package app_test

import (
	"ctRestClient/internal/app"
	"ctRestClient/internal/app/appfakes"
	"ctRestClient/internal/config"
	"ctRestClient/internal/csv/csvfakes"
	"ctRestClient/internal/dataprovider/dataproviderfakes"
	"ctRestClient/internal/logger/loggerfakes"
	"encoding/json"
	"os"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("InstanceProcessor", func() {

	var (
		groupExporter          *appfakes.FakeGroupExporter
		csvWriter              *csvfakes.FakeCSVFileWriter
		logger                 *loggerfakes.FakeLogger
		personDataProvider     *dataproviderfakes.FakeFileDataProvider
		blocklistsDataProvider *dataproviderfakes.FakeBlockListDataProvider
		instances              []config.Instance
		result                 []json.RawMessage
	)

	BeforeEach(func() {
		groupExporter = &appfakes.FakeGroupExporter{}
		csvWriter = &csvfakes.FakeCSVFileWriter{}
		logger = &loggerfakes.FakeLogger{}
		personDataProvider = &dataproviderfakes.FakeFileDataProvider{}
		blocklistsDataProvider = &dataproviderfakes.FakeBlockListDataProvider{}

		instances = []config.Instance{
			{
				Hostname:  "foo",
				TokenName: "THE_TOKEN",
				Groups: []config.Group{
					{
						Name:   "foo_group",
						Fields: []config.Field{{FieldName: ptr("id")}, {FieldName: ptr("firstName")}, {FieldName: ptr("lastName")}},
					},
				},
			},
		}

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
	})

	var _ = Describe("Process", func() {
		It("writes a csv", func() {
			groupExporter.ExportGroupMembersReturns(result, nil)
			csvWriter.WriteReturns(nil)

			exportGroupTask := app.NewExportGroupTask(groupExporter, csvWriter, os.TempDir(), personDataProvider, blocklistsDataProvider, logger)
			err := exportGroupTask.Execute(instances[0], nil)
			Expect(err).NotTo(HaveOccurred())

			path, header, content := csvWriter.WriteArgsForCall(0)
			Expect(path).To(ContainSubstring("foo_group.csv"))
			Expect(header).To(Equal([]string{"id", "firstName", "lastName"}))
			Expect(content).To(Equal([][]string{{"1", "foo_firstname", "foo_lastname"}, {"2", "bar_firstname", "bar_lastname"}}))
		})

		It("logs empty groups", func() {
			emptyGroupResult := []json.RawMessage{}
			groupExporter.ExportGroupMembersReturns(emptyGroupResult, nil)
			csvWriter.WriteReturns(nil)

			exportGroupTask := app.NewExportGroupTask(groupExporter, csvWriter, os.TempDir(), personDataProvider, blocklistsDataProvider, logger)
			err := exportGroupTask.Execute(instances[0], nil)
			Expect(err).NotTo(HaveOccurred())

			Expect(logger.InfoArgsForCall(1)).To(Equal("  processing group 'foo_group'"))
			Expect(logger.InfoArgsForCall(2)).To(Equal("      the group is empty"))
		})

		It("logs the group size", func() {
			groupExporter.ExportGroupMembersReturns(result, nil)
			csvWriter.WriteReturns(nil)

			exportGroupTask := app.NewExportGroupTask(groupExporter, csvWriter, os.TempDir(), personDataProvider, blocklistsDataProvider, logger)
			err := exportGroupTask.Execute(instances[0], nil)
			Expect(err).NotTo(HaveOccurred())

			Expect(logger.InfoArgsForCall(1)).To(Equal("  processing group 'foo_group'"))
			Expect(logger.InfoArgsForCall(2)).To(Equal("      the group has 2 persons"))
		})

		It("logs a warning for not active groups", func() {
			groupExporter.ExportGroupMembersReturns(nil, &app.GroupNotActiveError{GroupName: "foo_group"})

			exportGroupTask := app.NewExportGroupTask(groupExporter, csvWriter, os.TempDir(), personDataProvider, blocklistsDataProvider, logger)
			err := exportGroupTask.Execute(instances[0], nil)
			Expect(err).NotTo(HaveOccurred())

			Expect(logger.InfoArgsForCall(1)).To(Equal("  processing group 'foo_group'"))
			Expect(logger.WarnArgsForCall(0)).To(Equal("      skipping csv creation since the group is not active"))
		})

		It("returns an error if person data export fails", func() {
			groupExporter.ExportGroupMembersReturns(nil, errors.New("boom"))

			exportGroupTask := app.NewExportGroupTask(groupExporter, csvWriter, os.TempDir(), personDataProvider, blocklistsDataProvider, logger)
			err := exportGroupTask.Execute(instances[0], nil)
			Expect(err).NotTo(HaveOccurred())

			Expect(logger.InfoArgsForCall(1)).To(Equal("  processing group 'foo_group'"))
			Expect(logger.ErrorArgsForCall(0)).To(Equal("      failed to get person information: boom"))
		})

		It("logs an error if person data cannot be read", func() {
			result := []json.RawMessage{json.RawMessage(`[]`)}

			groupExporter.ExportGroupMembersReturns(result, nil)
			csvWriter.WriteReturns(nil)

			exportGroupTask := app.NewExportGroupTask(groupExporter, csvWriter, os.TempDir(), personDataProvider, blocklistsDataProvider, logger)
			err := exportGroupTask.Execute(instances[0], nil)
			Expect(err).NotTo(HaveOccurred())

			Expect(logger.InfoArgsForCall(1)).To(Equal("  processing group 'foo_group'"))
			Expect(logger.ErrorArgsForCall(0)).To(ContainSubstring("    failed to extract persons:"))
		})
	})
})
