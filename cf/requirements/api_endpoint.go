package requirements

import (
	"fmt"

	"errors"
	"github.com/theophoric/cf-cli/cf"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/cf/terminal"
)

type ApiEndpointRequirement struct {
	config core_config.Reader
}

func NewApiEndpointRequirement(config core_config.Reader) ApiEndpointRequirement {
	return ApiEndpointRequirement{config}
}

func (req ApiEndpointRequirement) Execute() error {
	if req.config.ApiEndpoint() == "" {
		loginTip := terminal.CommandColor(fmt.Sprintf(T("{{.CFName}} login", map[string]interface{}{"CFName": cf.Name})))
		apiTip := terminal.CommandColor(fmt.Sprintf(T("{{.CFName}} api", map[string]interface{}{"CFName": cf.Name})))
		return errors.New(T("No API endpoint set. Use '{{.LoginTip}}' or '{{.APITip}}' to target an endpoint.",
			map[string]interface{}{
				"LoginTip": loginTip,
				"APITip":   apiTip,
			}))
	}

	return nil
}
