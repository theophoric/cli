package requirements

import (
	"fmt"

	"errors"

	"github.com/theophoric/cf-cli/cf"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/cf/terminal"
)

type TargetedSpaceRequirement struct {
	config core_config.Reader
}

func NewTargetedSpaceRequirement(config core_config.Reader) TargetedSpaceRequirement {
	return TargetedSpaceRequirement{config}
}

func (req TargetedSpaceRequirement) Execute() error {
	if !req.config.HasOrganization() {
		message := fmt.Sprintf(T("No org and space targeted, use '{{.Command}}' to target an org and space", map[string]interface{}{"Command": terminal.CommandColor(cf.Name + " target -o ORG -s SPACE")}))
		return errors.New(message)
	}

	if !req.config.HasSpace() {
		message := fmt.Sprintf(T("No space targeted, use '{{.Command}}' to target a space", map[string]interface{}{"Command": terminal.CommandColor("cf target -s")}))
		return errors.New(message)
	}

	return nil
}
