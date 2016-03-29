package cloud_controller_gateway

import (
	"time"

	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/net"
	"github.com/theophoric/cf-cli/cf/trace/fakes"
	testterm "github.com/theophoric/cf-cli/testhelpers/terminal"
)

func NewTestCloudControllerGateway(configRepo core_config.Reader) net.Gateway {
	fakeLogger := new(fakes.FakePrinter)
	return net.NewCloudControllerGateway(configRepo, time.Now, &testterm.FakeUI{}, fakeLogger)
}
