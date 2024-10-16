package app_test

import (
    "ctRestClient/app"
    "ctRestClient/app/appfakes"
    "ctRestClient/config"
    "encoding/json"
    "errors"
    "os"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

var _ = Describe("InstanceProcessor", func() {

    var (
        groupExporter      *appfakes.FakeGroupExporter
        csvWriter          *appfakes.FakeCSVFileWriter
        logger             *appfakes.FakeLogger
        cfg                config.Config
        instancesProcessor app.InstancesProcessor
        result             []json.RawMessage
    )

    BeforeEach(func() {
        groupExporter = &appfakes.FakeGroupExporter{}
        csvWriter = &appfakes.FakeCSVFileWriter{}
        logger = &appfakes.FakeLogger{}

        cfg = config.Config{
            Instances: []config.Instance{
                {
                    Hostname:  "foo",
                    TokenName: "THE_TOKEN",
                    Groups: []config.Group{
                        {
                            Name:   "foo_group",
                            Fields: []string{"id", "firstName", "lastName"},
                        },
                    },
                },
            },
        }

        groupExporter.GetGroupNames2IDMappingReturns(
            map[string]int{"foo_group": 1, "bar_group": 2}, nil,
        )

        instancesProcessor = app.NewInstancesProcessor(cfg, os.TempDir(), logger)

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

        os.Setenv("THE_TOKEN", "the_token")
    })

    AfterEach(func() {
        os.Setenv("THE_TOKEN", "")
    })

    var _ = Describe("Process", func() {
        It("writes a csv", func() {
            groupExporter.ExportGroupMembersReturns(result, nil)
            csvWriter.WriteReturns(nil)

            instancesProcessor.Process(groupExporter, csvWriter)

            path, header, content := csvWriter.WriteArgsForCall(0)
            Expect(path).To(ContainSubstring("foo_group.csv"))
            Expect(header).To(Equal([]string{"id", "firstName", "lastName"}))
            Expect(content).To(Equal([][]string{{"1", "foo_firstname", "foo_lastname"}, {"2", "bar_firstname", "bar_lastname"}}))
        })

        It("returns an error if group name2id mappinfg cannot be created", func() {
            groupExporter.GetGroupNames2IDMappingReturns(nil, errors.New("boom"))

            instancesProcessor = app.NewInstancesProcessor(cfg, os.TempDir(), logger)
            err := instancesProcessor.Process(groupExporter, csvWriter)

            Expect(err.Error()).To(Equal("failed to resolve groupnames to ids, boom"))
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
                                Fields: []string{"id", "firstName", "lastName"},
                            },
                        },
                    },
                },
            }

            instancesProcessor = app.NewInstancesProcessor(cfg, os.TempDir(), logger)
            instancesProcessor.Process(groupExporter, csvWriter)

            message := logger.WarnArgsForCall(0)
            Expect(message).To(Equal("  skipping export, a token with name 'THE_UNKNOWN_TOKEN' is not set in the environment"))
        })

        It("logs an error if groupname cannot be found in the name2ID mapping", func() {
            cfg = config.Config{
                Instances: []config.Instance{
                    {
                        Hostname:  "foo",
                        TokenName: "THE_TOKEN",
                        Groups: []config.Group{
                            {
                                Name:   "missing_group",
                                Fields: []string{"id", "firstName", "lastName"},
                            },
                        },
                    },
                },
            }

            emptyGroupResult := []json.RawMessage{}
            groupExporter.ExportGroupMembersReturns(emptyGroupResult, nil)
            csvWriter.WriteReturns(nil)

            instancesProcessor = app.NewInstancesProcessor(cfg, os.TempDir(), logger)
            instancesProcessor.Process(groupExporter, csvWriter)

            Expect(logger.InfoArgsForCall(0)).To(Equal("processing instance 'foo'"))
            Expect(logger.InfoArgsForCall(1)).To(Equal("  processing group 'missing_group'"))
            Expect(logger.ErrorArgsForCall(0)).To(Equal("    could not find group to id mapping"))
        })

        It("logs empty groups", func() {
            emptyGroupResult := []json.RawMessage{}
            groupExporter.ExportGroupMembersReturns(emptyGroupResult, nil)
            csvWriter.WriteReturns(nil)

            instancesProcessor.Process(groupExporter, csvWriter)

            Expect(logger.InfoArgsForCall(0)).To(Equal("processing instance 'foo'"))
            Expect(logger.InfoArgsForCall(1)).To(Equal("  processing group 'foo_group'"))
            Expect(logger.InfoArgsForCall(2)).To(Equal("    the group is empty"))
        })

        It("logs the group size", func() {
            groupExporter.ExportGroupMembersReturns(result, nil)
            csvWriter.WriteReturns(nil)

            instancesProcessor.Process(groupExporter, csvWriter)

            Expect(logger.InfoArgsForCall(0)).To(Equal("processing instance 'foo'"))
            Expect(logger.InfoArgsForCall(1)).To(Equal("  processing group 'foo_group'"))
            Expect(logger.InfoArgsForCall(2)).To(Equal("    the group has 2 persons"))
        })

        It("returns an error if person data export fails", func() {
            groupExporter.ExportGroupMembersReturns(nil, errors.New("boom"))

            err := instancesProcessor.Process(groupExporter, csvWriter)

            Expect(err.Error()).To(Equal("failed to get person informations: boom"))
        })

        It("returns an error if json cannot be read", func() {
            result := []json.RawMessage{json.RawMessage(`[]`)}

            groupExporter.ExportGroupMembersReturns(result, nil)
            csvWriter.WriteReturns(nil)

            err := instancesProcessor.Process(groupExporter, csvWriter)

            Expect(err.Error()).To(ContainSubstring("failed to read person information raw json"))
        })
    })
})
