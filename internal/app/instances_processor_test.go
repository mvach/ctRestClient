package app_test

import (
	"ctRestClient/internal/app"
	"ctRestClient/internal/app/appfakes"
	"ctRestClient/internal/config"
	"ctRestClient/internal/logger/loggerfakes"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("InstanceProcessor", func() {

	var (
		instanceTask       *appfakes.FakeInstanceTask
		logger             *loggerfakes.FakeLogger
		keepassCli         *appfakes.FakeKeepassCli
		instancesProcessor app.InstancesProcessor
	)

	BeforeEach(func() {
		instanceTask = &appfakes.FakeInstanceTask{}
		logger = &loggerfakes.FakeLogger{}
		keepassCli = &appfakes.FakeKeepassCli{}
	})

	var _ = Describe("Process", func() {
		It("invokes the provided task", func() {
			instances := []config.Instance{
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

			keepassCli.GetPasswordReturns("the_token", nil)

			instancesProcessor = app.NewInstancesProcessor(instances, keepassCli, logger)
			err := instancesProcessor.Process(instanceTask)
			Expect(err).NotTo(HaveOccurred())

			Expect(instanceTask.ExecuteCallCount()).To(Equal(1))
			Expect(logger.InfoArgsForCall(2)).To(ContainSubstring("Processing instance 'foo'"))
		})

		It("logs a warning if a token is not in the environment", func() {
			instances := []config.Instance{
				{
					Hostname:  "foo",
					TokenName: "THE_UNKNOWN_TOKEN",
					Groups: []config.Group{
						{
							Name:   "foo_group",
							Fields: []config.Field{{FieldName: ptr("id")}, {FieldName: ptr("firstName")}, {FieldName: ptr("lastName")}},
						},
					},
				},
			}

			keepassCli.GetPasswordReturns("", errors.New("booom"))

			instancesProcessor = app.NewInstancesProcessor(instances, keepassCli, logger)
			err := instancesProcessor.Process(instanceTask)
			Expect(err).NotTo(HaveOccurred())

			message := logger.WarnArgsForCall(0)
			Expect(message).To(Equal("  skipping export, failed to get token with name 'THE_UNKNOWN_TOKEN' from Keepass. Err: booom"))
		})
	})
})
