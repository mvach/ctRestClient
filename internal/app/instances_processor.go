package app

import (
	"ctRestClient/internal/config"
	"ctRestClient/internal/httpclient"
	"ctRestClient/internal/logger"
	"fmt"
	"strings"
)

type InstancesProcessor interface {
	Process(
		instanceTask InstanceTask,
	) error
}

type instancesProcessor struct {
	instances  []config.Instance
	keepassCli KeepassCli
	logger     logger.Logger
}

func NewInstancesProcessor(
	instances []config.Instance,
	keepassCli KeepassCli,
	logger logger.Logger,

) InstancesProcessor {
	return instancesProcessor{
		instances:  instances,
		keepassCli: keepassCli,
		logger:     logger,
	}
}

func (p instancesProcessor) Process(
	instanceTask InstanceTask,
) error {
	for _, instance := range p.instances {

		p.logTitle(instance)

		token, err := p.keepassCli.GetPassword(instance.TokenName)
		if err != nil {
			p.logger.Warn(fmt.Sprintf("  skipping export, failed to get token with name '%s' from Keepass. Err: %v", instance.TokenName, err))
			continue
		}

		httpClient := httpclient.NewHTTPClient(instance.Hostname, token)

		instanceTask.Execute(instance, httpClient)
	}

	return nil
}

func (p instancesProcessor) logTitle(instance config.Instance) {
	boxLength := 70
	title := fmt.Sprintf("Processing instance '%s'", instance.Hostname)
	titleLength := len(title)
	border := strings.Repeat("-", boxLength)

	p.logger.Info("")
	p.logger.Info(fmt.Sprintf("+%s+", border))
	p.logger.Info(fmt.Sprintf("| %s "+strings.Repeat(" ", boxLength-titleLength-2)+"|", title))
	p.logger.Info(fmt.Sprintf("+%s+", border))
}
