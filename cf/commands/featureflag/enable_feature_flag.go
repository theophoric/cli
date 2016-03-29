package featureflag

import (
	"github.com/theophoric/cf-cli/cf/api/feature_flags"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/flags"
)

type EnableFeatureFlag struct {
	ui       terminal.UI
	config   core_config.ReadWriter
	flagRepo feature_flags.FeatureFlagRepository
}

func init() {
	command_registry.Register(&EnableFeatureFlag{})
}

func (cmd *EnableFeatureFlag) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "enable-feature-flag",
		Description: T("Enable the use of a feature so that users have access to and can use the feature"),
		Usage: []string{
			T("CF_NAME enable-feature-flag FEATURE_NAME"),
		},
	}
}

func (cmd *EnableFeatureFlag) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	if len(fc.Args()) != 1 {
		cmd.ui.Failed(T("Incorrect Usage. Requires an argument\n\n") + command_registry.Commands.CommandUsage("enable-feature-flag"))
	}

	reqs := []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
	}

	return reqs
}

func (cmd *EnableFeatureFlag) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.flagRepo = deps.RepoLocator.GetFeatureFlagRepository()
	return cmd
}

func (cmd *EnableFeatureFlag) Execute(c flags.FlagContext) {
	flag := c.Args()[0]

	cmd.ui.Say(T("Setting status of {{.FeatureFlag}} as {{.Username}}...", map[string]interface{}{
		"FeatureFlag": terminal.EntityNameColor(flag),
		"Username":    terminal.EntityNameColor(cmd.config.Username())}))

	apiErr := cmd.flagRepo.Update(flag, true)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
	}

	cmd.ui.Say("")
	cmd.ui.Ok()
	cmd.ui.Say("")
	cmd.ui.Say(T("Feature {{.FeatureFlag}} Enabled.", map[string]interface{}{
		"FeatureFlag": terminal.EntityNameColor(flag)}))
	return
}
