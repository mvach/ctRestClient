package app

import (
	"ctRestClient/internal/config"
	"ctRestClient/internal/httpclient"
)

//counterfeiter:generate . InstanceTask
type InstanceTask interface {
	Execute(
		instance config.Instance,
		httpClient httpclient.HTTPClient,
	) error
}
