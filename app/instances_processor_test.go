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
        csvWriter          *appfakes.FakeCSVWriter
        cfg                config.Config
        instancesProcessor app.InstancesProcessor
        result             []json.RawMessage
    )

    BeforeEach(func() {
        groupExporter = &appfakes.FakeGroupExporter{}
        csvWriter = &appfakes.FakeCSVWriter{}

        cfg = config.Config{
            Instances: []config.Instance{
                {
                    HostName: "foo",
                    Groups: []config.Group{
                        {
                            Name:   "foo_group",
                            Fields: []string{"id", "firstName", "lastName"},
                        },
                    },
                },
            },
        }
        
        instancesProcessor = app.NewInstancesProcessor(cfg, "token", os.TempDir())

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
            groupExporter.ExportPersonDataReturns(result, nil)
            csvWriter.WriteReturns(nil)

            instancesProcessor.Process(groupExporter, csvWriter)

            path, header, content := csvWriter.WriteArgsForCall(0)
            Expect(path).To(ContainSubstring("foo_group.csv"))
            Expect(header).To(Equal([]string{"id", "firstName", "lastName"}))
            Expect(content).To(Equal([][]string{{"1", "foo_firstname", "foo_lastname"}, {"2", "bar_firstname", "bar_lastname"}}))
        })

        It("returns an error if person data export fails", func() {
            groupExporter.ExportPersonDataReturns(nil, errors.New("boom"))

            err := instancesProcessor.Process(groupExporter, csvWriter)

            Expect(err.Error()).To(Equal("failed to get person informations: boom"))
        })

        It("returns an error if json cannot be read", func() {
            result := []json.RawMessage{json.RawMessage(`[]`)}

            groupExporter.ExportPersonDataReturns(result, nil)
            csvWriter.WriteReturns(nil)

            err := instancesProcessor.Process(groupExporter, csvWriter)

            Expect(err.Error()).To(ContainSubstring("failed to read person information raw json"))
        })
    })
})
