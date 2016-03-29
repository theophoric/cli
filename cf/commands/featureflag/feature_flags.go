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

type ListFeatureFlags struct {
	ui       terminal.UI
	config   core_config.ReadWriter
	flagRepo feature_flags.FeatureFlagRepository
}

func init() {
	command_registry.Register(&ListFeatureFlags{})
}

func (cmd *ListFeatureFlags) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "feature-flags",
		Description: T("Retrieve list of feature flags with status of each flag-able feature"),
		Usage: []string{
			T("CF_NAME feature-flags"),
		},
	}
}

func (cmd *ListFeatureFlags) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	usageReq := requirements.NewUsageRequirement(command_registry.CliCommandUsagePresenter(cmd),
		T("No argument required"),
		func() bool {
			return len(fc.Args()) != 0
		},
	)

	reqs := []requirements.Requirement{
		usageReq,
		requirementsFactory.NewLoginRequirement(),
	}

	return reqs
}

func (cmd *ListFeatureFlags) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.flagRepo = deps.RepoLocator.GetFeatureFlagRepository()
	return cmd
}

func (cmd *ListFeatureFlags) Execute(c flags.FlagContext) {
	cmd.ui.Say(T("Retrieving status of all flagged features as {{.Username}}...", map[string]interface{}{
		"Username": terminal.EntityNameColor(cmd.config.Username())}))

	flags, err := cmd.flagRepo.List()
	if err != nil {
		cmd.ui.Failed(err.Error())
		return
	}

	cmd.ui.Ok()
	cmd.ui.Say("")

	table := terminal.NewTable(cmd.ui, []string{T("Features"), T("State")})

	for _, flag := range flags {
		table.Add(
			flag.Name,
			cmd.flagBoolToString(flag.Enabled),
		)
	}

	table.Print()
	return
}

func (cmd ListFeatureFlags) flagBoolToString(enabled bool) string {
	if enabled {
		return "enabled"
	}
	return "disabled"
}
